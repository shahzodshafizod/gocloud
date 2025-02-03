// Code generated by MockGen. DO NOT EDIT.
// Source: router.go
//
// Generated by this command:
//
//	mockgen -source=router.go -package=mocks -destination=mocks/router.go
//

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	http "net/http"
	reflect "reflect"
	time "time"

	validator "github.com/go-playground/validator/v10"
	pkg "github.com/shahzodshafizod/gocloud/pkg"
	gomock "go.uber.org/mock/gomock"
)

// MockRouter is a mock of Router interface.
type MockRouter struct {
	ctrl     *gomock.Controller
	recorder *MockRouterMockRecorder
	isgomock struct{}
}

// MockRouterMockRecorder is the mock recorder for MockRouter.
type MockRouterMockRecorder struct {
	mock *MockRouter
}

// NewMockRouter creates a new mock instance.
func NewMockRouter(ctrl *gomock.Controller) *MockRouter {
	mock := &MockRouter{ctrl: ctrl}
	mock.recorder = &MockRouterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRouter) EXPECT() *MockRouterMockRecorder {
	return m.recorder
}

// DELETE mocks base method.
func (m *MockRouter) DELETE(path, actionName string, handler pkg.Handler, middlewares ...pkg.Middleware) {
	m.ctrl.T.Helper()
	varargs := []any{path, actionName, handler}
	for _, a := range middlewares {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "DELETE", varargs...)
}

// DELETE indicates an expected call of DELETE.
func (mr *MockRouterMockRecorder) DELETE(path, actionName, handler any, middlewares ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{path, actionName, handler}, middlewares...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DELETE", reflect.TypeOf((*MockRouter)(nil).DELETE), varargs...)
}

// GET mocks base method.
func (m *MockRouter) GET(path, actionName string, handler pkg.Handler, middlewares ...pkg.Middleware) {
	m.ctrl.T.Helper()
	varargs := []any{path, actionName, handler}
	for _, a := range middlewares {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "GET", varargs...)
}

// GET indicates an expected call of GET.
func (mr *MockRouterMockRecorder) GET(path, actionName, handler any, middlewares ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{path, actionName, handler}, middlewares...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GET", reflect.TypeOf((*MockRouter)(nil).GET), varargs...)
}

// POST mocks base method.
func (m *MockRouter) POST(path, actionName string, handler pkg.Handler, middlewares ...pkg.Middleware) {
	m.ctrl.T.Helper()
	varargs := []any{path, actionName, handler}
	for _, a := range middlewares {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "POST", varargs...)
}

// POST indicates an expected call of POST.
func (mr *MockRouterMockRecorder) POST(path, actionName, handler any, middlewares ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{path, actionName, handler}, middlewares...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "POST", reflect.TypeOf((*MockRouter)(nil).POST), varargs...)
}

// PUT mocks base method.
func (m *MockRouter) PUT(path, actionName string, handler pkg.Handler, middlewares ...pkg.Middleware) {
	m.ctrl.T.Helper()
	varargs := []any{path, actionName, handler}
	for _, a := range middlewares {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "PUT", varargs...)
}

// PUT indicates an expected call of PUT.
func (mr *MockRouterMockRecorder) PUT(path, actionName, handler any, middlewares ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{path, actionName, handler}, middlewares...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PUT", reflect.TypeOf((*MockRouter)(nil).PUT), varargs...)
}

