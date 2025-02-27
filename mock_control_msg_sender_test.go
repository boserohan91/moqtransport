// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/mengelbart/moqtransport (interfaces: ControlMsgSender)
//
// Generated by this command:
//
//	mockgen -build_flags=-tags=gomock -package moqtransport -self_package github.com/mengelbart/moqtransport -destination mock_control_msg_sender_test.go github.com/mengelbart/moqtransport ControlMsgSender
//
// Package moqtransport is a generated GoMock package.
package moqtransport

import (
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockControlMsgSender is a mock of ControlMsgSender interface.
type MockControlMsgSender struct {
	ctrl     *gomock.Controller
	recorder *MockControlMsgSenderMockRecorder
}

// MockControlMsgSenderMockRecorder is the mock recorder for MockControlMsgSender.
type MockControlMsgSenderMockRecorder struct {
	mock *MockControlMsgSender
}

// NewMockControlMsgSender creates a new mock instance.
func NewMockControlMsgSender(ctrl *gomock.Controller) *MockControlMsgSender {
	mock := &MockControlMsgSender{ctrl: ctrl}
	mock.recorder = &MockControlMsgSenderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockControlMsgSender) EXPECT() *MockControlMsgSenderMockRecorder {
	return m.recorder
}

// send mocks base method.
func (m *MockControlMsgSender) send(arg0 message) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "send", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// send indicates an expected call of send.
func (mr *MockControlMsgSenderMockRecorder) send(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "send", reflect.TypeOf((*MockControlMsgSender)(nil).send), arg0)
}
