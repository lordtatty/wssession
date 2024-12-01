// Code generated by mockery v2.46.3. DO NOT EDIT.

package wssession_mock

import (
	mock "github.com/stretchr/testify/mock"
	context "golang.org/x/net/context"

	redis "github.com/redis/go-redis/v9"

	time "time"
)

// MockRedisClient is an autogenerated mock type for the RedisClient type
type MockRedisClient struct {
	mock.Mock
}

type MockRedisClient_Expecter struct {
	mock *mock.Mock
}

func (_m *MockRedisClient) EXPECT() *MockRedisClient_Expecter {
	return &MockRedisClient_Expecter{mock: &_m.Mock}
}

// Get provides a mock function with given fields: ctx, key
func (_m *MockRedisClient) Get(ctx context.Context, key string) *redis.StringCmd {
	ret := _m.Called(ctx, key)

	if len(ret) == 0 {
		panic("no return value specified for Get")
	}

	var r0 *redis.StringCmd
	if rf, ok := ret.Get(0).(func(context.Context, string) *redis.StringCmd); ok {
		r0 = rf(ctx, key)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*redis.StringCmd)
		}
	}

	return r0
}

// MockRedisClient_Get_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Get'
type MockRedisClient_Get_Call struct {
	*mock.Call
}

// Get is a helper method to define mock.On call
//   - ctx context.Context
//   - key string
func (_e *MockRedisClient_Expecter) Get(ctx interface{}, key interface{}) *MockRedisClient_Get_Call {
	return &MockRedisClient_Get_Call{Call: _e.mock.On("Get", ctx, key)}
}

func (_c *MockRedisClient_Get_Call) Run(run func(ctx context.Context, key string)) *MockRedisClient_Get_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *MockRedisClient_Get_Call) Return(_a0 *redis.StringCmd) *MockRedisClient_Get_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockRedisClient_Get_Call) RunAndReturn(run func(context.Context, string) *redis.StringCmd) *MockRedisClient_Get_Call {
	_c.Call.Return(run)
	return _c
}

// Keys provides a mock function with given fields: ctx, pattern
func (_m *MockRedisClient) Keys(ctx context.Context, pattern string) *redis.StringSliceCmd {
	ret := _m.Called(ctx, pattern)

	if len(ret) == 0 {
		panic("no return value specified for Keys")
	}

	var r0 *redis.StringSliceCmd
	if rf, ok := ret.Get(0).(func(context.Context, string) *redis.StringSliceCmd); ok {
		r0 = rf(ctx, pattern)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*redis.StringSliceCmd)
		}
	}

	return r0
}

// MockRedisClient_Keys_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Keys'
type MockRedisClient_Keys_Call struct {
	*mock.Call
}

// Keys is a helper method to define mock.On call
//   - ctx context.Context
//   - pattern string
func (_e *MockRedisClient_Expecter) Keys(ctx interface{}, pattern interface{}) *MockRedisClient_Keys_Call {
	return &MockRedisClient_Keys_Call{Call: _e.mock.On("Keys", ctx, pattern)}
}

func (_c *MockRedisClient_Keys_Call) Run(run func(ctx context.Context, pattern string)) *MockRedisClient_Keys_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *MockRedisClient_Keys_Call) Return(_a0 *redis.StringSliceCmd) *MockRedisClient_Keys_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockRedisClient_Keys_Call) RunAndReturn(run func(context.Context, string) *redis.StringSliceCmd) *MockRedisClient_Keys_Call {
	_c.Call.Return(run)
	return _c
}

// Set provides a mock function with given fields: ctx, key, value, expiration
func (_m *MockRedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	ret := _m.Called(ctx, key, value, expiration)

	if len(ret) == 0 {
		panic("no return value specified for Set")
	}

	var r0 *redis.StatusCmd
	if rf, ok := ret.Get(0).(func(context.Context, string, interface{}, time.Duration) *redis.StatusCmd); ok {
		r0 = rf(ctx, key, value, expiration)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*redis.StatusCmd)
		}
	}

	return r0
}

// MockRedisClient_Set_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Set'
type MockRedisClient_Set_Call struct {
	*mock.Call
}

// Set is a helper method to define mock.On call
//   - ctx context.Context
//   - key string
//   - value interface{}
//   - expiration time.Duration
func (_e *MockRedisClient_Expecter) Set(ctx interface{}, key interface{}, value interface{}, expiration interface{}) *MockRedisClient_Set_Call {
	return &MockRedisClient_Set_Call{Call: _e.mock.On("Set", ctx, key, value, expiration)}
}

func (_c *MockRedisClient_Set_Call) Run(run func(ctx context.Context, key string, value interface{}, expiration time.Duration)) *MockRedisClient_Set_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(interface{}), args[3].(time.Duration))
	})
	return _c
}

func (_c *MockRedisClient_Set_Call) Return(_a0 *redis.StatusCmd) *MockRedisClient_Set_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockRedisClient_Set_Call) RunAndReturn(run func(context.Context, string, interface{}, time.Duration) *redis.StatusCmd) *MockRedisClient_Set_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockRedisClient creates a new instance of MockRedisClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockRedisClient(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockRedisClient {
	mock := &MockRedisClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}