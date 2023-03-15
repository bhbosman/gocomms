// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/bhbosman/gocomms/intf (interfaces: IConnectionReactor,IInitParams)

// Package intf is a generated GoMock package.
package intf

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	rxgo "github.com/reactivex/rxgo/v2"
)

// MockIConnectionReactor is a mock of IConnectionReactor interface.
type MockIConnectionReactor struct {
	ctrl     *gomock.Controller
	recorder *MockIConnectionReactorMockRecorder
}

// MockIConnectionReactorMockRecorder is the mock recorder for MockIConnectionReactor.
type MockIConnectionReactorMockRecorder struct {
	mock *MockIConnectionReactor
}

// NewMockIConnectionReactor creates a new mock instance.
func NewMockIConnectionReactor(ctrl *gomock.Controller) *MockIConnectionReactor {
	mock := &MockIConnectionReactor{ctrl: ctrl}
	mock.recorder = &MockIConnectionReactorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIConnectionReactor) EXPECT() *MockIConnectionReactorMockRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *MockIConnectionReactor) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockIConnectionReactorMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockIConnectionReactor)(nil).Close))
}

// Init mocks base method.
func (m *MockIConnectionReactor) Init(arg0 IInitParams) (rxgo.NextFunc, rxgo.ErrFunc, rxgo.CompletedFunc, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Init", arg0)
	ret0, _ := ret[0].(rxgo.NextFunc)
	ret1, _ := ret[1].(rxgo.ErrFunc)
	ret2, _ := ret[2].(rxgo.CompletedFunc)
	ret3, _ := ret[3].(error)
	return ret0, ret1, ret2, ret3
}

// Init indicates an expected call of Init.
func (mr *MockIConnectionReactorMockRecorder) Init(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Init", reflect.TypeOf((*MockIConnectionReactor)(nil).Init), arg0)
}

// Open mocks base method.
func (m *MockIConnectionReactor) Open() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Open")
	ret0, _ := ret[0].(error)
	return ret0
}

// Open indicates an expected call of Open.
func (mr *MockIConnectionReactorMockRecorder) Open() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Open", reflect.TypeOf((*MockIConnectionReactor)(nil).Open))
}

// argNames: []
// defaultArgs: []
// defaultArgsAsString:
// argTypes: []
// argString:
// rets: [error]
// retString: error
// retString:  error
// ia: map[]
// idRecv: mr
// 0
func (mr *MockIConnectionReactorMockRecorder) OnCloseDoAndReturn(
	f func() error) *gomock.Call {
	return mr.
		Close().
		DoAndReturn(f)
}

// 0
func (mr *MockIConnectionReactorMockRecorder) OnCloseDo(
	f func()) *gomock.Call {
	return mr.
		Close().
		DoAndReturn(f)
}

// retNames: [ret0]
// retArgs: [ret0 error]
// retArgs22: ret0 error
// 1
func (mr *MockIConnectionReactorMockRecorder) OnCloseReturn(ret0 error) *gomock.Call {
	return mr.
		Close().
		Return(ret0)
}

// argNames: [arg0]
// defaultArgs: [gomock.Any()]
// defaultArgsAsString: gomock.Any()
// argTypes: [IInitParams]
// argString: arg0 IInitParams
// rets: [rxgo.NextFunc rxgo.ErrFunc rxgo.CompletedFunc error]
// retString: rxgo.NextFunc, rxgo.ErrFunc, rxgo.CompletedFunc, error
// retString:  (rxgo.NextFunc, rxgo.ErrFunc, rxgo.CompletedFunc, error)
// ia: map[arg0:{}]
// idRecv: mr
// 1
func (mr *MockIConnectionReactorMockRecorder) OnInitDoAndReturn(
	arg0 interface{},
	f func(arg0 IInitParams) (rxgo.NextFunc, rxgo.ErrFunc, rxgo.CompletedFunc, error)) *gomock.Call {
	return mr.
		Init(arg0).
		DoAndReturn(f)
}

// 1
func (mr *MockIConnectionReactorMockRecorder) OnInitDo(
	arg0 interface{},
	f func(arg0 IInitParams)) *gomock.Call {
	return mr.
		Init(arg0).
		Do(f)
}

// 1
func (mr *MockIConnectionReactorMockRecorder) OnInitDoAndReturnDefault(
	f func(arg0 IInitParams) (rxgo.NextFunc, rxgo.ErrFunc, rxgo.CompletedFunc, error)) *gomock.Call {
	return mr.
		Init(gomock.Any()).
		DoAndReturn(f)
}

// 1
func (mr *MockIConnectionReactorMockRecorder) OnInitDoDefault(
	f func(arg0 IInitParams)) *gomock.Call {
	return mr.
		Init(gomock.Any()).
		Do(f)
}

