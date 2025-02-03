package notifications

import (
	"context"
	"errors"
	"testing"

	"github.com/shahzodshafizod/gocloud/pkg"
	"github.com/shahzodshafizod/gocloud/pkg/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type testConfig struct {
	ctrl         *gomock.Controller
	row          *mocks.MockRow
	postgres     *mocks.MockPostgres
	nosql        *mocks.MockNoSQL
	notification *mocks.MockNotification
	service      Service
}

func __SetupTestConfig(t *testing.T) *testConfig {
	ctrl := gomock.NewController(t)
	cfg := &testConfig{
		ctrl:         ctrl,
		row:          mocks.NewMockRow(ctrl),
		postgres:     mocks.NewMockPostgres(ctrl),
		nosql:        mocks.NewMockNoSQL(ctrl),
		notification: mocks.NewMockNotification(ctrl),
	}
	repository := NewRepository(cfg.postgres, cfg.nosql)
	cfg.service = NewService(repository, cfg.notification)
	return cfg
}

// go test -v -count=1 ./internal/notifications/ -run ^TestSendNotification$
func TestSendNotification(t *testing.T) {
	cfg := __SetupTestConfig(t)
	defer cfg.ctrl.Finish()

	ctx := context.Background()
	message := &Message{
		AgentID: "test agent id",
		Token:   "test token",
		Title:   "test title",
		Body:    "test body",
	}

	cfg.postgres.EXPECT().QueryRow(gomock.Any(), gomock.Any(), gomock.Any()).Return(cfg.row).AnyTimes()

	// Test case #1
	targetError := errors.New("getAgent error")
	cfg.row.EXPECT().Scan(gomock.Any(), gomock.Any()).Return(targetError)
	err := cfg.service.sendNotification(ctx, message)
	assert.True(t, errors.Is(err, targetError))

	cfg.row.EXPECT().Scan(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	// Test case #2
	targetError = errors.New("saveNotification error")
	cfg.nosql.EXPECT().Insert(gomock.Any(), gomock.Any(), gomock.Any()).Return("", targetError)
	err = cfg.service.sendNotification(ctx, message)
	assert.True(t, errors.Is(err, targetError))

	cfg.nosql.EXPECT().Insert(gomock.Any(), gomock.Any(), gomock.Any()).Return("", nil).AnyTimes()

	// Test case #3
	targetError = errors.New("SendPush error")
	cfg.notification.EXPECT().SendPush(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return("", targetError)
	err = cfg.service.sendNotification(ctx, message)
	assert.True(t, errors.Is(err, targetError))

	cfg.notification.EXPECT().SendPush(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return("", nil).AnyTimes()

	// Test case #4
	targetError = errors.New("updateNotification error")
	cfg.nosql.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, targetError)
	err = cfg.service.sendNotification(ctx, message)
	assert.True(t, errors.Is(err, targetError))

	cfg.nosql.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(pkg.Map{}, nil).AnyTimes()

	// Test case #5: Success
	err = cfg.service.sendNotification(ctx, message)
	assert.NoError(t, err)
}
