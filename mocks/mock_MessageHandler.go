// Code generated by mockery v2.46.3. DO NOT EDIT.

package wssession_mock

import (
	json "encoding/json"

	wssession "github.com/lordtatty/wssession"
	mock "github.com/stretchr/testify/mock"
)

// MockMessageHandler is an autogenerated mock type for the MessageHandler type
type MockMessageHandler[T any] struct {
	mock.Mock
}

type MockMessageHandler_Expecter[T any] struct {
	mock *mock.Mock
}

func (_m *MockMessageHandler[T]) EXPECT() *MockMessageHandler_Expecter[T] {
	return &MockMessageHandler_Expecter[T]{mock: &_m.Mock}
}

// WSHandle provides a mock function with given fields: w, msg, customState
func (_m *MockMessageHandler[T]) WSHandle(w wssession.Writer, msg json.RawMessage, customState *T) error {
	ret := _m.Called(w, msg, customState)

	if len(ret) == 0 {
		panic("no return value specified for WSHandle")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(wssession.Writer, json.RawMessage, *T) error); ok {
		r0 = rf(w, msg, customState)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockMessageHandler_WSHandle_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'WSHandle'
type MockMessageHandler_WSHandle_Call[T any] struct {
	*mock.Call
}

// WSHandle is a helper method to define mock.On call
//   - w wssession.Writer
//   - msg json.RawMessage
//   - customState *T
func (_e *MockMessageHandler_Expecter[T]) WSHandle(w interface{}, msg interface{}, customState interface{}) *MockMessageHandler_WSHandle_Call[T] {
	return &MockMessageHandler_WSHandle_Call[T]{Call: _e.mock.On("WSHandle", w, msg, customState)}
}

func (_c *MockMessageHandler_WSHandle_Call[T]) Run(run func(w wssession.Writer, msg json.RawMessage, customState *T)) *MockMessageHandler_WSHandle_Call[T] {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(wssession.Writer), args[1].(json.RawMessage), args[2].(*T))
	})
	return _c
}

func (_c *MockMessageHandler_WSHandle_Call[T]) Return(_a0 error) *MockMessageHandler_WSHandle_Call[T] {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockMessageHandler_WSHandle_Call[T]) RunAndReturn(run func(wssession.Writer, json.RawMessage, *T) error) *MockMessageHandler_WSHandle_Call[T] {
	_c.Call.Return(run)
	return _c
}

// NewMockMessageHandler creates a new instance of MockMessageHandler. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockMessageHandler[T any](t interface {
	mock.TestingT
	Cleanup(func())
}) *MockMessageHandler[T] {
	mock := &MockMessageHandler[T]{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