// retNames: [ret0 ret1 ret2 ret3]
// retArgs: [ret0 rxgo.NextFunc ret1 rxgo.ErrFunc ret2 rxgo.CompletedFunc ret3 error]
// retArgs22: ret0 rxgo.NextFunc,ret1 rxgo.ErrFunc,ret2 rxgo.CompletedFunc,ret3 error
// 1
func (mr *MockIConnectionReactorMockRecorder) OnInitReturn(
	arg0 interface{},
	ret0 rxgo.NextFunc, ret1 rxgo.ErrFunc, ret2 rxgo.CompletedFunc, ret3 error) *gomock.Call {
	return mr.
		Init(arg0).
		Return(ret0, ret1, ret2, ret3)
}

// 1
func (mr *MockIConnectionReactorMockRecorder) OnInitReturnDefault(
	ret0 rxgo.NextFunc, ret1 rxgo.ErrFunc, ret2 rxgo.CompletedFunc, ret3 error) *gomock.Call {
	return mr.
		Init(gomock.Any()).
		Return(ret0, ret1, ret2, ret3)
}

// argNames: []
// defaultArgs: []
// defaultArgsAsString:
// argTypes: []
// argString:
// rets: [error]
// retString: error
// retString:  error
// ia: map[]
// idRecv: mr
// 0
func (mr *MockIConnectionReactorMockRecorder) OnOpenDoAndReturn(
	f func() error) *gomock.Call {
	return mr.
		Open().
		DoAndReturn(f)
}

// 0
func (mr *MockIConnectionReactorMockRecorder) OnOpenDo(
	f func()) *gomock.Call {
	return mr.
		Open().
		DoAndReturn(f)
}

// retNames: [ret0]
// retArgs: [ret0 error]
// retArgs22: ret0 error
// 1
func (mr *MockIConnectionReactorMockRecorder) OnOpenReturn(ret0 error) *gomock.Call {
	return mr.
		Open().
		Return(ret0)
}

// MockIInitParams is a mock of IInitParams interface.
type MockIInitParams struct {
	ctrl     *gomock.Controller
	recorder *MockIInitParamsMockRecorder
}

// MockIInitParamsMockRecorder is the mock recorder for MockIInitParams.
type MockIInitParamsMockRecorder struct {
	mock *MockIInitParams
}

// NewMockIInitParams creates a new mock instance.
func NewMockIInitParams(ctrl *gomock.Controller) *MockIInitParams {
	mock := &MockIInitParams{ctrl: ctrl}
	mock.recorder = &MockIInitParamsMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIInitParams) EXPECT() *MockIInitParamsMockRecorder {
	return m.recorder
}

// NextFuncInBoundChannel mocks base method.
func (m *MockIInitParams) NextFuncInBoundChannel() rxgo.NextFunc {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NextFuncInBoundChannel")
	ret0, _ := ret[0].(rxgo.NextFunc)
	return ret0
}

// NextFuncInBoundChannel indicates an expected call of NextFuncInBoundChannel.
func (mr *MockIInitParamsMockRecorder) NextFuncInBoundChannel() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NextFuncInBoundChannel", reflect.TypeOf((*MockIInitParams)(nil).NextFuncInBoundChannel))
}

// NextFuncOutBoundChannel mocks base method.
func (m *MockIInitParams) NextFuncOutBoundChannel() rxgo.NextFunc {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NextFuncOutBoundChannel")
	ret0, _ := ret[0].(rxgo.NextFunc)
	return ret0
}

// NextFuncOutBoundChannel indicates an expected call of NextFuncOutBoundChannel.
func (mr *MockIInitParamsMockRecorder) NextFuncOutBoundChannel() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NextFuncOutBoundChannel", reflect.TypeOf((*MockIInitParams)(nil).NextFuncOutBoundChannel))
}

// OnSendToConnection mocks base method.
func (m *MockIInitParams) OnSendToConnection() rxgo.NextFunc {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "OnSendToConnection")
	ret0, _ := ret[0].(rxgo.NextFunc)
	return ret0
}

// OnSendToConnection indicates an expected call of OnSendToConnection.
func (mr *MockIInitParamsMockRecorder) OnSendToConnection() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OnSendToConnection", reflect.TypeOf((*MockIInitParams)(nil).OnSendToConnection))
}

// OnSendToReactor mocks base method.
func (m *MockIInitParams) OnSendToReactor() rxgo.NextFunc {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "OnSendToReactor")
	ret0, _ := ret[0].(rxgo.NextFunc)
	return ret0
}

