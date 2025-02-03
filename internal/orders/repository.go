package orders

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/shahzodshafizod/gocloud/pkg"
)

// TODO: NoSQL for payments and orders_history (status changes, delivery updates, ect...)

type Repository interface {
	getBank(context.Context, string) (*bank, error)
	createOrder(context.Context, *Order) (int64, error)
	payOrder(context.Context, *PayRequest) (string, *PaidOrder, error)
	pickupOrder(ctx context.Context, orderID int64, pickupAddress string) error
	assignOrder(context.Context, *AssignRequest) (string, error)
}

type repository struct {
	postgres pkg.Postgres
	nosql    pkg.NoSQL
}

func NewRepository(postgres pkg.Postgres, nosql pkg.NoSQL) Repository {
	return &repository{
		postgres: postgres,
		nosql:    nosql,
	}
}

func (r *repository) getBank(ctx context.Context, id string) (*bank, error) {
	query := `SELECT id, webcheckout_url FROM banks WHERE id = $1 AND active`
	var bank = &bank{}
	err := r.postgres.QueryRow(ctx, query, id).Scan(&bank.ID, &bank.WebcheckoutURL)
	if err != nil {
		return nil, errors.Wrap(err, "r.postgres.QueryRow")
	}
	return bank, nil
}

func (r *repository) createOrder(ctx context.Context, order *Order) (int64, error) {
	query := `INSERT INTO orders (
		order_id
		, customer_id
		, customer_name
		, customer_phone
		, customer_notif_token
		, delivery_address
		, partner_id
		, partner_title
		, partner_brand
		, total_amount
		, paytype
		, products
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	RETURNING id`
	var id int64
	err := r.postgres.QueryRow(ctx, query,
		order.OrderID,
		order.CustomerID,
		order.CustomerName,
		order.CustomerPhone,
		order.CustomerNotifToken,
		order.DeliveryAddress,
		order.PartnerID,
		order.PartnerTitle,
		order.PartnerBrand,
		order.TotalAmount,
		order.Paytype,
		order.Products,
	).Scan(&id)
	if err != nil {
		return 0, errors.Wrap(err, "r.postgres.QueryRow")
	}

	return id, nil
}

func (r *repository) payOrder(ctx context.Context, req *PayRequest) (string, *PaidOrder, error) {
	tx, err := r.postgres.Begin(ctx)
	if err != nil {
		return "", nil, errors.Wrap(err, "r.postgres.Begin")
	}

	query := `UPDATE orders
	SET
		paid_amount = $1
		, updated_at = now()
		, status = 'paid'
	WHERE
		id = $2
		AND status = 'pending'
		AND total_amount <= $1
	RETURNING products, partner_id`

	var order = &PaidOrder{
		OrderID: req.OrderID,
	}
	err = tx.QueryRow(ctx, query,
		req.PaidAmount,
		req.OrderID,
	).Scan(
		&order.Products,
		&order.PartnerID,
	)
	if err != nil {
		tx.Rollback(ctx)
		return "", nil, errors.Wrap(err, "UPDATE orders")
	}

	var payItems = pkg.Map{
		"payment_id": req.PaymentID,
		"order_id":   req.OrderID,
		"amount":     req.PaidAmount,
		"status":     "paid",
		"created_at": time.Now(),
	}
	paymentID, err := r.nosql.Insert(ctx, "payments", payItems)
	if err != nil {
		tx.Rollback(ctx)
		return "", nil, errors.Wrap(err, "r.nosql.Insert")
	}

	err = tx.Commit(ctx)
	if err != nil {
		// TODO: What to do with the inserted payment into NoSQL?
		tx.Rollback(ctx)
		return "", nil, errors.Wrap(err, "tx.Commit")
	}

	return paymentID, order, nil
}

func (r *repository) pickupOrder(ctx context.Context, orderID int64, pickupAddress string) error {
	query := `UPDATE orders
	SET
		pickup_address = $1
		, updated_at = now()
		, status = 'ready'
	WHERE id = $2 AND status = 'paid'`
	err := r.postgres.Exec(ctx, query, pickupAddress, orderID)
	if err != nil {
		return errors.Wrap(err, "r.postgres.Exec")
	}
	return nil
}

func (r *repository) assignOrder(ctx context.Context, req *AssignRequest) (string, error) {
	query := `UPDATE orders
	SET
		deliverer_id = $1
		, updated_at = now()
		, status = 'delivering'
	WHERE id = $2 AND status = 'ready'
	RETURNING customer_notif_token`
	var customerNotifToken string
	err := r.postgres.QueryRow(ctx, query,
		req.DelivererID,
		req.OrderID,
	).Scan(&customerNotifToken)
	if err != nil {
		return "", errors.Wrap(err, "r.postgres.QueryRow")
	}
	return customerNotifToken, nil
}
