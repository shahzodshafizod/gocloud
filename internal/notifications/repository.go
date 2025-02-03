package notifications

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/shahzodshafizod/gocloud/pkg"
)

type Repository interface {
	getAgent(context.Context, string) (*agent, error)
	saveNotification(context.Context, *agent, *Message) (string, error)
	updateNotification(ctx context.Context, id string, notifID string) error
}

type repository struct {
	tableName string
	postgres  pkg.Postgres
	nosql     pkg.NoSQL
}

func NewRepository(postgres pkg.Postgres, nosql pkg.NoSQL) Repository {
	return &repository{
		tableName: "notifications",
		postgres:  postgres,
		nosql:     nosql,
	}
}

func (r *repository) getAgent(ctx context.Context, id string) (*agent, error) {
	query := `SELECT name, priority
	FROM agents WHERE id = $1`
	var agent = &agent{ID: id}
	err := r.postgres.QueryRow(ctx, query, id).Scan(
		&agent.Name,
		&agent.Priority,
	)
	if err != nil {
		return nil, errors.Wrap(err, "r.postgres.QueryRow")
	}
	return agent, nil
}

func (r *repository) saveNotification(ctx context.Context, agent *agent, msg *Message) (string, error) {
	items := pkg.Map{
		"agent_id":   agent.ID,
		"agent_name": agent.Name,
		"token_id":   msg.Token,
		"title":      msg.Title,
		"body":       msg.Body,
		"priority":   agent.Priority,
		"created_at": time.Now(),
	}
	id, err := r.nosql.Insert(ctx, r.tableName, items)
	if err != nil {
		return "", errors.Wrap(err, "r.nosql.Insert")
	}
	return id, nil
}

func (r *repository) updateNotification(ctx context.Context, id string, notifID string) error {
	_, err := r.nosql.Update(ctx, r.tableName,
		pkg.Map{"id": id},
		pkg.Map{"notif_id": notifID},
	)
	if err != nil {
		return errors.Wrap(err, "r.nosql.Update")
	}
	return nil
}