// Serve mocks base method.
func (m *MockRouter) Serve(ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Serve", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// Serve indicates an expected call of Serve.
func (mr *MockRouterMockRecorder) Serve(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Serve", reflect.TypeOf((*MockRouter)(nil).Serve), ctx)
}

// Shutdown mocks base method.
func (m *MockRouter) Shutdown(ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Shutdown", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// Shutdown indicates an expected call of Shutdown.
func (mr *MockRouterMockRecorder) Shutdown(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Shutdown", reflect.TypeOf((*MockRouter)(nil).Shutdown), ctx)
}

// WrapHandler mocks base method.
func (m *MockRouter) WrapHandler(handler http.HandlerFunc) pkg.Handler {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WrapHandler", handler)
	ret0, _ := ret[0].(pkg.Handler)
	return ret0
}

// WrapHandler indicates an expected call of WrapHandler.
func (mr *MockRouterMockRecorder) WrapHandler(handler any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WrapHandler", reflect.TypeOf((*MockRouter)(nil).WrapHandler), handler)
}

// MockContext is a mock of Context interface.
type MockContext struct {
	ctrl     *gomock.Controller
	recorder *MockContextMockRecorder
	isgomock struct{}
}

// MockContextMockRecorder is the mock recorder for MockContext.
type MockContextMockRecorder struct {
	mock *MockContext
}

// NewMockContext creates a new mock instance.
func NewMockContext(ctrl *gomock.Controller) *MockContext {
	mock := &MockContext{ctrl: ctrl}
	mock.recorder = &MockContextMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockContext) EXPECT() *MockContextMockRecorder {
	return m.recorder
}

// GetCookie mocks base method.
func (m *MockContext) GetCookie(name string) string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCookie", name)
	ret0, _ := ret[0].(string)
	return ret0
}

// GetCookie indicates an expected call of GetCookie.
func (mr *MockContextMockRecorder) GetCookie(name any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCookie", reflect.TypeOf((*MockContext)(nil).GetCookie), name)
}

// GetFormValue mocks base method.
func (m *MockContext) GetFormValue(key string) string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFormValue", key)
	ret0, _ := ret[0].(string)
	return ret0
}

// GetFormValue indicates an expected call of GetFormValue.
func (mr *MockContextMockRecorder) GetFormValue(key any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFormValue", reflect.TypeOf((*MockContext)(nil).GetFormValue), key)
}

// GetHeader mocks base method.
func (m *MockContext) GetHeader(key string) string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetHeader", key)
	ret0, _ := ret[0].(string)
	return ret0
}

// GetHeader indicates an expected call of GetHeader.
func (mr *MockContextMockRecorder) GetHeader(key any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetHeader", reflect.TypeOf((*MockContext)(nil).GetHeader), key)
}

// GetParam mocks base method.
func (m *MockContext) GetParam(key string) string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetParam", key)
	ret0, _ := ret[0].(string)
	return ret0
}

// GetParam indicates an expected call of GetParam.
func (mr *MockContextMockRecorder) GetParam(key any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetParam", reflect.TypeOf((*MockContext)(nil).GetParam), key)
}

// GetQueryValue mocks base method.
func (m *MockContext) GetQueryValue(key string) string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetQueryValue", key)
	ret0, _ := ret[0].(string)
	return ret0
}

// GetQueryValue indicates an expected call of GetQueryValue.
func (mr *MockContextMockRecorder) GetQueryValue(key any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetQueryValue", reflect.TypeOf((*MockContext)(nil).GetQueryValue), key)
}

// GetRequestContext mocks base method.
func (m *MockContext) GetRequestContext() context.Context {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRequestContext")
	ret0, _ := ret[0].(context.Context)
	return ret0
}

// GetRequestContext indicates an expected call of GetRequestContext.
func (mr *MockContextMockRecorder) GetRequestContext() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRequestContext", reflect.TypeOf((*MockContext)(nil).GetRequestContext))
}

// GetValue mocks base method.
func (m *MockContext) GetValue(key string) any {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetValue", key)
	ret0, _ := ret[0].(any)
	return ret0
}

// GetValue indicates an expected call of GetValue.
func (mr *MockContextMockRecorder) GetValue(key any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetValue", reflect.TypeOf((*MockContext)(nil).GetValue), key)
}

// OpenFormFile mocks base method.
func (m *MockContext) OpenFormFile(key string) (pkg.File, pkg.FileInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "OpenFormFile", key)
	ret0, _ := ret[0].(pkg.File)
	ret1, _ := ret[1].(pkg.FileInfo)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// OpenFormFile indicates an expected call of OpenFormFile.
func (mr *MockContextMockRecorder) OpenFormFile(key any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OpenFormFile", reflect.TypeOf((*MockContext)(nil).OpenFormFile), key)
}

// ParseBody mocks base method.
func (m *MockContext) ParseBody(v any) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ParseBody", v)
	ret0, _ := ret[0].(error)
	return ret0
}

// ParseBody indicates an expected call of ParseBody.
func (mr *MockContextMockRecorder) ParseBody(v any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ParseBody", reflect.TypeOf((*MockContext)(nil).ParseBody), v)
}

