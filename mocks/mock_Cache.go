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
func (_m *MockCache) Add(connID string, r wssession.ResponseMsg) {
	_m.Called(connID, r)
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

func (_c *MockCache_Add_Call) Return() *MockCache_Add_Call {
	_c.Call.Return()
	return _c
}

func (_c *MockCache_Add_Call) RunAndReturn(run func(string, wssession.ResponseMsg)) *MockCache_Add_Call {
	_c.Call.Return(run)
	return _c
}

// Items provides a mock function with given fields: connID
func (_m *MockCache) Items(connID string) []*wssession.ResponseMsg {
	ret := _m.Called(connID)

	if len(ret) == 0 {
		panic("no return value specified for Items")
	}

	var r0 []*wssession.ResponseMsg
	if rf, ok := ret.Get(0).(func(string) []*wssession.ResponseMsg); ok {
		r0 = rf(connID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*wssession.ResponseMsg)
		}
	}

	return r0
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

func (_c *MockCache_Items_Call) Return(_a0 []*wssession.ResponseMsg) *MockCache_Items_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockCache_Items_Call) RunAndReturn(run func(string) []*wssession.ResponseMsg) *MockCache_Items_Call {
	_c.Call.Return(run)
	return _c
}

// Len provides a mock function with given fields: connID
func (_m *MockCache) Len(connID string) int {
	ret := _m.Called(connID)

	if len(ret) == 0 {
		panic("no return value specified for Len")
	}

	var r0 int
	if rf, ok := ret.Get(0).(func(string) int); ok {
		r0 = rf(connID)
	} else {
		r0 = ret.Get(0).(int)
	}

	return r0
}

// MockCache_Len_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Len'
type MockCache_Len_Call struct {
	*mock.Call
}

// Len is a helper method to define mock.On call
//   - connID string
func (_e *MockCache_Expecter) Len(connID interface{}) *MockCache_Len_Call {
	return &MockCache_Len_Call{Call: _e.mock.On("Len", connID)}
}

func (_c *MockCache_Len_Call) Run(run func(connID string)) *MockCache_Len_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockCache_Len_Call) Return(_a0 int) *MockCache_Len_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockCache_Len_Call) RunAndReturn(run func(string) int) *MockCache_Len_Call {
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
