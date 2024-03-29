// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/bhbosman/gocommon/pubSub (interfaces: IPubSub)

// Package fullMarketDataManagerService is a generated GoMock package.
package pubSub

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockIPubSub is a mock of IPubSub interface.
type MockIPubSub struct {
	ctrl     *gomock.Controller
	recorder *MockIPubSubMockRecorder
}

// MockIPubSubMockRecorder is the mock recorder for MockIPubSub.
type MockIPubSubMockRecorder struct {
	mock *MockIPubSub
}

// NewMockIPubSub creates a new mock instance.
func NewMockIPubSub(ctrl *gomock.Controller) *MockIPubSub {
	mock := &MockIPubSub{ctrl: ctrl}
	mock.recorder = &MockIPubSubMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIPubSub) EXPECT() *MockIPubSubMockRecorder {
	return m.recorder
}

// Pub mocks base method.
func (m *MockIPubSub) Pub(arg0 interface{}, arg1 ...string) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Pub", varargs...)
}

// Pub indicates an expected call of Pub.
func (mr *MockIPubSubMockRecorder) Pub(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Pub", reflect.TypeOf((*MockIPubSub)(nil).Pub), varargs...)
}

// PubWithContext mocks base method.
func (m *MockIPubSub) PubWithContext(arg0 interface{}, arg1 ...string) bool {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "PubWithContext", varargs...)
	ret0, _ := ret[0].(bool)
	return ret0
}

// PubWithContext indicates an expected call of PubWithContext.
func (mr *MockIPubSubMockRecorder) PubWithContext(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PubWithContext", reflect.TypeOf((*MockIPubSub)(nil).PubWithContext), varargs...)
}

// argNames: [arg0 arg1]
// defaultArgs: [gomock.Any() gomock.Any()]
// defaultArgsAsString: gomock.Any(),gomock.Any()
// argTypes: [interface{} ...string]
// argString: arg0 interface{}, arg1 ...string
// rets: []
// retString:
// retString:
// ia: map[arg0:{} arg1:{}]
// idRecv: mr
// 1
func (mr *MockIPubSubMockRecorder) OnPubDoAndReturn(
	arg0, arg1 interface{},
	f func(arg0 interface{}, arg1 ...string)) *gomock.Call {
	return mr.
		Pub(arg0, arg1).
		DoAndReturn(f)
}

// 1
func (mr *MockIPubSubMockRecorder) OnPubDo(
	arg0, arg1 interface{},
	f func(arg0 interface{}, arg1 ...string)) *gomock.Call {
	return mr.
		Pub(arg0, arg1).
		Do(f)
}

// 1
func (mr *MockIPubSubMockRecorder) OnPubDoAndReturnDefault(
	f func(arg0 interface{}, arg1 ...string)) *gomock.Call {
	return mr.
		Pub(gomock.Any(), gomock.Any()).
		DoAndReturn(f)
}

// 1
func (mr *MockIPubSubMockRecorder) OnPubDoDefault(
	f func(arg0 interface{}, arg1 ...string)) *gomock.Call {
	return mr.
		Pub(gomock.Any(), gomock.Any()).
		Do(f)
}

// argNames: [arg0 arg1]
// defaultArgs: [gomock.Any() gomock.Any()]
// defaultArgsAsString: gomock.Any(),gomock.Any()
// argTypes: [interface{} ...string]
// argString: arg0 interface{}, arg1 ...string
// rets: [bool]
// retString: bool
// retString:  bool
// ia: map[arg0:{} arg1:{}]
// idRecv: mr
// 1
func (mr *MockIPubSubMockRecorder) OnPubWithContextDoAndReturn(
	arg0, arg1 interface{},
	f func(arg0 interface{}, arg1 ...string) bool) *gomock.Call {
	return mr.
		PubWithContext(arg0, arg1).
		DoAndReturn(f)
}

// 1
func (mr *MockIPubSubMockRecorder) OnPubWithContextDo(
	arg0, arg1 interface{},
	f func(arg0 interface{}, arg1 ...string)) *gomock.Call {
	return mr.
		PubWithContext(arg0, arg1).
		Do(f)
}

// 1
func (mr *MockIPubSubMockRecorder) OnPubWithContextDoAndReturnDefault(
	f func(arg0 interface{}, arg1 ...string) bool) *gomock.Call {
	return mr.
		PubWithContext(gomock.Any(), gomock.Any()).
		DoAndReturn(f)
}

// 1
func (mr *MockIPubSubMockRecorder) OnPubWithContextDoDefault(
	f func(arg0 interface{}, arg1 ...string)) *gomock.Call {
	return mr.
		PubWithContext(gomock.Any(), gomock.Any()).
		Do(f)
}

// retNames: [ret0]
// retArgs: [ret0 bool]
// retArgs22: ret0 bool
// 1
func (mr *MockIPubSubMockRecorder) OnPubWithContextReturn(
	arg0, arg1 interface{},
	ret0 bool) *gomock.Call {
	return mr.
		PubWithContext(arg0, arg1).
		Return(ret0)
}

// 1
func (mr *MockIPubSubMockRecorder) OnPubWithContextReturnDefault(
	ret0 bool) *gomock.Call {
	return mr.
		PubWithContext(gomock.Any(), gomock.Any()).
		Return(ret0)
}
