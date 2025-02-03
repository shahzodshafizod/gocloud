package products

import (
	"context"

	"github.com/pkg/errors"
	"github.com/shahzodshafizod/gocloud/pkg"
)

type Repository interface {
	getPartnerProducts(context.Context) (*GetAllResponse, error)
}

type repository struct {
	postgres pkg.Postgres
}

func NewRepository(postgres pkg.Postgres) Repository {
	return &repository{
		postgres: postgres,
	}
}

func (r *repository) getPartnerProducts(ctx context.Context) (*GetAllResponse, error) {
	var query = `SELECT
		pts.id
		, pts.title
		, pts.brand

		, pds.id
		, pds.title
		, pds.description
		, pds.picture_url
		, ava.price
	FROM available ava
	INNER JOIN partners pts ON pts.id = ava.partner_id
	INNER JOIN products pds ON pds.id = ava.product_id
	WHERE ava.active AND pts.verified AND pts.enabled
	ORDER BY pts.id, pds.id`

	rows, err := r.postgres.Query(ctx, query)
	if err != nil {
		return nil, errors.Wrap(err, "r.postgres.Query")
	}
	defer rows.Close()

	var prevID int32 = -1
	var partner *Partner
	var resp = &GetAllResponse{Partners: make([]*Partner, 0)}

	for rows.Next() {
		var id int32
		var title, brand string
		var product = &PartnerProduct{}
		err = rows.Scan(
			&id,
			&title,
			&brand,
			&product.ID,
			&product.Title,
			&product.Description,
			&product.PictureURL,
			&product.Price,
		)
		if err != nil {
			return nil, errors.Wrap(err, "rows.Scan")
		}
		if id != prevID {
			resp.Partners = append(resp.Partners, partner)
			partner = &Partner{
				ID:       id,
				Title:    title,
				Brand:    brand,
				Products: []*PartnerProduct{product},
			}
		} else {
			partner.Products = append(partner.Products, product)
		}
		prevID = id
	}

	resp.Partners = append(resp.Partners, partner)
	resp.Partners = resp.Partners[1:]
	return resp, nil
}
