// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/bhbosman/gocommon/Services/interfaces (interfaces: IUniqueReferenceService)

// Package interfaces is a generated GoMock package.
package interfaces

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockIUniqueReferenceService is a mock of IUniqueReferenceService interface.
type MockIUniqueReferenceService struct {
	ctrl     *gomock.Controller
	recorder *MockIUniqueReferenceServiceMockRecorder
}

// MockIUniqueReferenceServiceMockRecorder is the mock recorder for MockIUniqueReferenceService.
type MockIUniqueReferenceServiceMockRecorder struct {
	mock *MockIUniqueReferenceService
}

// NewMockIUniqueReferenceService creates a new mock instance.
func NewMockIUniqueReferenceService(ctrl *gomock.Controller) *MockIUniqueReferenceService {
	mock := &MockIUniqueReferenceService{ctrl: ctrl}
	mock.recorder = &MockIUniqueReferenceServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIUniqueReferenceService) EXPECT() *MockIUniqueReferenceServiceMockRecorder {
	return m.recorder
}

// Next mocks base method.
func (m *MockIUniqueReferenceService) Next(arg0 string) string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Next", arg0)
	ret0, _ := ret[0].(string)
	return ret0
}

// Next indicates an expected call of Next.
func (mr *MockIUniqueReferenceServiceMockRecorder) Next(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Next", reflect.TypeOf((*MockIUniqueReferenceService)(nil).Next), arg0)
}

// argNames: [arg0]
// defaultArgs: [gomock.Any()]
// defaultArgsAsString: gomock.Any()
// argTypes: [string]
// argString: arg0 string
// rets: [string]
// retString: string
// retString:  string
// ia: map[arg0:{}]
// idRecv: mr
// 1
func (mr *MockIUniqueReferenceServiceMockRecorder) OnNextDoAndReturn(
	arg0 interface{},
	f func(arg0 string) string) *gomock.Call {
	return mr.
		Next(arg0).
		DoAndReturn(f)
}

// 1
func (mr *MockIUniqueReferenceServiceMockRecorder) OnNextDo(
	arg0 interface{},
	f func(arg0 string)) *gomock.Call {
	return mr.
		Next(arg0).
		Do(f)
}

// 1
func (mr *MockIUniqueReferenceServiceMockRecorder) OnNextDoAndReturnDefault(
	f func(arg0 string) string) *gomock.Call {
	return mr.
		Next(gomock.Any()).
		DoAndReturn(f)
}

// 1
func (mr *MockIUniqueReferenceServiceMockRecorder) OnNextDoDefault(
	f func(arg0 string)) *gomock.Call {
	return mr.
		Next(gomock.Any()).
		Do(f)
}

// retNames: [ret0]
// retArgs: [ret0 string]
// retArgs22: ret0 string
// 1
func (mr *MockIUniqueReferenceServiceMockRecorder) OnNextReturn(
	arg0 interface{},
	ret0 string) *gomock.Call {
	return mr.
		Next(arg0).
		Return(ret0)
}

// 1
func (mr *MockIUniqueReferenceServiceMockRecorder) OnNextReturnDefault(
	ret0 string) *gomock.Call {
	return mr.
		Next(gomock.Any()).
		Return(ret0)
}
