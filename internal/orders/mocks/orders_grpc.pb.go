// Code generated by MockGen. DO NOT EDIT.
// Source: internal/orders/orders_grpc.pb.go
//
// Generated by this command:
//
//	mockgen -source=internal/orders/orders_grpc.pb.go -destination=internal/orders/mocks/orders_grpc.pb.go -package=mocks
//

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	orders "github.com/shahzodshafizod/gocloud/internal/orders"
	gomock "go.uber.org/mock/gomock"
	grpc "google.golang.org/grpc"
)

// MockOrdersClient is a mock of OrdersClient interface.
type MockOrdersClient struct {
	ctrl     *gomock.Controller
	recorder *MockOrdersClientMockRecorder
	isgomock struct{}
}

// MockOrdersClientMockRecorder is the mock recorder for MockOrdersClient.
type MockOrdersClientMockRecorder struct {
	mock *MockOrdersClient
}

// NewMockOrdersClient creates a new mock instance.
func NewMockOrdersClient(ctrl *gomock.Controller) *MockOrdersClient {
	mock := &MockOrdersClient{ctrl: ctrl}
	mock.recorder = &MockOrdersClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockOrdersClient) EXPECT() *MockOrdersClientMockRecorder {
	return m.recorder
}

// AssignOrder mocks base method.
func (m *MockOrdersClient) AssignOrder(ctx context.Context, in *orders.AssignRequest, opts ...grpc.CallOption) (*orders.AssignResponse, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "AssignOrder", varargs...)
	ret0, _ := ret[0].(*orders.AssignResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AssignOrder indicates an expected call of AssignOrder.
func (mr *MockOrdersClientMockRecorder) AssignOrder(ctx, in any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AssignOrder", reflect.TypeOf((*MockOrdersClient)(nil).AssignOrder), varargs...)
}

// CreateOrder mocks base method.
func (m *MockOrdersClient) CreateOrder(ctx context.Context, in *orders.Order, opts ...grpc.CallOption) (*orders.CreateResponse, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "CreateOrder", varargs...)
	ret0, _ := ret[0].(*orders.CreateResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateOrder indicates an expected call of CreateOrder.
func (mr *MockOrdersClientMockRecorder) CreateOrder(ctx, in any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateOrder", reflect.TypeOf((*MockOrdersClient)(nil).CreateOrder), varargs...)
}

// PayOrder mocks base method.
func (m *MockOrdersClient) PayOrder(ctx context.Context, in *orders.PayRequest, opts ...grpc.CallOption) (*orders.PayResponse, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "PayOrder", varargs...)
	ret0, _ := ret[0].(*orders.PayResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PayOrder indicates an expected call of PayOrder.
func (mr *MockOrdersClientMockRecorder) PayOrder(ctx, in any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PayOrder", reflect.TypeOf((*MockOrdersClient)(nil).PayOrder), varargs...)
}

// MockOrdersServer is a mock of OrdersServer interface.
type MockOrdersServer struct {
	ctrl     *gomock.Controller
	recorder *MockOrdersServerMockRecorder
	isgomock struct{}
}

// MockOrdersServerMockRecorder is the mock recorder for MockOrdersServer.
type MockOrdersServerMockRecorder struct {
	mock *MockOrdersServer
}

// NewMockOrdersServer creates a new mock instance.
func NewMockOrdersServer(ctrl *gomock.Controller) *MockOrdersServer {
	mock := &MockOrdersServer{ctrl: ctrl}
	mock.recorder = &MockOrdersServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockOrdersServer) EXPECT() *MockOrdersServerMockRecorder {
	return m.recorder
}

// AssignOrder mocks base method.
func (m *MockOrdersServer) AssignOrder(arg0 context.Context, arg1 *orders.AssignRequest) (*orders.AssignResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AssignOrder", arg0, arg1)
	ret0, _ := ret[0].(*orders.AssignResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AssignOrder indicates an expected call of AssignOrder.
func (mr *MockOrdersServerMockRecorder) AssignOrder(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AssignOrder", reflect.TypeOf((*MockOrdersServer)(nil).AssignOrder), arg0, arg1)
}

// CreateOrder mocks base method.
func (m *MockOrdersServer) CreateOrder(arg0 context.Context, arg1 *orders.Order) (*orders.CreateResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateOrder", arg0, arg1)
	ret0, _ := ret[0].(*orders.CreateResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateOrder indicates an expected call of CreateOrder.
func (mr *MockOrdersServerMockRecorder) CreateOrder(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateOrder", reflect.TypeOf((*MockOrdersServer)(nil).CreateOrder), arg0, arg1)
}

// PayOrder mocks base method.
func (m *MockOrdersServer) PayOrder(arg0 context.Context, arg1 *orders.PayRequest) (*orders.PayResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PayOrder", arg0, arg1)
	ret0, _ := ret[0].(*orders.PayResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PayOrder indicates an expected call of PayOrder.
func (mr *MockOrdersServerMockRecorder) PayOrder(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PayOrder", reflect.TypeOf((*MockOrdersServer)(nil).PayOrder), arg0, arg1)
}

// mustEmbedUnimplementedOrdersServer mocks base method.
func (m *MockOrdersServer) mustEmbedUnimplementedOrdersServer() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "mustEmbedUnimplementedOrdersServer")
}

// mustEmbedUnimplementedOrdersServer indicates an expected call of mustEmbedUnimplementedOrdersServer.
func (mr *MockOrdersServerMockRecorder) mustEmbedUnimplementedOrdersServer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "mustEmbedUnimplementedOrdersServer", reflect.TypeOf((*MockOrdersServer)(nil).mustEmbedUnimplementedOrdersServer))
}

// MockUnsafeOrdersServer is a mock of UnsafeOrdersServer interface.
type MockUnsafeOrdersServer struct {
	ctrl     *gomock.Controller
	recorder *MockUnsafeOrdersServerMockRecorder
	isgomock struct{}
}

// MockUnsafeOrdersServerMockRecorder is the mock recorder for MockUnsafeOrdersServer.
type MockUnsafeOrdersServerMockRecorder struct {
	mock *MockUnsafeOrdersServer
}

// NewMockUnsafeOrdersServer creates a new mock instance.
func NewMockUnsafeOrdersServer(ctrl *gomock.Controller) *MockUnsafeOrdersServer {
	mock := &MockUnsafeOrdersServer{ctrl: ctrl}
	mock.recorder = &MockUnsafeOrdersServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUnsafeOrdersServer) EXPECT() *MockUnsafeOrdersServerMockRecorder {
	return m.recorder
}

// mustEmbedUnimplementedOrdersServer mocks base method.
func (m *MockUnsafeOrdersServer) mustEmbedUnimplementedOrdersServer() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "mustEmbedUnimplementedOrdersServer")
}

// mustEmbedUnimplementedOrdersServer indicates an expected call of mustEmbedUnimplementedOrdersServer.
func (mr *MockUnsafeOrdersServerMockRecorder) mustEmbedUnimplementedOrdersServer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "mustEmbedUnimplementedOrdersServer", reflect.TypeOf((*MockUnsafeOrdersServer)(nil).mustEmbedUnimplementedOrdersServer))
}
