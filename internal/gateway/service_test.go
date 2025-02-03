package gateway

import (
	"context"
	"errors"
	"testing"

	"github.com/shahzodshafizod/gocloud/internal/orders"
	ordersmocks "github.com/shahzodshafizod/gocloud/internal/orders/mocks"
	"github.com/shahzodshafizod/gocloud/internal/partners"
	partnersmocks "github.com/shahzodshafizod/gocloud/internal/partners/mocks"
	"github.com/shahzodshafizod/gocloud/internal/products"
	"github.com/shahzodshafizod/gocloud/pkg"
	"github.com/shahzodshafizod/gocloud/pkg/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type testConfig struct {
	ctrl           *gomock.Controller
	authManager    *mocks.MockAuth
	partnersClient *partnersmocks.MockPartnersClient
	ordersClient   *ordersmocks.MockOrdersClient
	cache          *mocks.MockCache
	queue          *mocks.MockQueue
	storage        *mocks.MockStorage
	service        Service
}

func __SetupTestConfig(t *testing.T) *testConfig {
	ctrl := gomock.NewController(t)
	cfg := &testConfig{
		ctrl:           ctrl,
		authManager:    mocks.NewMockAuth(ctrl),
		partnersClient: partnersmocks.NewMockPartnersClient(ctrl),
		ordersClient:   ordersmocks.NewMockOrdersClient(ctrl),
		cache:          mocks.NewMockCache(ctrl),
		queue:          mocks.NewMockQueue(ctrl),
		storage:        mocks.NewMockStorage(ctrl),
	}
	cfg.service = NewService(cfg.authManager, cfg.partnersClient,
		cfg.ordersClient, cfg.cache, cfg.queue, cfg.storage)
	return cfg
}

