package jetstream

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/pkg/errors"
	"github.com/shahzodshafizod/gocloud/pkg"
)

type queue struct {
	nc          *nats.Conn
	js          jetstream.JetStream
	stream      jetstream.Stream
	streamname  string
	subjects    []string
	durableName string
	doneCh      chan bool
	subscribers []*pkg.Subscribe
	tracer      pkg.Tracer
}

func NewQueue(serviceName string, tracer pkg.Tracer) (pkg.Queue, error) {
	queue := &queue{
		streamname: os.Getenv("JETSTREAM_STREAM_NAME"),
		subjects: []string{
			os.Getenv("JETSTREAM_ORDERS_TOPIC"),
			os.Getenv("JETSTREAM_PRODUCTS_TOPIC"),
			os.Getenv("JETSTREAM_PARTNERS_TOPIC"),
		},
		durableName: os.Getenv("JETSTREAM_DURABLE"),
		doneCh:      make(chan bool),
		subscribers: make([]*pkg.Subscribe, 0),
		tracer:      tracer,
	}

	// Connect Options.
	opts, err := queue.setupConnOptions(serviceName)
	if err != nil {
		return nil, errors.Wrap(err, "queue.setupConnOptions")
	}

	// Connect to NATS
	queue.nc, err = nats.Connect(os.Getenv("JETSTREAM_NATS_URL"), opts...)
	if err != nil {
		return nil, errors.Wrap(err, "nats.Connect")
	}

	err = queue.connect()
	if err != nil {
		return nil, errors.Wrap(err, "queue.connect")
	}

	return queue, nil
}

func (q *queue) Close() error {
	q.closeGoroutines()
	if q.nc != nil {
		q.nc.Drain()
	}
	return nil
}

func (q *queue) connect() error {
	ctx, span := q.tracer.StartFromContext(context.Background(), "queue.connect")
	var err error
	q.js, err = jetstream.New(q.nc)
	if err != nil {
		return errors.Wrap(err, "jetstream.New")
	}

	q.stream, _ = q.js.Stream(ctx, q.streamname)
	if q.stream == nil {
		// stream not found, create it
		span.RecordError(errors.New("stream not found, create it: " + q.streamname))

		var err error
		q.stream, err = q.js.CreateStream(ctx, jetstream.StreamConfig{
			Name:      q.streamname,
			Subjects:  q.subjects,
			Storage:   jetstream.FileStorage,
			Replicas:  1,
			Retention: jetstream.InterestPolicy, // jetstream.WorkQueuePolicy,
		})
		if err != nil {
			return errors.Wrap(err, "q.js.CreateStream")
		}
	}

	for _, subscr := range q.subscribers {
		go q.subscribe(context.Background(), subscr)
	}

	return nil
}

func (q *queue) closeGoroutines() {
	q.doneCh <- true
	time.Sleep(time.Millisecond)
	<-q.doneCh
}

func (q *queue) Subscribe(ctx context.Context, subscr *pkg.Subscribe) error {
	go q.subscribe(ctx, subscr)
	q.subscribers = append(q.subscribers, subscr)
	return nil
}

func (q *queue) subscribe(ctx context.Context, subscr *pkg.Subscribe) {
	ctx, span := q.tracer.StartFromContext(ctx, "queue.subscribe")
	defer span.End()
	consumer, err := q.stream.CreateConsumer(ctx, jetstream.ConsumerConfig{
		// Durable:       c.durableName,
		FilterSubject: subscr.Topic,
		AckPolicy:     jetstream.AckExplicitPolicy,
	})
	if err != nil {
		span.RecordError(errors.Wrap(err, "q.stream.CreateConsumer"))
		return
	}
	cc, err := consumer.Consume(func(m jetstream.Msg) {
		m.Ack()
		err = subscr.Callback(ctx,
			q.tracer.Extract(ctx, mapCarrier(m.Headers())),
			&message{m.Data(), m.Reply()},
		)
		if err != nil {
			span.RecordError(errors.Wrap(err, "subscr.Callback"))
		}
	})
	if err != nil {
		span.RecordError(errors.Wrap(err, "consumer.Consume"))
		return
	}
	defer cc.Stop()
	<-q.doneCh
	q.doneCh <- true
}

