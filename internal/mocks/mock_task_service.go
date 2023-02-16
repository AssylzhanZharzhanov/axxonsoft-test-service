// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/AssylzhanZharzhanov/axxonsoft-test-service/internal/domain (interfaces: TaskService)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	domain "github.com/AssylzhanZharzhanov/axxonsoft-test-service/internal/domain"
	gomock "github.com/golang/mock/gomock"
)

// MockTaskService is a mock of TaskService interface.
type MockTaskService struct {
	ctrl     *gomock.Controller
	recorder *MockTaskServiceMockRecorder
}

// MockTaskServiceMockRecorder is the mock recorder for MockTaskService.
type MockTaskServiceMockRecorder struct {
	mock *MockTaskService
}

// NewMockTaskService creates a new mock instance.
func NewMockTaskService(ctrl *gomock.Controller) *MockTaskService {
	mock := &MockTaskService{ctrl: ctrl}
	mock.recorder = &MockTaskServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTaskService) EXPECT() *MockTaskServiceMockRecorder {
	return m.recorder
}

// CreateTask mocks base method.
func (m *MockTaskService) CreateTask(arg0 context.Context, arg1 *domain.Task) (domain.TaskID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateTask", arg0, arg1)
	ret0, _ := ret[0].(domain.TaskID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateTask indicates an expected call of CreateTask.
func (mr *MockTaskServiceMockRecorder) CreateTask(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateTask", reflect.TypeOf((*MockTaskService)(nil).CreateTask), arg0, arg1)
}

// GetTask mocks base method.
func (m *MockTaskService) GetTask(arg0 context.Context, arg1 domain.TaskID) (*domain.Task, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTask", arg0, arg1)
	ret0, _ := ret[0].(*domain.Task)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTask indicates an expected call of GetTask.
func (mr *MockTaskServiceMockRecorder) GetTask(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTask", reflect.TypeOf((*MockTaskService)(nil).GetTask), arg0, arg1)
}

// ListTasks mocks base method.
func (m *MockTaskService) ListTasks(arg0 context.Context, arg1 domain.TaskSearchCriteria) ([]*domain.Task, domain.Total, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListTasks", arg0, arg1)
	ret0, _ := ret[0].([]*domain.Task)
	ret1, _ := ret[1].(domain.Total)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// ListTasks indicates an expected call of ListTasks.
func (mr *MockTaskServiceMockRecorder) ListTasks(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListTasks", reflect.TypeOf((*MockTaskService)(nil).ListTasks), arg0, arg1)
}

// UpdateTask mocks base method.
func (m *MockTaskService) UpdateTask(arg0 context.Context, arg1 *domain.Task) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateTask", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateTask indicates an expected call of UpdateTask.
func (mr *MockTaskServiceMockRecorder) UpdateTask(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateTask", reflect.TypeOf((*MockTaskService)(nil).UpdateTask), arg0, arg1)
}
