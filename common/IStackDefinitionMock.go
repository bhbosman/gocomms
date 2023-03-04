// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/bhbosman/gocomms/common (interfaces: IStackDefinition)

// Package common is a generated GoMock package.
package common

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockIStackDefinition is a mock of IStackDefinition interface.
type MockIStackDefinition struct {
	ctrl     *gomock.Controller
	recorder *MockIStackDefinitionMockRecorder
}

// MockIStackDefinitionMockRecorder is the mock recorder for MockIStackDefinition.
type MockIStackDefinitionMockRecorder struct {
	mock *MockIStackDefinition
}

// NewMockIStackDefinition creates a new mock instance.
func NewMockIStackDefinition(ctrl *gomock.Controller) *MockIStackDefinition {
	mock := &MockIStackDefinition{ctrl: ctrl}
	mock.recorder = &MockIStackDefinitionMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIStackDefinition) EXPECT() *MockIStackDefinitionMockRecorder {
	return m.recorder
}

// Inbound mocks base method.
func (m *MockIStackDefinition) Inbound() BoundResult {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Inbound")
	ret0, _ := ret[0].(BoundResult)
	return ret0
}

// Inbound indicates an expected call of Inbound.
func (mr *MockIStackDefinitionMockRecorder) Inbound() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Inbound", reflect.TypeOf((*MockIStackDefinition)(nil).Inbound))
}

// Name mocks base method.
func (m *MockIStackDefinition) Name() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Name")
	ret0, _ := ret[0].(string)
	return ret0
}

// Name indicates an expected call of Name.
func (mr *MockIStackDefinitionMockRecorder) Name() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Name", reflect.TypeOf((*MockIStackDefinition)(nil).Name))
}

// Outbound mocks base method.
func (m *MockIStackDefinition) Outbound() BoundResult {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Outbound")
	ret0, _ := ret[0].(BoundResult)
	return ret0
}

// Outbound indicates an expected call of Outbound.
func (mr *MockIStackDefinitionMockRecorder) Outbound() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Outbound", reflect.TypeOf((*MockIStackDefinition)(nil).Outbound))
}

// StackState mocks base method.
func (m *MockIStackDefinition) StackState() *StackState {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StackState")
	ret0, _ := ret[0].(*StackState)
	return ret0
}

// StackState indicates an expected call of StackState.
func (mr *MockIStackDefinitionMockRecorder) StackState() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StackState", reflect.TypeOf((*MockIStackDefinition)(nil).StackState))
}

// argNames: []
// defaultArgs: []
// defaultArgsAsString:
// argTypes: []
// argString:
// rets: [BoundResult]
// retString: BoundResult
// retString:  BoundResult
// ia: map[]
// idRecv: mr
// 0
func (mr *MockIStackDefinitionMockRecorder) OnInboundDoAndReturn(
	f func() BoundResult) *gomock.Call {
	return mr.
		Inbound().
		DoAndReturn(f)
}

// 0
func (mr *MockIStackDefinitionMockRecorder) OnInboundDo(
	f func()) *gomock.Call {
	return mr.
		Inbound().
		DoAndReturn(f)
}

// retNames: [ret0]
// retArgs: [ret0 BoundResult]
// retArgs22: ret0 BoundResult
// 1
func (mr *MockIStackDefinitionMockRecorder) OnInboundReturn(ret0 BoundResult) *gomock.Call {
	return mr.
		Inbound().
		Return(ret0)
}

// argNames: []
// defaultArgs: []
// defaultArgsAsString:
// argTypes: []
// argString:
// rets: [string]
// retString: string
// retString:  string
// ia: map[]
// idRecv: mr
// 0
func (mr *MockIStackDefinitionMockRecorder) OnNameDoAndReturn(
	f func() string) *gomock.Call {
	return mr.
		Name().
		DoAndReturn(f)
}

// 0
func (mr *MockIStackDefinitionMockRecorder) OnNameDo(
	f func()) *gomock.Call {
	return mr.
		Name().
		DoAndReturn(f)
}

// retNames: [ret0]
// retArgs: [ret0 string]
// retArgs22: ret0 string
// 1
func (mr *MockIStackDefinitionMockRecorder) OnNameReturn(ret0 string) *gomock.Call {
	return mr.
		Name().
		Return(ret0)
}

// argNames: []
// defaultArgs: []
// defaultArgsAsString:
// argTypes: []
// argString:
// rets: [BoundResult]
// retString: BoundResult
// retString:  BoundResult
// ia: map[]
// idRecv: mr
// 0
func (mr *MockIStackDefinitionMockRecorder) OnOutboundDoAndReturn(
	f func() BoundResult) *gomock.Call {
	return mr.
		Outbound().
		DoAndReturn(f)
}

// 0
func (mr *MockIStackDefinitionMockRecorder) OnOutboundDo(
	f func()) *gomock.Call {
	return mr.
		Outbound().
		DoAndReturn(f)
}

// retNames: [ret0]
// retArgs: [ret0 BoundResult]
// retArgs22: ret0 BoundResult
// 1
func (mr *MockIStackDefinitionMockRecorder) OnOutboundReturn(ret0 BoundResult) *gomock.Call {
	return mr.
		Outbound().
		Return(ret0)
}

// argNames: []
// defaultArgs: []
// defaultArgsAsString:
// argTypes: []
// argString:
// rets: [*StackState]
// retString: *StackState
// retString:  *StackState
// ia: map[]
// idRecv: mr
// 0
func (mr *MockIStackDefinitionMockRecorder) OnStackStateDoAndReturn(
	f func() *StackState) *gomock.Call {
	return mr.
		StackState().
		DoAndReturn(f)
}

// 0
func (mr *MockIStackDefinitionMockRecorder) OnStackStateDo(
	f func()) *gomock.Call {
	return mr.
		StackState().
		DoAndReturn(f)
}

// retNames: [ret0]
// retArgs: [ret0 *StackState]
// retArgs22: ret0 *StackState
// 1
func (mr *MockIStackDefinitionMockRecorder) OnStackStateReturn(ret0 *StackState) *gomock.Call {
	return mr.
		StackState().
		Return(ret0)
}
