package notifications

import (
	"context"

	"github.com/pkg/errors"
	"github.com/shahzodshafizod/gocloud/pkg"
)

type Service interface {
	sendNotification(context.Context, *Message) error
}

type service struct {
	repository   Repository
	notification pkg.Notification
}

func NewService(repository Repository, notification pkg.Notification) Service {
	return &service{
		repository:   repository,
		notification: notification,
	}
}

func (s *service) sendNotification(ctx context.Context, msg *Message) error {
	agent, err := s.repository.getAgent(ctx, msg.AgentID)
	if err != nil {
		return errors.Wrap(err, "s.repository.getAgent")
	}

	id, err := s.repository.saveNotification(ctx, agent, msg)
	if err != nil {
		return errors.Wrap(err, "s.repository.saveNotification")
	}

	respID, err := s.notification.SendPush(ctx, msg.Token, msg.Body, agent.Priority)
	if err != nil {
		return errors.Wrap(err, "s.notification.SendPush")
	}

	err = s.repository.updateNotification(ctx, id, respID)
	if err != nil {
		return errors.Wrap(err, "s.repository.updateNotification")
	}

	return nil
}