// go test -v -count=1 ./internal/gateway/ -run ^TestSignUp$
func TestSignUp(t *testing.T) {
	cfg := __SetupTestConfig(t)
	defer cfg.ctrl.Finish()

	ctx := context.Background()
	signUp := &signUp{}

	// Test case #1
	targetError := errors.New("s.authManager.SignUp error")
	cfg.authManager.EXPECT().SignUp(gomock.Any(), gomock.Any()).Return("", targetError)
	_, err := cfg.service.signUp(ctx, signUp)
	assert.True(t, errors.Is(err, targetError))

	cfg.authManager.EXPECT().SignUp(gomock.Any(), gomock.Any()).Return("", nil).AnyTimes()

	// Test case #2: Success
	resp, err := cfg.service.signUp(ctx, signUp)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

// go test -v -count=1 ./internal/gateway/ -run ^TestConfirmSignUp$
func TestConfirmSignUp(t *testing.T) {
	cfg := __SetupTestConfig(t)
	defer cfg.ctrl.Finish()

	ctx := context.Background()
	confirm := &confirmSignUp{}

	// Test case #1
	targetError := errors.New("s.authManager.ConfirmSignUp")
	cfg.authManager.EXPECT().ConfirmSignUp(gomock.Any(), gomock.Any()).Return(targetError)
	err := cfg.service.confirmSignUp(ctx, confirm)
	assert.True(t, errors.Is(err, targetError))

	cfg.authManager.EXPECT().ConfirmSignUp(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	err = cfg.service.confirmSignUp(ctx, confirm)
	assert.NoError(t, err)
}

// go test -v -count=1 ./internal/gateway/ -run ^TestSignIn$
func TestSignIn(t *testing.T) {
	cfg := __SetupTestConfig(t)
	defer cfg.ctrl.Finish()

	ctx := context.Background()
	signin := &signIn{}

	// Test case #1
	targetError := errors.New("s.authManager.SignIn error")
	cfg.authManager.EXPECT().SignIn(gomock.Any(), gomock.Any()).Return(nil, nil, targetError)
	_, err := cfg.service.signIn(ctx, signin)
	assert.True(t, errors.Is(err, targetError))

	cfg.authManager.EXPECT().SignIn(gomock.Any(), gomock.Any()).Return(nil, &pkg.Token{}, nil).AnyTimes()

	// Test case #2: Success
	token, err := cfg.service.signIn(ctx, signin)
	assert.NoError(t, err)
	assert.NotNil(t, token)
}

// go test -v -count=1 ./internal/gateway/ -run ^TestCheckToken$
func TestCheckToken(t *testing.T) {
	cfg := __SetupTestConfig(t)
	defer cfg.ctrl.Finish()

	ctx := context.Background()

	// Test case #1
	targetError := errors.New("s.authManager.CheckToken error")
	cfg.authManager.EXPECT().CheckToken(gomock.Any(), gomock.Any()).Return(nil, targetError)
	_, err := cfg.service.checkToken(ctx, "")
	assert.True(t, errors.Is(err, targetError))

	cfg.authManager.EXPECT().CheckToken(gomock.Any(), gomock.Any()).Return(&pkg.User{}, nil).AnyTimes()

	// Test case #2: Success
	user, err := cfg.service.checkToken(ctx, "")
	assert.NoError(t, err)
	assert.NotNil(t, user)
}

// go test -v -count=1 ./internal/gateway/ -run ^TestRefreshToken$
func TestRefreshToken(t *testing.T) {
	cfg := __SetupTestConfig(t)
	defer cfg.ctrl.Finish()

	ctx := context.Background()
	refresh := &refreshToken{}

	// Test case #1
	targetError := errors.New("s.authManager.RefreshToken error")
	cfg.authManager.EXPECT().RefreshToken(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, targetError)
	_, err := cfg.service.refreshToken(ctx, "", refresh)
	assert.True(t, errors.Is(err, targetError))

	cfg.authManager.EXPECT().RefreshToken(gomock.Any(), gomock.Any(), gomock.Any()).Return(&pkg.Token{}, nil).AnyTimes()

	// Test case #2: Success
	token, err := cfg.service.refreshToken(ctx, "", refresh)
	assert.NoError(t, err)
	assert.NotNil(t, token)
}

// go test -v -count=1 ./internal/gateway/ -run ^TestUpdateUser$
func TestUpdateUser(t *testing.T) {
	cfg := __SetupTestConfig(t)
	defer cfg.ctrl.Finish()

	ctx := context.Background()
	user := &user{}
	update := &updateUser{
		FirstName: "test first name",
		LastName:  "test last name",
		Email:     "test email",
		Phone:     "test phone",
		// BirthDate:  "test birth date",
		NotifToken: "test notification token",
	}

	// Test case #1
	targetError := errors.New("s.authManager.UpdateUser error")
	cfg.authManager.EXPECT().UpdateUser(gomock.Any(), gomock.Any(), gomock.Any()).Return(targetError)
	err := cfg.service.updateUser(ctx, user, update)
	assert.True(t, errors.Is(err, targetError))

	cfg.authManager.EXPECT().UpdateUser(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	// Test case #2: Succes
	err = cfg.service.updateUser(ctx, user, update)
	assert.NoError(t, err)
}

// go test -v -count=1 ./internal/gateway/ -run ^TestConfirmChangeEmail$
func TestConfirmChangeEmail(t *testing.T) {
	cfg := __SetupTestConfig(t)
	defer cfg.ctrl.Finish()

	ctx := context.Background()
	confirm := &confirmChangeEmail{}

	// Test case #1
	targetError := errors.New("s.authManager.ConfirmChangeEmail error")
	cfg.authManager.EXPECT().ConfirmChangeEmail(gomock.Any(), gomock.Any(), gomock.Any()).Return(targetError)
	err := cfg.service.confirmChangeEmail(ctx, confirm)
	assert.True(t, errors.Is(err, targetError))

	cfg.authManager.EXPECT().ConfirmChangeEmail(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	// Test case #2
	err = cfg.service.confirmChangeEmail(ctx, confirm)
	assert.NoError(t, err)
}

// // go test -v -count=1 ./internal/gateway/ -run ^TestUpdateAvatar$
// func TestUpdateAvatar(t *testing.T) {
// 	cfg := __SetupTestConfig(t)
// 	defer cfg.ctrl.Finish()

// 	ctx := context.Background()
// 	user := &user{PhotoURL: "old user photo URL"}
// 	var mfile = mocks.NewMockFile(cfg.ctrl)
// 	var info = mocks.NewMockFileInfo(cfg.ctrl)

// 	mfile.EXPECT().Read(gomock.Any()).Return(0, nil).AnyTimes()

// 	// Test case #1: image.Decode error
// 	err := cfg.service.updateAvatar(ctx, user, mfile, info)
// 	assert.Error(t, err)

// 	file, err := os.Open("../../assets/picture.jpg")
// 	assert.NoError(t, err)
// 	info.EXPECT().ContentType().Return("").AnyTimes()

// 	// Test case #2
// 	targetError := errors.New("s.storage.Upload error")
// 	cfg.storage.EXPECT().Upload(gomock.Any(), gomock.Any()).Return("", targetError)
// 	err = cfg.service.updateAvatar(ctx, user, file, info)
// 	assert.True(t, errors.Is(err, targetError))

// 	cfg.storage.EXPECT().Upload(gomock.Any(), gomock.Any()).Return("", nil).AnyTimes()
// 	file, err = os.Open("../../assets/picture.jpg")
// 	assert.NoError(t, err)

// 	// Test case #3
// 	targetError = errors.New("s.authManager.UpdateUser error")
// 	cfg.authManager.EXPECT().UpdateUser(gomock.Any(), gomock.Any(), gomock.Any()).Return(targetError)
// 	cfg.storage.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(nil)
// 	err = cfg.service.updateAvatar(ctx, user, file, info)
// 	assert.True(t, errors.Is(err, targetError))

// 	file, err = os.Open("../../assets/picture.jpg")
// 	assert.NoError(t, err)

// 	// Test case #4
// 	targetError = errors.New("s.storage.Delete error")
// 	cfg.authManager.EXPECT().UpdateUser(gomock.Any(), gomock.Any(), gomock.Any()).Return(targetError)
// 	cfg.storage.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(targetError)
// 	err = cfg.service.updateAvatar(ctx, user, file, info)
// 	assert.True(t, errors.Is(err, targetError))

// 	file, err = os.Open("../../assets/picture.jpg")
// 	assert.NoError(t, err)
// 	cfg.authManager.EXPECT().UpdateUser(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

// 	// Test case #5
// 	targetError = errors.New("s.storage.Delete error")
// 	cfg.storage.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(targetError)
// 	err = cfg.service.updateAvatar(ctx, user, file, info)
// 	assert.True(t, errors.Is(err, targetError))

// 	file, err = os.Open("../../assets/picture.jpg")
// 	assert.NoError(t, err)
// 	cfg.storage.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

// 	// Test case #6: Success
// 	err = cfg.service.updateAvatar(ctx, user, file, info)
// 	assert.NoError(t, err)
// }

// go test -v -count=1 ./internal/gateway/ -run ^TestForgotPassword$
func TestForgotPassword(t *testing.T) {
	cfg := __SetupTestConfig(t)
	defer cfg.ctrl.Finish()

	ctx := context.Background()

	// Test case #1
	targetError := errors.New("s.authManager.ForgotPassword error")
	cfg.authManager.EXPECT().ForgotPassword(gomock.Any(), gomock.Any()).Return(targetError)
	err := cfg.service.forgotPassword(ctx, "")
	assert.True(t, errors.Is(err, targetError))

	cfg.authManager.EXPECT().ForgotPassword(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	// Test case #2: Success
	err = cfg.service.forgotPassword(ctx, "")
	assert.NoError(t, err)
}

// go test -v -count=1 ./internal/gateway/ -run ^TestResetPassword$
func TestResetPassword(t *testing.T) {
	cfg := __SetupTestConfig(t)
	defer cfg.ctrl.Finish()

	ctx := context.Background()
	reset := &resetPassword{}

	// Test case #1
	targetError := errors.New("s.authManager.ResetPassword error")
	cfg.authManager.EXPECT().ResetPassword(gomock.Any(), gomock.Any()).Return(targetError)
	err := cfg.service.resetPassword(ctx, reset)
	assert.True(t, errors.Is(err, targetError))

	cfg.authManager.EXPECT().ResetPassword(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	// Test case #2: Success
	err = cfg.service.resetPassword(ctx, reset)
	assert.NoError(t, err)
}

// go test -v -count=1 ./internal/gateway/ -run ^TestChangePassword$
func TestChangePassword(t *testing.T) {
	cfg := __SetupTestConfig(t)
	defer cfg.ctrl.Finish()

	ctx := context.Background()
	user := &user{}
	change := &changePassword{}

	// Test case #1
	targetError := errors.New("s.authManager.ChangePassword error")
	cfg.authManager.EXPECT().ChangePassword(gomock.Any(), gomock.Any()).Return(targetError)
	err := cfg.service.changePassword(ctx, user, change)
	assert.True(t, errors.Is(err, targetError))

	cfg.authManager.EXPECT().ChangePassword(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	// Test case #2: Success
	err = cfg.service.changePassword(ctx, user, change)
	assert.NoError(t, err)
}

// go test -v -count=1 ./internal/gateway/ -run ^TestSignOut$
func TestSignOut(t *testing.T) {
	cfg := __SetupTestConfig(t)
	defer cfg.ctrl.Finish()

	ctx := context.Background()
	refresh := &refreshToken{}

	// Test case #1
	targetError := errors.New("s.authManager.SignOut error")
	cfg.authManager.EXPECT().SignOut(gomock.Any(), gomock.Any(), gomock.Any()).Return(targetError)
	err := cfg.service.signOut(ctx, "", refresh)
	assert.True(t, errors.Is(err, targetError))

	cfg.authManager.EXPECT().SignOut(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	// Test case #2: Success
	err = cfg.service.signOut(ctx, "", refresh)
	assert.NoError(t, err)
}

// go test -v -count=1 ./internal/gateway/ -run ^TestDeleteUser$
func TestDeleteUser(t *testing.T) {
	cfg := __SetupTestConfig(t)
	defer cfg.ctrl.Finish()

	ctx := context.Background()

	// Test case #1
	targetError := errors.New("s.authManager.DeleteUser error")
	cfg.authManager.EXPECT().DeleteUser(gomock.Any(), gomock.Any()).Return(targetError)
	err := cfg.service.deleteUser(ctx, "")
	assert.True(t, errors.Is(err, targetError))

	cfg.authManager.EXPECT().DeleteUser(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	// Test case #2: Success
	err = cfg.service.deleteUser(ctx, "")
	assert.NoError(t, err)
}

// go test -count=1 -v ./internal/gateway/ -run ^TestGetPartnerProducts$
func TestGetPartnerProducts(t *testing.T) {
	cfg := __SetupTestConfig(t)
	defer cfg.ctrl.Finish()

	ctx := context.Background()

	// Test case 1: Error
	targetError := errors.New("s.partners.GetPartnerProducts error")
	cfg.partnersClient.EXPECT().GetPartnerProducts(gomock.Any(), gomock.Any()).Return(nil, targetError)
	_, err := cfg.service.getPartnerProducts(ctx)
	assert.True(t, errors.Is(err, targetError))

	cfg.partnersClient.EXPECT().GetPartnerProducts(gomock.Any(), gomock.Any()).Return(&products.GetAllResponse{}, nil).AnyTimes()

	// Test case 2: Success
	resp, err := cfg.service.getPartnerProducts(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

// go test -count=1 -v ./internal/gateway/ -run ^TestCheckOrder$
func TestCheckOrder(t *testing.T) {
	cfg := __SetupTestConfig(t)
	defer cfg.ctrl.Finish()

	ctx := context.Background()
	user := &user{}
	req := &checkRequest{
		Products: []*product{
			{ID: 1, Quantity: 1},
			{ID: 2, Quantity: 2},
		},
	}

	// Test case 1: Order is already checked
	cfg.cache.EXPECT().GetStruct(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
	err := cfg.service.checkOrder(ctx, user, req)
	assert.NoError(t, err)

	cfg.cache.EXPECT().GetStruct(gomock.Any(), gomock.Any(), gomock.Any()).Return(assert.AnError).AnyTimes()

	// Test case 2
	targetError := errors.New("s.partners.CheckPartnerProducts error")
	cfg.partnersClient.EXPECT().CheckPartnerProducts(gomock.Any(), gomock.Any()).Return(nil, targetError)
	err = cfg.service.checkOrder(ctx, user, req)
	assert.True(t, errors.Is(err, targetError))

	cfg.partnersClient.EXPECT().CheckPartnerProducts(gomock.Any(), gomock.Any()).Return(&partners.CheckResponse{}, nil).AnyTimes()

	// Test case 3
	targetError = errors.New("s.cache.SaveStruct error")
	cfg.cache.EXPECT().SaveStruct(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(targetError)
	err = cfg.service.checkOrder(ctx, user, req)
	assert.True(t, errors.Is(err, targetError))

	cfg.cache.EXPECT().SaveStruct(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	// Test case 4: success
	err = cfg.service.checkOrder(ctx, user, req)
	assert.NoError(t, err)
}

// go test -count=1 -v ./internal/gateway/ -run ^TestConfirmOrder$
func TestConfirmOrder(t *testing.T) {
	cfg := __SetupTestConfig(t)
	defer cfg.ctrl.Finish()

	ctx := context.Background()
	user := &user{}
	confirm := &confirmRequest{}

	// Test case 1
	targetError := errors.New("s.cache.GetStruct error")
	cfg.cache.EXPECT().GetStruct(gomock.Any(), gomock.Any(), gomock.Any()).Return(targetError)
	_, err := cfg.service.confirmOrder(ctx, user, confirm)
	assert.True(t, errors.Is(err, targetError))

	cfg.cache.EXPECT().GetStruct(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	// Test case #2
	targetError = errors.New("s.orders.CreateOrder error")
	cfg.ordersClient.EXPECT().CreateOrder(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, targetError)
	_, err = cfg.service.confirmOrder(ctx, user, confirm)
	assert.True(t, errors.Is(err, targetError))

	cfg.ordersClient.EXPECT().CreateOrder(gomock.Any(), gomock.Any(), gomock.Any()).Return(&orders.CreateResponse{}, nil)
	cfg.cache.EXPECT().Del(gomock.Any(), gomock.Any()).Return(nil)

	// Test case 3: Success
	resp, err := cfg.service.confirmOrder(ctx, user, confirm)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

// go test -count=1 -v ./internal/gateway/ -run ^TestPayOrder$
func TestPayOrder(t *testing.T) {
	cfg := __SetupTestConfig(t)
	defer cfg.ctrl.Finish()

	ctx := context.Background()
	req := &payRequest{}

	// Test case 1
	targetError := errors.New("s.orders.PayOrder error")
	cfg.ordersClient.EXPECT().PayOrder(gomock.Any(), gomock.Any()).Return(nil, targetError)
	_, err := cfg.service.payOrder(ctx, req)
	assert.True(t, errors.Is(err, targetError))

	cfg.ordersClient.EXPECT().PayOrder(gomock.Any(), gomock.Any()).Return(&orders.PayResponse{}, nil).AnyTimes()

	// Test case 2: Success
	resp, err := cfg.service.payOrder(ctx, req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

// go test -count=1 -v ./internal/gateway/ -run ^TestPickUpOrder$
func TestPickUpOrder(t *testing.T) {
	cfg := __SetupTestConfig(t)
	defer cfg.ctrl.Finish()

	ctx := context.Background()
	req := &pickupRequest{}

	// Test case #1
	targetError := errors.New("s.queue.Publish error")
	cfg.queue.EXPECT().Publish(gomock.Any(), gomock.Any(), gomock.Any()).Return(targetError)
	err := cfg.service.pickUpOrder(ctx, req)
	assert.True(t, errors.Is(err, targetError))

	cfg.queue.EXPECT().Publish(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	// Test case #2: Success
	err = cfg.service.pickUpOrder(ctx, req)
	assert.NoError(t, err)
}

// go test -count=1 -v ./internal/gateway/ -run ^TestAssignOrder$
func TestAssignOrder(t *testing.T) {
	cfg := __SetupTestConfig(t)
	defer cfg.ctrl.Finish()

	ctx := context.Background()
	user := &user{}
	req := &assignRequest{}

	// Test case #1
	targetError := errors.New("s.orders.AssignOrder error")
	cfg.ordersClient.EXPECT().AssignOrder(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, targetError)
	err := cfg.service.assignOrder(ctx, user, req)
	assert.True(t, errors.Is(err, targetError))

	cfg.ordersClient.EXPECT().AssignOrder(gomock.Any(), gomock.Any(), gomock.Any()).Return(&orders.AssignResponse{}, nil).AnyTimes()

	// Test case #2: Success
	err = cfg.service.assignOrder(ctx, user, req)
	assert.NoError(t, err)
}