// OnSendToReactor indicates an expected call of OnSendToReactor.
func (mr *MockIInitParamsMockRecorder) OnSendToReactor() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OnSendToReactor", reflect.TypeOf((*MockIInitParams)(nil).OnSendToReactor))
}

// argNames: []
// defaultArgs: []
// defaultArgsAsString:
// argTypes: []
// argString:
// rets: [rxgo.NextFunc]
// retString: rxgo.NextFunc
// retString:  rxgo.NextFunc
// ia: map[]
// idRecv: mr
// 0
func (mr *MockIInitParamsMockRecorder) OnNextFuncInBoundChannelDoAndReturn(
	f func() rxgo.NextFunc) *gomock.Call {
	return mr.
		NextFuncInBoundChannel().
		DoAndReturn(f)
}

// 0
func (mr *MockIInitParamsMockRecorder) OnNextFuncInBoundChannelDo(
	f func()) *gomock.Call {
	return mr.
		NextFuncInBoundChannel().
		DoAndReturn(f)
}

// retNames: [ret0]
// retArgs: [ret0 rxgo.NextFunc]
// retArgs22: ret0 rxgo.NextFunc
// 1
func (mr *MockIInitParamsMockRecorder) OnNextFuncInBoundChannelReturn(ret0 rxgo.NextFunc) *gomock.Call {
	return mr.
		NextFuncInBoundChannel().
		Return(ret0)
}

// argNames: []
// defaultArgs: []
// defaultArgsAsString:
// argTypes: []
// argString:
// rets: [rxgo.NextFunc]
// retString: rxgo.NextFunc
// retString:  rxgo.NextFunc
// ia: map[]
// idRecv: mr
// 0
func (mr *MockIInitParamsMockRecorder) OnNextFuncOutBoundChannelDoAndReturn(
	f func() rxgo.NextFunc) *gomock.Call {
	return mr.
		NextFuncOutBoundChannel().
		DoAndReturn(f)
}

// 0
func (mr *MockIInitParamsMockRecorder) OnNextFuncOutBoundChannelDo(
	f func()) *gomock.Call {
	return mr.
		NextFuncOutBoundChannel().
		DoAndReturn(f)
}

// retNames: [ret0]
// retArgs: [ret0 rxgo.NextFunc]
// retArgs22: ret0 rxgo.NextFunc
// 1
func (mr *MockIInitParamsMockRecorder) OnNextFuncOutBoundChannelReturn(ret0 rxgo.NextFunc) *gomock.Call {
	return mr.
		NextFuncOutBoundChannel().
		Return(ret0)
}

// argNames: []
// defaultArgs: []
// defaultArgsAsString:
// argTypes: []
// argString:
// rets: [rxgo.NextFunc]
// retString: rxgo.NextFunc
// retString:  rxgo.NextFunc
// ia: map[]
// idRecv: mr
// 0
func (mr *MockIInitParamsMockRecorder) OnOnSendToConnectionDoAndReturn(
	f func() rxgo.NextFunc) *gomock.Call {
	return mr.
		OnSendToConnection().
		DoAndReturn(f)
}

// 0
func (mr *MockIInitParamsMockRecorder) OnOnSendToConnectionDo(
	f func()) *gomock.Call {
	return mr.
		OnSendToConnection().
		DoAndReturn(f)
}

// retNames: [ret0]
// retArgs: [ret0 rxgo.NextFunc]
// retArgs22: ret0 rxgo.NextFunc
// 1
func (mr *MockIInitParamsMockRecorder) OnOnSendToConnectionReturn(ret0 rxgo.NextFunc) *gomock.Call {
	return mr.
		OnSendToConnection().
		Return(ret0)
}

// argNames: []
// defaultArgs: []
// defaultArgsAsString:
// argTypes: []
// argString:
// rets: [rxgo.NextFunc]
// retString: rxgo.NextFunc
// retString:  rxgo.NextFunc
// ia: map[]
// idRecv: mr
// 0
func (mr *MockIInitParamsMockRecorder) OnOnSendToReactorDoAndReturn(
	f func() rxgo.NextFunc) *gomock.Call {
	return mr.
		OnSendToReactor().
		DoAndReturn(f)
}

// 0
func (mr *MockIInitParamsMockRecorder) OnOnSendToReactorDo(
	f func()) *gomock.Call {
	return mr.
		OnSendToReactor().
		DoAndReturn(f)
}

// retNames: [ret0]
// retArgs: [ret0 rxgo.NextFunc]
// retArgs22: ret0 rxgo.NextFunc
// 1
func (mr *MockIInitParamsMockRecorder) OnOnSendToReactorReturn(ret0 rxgo.NextFunc) *gomock.Call {
	return mr.
		OnSendToReactor().
		Return(ret0)
}