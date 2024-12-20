// Code generated by mockery v2.46.3. DO NOT EDIT.

package wssession_mock

import (
	wssession "github.com/lordtatty/wssession"
	mock "github.com/stretchr/testify/mock"
)

// MockCache is an autogenerated mock type for the Cache type
type MockCache struct {
	mock.Mock
}

type MockCache_Expecter struct {
	mock *mock.Mock
}

func (_m *MockCache) EXPECT() *MockCache_Expecter {
	return &MockCache_Expecter{mock: &_m.Mock}
}

// Add provides a mock function with given fields: connID, r
func (_m *MockCache) Add(connID string, r wssession.ResponseMsg) error {
	ret := _m.Called(connID, r)

	if len(ret) == 0 {
		panic("no return value specified for Add")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string, wssession.ResponseMsg) error); ok {
		r0 = rf(connID, r)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockCache_Add_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Add'
type MockCache_Add_Call struct {
	*mock.Call
}

// Add is a helper method to define mock.On call
//   - connID string
//   - r wssession.ResponseMsg
func (_e *MockCache_Expecter) Add(connID interface{}, r interface{}) *MockCache_Add_Call {
	return &MockCache_Add_Call{Call: _e.mock.On("Add", connID, r)}
}

func (_c *MockCache_Add_Call) Run(run func(connID string, r wssession.ResponseMsg)) *MockCache_Add_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(wssession.ResponseMsg))
	})
	return _c
}

func (_c *MockCache_Add_Call) Return(_a0 error) *MockCache_Add_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockCache_Add_Call) RunAndReturn(run func(string, wssession.ResponseMsg) error) *MockCache_Add_Call {
	_c.Call.Return(run)
	return _c
}

// Items provides a mock function with given fields: connID
func (_m *MockCache) Items(connID string) ([]*wssession.ResponseMsg, error) {
	ret := _m.Called(connID)

	if len(ret) == 0 {
		panic("no return value specified for Items")
	}

	var r0 []*wssession.ResponseMsg
	var r1 error
	if rf, ok := ret.Get(0).(func(string) ([]*wssession.ResponseMsg, error)); ok {
		return rf(connID)
	}
	if rf, ok := ret.Get(0).(func(string) []*wssession.ResponseMsg); ok {
		r0 = rf(connID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*wssession.ResponseMsg)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(connID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockCache_Items_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Items'
type MockCache_Items_Call struct {
	*mock.Call
}

// Items is a helper method to define mock.On call
//   - connID string
func (_e *MockCache_Expecter) Items(connID interface{}) *MockCache_Items_Call {
	return &MockCache_Items_Call{Call: _e.mock.On("Items", connID)}
}

func (_c *MockCache_Items_Call) Run(run func(connID string)) *MockCache_Items_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockCache_Items_Call) Return(_a0 []*wssession.ResponseMsg, _a1 error) *MockCache_Items_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockCache_Items_Call) RunAndReturn(run func(string) ([]*wssession.ResponseMsg, error)) *MockCache_Items_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockCache creates a new instance of MockCache. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockCache(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockCache {
	mock := &MockCache{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
