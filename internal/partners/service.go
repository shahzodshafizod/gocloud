package partners

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
	"github.com/shahzodshafizod/gocloud/internal/orders"
	"github.com/shahzodshafizod/gocloud/internal/products"
	"github.com/shahzodshafizod/gocloud/pkg"
)

type Service interface {
	getPartnerProducts(context.Context) (*products.GetAllResponse, error)
	checkPartnerProducts(context.Context, *CheckRequest) (*CheckResponse, error)
	sendToPartner(context.Context, *orders.PaidOrder) error
}

type service struct {
	repository Repository
	httpClient pkg.HTTPClient
	products   products.Service
}

func NewService(repository Repository, httpClient pkg.HTTPClient, products products.Service) Service {
	return &service{
		repository: repository,
		httpClient: httpClient,
		products:   products,
	}
}

func (s *service) getPartnerProducts(ctx context.Context) (*products.GetAllResponse, error) {
	resp, err := s.products.GetPartnerProducts(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "s.products.GetPartnerProducts")
	}
	return resp, nil
}

func (s *service) checkPartnerProducts(ctx context.Context, req *CheckRequest) (*CheckResponse, error) {
	resp, err := s.repository.checkPartnerProducts(ctx, req)
	if err != nil {
		return nil, errors.Wrap(err, "s.repository.checkPartnerProducts")
	}
	return resp, nil
}

func (s *service) sendToPartner(ctx context.Context, order *orders.PaidOrder) error {
	apiURL, err := s.repository.getPartnerApiURL(ctx, order.PartnerID)
	if err != nil {
		return errors.Wrap(err, "s.repository.getPartnerApiURL")
	}

	body, err := json.Marshal(order)
	if err != nil {
		return errors.Wrap(err, "json.Marshal")
	}

	_, err = s.httpClient.SendRequest(http.MethodPost, apiURL, bytes.NewReader(body))
	if err != nil {
		return errors.Wrap(err, "s.httpClient.SendRequest")
	}

	return nil
}
