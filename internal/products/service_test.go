package products

import (
	"context"
	"errors"
	"testing"

	"github.com/shahzodshafizod/gocloud/pkg/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type testConfig struct {
	ctrl     *gomock.Controller
	rows     *mocks.MockRows
	postgres *mocks.MockPostgres
	service  Service
}

func __SetupTestConfig(t *testing.T) *testConfig {
	ctrl := gomock.NewController(t)
	cfg := &testConfig{
		ctrl:     ctrl,
		rows:     mocks.NewMockRows(ctrl),
		postgres: mocks.NewMockPostgres(ctrl),
	}
	repository := NewRepository(cfg.postgres)
	cfg.service = NewService(repository)
	return cfg
}

// go test -v -count=1 ./internal/products/ -run TestGetPartnerProducts
func TestGetPartnerProducts(t *testing.T) {
	var cfg = __SetupTestConfig(t)
	defer cfg.ctrl.Finish()

	ctx := context.Background()

	// Test case #1
	targetError := errors.New("r.postgres.Query")
	cfg.postgres.EXPECT().Query(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, targetError)
	_, err := cfg.service.GetPartnerProducts(ctx)
	assert.True(t, errors.Is(err, targetError))

	cfg.postgres.EXPECT().Query(gomock.Any(), gomock.Any(), gomock.Any()).Return(cfg.rows, nil).AnyTimes()
	cfg.rows.EXPECT().Close().AnyTimes()

	// Test case #2: Empty response
	cfg.rows.EXPECT().Next().Return(false)
	resp, err := cfg.service.GetPartnerProducts(ctx)
	assert.NoError(t, err)
	assert.True(t, len(resp.Partners) == 0)

	// Test case #3
	targetError = errors.New("Scan error")
	cfg.rows.EXPECT().Next().Return(true).Times(1)
	cfg.rows.EXPECT().Scan(gomock.Any()).Return(targetError).Times(1)
	_, err = cfg.service.GetPartnerProducts(ctx)
	assert.True(t, errors.Is(err, targetError))

	cfg.rows.EXPECT().Next().Return(true).Times(2)
	cfg.rows.EXPECT().Scan(gomock.Any()).DoAndReturn(func(dest ...any) error {
		*dest[0].(*int32) = 1
		*dest[1].(*string) = "test partner title"
		*dest[2].(*string) = "test partner brand"
		*dest[3].(*int32) = 1
		*dest[4].(*string) = "test product title"
		*dest[5].(*string) = "test product description"
		*dest[6].(*string) = "test product picture URL"
		*dest[7].(*int32) = 100
		return nil
	}).Times(2)
	cfg.rows.EXPECT().Next().Return(false).Times(1)

	// Test case #4
	resp, err = cfg.service.GetPartnerProducts(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
}