func (q *queue) Publish(ctx context.Context, subject string, msgData []byte) error {
	q.tracer.Inject(ctx, mapCarrier(nats.Header{}))
	_, err := q.js.Publish(ctx, subject, msgData)
	if err != nil {
		return errors.Wrap(err, "q.js.Publish")
	}
	return nil
}

func (q *queue) Request(ctx context.Context, subject string, request []byte) ([]byte, error) {
	q.tracer.Inject(ctx, mapCarrier(nats.Header{}))
	msg, err := q.nc.RequestWithContext(ctx, subject, request)
	if err != nil {
		return nil, errors.Wrap(err, "q.nc.RequestWithContext")
	}
	return msg.Data, nil
}

func (q *queue) setupConnOptions(connectionName string) ([]nats.Option, error) {
	const (
		totalWait      = 10 * time.Minute
		reconnectDelay = time.Second * 2
	)

	_, span := q.tracer.StartFromContext(context.Background(), "queue.setupConnOptions")
	defer span.End()

	var userCreds = ""     // User Credentials File
	var nkeyFile = ""      // NKey Seed File
	var tlsClientCert = "" // TLS client certificate file
	var tlsClientKey = ""  // Private key file for client certificate
	var tlsCACert = ""     // CA certificate to verify peer against

	opts := []nats.Option{nats.Name(connectionName)}
	opts = append(opts, nats.ReconnectWait(reconnectDelay))
	opts = append(opts, nats.MaxReconnects(int(totalWait/reconnectDelay)))
	opts = append(opts, nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
		span.RecordError(fmt.Errorf("disconnected due to:%s, will attempt reconnects for %.0fm", err, totalWait.Minutes()))
	}))
	opts = append(opts, nats.ReconnectHandler(func(nc *nats.Conn) {
		q.closeGoroutines()
		err := q.connect()
		if err != nil {
			span.RecordError(errors.Wrap(err, "nats.ReconnectHandler c.connect"))
		}
		span.RecordError(errors.New("Reconnected: " + nc.ConnectedUrl()))
	}))
	opts = append(opts, nats.ClosedHandler(func(nc *nats.Conn) {
		span.RecordError(errors.Wrap(nc.LastError(), "nats.ClosedHandler"))
	}))
	if userCreds != "" && nkeyFile != "" {
		return nil, errors.New("specify -seed or -creds")
	}
	// Use UserCredentials
	if userCreds != "" {
		opts = append(opts, nats.UserCredentials(userCreds))
	}
	// Use TLS client authentication
	if tlsClientCert != "" && tlsClientKey != "" {
		opts = append(opts, nats.ClientCert(tlsClientCert, tlsClientKey))
	}
	// Use specific CA certificate
	if tlsCACert != "" {
		opts = append(opts, nats.RootCAs(tlsCACert))
	}
	// Use Nkey authentication.
	if nkeyFile != "" {
		opt, err := nats.NkeyOptionFromSeed(nkeyFile)
		if err != nil {
			return nil, errors.Wrap(err, "nats.NkeyOptionFromSeed")
		}
		opts = append(opts, opt)
	}
	return opts, nil
}

type message struct {
	body    []byte
	replyTo string
}

func (m *message) Body() []byte {
	return m.body
}

func (m *message) ReplyTo() string {
	return m.replyTo
}

type mapCarrier nats.Header

var _ pkg.TextMapCarrier = mapCarrier{}

func (m mapCarrier) Get(key string) string {
	if value, ok := m[key]; ok && len(value) > 0 {
		return value[0]
	}
	return ""
}

func (m mapCarrier) Set(key string, value string) {
	m[key] = []string{value}
}

func (m mapCarrier) Keys() []string {
	var keys = make([]string, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	return keys
}
