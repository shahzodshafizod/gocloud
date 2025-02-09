// Code generated by MockGen. DO NOT EDIT.
// Source: auth.go
//
// Generated by this command:
//
//	mockgen -source=auth.go -package=mocks -destination=mocks/auth.go
//

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	pkg "github.com/shahzodshafizod/gocloud/pkg"
	gomock "go.uber.org/mock/gomock"
)

// MockAuth is a mock of Auth interface.
type MockAuth struct {
	ctrl     *gomock.Controller
	recorder *MockAuthMockRecorder
	isgomock struct{}
}

// MockAuthMockRecorder is the mock recorder for MockAuth.
type MockAuthMockRecorder struct {
	mock *MockAuth
}

// NewMockAuth creates a new mock instance.
func NewMockAuth(ctrl *gomock.Controller) *MockAuth {
	mock := &MockAuth{ctrl: ctrl}
	mock.recorder = &MockAuthMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAuth) EXPECT() *MockAuthMockRecorder {
	return m.recorder
}

// ChangePassword mocks base method.
func (m *MockAuth) ChangePassword(ctx context.Context, change *pkg.ChangePassword) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ChangePassword", ctx, change)
	ret0, _ := ret[0].(error)
	return ret0
}

// ChangePassword indicates an expected call of ChangePassword.
func (mr *MockAuthMockRecorder) ChangePassword(ctx, change any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ChangePassword", reflect.TypeOf((*MockAuth)(nil).ChangePassword), ctx, change)
}

// CheckToken mocks base method.
func (m *MockAuth) CheckToken(ctx context.Context, accessToken string) (*pkg.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckToken", ctx, accessToken)
	ret0, _ := ret[0].(*pkg.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckToken indicates an expected call of CheckToken.
func (mr *MockAuthMockRecorder) CheckToken(ctx, accessToken any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckToken", reflect.TypeOf((*MockAuth)(nil).CheckToken), ctx, accessToken)
}

// ConfirmChangeEmail mocks base method.
func (m *MockAuth) ConfirmChangeEmail(ctx context.Context, accessToken string, verifyEmail *pkg.VerifyEmail) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ConfirmChangeEmail", ctx, accessToken, verifyEmail)
	ret0, _ := ret[0].(error)
	return ret0
}

// ConfirmChangeEmail indicates an expected call of ConfirmChangeEmail.
func (mr *MockAuthMockRecorder) ConfirmChangeEmail(ctx, accessToken, verifyEmail any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ConfirmChangeEmail", reflect.TypeOf((*MockAuth)(nil).ConfirmChangeEmail), ctx, accessToken, verifyEmail)
}

// ConfirmSignUp mocks base method.
func (m *MockAuth) ConfirmSignUp(ctx context.Context, verify *pkg.VerifyEmail) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ConfirmSignUp", ctx, verify)
	ret0, _ := ret[0].(error)
	return ret0
}

// ConfirmSignUp indicates an expected call of ConfirmSignUp.
func (mr *MockAuthMockRecorder) ConfirmSignUp(ctx, verify any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ConfirmSignUp", reflect.TypeOf((*MockAuth)(nil).ConfirmSignUp), ctx, verify)
}

// DeleteUser mocks base method.
func (m *MockAuth) DeleteUser(ctx context.Context, userID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteUser", ctx, userID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteUser indicates an expected call of DeleteUser.
func (mr *MockAuthMockRecorder) DeleteUser(ctx, userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteUser", reflect.TypeOf((*MockAuth)(nil).DeleteUser), ctx, userID)
}

// ForgotPassword mocks base method.
func (m *MockAuth) ForgotPassword(ctx context.Context, userID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ForgotPassword", ctx, userID)
	ret0, _ := ret[0].(error)
	return ret0
}

// ForgotPassword indicates an expected call of ForgotPassword.
func (mr *MockAuthMockRecorder) ForgotPassword(ctx, userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ForgotPassword", reflect.TypeOf((*MockAuth)(nil).ForgotPassword), ctx, userID)
}

// RefreshToken mocks base method.
func (m *MockAuth) RefreshToken(ctx context.Context, userID, refreshToken string) (*pkg.Token, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RefreshToken", ctx, userID, refreshToken)
	ret0, _ := ret[0].(*pkg.Token)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RefreshToken indicates an expected call of RefreshToken.
func (mr *MockAuthMockRecorder) RefreshToken(ctx, userID, refreshToken any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RefreshToken", reflect.TypeOf((*MockAuth)(nil).RefreshToken), ctx, userID, refreshToken)
}

// ResetPassword mocks base method.
func (m *MockAuth) ResetPassword(ctx context.Context, reset *pkg.ResetPassword) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ResetPassword", ctx, reset)
	ret0, _ := ret[0].(error)
	return ret0
}

// ResetPassword indicates an expected call of ResetPassword.
func (mr *MockAuthMockRecorder) ResetPassword(ctx, reset any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ResetPassword", reflect.TypeOf((*MockAuth)(nil).ResetPassword), ctx, reset)
}

// SignIn mocks base method.
func (m *MockAuth) SignIn(ctx context.Context, signIn *pkg.SignIn) (*pkg.User, *pkg.Token, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SignIn", ctx, signIn)
	ret0, _ := ret[0].(*pkg.User)
	ret1, _ := ret[1].(*pkg.Token)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// SignIn indicates an expected call of SignIn.
func (mr *MockAuthMockRecorder) SignIn(ctx, signIn any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SignIn", reflect.TypeOf((*MockAuth)(nil).SignIn), ctx, signIn)
}

// SignOut mocks base method.
func (m *MockAuth) SignOut(ctx context.Context, userID, refreshToken string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SignOut", ctx, userID, refreshToken)
	ret0, _ := ret[0].(error)
	return ret0
}

// SignOut indicates an expected call of SignOut.
func (mr *MockAuthMockRecorder) SignOut(ctx, userID, refreshToken any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SignOut", reflect.TypeOf((*MockAuth)(nil).SignOut), ctx, userID, refreshToken)
}

// SignUp mocks base method.
func (m *MockAuth) SignUp(ctx context.Context, signUp *pkg.SignUp) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SignUp", ctx, signUp)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SignUp indicates an expected call of SignUp.
func (mr *MockAuthMockRecorder) SignUp(ctx, signUp any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SignUp", reflect.TypeOf((*MockAuth)(nil).SignUp), ctx, signUp)
}

// UpdateUser mocks base method.
func (m *MockAuth) UpdateUser(ctx context.Context, accessToken string, user *pkg.UpdateUser) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUser", ctx, accessToken, user)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateUser indicates an expected call of UpdateUser.
func (mr *MockAuthMockRecorder) UpdateUser(ctx, accessToken, user any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUser", reflect.TypeOf((*MockAuth)(nil).UpdateUser), ctx, accessToken, user)
}
