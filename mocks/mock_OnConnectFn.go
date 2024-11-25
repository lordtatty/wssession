// Code generated by mockery v2.46.3. DO NOT EDIT.

package wssession_mock

import (
	wssession "github.com/lordtatty/wssession"
	mock "github.com/stretchr/testify/mock"
)

// MockOnConnectFn is an autogenerated mock type for the OnConnectFn type
type MockOnConnectFn struct {
	mock.Mock
}

type MockOnConnectFn_Expecter struct {
	mock *mock.Mock
}

func (_m *MockOnConnectFn) EXPECT() *MockOnConnectFn_Expecter {
	return &MockOnConnectFn_Expecter{mock: &_m.Mock}
}

// Execute provides a mock function with given fields: m
func (_m *MockOnConnectFn) Execute(m wssession.ReceivedMsg) error {
	ret := _m.Called(m)

	if len(ret) == 0 {
		panic("no return value specified for Execute")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(wssession.ReceivedMsg) error); ok {
		r0 = rf(m)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockOnConnectFn_Execute_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Execute'
type MockOnConnectFn_Execute_Call struct {
	*mock.Call
}

// Execute is a helper method to define mock.On call
//   - m wssession.ReceivedMsg
func (_e *MockOnConnectFn_Expecter) Execute(m interface{}) *MockOnConnectFn_Execute_Call {
	return &MockOnConnectFn_Execute_Call{Call: _e.mock.On("Execute", m)}
}

func (_c *MockOnConnectFn_Execute_Call) Run(run func(m wssession.ReceivedMsg)) *MockOnConnectFn_Execute_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(wssession.ReceivedMsg))
	})
	return _c
}

func (_c *MockOnConnectFn_Execute_Call) Return(_a0 error) *MockOnConnectFn_Execute_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockOnConnectFn_Execute_Call) RunAndReturn(run func(wssession.ReceivedMsg) error) *MockOnConnectFn_Execute_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockOnConnectFn creates a new instance of MockOnConnectFn. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockOnConnectFn(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockOnConnectFn {
	mock := &MockOnConnectFn{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
