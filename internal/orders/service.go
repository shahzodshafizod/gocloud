package orders

import (
	"context"
	"encoding/json"
	"os"

	"github.com/pkg/errors"
	"github.com/shahzodshafizod/gocloud/internal/notifications"
	"github.com/shahzodshafizod/gocloud/pkg"
)

type Service interface {
	createOrder(context.Context, *Order) (*CreateResponse, error)
	payOrder(context.Context, *PayRequest) (*PayResponse, error)
	pickupOrder(context.Context, *pickupRequest) error
	assignOrder(ctx context.Context, req *AssignRequest) (*AssignResponse, error)
	// updateOrderAddress(orderID int64, updatedAddress *UpdatedAddress) error
}

type service struct {
	repository          Repository
	queue               pkg.Queue
	paymentCallbackURL  string
	readyCallbackURL    string
	notificationAgentID string
}

func NewService(repository Repository, queue pkg.Queue) Service {
	return &service{
		repository:          repository,
		queue:               queue,
		paymentCallbackURL:  os.Getenv("PAYMENT_CALLBACK_URL"),
		readyCallbackURL:    os.Getenv("READY_CALLBACK_URL"),
		notificationAgentID: os.Getenv("NOTIFICATION_AGENT_ID"),
	}
}

func (s *service) createOrder(ctx context.Context, order *Order) (*CreateResponse, error) {
	bank, err := s.repository.getBank(ctx, order.Paytype)
	if err != nil {
		return nil, errors.Wrap(err, "s.repository.getBank")
	}

	id, err := s.repository.createOrder(ctx, order)
	if err != nil {
		return nil, errors.Wrap(err, "s.repository.createOrder")
	}

	resp := &CreateResponse{
		OrderID:        id,
		WebcheckoutURL: bank.WebcheckoutURL,
		CallbackURL:    s.paymentCallbackURL,
	}
	return resp, nil
}

func (s *service) payOrder(ctx context.Context, req *PayRequest) (*PayResponse, error) {
	id, order, err := s.repository.payOrder(ctx, req)
	if err != nil {
		return nil, errors.Wrap(err, "s.repository.payOrder")
	}
	order.CallbackURL = s.readyCallbackURL

	data, err := json.Marshal(order)
	if err != nil {
		return nil, errors.Wrap(err, "json.Marshal")
	}
	err = s.queue.Publish(ctx, "orders.paid", data)
	if err != nil {
		return nil, errors.Wrap(err, "s.queue.Publish")
	}

	return &PayResponse{PaymentID: id}, nil
}

func (s *service) pickupOrder(ctx context.Context, req *pickupRequest) error {
	err := s.repository.pickupOrder(ctx, req.OrderID, req.PickupAddress)
	if err != nil {
		return errors.Wrap(err, "s.repository.pickupOrder")
	}
	return nil
}

func (s *service) assignOrder(ctx context.Context, req *AssignRequest) (*AssignResponse, error) {
	customerNotifToken, err := s.repository.assignOrder(ctx, req)
	if err != nil {
		return nil, errors.Wrap(err, "s.repository.assignOrder")
	}

	data, err := json.Marshal(notifications.Message{
		AgentID: s.notificationAgentID,
		Token:   customerNotifToken,
		Title:   "Your Order Is on Its Way!",
		Body:    "Good news! Your order is on its way and will arrive soon. Keep an eye out for the delivery!",
	})
	if err != nil {
		return nil, errors.Wrap(err, "json.Marshal")
	}

	err = s.queue.Publish(ctx, "orders.delivering", data)
	if err != nil {
		return nil, errors.Wrap(err, "s.queue.Publish")
	}

	return &AssignResponse{}, nil
}
