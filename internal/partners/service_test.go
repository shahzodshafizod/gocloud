package partners

import (
	"context"
	"errors"
	"testing"

	"github.com/shahzodshafizod/gocloud/internal/orders"
	"github.com/shahzodshafizod/gocloud/internal/products"
	productsmocks "github.com/shahzodshafizod/gocloud/internal/products/mocks"
	"github.com/shahzodshafizod/gocloud/pkg/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type testConfig struct {
	ctrl        *gomock.Controller
	row         *mocks.MockRow
	postgres    *mocks.MockPostgres
	httpClient  *mocks.MockHTTPClient
	productsSvc *productsmocks.MockService
	service     Service
}

func __SetupTestConfig(t *testing.T) *testConfig {
	ctrl := gomock.NewController(t)
	cfg := &testConfig{
		ctrl:        ctrl,
		row:         mocks.NewMockRow(ctrl),
		postgres:    mocks.NewMockPostgres(ctrl),
		httpClient:  mocks.NewMockHTTPClient(ctrl),
		productsSvc: productsmocks.NewMockService(ctrl),
	}
	repository := NewRepository(cfg.postgres)
	cfg.service = NewService(repository, cfg.httpClient, cfg.productsSvc)
	return cfg
}

// go test -v -count=1 ./internal/partners/ -run ^TestGetPartnerProducts$
func TestGetPartnerProducts(t *testing.T) {
	cfg := __SetupTestConfig(t)
	defer cfg.ctrl.Finish()

	ctx := context.Background()

	// Test case #1
	targetError := errors.New("products GetPartnerProducts error")
	cfg.productsSvc.EXPECT().GetPartnerProducts(gomock.Any()).Return(nil, targetError)
	_, err := cfg.service.getPartnerProducts(ctx)
	assert.True(t, errors.Is(err, targetError))

	cfg.productsSvc.EXPECT().GetPartnerProducts(gomock.Any()).Return(&products.GetAllResponse{}, nil).AnyTimes()

	// Test case #2: Success
	resp, err := cfg.service.getPartnerProducts(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

// go test -v -count=1 ./internal/partners/ -run ^TestCheckPartnerProducts$
func TestCheckPartnerProducts(t *testing.T) {
	cfg := __SetupTestConfig(t)
	defer cfg.ctrl.Finish()

	ctx := context.Background()
	req := &CheckRequest{
		PartnerID:   1,
		TotalAmount: 300,
		Products: []*orders.Product{
			{Quantity: 2, Price: 100},
			{Quantity: 1, Price: 100},
		},
	}

	cfg.postgres.EXPECT().QueryRow(gomock.Any(), gomock.Any(), gomock.Any()).Return(cfg.row).AnyTimes()

	// Test case #1
	targetError := errors.New("Scan error")
	cfg.row.EXPECT().Scan(gomock.Any()).Return(targetError)
	_, err := cfg.service.checkPartnerProducts(ctx, req)
	assert.True(t, errors.Is(err, targetError))

	// Test case #2: Wrong calculated Total Amount
	cfg.row.EXPECT().Scan(gomock.Any()).DoAndReturn(func(dest ...any) error {
		*dest[0].(*string) = "test partner title"
		*dest[1].(*string) = "test partner brand"
		*dest[2].(*int32) = 200
		return nil
	}).Times(2)
	_, err = cfg.service.checkPartnerProducts(ctx, req)
	assert.Error(t, err)

	// Test case #3: Success
	cfg.row.EXPECT().Scan(gomock.Any()).DoAndReturn(func(dest ...any) error {
		*dest[0].(*string) = "test partner title"
		*dest[1].(*string) = "test partner brand"
		*dest[2].(*int32) = 100
		return nil
	}).AnyTimes()
	resp, err := cfg.service.checkPartnerProducts(ctx, req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

// go test -v -count=1 ./internal/partners/ -run ^TestSendToPartner$
func TestSendToPartner(t *testing.T) {
	cfg := __SetupTestConfig(t)
	defer cfg.ctrl.Finish()

	ctx := context.Background()
	order := &orders.PaidOrder{
		OrderID:     100000,
		PartnerID:   1,
		CallbackURL: "test callback URL",
		Products: []*orders.Product{
			{ID: 1, Quantity: 2, Price: 100},
			{ID: 2, Quantity: 1, Price: 100},
		},
	}

	cfg.postgres.EXPECT().QueryRow(gomock.Any(), gomock.Any(), gomock.Any()).Return(cfg.row).AnyTimes()

	// Test case #1
	targetError := errors.New("getPartnerApiURL error")
	cfg.row.EXPECT().Scan(gomock.Any()).Return(targetError)
	err := cfg.service.sendToPartner(ctx, order)
	assert.True(t, errors.Is(err, targetError))

	cfg.row.EXPECT().Scan(gomock.Any()).Return(nil).AnyTimes()

	// Test case #2
	targetError = errors.New("httpClient Send error")
	cfg.httpClient.EXPECT().SendRequest(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, targetError)
	err = cfg.service.sendToPartner(ctx, order)
	assert.True(t, errors.Is(err, targetError))

	cfg.httpClient.EXPECT().SendRequest(gomock.Any(), gomock.Any(), gomock.Any()).Return([]byte{}, nil).AnyTimes()

	// Test case #3: Success
	err = cfg.service.sendToPartner(ctx, order)
	assert.NoError(t, err)
}