// Redirect mocks base method.
func (m *MockContext) Redirect(url string, code int) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Redirect", url, code)
}

// Redirect indicates an expected call of Redirect.
func (mr *MockContextMockRecorder) Redirect(url, code any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Redirect", reflect.TypeOf((*MockContext)(nil).Redirect), url, code)
}

// Respond mocks base method.
func (m *MockContext) Respond(r pkg.Response) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Respond", r)
}

// Respond indicates an expected call of Respond.
func (mr *MockContextMockRecorder) Respond(r any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Respond", reflect.TypeOf((*MockContext)(nil).Respond), r)
}

// SaveCookie mocks base method.
func (m *MockContext) SaveCookie(name, value string, expiresIn time.Duration) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SaveCookie", name, value, expiresIn)
}

// SaveCookie indicates an expected call of SaveCookie.
func (mr *MockContextMockRecorder) SaveCookie(name, value, expiresIn any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveCookie", reflect.TypeOf((*MockContext)(nil).SaveCookie), name, value, expiresIn)
}

// SaveValue mocks base method.
func (m *MockContext) SaveValue(key string, value any) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SaveValue", key, value)
}

// SaveValue indicates an expected call of SaveValue.
func (mr *MockContextMockRecorder) SaveValue(key, value any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveValue", reflect.TypeOf((*MockContext)(nil).SaveValue), key, value)
}

// ServeFile mocks base method.
func (m *MockContext) ServeFile(filename string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ServeFile", filename)
}

// ServeFile indicates an expected call of ServeFile.
func (mr *MockContextMockRecorder) ServeFile(filename any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ServeFile", reflect.TypeOf((*MockContext)(nil).ServeFile), filename)
}

// StartSpan mocks base method.
func (m *MockContext) StartSpan() (context.Context, pkg.Span) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StartSpan")
	ret0, _ := ret[0].(context.Context)
	ret1, _ := ret[1].(pkg.Span)
	return ret0, ret1
}

// StartSpan indicates an expected call of StartSpan.
func (mr *MockContextMockRecorder) StartSpan() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StartSpan", reflect.TypeOf((*MockContext)(nil).StartSpan))
}

// ValidateStruct mocks base method.
func (m *MockContext) ValidateStruct(v any) []string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ValidateStruct", v)
	ret0, _ := ret[0].([]string)
	return ret0
}

// ValidateStruct indicates an expected call of ValidateStruct.
func (mr *MockContextMockRecorder) ValidateStruct(v any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidateStruct", reflect.TypeOf((*MockContext)(nil).ValidateStruct), v)
}

// ValidateVar mocks base method.
func (m *MockContext) ValidateVar(v any, tag string) string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ValidateVar", v, tag)
	ret0, _ := ret[0].(string)
	return ret0
}

// ValidateVar indicates an expected call of ValidateVar.
func (mr *MockContextMockRecorder) ValidateVar(v, tag any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidateVar", reflect.TypeOf((*MockContext)(nil).ValidateVar), v, tag)
}

// MockFile is a mock of File interface.
type MockFile struct {
	ctrl     *gomock.Controller
	recorder *MockFileMockRecorder
	isgomock struct{}
}

// MockFileMockRecorder is the mock recorder for MockFile.
type MockFileMockRecorder struct {
	mock *MockFile
}

// NewMockFile creates a new mock instance.
func NewMockFile(ctrl *gomock.Controller) *MockFile {
	mock := &MockFile{ctrl: ctrl}
	mock.recorder = &MockFileMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockFile) EXPECT() *MockFileMockRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *MockFile) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockFileMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockFile)(nil).Close))
}

// Read mocks base method.
func (m *MockFile) Read(p []byte) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Read", p)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Read indicates an expected call of Read.
func (mr *MockFileMockRecorder) Read(p any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Read", reflect.TypeOf((*MockFile)(nil).Read), p)
}

// ReadAt mocks base method.
func (m *MockFile) ReadAt(p []byte, off int64) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadAt", p, off)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReadAt indicates an expected call of ReadAt.
func (mr *MockFileMockRecorder) ReadAt(p, off any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadAt", reflect.TypeOf((*MockFile)(nil).ReadAt), p, off)
}

