package rabbitmq

import (
	"context"
	"os"
	"time"

	"github.com/pkg/errors"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/shahzodshafizod/gocloud/pkg"
)

type queue struct {
	amqpurl        string // AMQP Server URI
	connection     *amqp.Connection
	channel        *amqp.Channel
	exchange       string // Durable, non-auto-deleted AMQP exchange name
	routingKey     string // AMQP routing key
	subscribers    []*pkg.Subscribe
	doneCh         chan bool
	closedCh       chan *amqp.Error
	tracer         pkg.Tracer
	reconnectDelay time.Duration
}

func NewQueue(serviceName string, tracer pkg.Tracer) (pkg.Queue, error) {

	queue := &queue{
		amqpurl:        os.Getenv("RABBITMQ_AMQP_URL"),
		exchange:       os.Getenv("RABBITMQ_EXCHANGE"),
		routingKey:     os.Getenv("RABBITMQ_ROUTING_KEY"),
		subscribers:    make([]*pkg.Subscribe, 0),
		doneCh:         make(chan bool),
		tracer:         tracer,
		reconnectDelay: time.Second * 2,
	}

	err := queue.connect()
	if err != nil {
		return nil, errors.Wrap(err, "c.connect")
	}

	go queue.reconnect()

	return queue, nil
}

func (q *queue) Close() error {
	q.doneCh <- true
	q.close()
	return nil
}

func (q *queue) connect() error {
	// Create a new RabbitMQ connection.
	var err error
	q.connection, err = amqp.Dial(q.amqpurl)
	if err != nil {
		return errors.Wrap(err, "amqp.Dial")
	}

	// Opening a channel to our RabbitMQ instance over the connection we have already established.
	q.channel, err = q.connection.Channel()
	if err != nil {
		return errors.Wrap(err, "q.connection.Channel")
	}

	// got Channel, declaring Exchange {exchange}
	err = q.channel.ExchangeDeclare(
		q.exchange,         // name
		amqp.ExchangeTopic, // type (kind)
		true,               // durable
		false,              // auto-delete, delete when complete
		false,              // internal
		false,              // noWait
		nil,                // arguments
	)
	if err != nil {
		return errors.Wrap(err, "q.channel.ExchangeDeclare")
	}

	q.closedCh = make(chan *amqp.Error)
	q.connection.NotifyClose(q.closedCh)

	return nil
}

func (q *queue) reconnect() {
	ctx, span := q.tracer.StartFromContext(context.Background(), "queue.reconnect")
	defer span.End()
	for {
		select {
		case <-q.doneCh:
			return
		case <-q.closedCh:
			if q.connection == nil || q.connection.IsClosed() {
				q.close()
				time.Sleep(q.reconnectDelay)
				err := q.connect()
				if err != nil {
					span.RecordError(errors.Wrap(err, "queue.reconnect q.connect"))
				} else {
					for _, subscriber := range q.subscribers {
						q.subscribe(ctx, subscriber)
					}
				}
			}
		}
	}
}

func (q *queue) close() {
	_, span := q.tracer.StartFromContext(context.Background(), "queue.close")
	defer span.End()
	if q.channel != nil && !q.channel.IsClosed() {
		// will close() the deliveries channel
		for _, subscr := range q.subscribers {
			err := q.channel.Cancel(subscr.Topic, true)
			if err != nil {
				span.RecordError(errors.Wrap(err, "queue.close q.channel.Cancel "+subscr.Topic))
			}
		}
		err := q.channel.Close()
		if err != nil {
			span.RecordError(errors.Wrap(err, "queue.close q.channel.Close"))
		}
	}
	if q.connection != nil && !q.connection.IsClosed() {
		err := q.connection.Close()
		if err != nil {
			span.RecordError(errors.Wrap(err, "queue.close q.connection.Close"))
		}
	}
}

func (q *queue) Subscribe(ctx context.Context, subscr *pkg.Subscribe) error {
	err := q.subscribe(ctx, subscr)
	if err != nil {
		return errors.Wrap(err, "q.subscribe")
	}
	q.subscribers = append(q.subscribers, subscr)
	return nil
}

