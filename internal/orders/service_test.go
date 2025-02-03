package orders

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
	row      *mocks.MockRow
	postgres *mocks.MockPostgres
	nosql    *mocks.MockNoSQL
	queue    *mocks.MockQueue
	service  Service
}

func __SetupTestConfig(t *testing.T) *testConfig {
	ctrl := gomock.NewController(t)
	cfg := &testConfig{
		ctrl:     ctrl,
		row:      mocks.NewMockRow(ctrl),
		postgres: mocks.NewMockPostgres(ctrl),
		nosql:    mocks.NewMockNoSQL(ctrl),
		queue:    mocks.NewMockQueue(ctrl),
	}
	repository := NewRepository(cfg.postgres, cfg.nosql)
	cfg.service = NewService(repository, cfg.queue)
	return cfg
}

// go test -v -count=1 ./internal/orders/ -run ^TestCreateOrder$
func TestCreateOrder(t *testing.T) {
	cfg := __SetupTestConfig(t)
	defer cfg.ctrl.Finish()

	ctx := context.Background()
	order := &Order{}

	cfg.postgres.EXPECT().QueryRow(gomock.Any(), gomock.Any(), gomock.Any()).Return(cfg.row).AnyTimes()

	// Test case #1
	targetError := errors.New("getBank error")
	cfg.row.EXPECT().Scan(gomock.Any()).Return(targetError)
	_, err := cfg.service.createOrder(ctx, order)
	assert.True(t, errors.Is(err, targetError))

	// Test case #2
	targetError = errors.New("createOrder error")
	cfg.row.EXPECT().Scan(gomock.Any()).Return(nil)
	cfg.row.EXPECT().Scan(gomock.Any()).Return(targetError)
	_, err = cfg.service.createOrder(ctx, order)
	assert.True(t, errors.Is(err, targetError))

	cfg.row.EXPECT().Scan(gomock.Any()).Return(nil).AnyTimes()

	// Test case #3: Success
	resp, err := cfg.service.createOrder(ctx, order)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

// go test -v -count=1 ./internal/orders/ -run ^TestPayOrder$
func TestPayOrder(t *testing.T) {
	cfg := __SetupTestConfig(t)
	defer cfg.ctrl.Finish()

	ctx := context.Background()
	req := &PayRequest{
		OrderID:    1,
		PaymentID:  "8baec4ff-e08a-4d4a-bf5e-1e8dc4dc9f55",
		PaidAmount: 300,
	}

	tx := mocks.NewMockTx(cfg.ctrl)

	// Test case #1
	targetError := errors.New("payOrder Begin tx error")
	cfg.postgres.EXPECT().Begin(gomock.Any()).Return(nil, targetError)
	_, err := cfg.service.payOrder(ctx, req)
	assert.True(t, errors.Is(err, targetError))

	cfg.postgres.EXPECT().Begin(gomock.Any()).Return(tx, nil).AnyTimes()
	tx.EXPECT().QueryRow(gomock.Any(), gomock.Any(), gomock.Any()).Return(cfg.row).AnyTimes()
	tx.EXPECT().Rollback(gomock.Any()).AnyTimes()

	// Test case #2
	targetError = errors.New("payOrder Scan error")
	cfg.row.EXPECT().Scan(gomock.Any()).Return(targetError)
	_, err = cfg.service.payOrder(ctx, req)
	assert.True(t, errors.Is(err, targetError))

	cfg.row.EXPECT().Scan(gomock.Any()).Return(nil).AnyTimes()

	// Test case #3
	targetError = errors.New("payOrder nosql.Insert error")
	cfg.nosql.EXPECT().Insert(gomock.Any(), gomock.Any(), gomock.Any()).Return("", targetError)
	_, err = cfg.service.payOrder(ctx, req)
	assert.True(t, errors.Is(err, targetError))

	cfg.nosql.EXPECT().Insert(gomock.Any(), gomock.Any(), gomock.Any()).Return("", nil).AnyTimes()

	// Test case #4
	targetError = errors.New("payOrder Commit error")
	tx.EXPECT().Commit(gomock.Any()).Return(targetError)
	_, err = cfg.service.payOrder(ctx, req)
	assert.True(t, errors.Is(err, targetError))

	tx.EXPECT().Commit(gomock.Any()).Return(nil).AnyTimes()

	// Test case #5
	targetError = errors.New("Publish error")
	cfg.queue.EXPECT().Publish(gomock.Any(), gomock.Any(), gomock.Any()).Return(targetError)
	_, err = cfg.service.payOrder(ctx, req)
	assert.True(t, errors.Is(err, targetError))

	cfg.queue.EXPECT().Publish(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	// Test case #6: Success
	resp, err := cfg.service.payOrder(ctx, req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

// go test -v -count=1 ./internal/orders/ -run ^TestPickupOrder$
func TestPickupOrder(t *testing.T) {
	cfg := __SetupTestConfig(t)
	defer cfg.ctrl.Finish()

	ctx := context.Background()
	req := &pickupRequest{
		OrderID:       1,
		PickupAddress: "London, Big Ben",
	}

	// Test case #1
	targetError := errors.New("pickupOrder error")
	cfg.postgres.EXPECT().Exec(gomock.Any(), gomock.Any(), gomock.Any()).Return(targetError)
	err := cfg.service.pickupOrder(ctx, req)
	assert.True(t, errors.Is(err, targetError))

	cfg.postgres.EXPECT().Exec(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	// Test case #2: Success
	err = cfg.service.pickupOrder(ctx, req)
	assert.NoError(t, err)
}

// go test -v -count=1 ./internal/orders/ -run ^TestAssignOrder$
func TestAssignOrder(t *testing.T) {
	cfg := __SetupTestConfig(t)
	defer cfg.ctrl.Finish()

	ctx := context.Background()
	req := &AssignRequest{
		OrderID:     1,
		DelivererID: "0258dc6d-cc4f-418d-a9a1-62a474d86bb2",
	}

	cfg.postgres.EXPECT().QueryRow(gomock.Any(), gomock.Any(), gomock.Any()).Return(cfg.row).AnyTimes()

	// Test case #1
	targetError := errors.New("assignOrder error")
	cfg.row.EXPECT().Scan(gomock.Any()).Return(targetError)
	_, err := cfg.service.assignOrder(ctx, req)
	assert.True(t, errors.Is(err, targetError))

	cfg.row.EXPECT().Scan(gomock.Any()).Return(nil).AnyTimes()

	// Test case #2
	targetError = errors.New("Publish error")
	cfg.queue.EXPECT().Publish(gomock.Any(), gomock.Any(), gomock.Any()).Return(targetError)
	_, err = cfg.service.assignOrder(ctx, req)
	assert.True(t, errors.Is(err, targetError))

	cfg.queue.EXPECT().Publish(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	// Test case #3: Success
	resp, err := cfg.service.assignOrder(ctx, req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
}