// Seek mocks base method.
func (m *MockFile) Seek(offset int64, whence int) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Seek", offset, whence)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Seek indicates an expected call of Seek.
func (mr *MockFileMockRecorder) Seek(offset, whence any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Seek", reflect.TypeOf((*MockFile)(nil).Seek), offset, whence)
}

// MockFileInfo is a mock of FileInfo interface.
type MockFileInfo struct {
	ctrl     *gomock.Controller
	recorder *MockFileInfoMockRecorder
	isgomock struct{}
}

// MockFileInfoMockRecorder is the mock recorder for MockFileInfo.
type MockFileInfoMockRecorder struct {
	mock *MockFileInfo
}

// NewMockFileInfo creates a new mock instance.
func NewMockFileInfo(ctrl *gomock.Controller) *MockFileInfo {
	mock := &MockFileInfo{ctrl: ctrl}
	mock.recorder = &MockFileInfoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockFileInfo) EXPECT() *MockFileInfoMockRecorder {
	return m.recorder
}

// ContentType mocks base method.
func (m *MockFileInfo) ContentType() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ContentType")
	ret0, _ := ret[0].(string)
	return ret0
}

// ContentType indicates an expected call of ContentType.
func (mr *MockFileInfoMockRecorder) ContentType() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ContentType", reflect.TypeOf((*MockFileInfo)(nil).ContentType))
}

// FileName mocks base method.
func (m *MockFileInfo) FileName() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FileName")
	ret0, _ := ret[0].(string)
	return ret0
}

// FileName indicates an expected call of FileName.
func (mr *MockFileInfoMockRecorder) FileName() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FileName", reflect.TypeOf((*MockFileInfo)(nil).FileName))
}

// FileSize mocks base method.
func (m *MockFileInfo) FileSize() int64 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FileSize")
	ret0, _ := ret[0].(int64)
	return ret0
}

// FileSize indicates an expected call of FileSize.
func (mr *MockFileInfoMockRecorder) FileSize() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FileSize", reflect.TypeOf((*MockFileInfo)(nil).FileSize))
}

// MockResponse is a mock of Response interface.
type MockResponse struct {
	ctrl     *gomock.Controller
	recorder *MockResponseMockRecorder
	isgomock struct{}
}

// MockResponseMockRecorder is the mock recorder for MockResponse.
type MockResponseMockRecorder struct {
	mock *MockResponse
}

// NewMockResponse creates a new mock instance.
func NewMockResponse(ctrl *gomock.Controller) *MockResponse {
	mock := &MockResponse{ctrl: ctrl}
	mock.recorder = &MockResponseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockResponse) EXPECT() *MockResponseMockRecorder {
	return m.recorder
}

// GetCode mocks base method.
func (m *MockResponse) GetCode() int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCode")
	ret0, _ := ret[0].(int)
	return ret0
}

// GetCode indicates an expected call of GetCode.
func (mr *MockResponseMockRecorder) GetCode() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCode", reflect.TypeOf((*MockResponse)(nil).GetCode))
}

// MockValidator is a mock of Validator interface.
type MockValidator struct {
	ctrl     *gomock.Controller
	recorder *MockValidatorMockRecorder
	isgomock struct{}
}

// MockValidatorMockRecorder is the mock recorder for MockValidator.
type MockValidatorMockRecorder struct {
	mock *MockValidator
}

// NewMockValidator creates a new mock instance.
func NewMockValidator(ctrl *gomock.Controller) *MockValidator {
	mock := &MockValidator{ctrl: ctrl}
	mock.recorder = &MockValidatorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockValidator) EXPECT() *MockValidatorMockRecorder {
	return m.recorder
}

// GetFunc mocks base method.
func (m *MockValidator) GetFunc() validator.FuncCtx {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFunc")
	ret0, _ := ret[0].(validator.FuncCtx)
	return ret0
}

// GetFunc indicates an expected call of GetFunc.
func (mr *MockValidatorMockRecorder) GetFunc() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFunc", reflect.TypeOf((*MockValidator)(nil).GetFunc))
}

// GetTag mocks base method.
func (m *MockValidator) GetTag() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTag")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetTag indicates an expected call of GetTag.
func (mr *MockValidatorMockRecorder) GetTag() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTag", reflect.TypeOf((*MockValidator)(nil).GetTag))
}
