package products

import (
	"context"

	"github.com/pkg/errors"
)

type Service interface {
	GetPartnerProducts(context.Context) (*GetAllResponse, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &service{
		repository: repository,
	}
}

func (s *service) GetPartnerProducts(ctx context.Context) (*GetAllResponse, error) {
	resp, err := s.repository.getPartnerProducts(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "s.repository.getPartnerProducts")
	}
	return resp, nil
}
