// Code generated by mockery v2.46.3. DO NOT EDIT.

package wssession_mock

import (
	wssession "github.com/lordtatty/wssession"
	mock "github.com/stretchr/testify/mock"
)

// MockSessionGetter is an autogenerated mock type for the SessionGetter type
type MockSessionGetter struct {
	mock.Mock
}

type MockSessionGetter_Expecter struct {
	mock *mock.Mock
}

func (_m *MockSessionGetter) EXPECT() *MockSessionGetter_Expecter {
	return &MockSessionGetter_Expecter{mock: &_m.Mock}
}

// Get provides a mock function with given fields: connID, conn
func (_m *MockSessionGetter) Get(connID string, conn wssession.WebsocketConn) (*wssession.Session, error) {
	ret := _m.Called(connID, conn)

	if len(ret) == 0 {
		panic("no return value specified for Get")
	}

	var r0 *wssession.Session
	var r1 error
	if rf, ok := ret.Get(0).(func(string, wssession.WebsocketConn) (*wssession.Session, error)); ok {
		return rf(connID, conn)
	}
	if rf, ok := ret.Get(0).(func(string, wssession.WebsocketConn) *wssession.Session); ok {
		r0 = rf(connID, conn)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*wssession.Session)
		}
	}

	if rf, ok := ret.Get(1).(func(string, wssession.WebsocketConn) error); ok {
		r1 = rf(connID, conn)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockSessionGetter_Get_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Get'
type MockSessionGetter_Get_Call struct {
	*mock.Call
}

// Get is a helper method to define mock.On call
//   - connID string
//   - conn wssession.WebsocketConn
func (_e *MockSessionGetter_Expecter) Get(connID interface{}, conn interface{}) *MockSessionGetter_Get_Call {
	return &MockSessionGetter_Get_Call{Call: _e.mock.On("Get", connID, conn)}
}

func (_c *MockSessionGetter_Get_Call) Run(run func(connID string, conn wssession.WebsocketConn)) *MockSessionGetter_Get_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(wssession.WebsocketConn))
	})
	return _c
}

func (_c *MockSessionGetter_Get_Call) Return(_a0 *wssession.Session, _a1 error) *MockSessionGetter_Get_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockSessionGetter_Get_Call) RunAndReturn(run func(string, wssession.WebsocketConn) (*wssession.Session, error)) *MockSessionGetter_Get_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockSessionGetter creates a new instance of MockSessionGetter. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockSessionGetter(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockSessionGetter {
	mock := &MockSessionGetter{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}