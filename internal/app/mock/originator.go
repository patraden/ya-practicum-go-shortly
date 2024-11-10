// Code generated by MockGen. DO NOT EDIT.
// Source: internal/app/memento/originator.go
//
// Generated by this command:
//
//	mockgen -source=internal/app/memento/originator.go -destination=internal/app/mock/originator.go -package=mock Originator
//

// Package mock is a generated GoMock package.
package mock

import (
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"

	memento "github.com/patraden/ya-practicum-go-shortly/internal/app/memento"
)

// MockOriginator is a mock of Originator interface.
type MockOriginator struct {
	ctrl     *gomock.Controller
	recorder *MockOriginatorMockRecorder
	isgomock struct{}
}

// MockOriginatorMockRecorder is the mock recorder for MockOriginator.
type MockOriginatorMockRecorder struct {
	mock *MockOriginator
}

// NewMockOriginator creates a new mock instance.
func NewMockOriginator(ctrl *gomock.Controller) *MockOriginator {
	mock := &MockOriginator{ctrl: ctrl}
	mock.recorder = &MockOriginatorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockOriginator) EXPECT() *MockOriginatorMockRecorder {
	return m.recorder
}

// CreateMemento mocks base method.
func (m *MockOriginator) CreateMemento() (*memento.Memento, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateMemento")
	ret0, _ := ret[0].(*memento.Memento)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateMemento indicates an expected call of CreateMemento.
func (mr *MockOriginatorMockRecorder) CreateMemento() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateMemento", reflect.TypeOf((*MockOriginator)(nil).CreateMemento))
}

// RestoreMemento mocks base method.
func (m_2 *MockOriginator) RestoreMemento(m *memento.Memento) error {
	m_2.ctrl.T.Helper()
	ret := m_2.ctrl.Call(m_2, "RestoreMemento", m)
	ret0, _ := ret[0].(error)
	return ret0
}

// RestoreMemento indicates an expected call of RestoreMemento.
func (mr *MockOriginatorMockRecorder) RestoreMemento(m any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RestoreMemento", reflect.TypeOf((*MockOriginator)(nil).RestoreMemento), m)
}