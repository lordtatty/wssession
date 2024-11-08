// Code generated by mockery v2.46.3. DO NOT EDIT.

package wssession_mock

import (
	json "encoding/json"
	time "time"

	mock "github.com/stretchr/testify/mock"
)

// MockWriter is an autogenerated mock type for the Writer type
type MockWriter struct {
	mock.Mock
}

type MockWriter_Expecter struct {
	mock *mock.Mock
}

func (_m *MockWriter) EXPECT() *MockWriter_Expecter {
	return &MockWriter_Expecter{mock: &_m.Mock}
}

// SendJSON provides a mock function with given fields: msgType, j
func (_m *MockWriter) SendJSON(msgType string, j any) error {
	ret := _m.Called(msgType, j)

	if len(ret) == 0 {
		panic("no return value specified for SendJSON")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string, any) error); ok {
		r0 = rf(msgType, j)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockWriter_SendJSON_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SendJSON'
type MockWriter_SendJSON_Call struct {
	*mock.Call
}

// SendJSON is a helper method to define mock.On call
//   - msgType string
//   - j any
func (_e *MockWriter_Expecter) SendJSON(msgType interface{}, j interface{}) *MockWriter_SendJSON_Call {
	return &MockWriter_SendJSON_Call{Call: _e.mock.On("SendJSON", msgType, j)}
}

func (_c *MockWriter_SendJSON_Call) Run(run func(msgType string, j any)) *MockWriter_SendJSON_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(any))
	})
	return _c
}

func (_c *MockWriter_SendJSON_Call) Return(_a0 error) *MockWriter_SendJSON_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockWriter_SendJSON_Call) RunAndReturn(run func(string, any) error) *MockWriter_SendJSON_Call {
	_c.Call.Return(run)
	return _c
}

// SendJSONAndWait provides a mock function with given fields: msgType, j, timeout
func (_m *MockWriter) SendJSONAndWait(msgType string, j any, timeout time.Duration) (*json.RawMessage, error) {
	ret := _m.Called(msgType, j, timeout)

	if len(ret) == 0 {
		panic("no return value specified for SendJSONAndWait")
	}

	var r0 *json.RawMessage
	var r1 error
	if rf, ok := ret.Get(0).(func(string, any, time.Duration) (*json.RawMessage, error)); ok {
		return rf(msgType, j, timeout)
	}
	if rf, ok := ret.Get(0).(func(string, any, time.Duration) *json.RawMessage); ok {
		r0 = rf(msgType, j, timeout)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*json.RawMessage)
		}
	}

	if rf, ok := ret.Get(1).(func(string, any, time.Duration) error); ok {
		r1 = rf(msgType, j, timeout)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockWriter_SendJSONAndWait_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SendJSONAndWait'
type MockWriter_SendJSONAndWait_Call struct {
	*mock.Call
}

// SendJSONAndWait is a helper method to define mock.On call
//   - msgType string
//   - j any
//   - timeout time.Duration
func (_e *MockWriter_Expecter) SendJSONAndWait(msgType interface{}, j interface{}, timeout interface{}) *MockWriter_SendJSONAndWait_Call {
	return &MockWriter_SendJSONAndWait_Call{Call: _e.mock.On("SendJSONAndWait", msgType, j, timeout)}
}

func (_c *MockWriter_SendJSONAndWait_Call) Run(run func(msgType string, j any, timeout time.Duration)) *MockWriter_SendJSONAndWait_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(any), args[2].(time.Duration))
	})
	return _c
}

func (_c *MockWriter_SendJSONAndWait_Call) Return(_a0 *json.RawMessage, _a1 error) *MockWriter_SendJSONAndWait_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockWriter_SendJSONAndWait_Call) RunAndReturn(run func(string, any, time.Duration) (*json.RawMessage, error)) *MockWriter_SendJSONAndWait_Call {
	_c.Call.Return(run)
	return _c
}

// SendStr provides a mock function with given fields: msgType, msg
func (_m *MockWriter) SendStr(msgType string, msg string) error {
	ret := _m.Called(msgType, msg)

	if len(ret) == 0 {
		panic("no return value specified for SendStr")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string) error); ok {
		r0 = rf(msgType, msg)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockWriter_SendStr_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SendStr'
type MockWriter_SendStr_Call struct {
	*mock.Call
}

// SendStr is a helper method to define mock.On call
//   - msgType string
//   - msg string
func (_e *MockWriter_Expecter) SendStr(msgType interface{}, msg interface{}) *MockWriter_SendStr_Call {
	return &MockWriter_SendStr_Call{Call: _e.mock.On("SendStr", msgType, msg)}
}

func (_c *MockWriter_SendStr_Call) Run(run func(msgType string, msg string)) *MockWriter_SendStr_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(string))
	})
	return _c
}

func (_c *MockWriter_SendStr_Call) Return(_a0 error) *MockWriter_SendStr_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockWriter_SendStr_Call) RunAndReturn(run func(string, string) error) *MockWriter_SendStr_Call {
	_c.Call.Return(run)
	return _c
}

// SendStrAndWait provides a mock function with given fields: msgType, msg, timeout
func (_m *MockWriter) SendStrAndWait(msgType string, msg string, timeout time.Duration) (*json.RawMessage, error) {
	ret := _m.Called(msgType, msg, timeout)

	if len(ret) == 0 {
		panic("no return value specified for SendStrAndWait")
	}

	var r0 *json.RawMessage
	var r1 error
	if rf, ok := ret.Get(0).(func(string, string, time.Duration) (*json.RawMessage, error)); ok {
		return rf(msgType, msg, timeout)
	}
	if rf, ok := ret.Get(0).(func(string, string, time.Duration) *json.RawMessage); ok {
		r0 = rf(msgType, msg, timeout)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*json.RawMessage)
		}
	}

	if rf, ok := ret.Get(1).(func(string, string, time.Duration) error); ok {
		r1 = rf(msgType, msg, timeout)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockWriter_SendStrAndWait_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SendStrAndWait'
type MockWriter_SendStrAndWait_Call struct {
	*mock.Call
}

// SendStrAndWait is a helper method to define mock.On call
//   - msgType string
//   - msg string
//   - timeout time.Duration
func (_e *MockWriter_Expecter) SendStrAndWait(msgType interface{}, msg interface{}, timeout interface{}) *MockWriter_SendStrAndWait_Call {
	return &MockWriter_SendStrAndWait_Call{Call: _e.mock.On("SendStrAndWait", msgType, msg, timeout)}
}

func (_c *MockWriter_SendStrAndWait_Call) Run(run func(msgType string, msg string, timeout time.Duration)) *MockWriter_SendStrAndWait_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(string), args[2].(time.Duration))
	})
	return _c
}

func (_c *MockWriter_SendStrAndWait_Call) Return(_a0 *json.RawMessage, _a1 error) *MockWriter_SendStrAndWait_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockWriter_SendStrAndWait_Call) RunAndReturn(run func(string, string, time.Duration) (*json.RawMessage, error)) *MockWriter_SendStrAndWait_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockWriter creates a new instance of MockWriter. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockWriter(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockWriter {
	mock := &MockWriter{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
