package partners

import (
	"context"

	"github.com/pkg/errors"
	"github.com/shahzodshafizod/gocloud/pkg"
)

type Repository interface {
	checkPartnerProducts(context.Context, *CheckRequest) (*CheckResponse, error)
	getPartnerApiURL(context.Context, int) (string, error)
}

type repository struct {
	postgres pkg.Postgres
}

func NewRepository(postgres pkg.Postgres) Repository {
	return &repository{
		postgres: postgres,
	}
}

func (r *repository) checkPartnerProducts(ctx context.Context, req *CheckRequest) (*CheckResponse, error) {
	var query = `
	SELECT
		pts.title
		, pts.brand
		, ava.price
	FROM available ava
	INNER JOIN partners pts ON pts.id = ava.partner_id AND pts.enabled
	INNER JOIN products pds ON pds.id = ava.product_id
	WHERE ava.active AND ava.product_id = $1 AND ava.partner_id = $2`
	var err error
	var title, brand string
	var totalAmount int64 = 0
	for idx := range req.Products {
		err = r.postgres.QueryRow(ctx, query,
			req.Products[idx].ID,
			req.PartnerID,
		).Scan(
			&title,
			&brand,
			&req.Products[idx].Price,
		)
		if err != nil {
			return nil, errors.Wrap(err, "r.postgres.QueryRow.Scan")
		}
		totalAmount += int64(req.Products[idx].Price * req.Products[idx].Quantity)
	}
	if totalAmount != req.TotalAmount {
		return nil, errors.New("incorrectly calculated TotalAmount")
	}
	return &CheckResponse{
		PartnerTitle: title,
		PartnerBrand: brand,
		Products:     req.Products,
	}, nil
}

func (r *repository) getPartnerApiURL(ctx context.Context, partnerID int) (string, error) {
	query := `SELECT api_url FROM partners WHERE id = $1`
	var apiURL string
	err := r.postgres.QueryRow(ctx, query, partnerID).Scan(&apiURL)
	if err != nil {
		return "", errors.Wrap(err, "r.postgres.QueryRow")
	}
	return apiURL, nil
}