func (q *queue) Publish(ctx context.Context, queuename string, body []byte) error {

	// // Explanation of commenting: publishers don't own queues
	// err := q.declareAndBindQueue(subject)
	// if err != nil {
	// 	return errors.Wrap(err, "q.declareAndBindQueue")
	// }

	msg := amqp.Publishing{
		Headers:         amqp.Table{},
		ContentType:     "application/json",
		ContentEncoding: "",
		DeliveryMode:    amqp.Persistent,
		Priority:        0,
		AppId:           "",
		Body:            body,
	}

	q.tracer.Inject(ctx, mapCarrier(msg.Headers))
	err := q.channel.PublishWithContext(
		ctx,
		"",        // amqp.ExchangeDirect,
		queuename, // q.routingKey,
		false,
		false,
		msg,
	)
	if err != nil {
		return errors.Wrap(err, "q.channel.PublishWithContext")
	}

	return nil
}

func (q *queue) Request(ctx context.Context, queuename string, request []byte) ([]byte, error) {
	delivery, err := q.channel.ConsumeWithContext(
		ctx,
		"amq.rabbitmq.reply-to",
		"amq.rabbitmq.reply-to",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, errors.Wrap(err, "q.channel.ConsumeWithContext")
	}

	pubMsg := amqp.Publishing{
		Headers:      amqp.Table{},
		ContentType:  "application/json",
		DeliveryMode: amqp.Persistent,
		Priority:     0,
		Body:         request,
		ReplyTo:      "amq.rabbitmq.reply-to",
	}

	q.tracer.Inject(ctx, mapCarrier(pubMsg.Headers))
	err = q.channel.PublishWithContext(
		ctx,
		q.exchange,
		queuename,
		false,
		false,
		pubMsg,
	)
	if err != nil {
		return nil, errors.Wrap(err, "q.channel.PublishWithContext")
	}

	msg := <-delivery
	return msg.Body, nil
}

func (q *queue) subscribe(ctx context.Context, subscr *pkg.Subscribe) error {
	err := q.declareAndBindQueue(subscr.Topic)
	if err != nil {
		return errors.Wrap(err, "q.declareAndBindQueue")
	}

	go q.consume(ctx, subscr)
	return nil
}

func (q *queue) declareAndBindQueue(queueName string) error {
	// declared Exchange, declaring Queue {queueName}
	_, err := q.channel.QueueDeclare(
		queueName, // name of the queue
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // noWait
		nil,       // arguments
	)
	if err != nil {
		return errors.Wrap(err, "q.channel.QueueDeclare")
	}

	// declared Queue ({queueName}, binding to Exchange (key {routingkey})
	return q.channel.QueueBind(
		queueName,  // name of the queue
		queueName,  // c.routingKey, // bindingKey
		q.exchange, // sourceExchange
		false,      // noWait
		nil,        // arguments
	)
}

func (q *queue) consume(ctx context.Context, subscr *pkg.Subscribe) {
	_, span := q.tracer.StartFromContext(ctx, "queue.consume")
	defer span.End()
	// Queue bound to Exchange, starting Consume
	delivery, err := q.channel.ConsumeWithContext(
		ctx,
		subscr.Topic, // queue name
		subscr.Topic, // The consumer identity will be included in every Delivery in the ConsumerTag field
		false,        // autoAck
		false,        // exclusive
		false,        // noLocal
		false,        // noWait
		nil,          // arguments
	)
	if err != nil {
		span.RecordError(errors.Wrap(err, "queue.consume q.channel.ConsumeWithContext"))
		return
	}

	for d := range delivery {
		err = subscr.Callback(ctx,
			q.tracer.Extract(ctx, mapCarrier(d.Headers)),
			&message{d.Body, d.ReplyTo},
		)
		if err != nil {
			span.RecordError(errors.Wrap(err, "queue.consume callback"))
		}
		d.Ack(false)
	}
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

type mapCarrier amqp.Table

var _ pkg.TextMapCarrier = mapCarrier{}

func (m mapCarrier) Get(key string) string {
	if value, ok := m[key].(string); ok {
		return value
	}
	return ""
}

func (m mapCarrier) Set(key string, value string) {
	m[key] = value
}

func (m mapCarrier) Keys() []string {
	var keys = make([]string, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	return keys
}
