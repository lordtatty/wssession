package wssession_test

import (
	"fmt"
	"sort"
	"testing"
	"time"

	"github.com/lordtatty/wssession"
	mocks "github.com/lordtatty/wssession/mocks"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/net/context"
)

func TestRedisCache_ImplementsCache(t *testing.T) {
	// The goal of this test is to make sure that RedisCache can be passed into ServeSession
	// (which expects a Cache interface)
	t.Parallel()
	assert := assert.New(t)
	mConn := mocks.MockWebsocketConn{}
	defer mConn.AssertExpectations(t)
	mConn.EXPECT().ReadMessage().Return(0, nil, nil) // not sending a connection message first will immediately end and return the session

	sut := &wssession.RedisCache{}
	s := wssession.Mgr{}
	s.ServeSession(&mConn, sut)
	assert.Implements((*wssession.Cache)(nil), sut)
	assert.True(true, "this will always be true, if sut does not implement Cache, the code won't compile")
}

func TestRedisCache_Add(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	fixedTime := time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name        string
		msg         wssession.ResponseMsg
		connID      string
		redisKey    string
		redisVal    string
		cacheTTL    time.Duration
		expectError bool
		errorMsg    string // expected error message
		redisError  error
		skipRedis   bool // skip redis mock setup if we expect early error
	}{
		{
			name: "successful storage with default TTL",
			msg: wssession.ResponseMsg{
				ID:      "1",
				ConnID:  "conn-1",
				Type:    "test",
				Message: "hello",
			},
			connID:   "0001",
			redisKey: "cache:0001:1609459200000000000",
			redisVal: `{"id":"1","conn_id":"conn-1","type":"test","message":"hello"}`,
			cacheTTL: 0, // test default TTL
		},
		{
			name: "successful storage with custom TTL",
			msg: wssession.ResponseMsg{
				ID:      "2",
				ConnID:  "conn-2",
				Type:    "test",
				Message: "world",
			},
			connID:   "0002",
			redisKey: "cache:0002:1609459200000000000",
			redisVal: `{"id":"2","conn_id":"conn-2","type":"test","message":"world"}`,
			cacheTTL: time.Hour,
		},
		{
			name: "redis set operation fails",
			msg: wssession.ResponseMsg{
				ID:      "3",
				ConnID:  "conn-3",
				Type:    "test",
				Message: "error test",
			},
			connID:      "0003",
			redisKey:    "cache:0003:1609459200000000000",
			redisVal:    `{"id":"3","conn_id":"conn-3","type":"test","message":"error test"}`,
			redisError:  fmt.Errorf("redis connection failed"),
			expectError: true,
			errorMsg:    "failed to set message",
		},
		{
			name: "empty connection ID",
			msg: wssession.ResponseMsg{
				ID:      "4",
				ConnID:  "conn-4",
				Type:    "test",
				Message: "test",
			},
			connID:      "",
			expectError: true,
			errorMsg:    "connection ID cannot be empty",
			skipRedis:   true, // Redis mock won't be called due to early error
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mRClient := &mocks.MockRedisClient{}
			defer mRClient.AssertExpectations(t)

			if !tt.skipRedis {
				expectedTTL := time.Minute
				if tt.cacheTTL != 0 {
					expectedTTL = tt.cacheTTL
				}

				statusCmd := &redis.StatusCmd{}
				if tt.redisError != nil {
					statusCmd.SetErr(tt.redisError)
				}
				mRClient.EXPECT().Set(
					mock.Anything,
					tt.redisKey,
					tt.redisVal,
					expectedTTL,
				).Return(statusCmd)
			}

			sut := &wssession.RedisCache{
				Client: mRClient,
				TTL:    tt.cacheTTL,
				TimestampFn: func() time.Time {
					return fixedTime
				},
			}

			err := sut.Add(tt.connID, tt.msg)
			if tt.expectError {
				assert.Error(err)
				if tt.errorMsg != "" {
					assert.Contains(err.Error(), tt.errorMsg)
				}
			} else {
				assert.NoError(err)
			}
		})
	}
}

func TestRedisCache_Items(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	ctx := context.Background()

	tests := []struct {
		name        string
		desc        string
		connID      string
		keys        []string
		messages    map[string]string
		getErrors   map[string]error
		keysError   error
		expectError bool
		expected    []*wssession.ResponseMsg
	}{
		{
			name:   "successful retrieval in order",
			desc:   "Messages should be returned in chronological order based on timestamp in key",
			connID: "conn1",
			// Keys as returned by Redis - could be any order
			keys: []string{
				"cache:conn1:3", // timestamp 3
				"cache:conn1:1", // timestamp 1
				"cache:conn1:2", // timestamp 2
			},
			// Messages should match the keys, will be retrieved in timestamp order
			messages: map[string]string{
				"cache:conn1:1": `{"id":"1","conn_id":"conn1","type":"test","message":"first"}`,
				"cache:conn1:2": `{"id":"2","conn_id":"conn1","type":"test","message":"second"}`,
				"cache:conn1:3": `{"id":"3","conn_id":"conn1","type":"test","message":"third"}`,
			},
			// Results should be ordered by the timestamp in the key
			expected: []*wssession.ResponseMsg{
				{ID: "1", ConnID: "conn1", Type: "test", Message: "first"},
				{ID: "2", ConnID: "conn1", Type: "test", Message: "second"},
				{ID: "3", ConnID: "conn1", Type: "test", Message: "third"},
			},
		},
		{
			name:        "keys command fails",
			desc:        "Should return error if Redis Keys command fails",
			connID:      "conn2",
			keysError:   fmt.Errorf("redis connection failed"),
			expectError: true,
		},
		{
			name:   "get command fails with non-nil error",
			desc:   "Should return error if Redis Get fails with anything other than key not found",
			connID: "conn3",
			keys: []string{
				"cache:conn3:1",
			},
			getErrors: map[string]error{
				"cache:conn3:1": fmt.Errorf("redis get failed"),
			},
			expectError: true,
		},
		{
			name:   "expired keys are correctly skipped",
			desc:   "Should skip expired keys and continue processing remaining keys",
			connID: "conn4",
			// Keys as returned by Redis
			keys: []string{
				"cache:conn4:2", // timestamp 2 (valid)
				"cache:conn4:1", // timestamp 1 (expired)
				"cache:conn4:4", // timestamp 4 (valid)
				"cache:conn4:3", // timestamp 3 (expired)
			},
			messages: map[string]string{
				"cache:conn4:2": `{"id":"2","conn_id":"conn4","type":"test","message":"valid1"}`,
				"cache:conn4:4": `{"id":"4","conn_id":"conn4","type":"test","message":"valid2"}`,
			},
			getErrors: map[string]error{
				"cache:conn4:1": redis.Nil,
				"cache:conn4:3": redis.Nil,
			},
			// Results should be ordered by timestamp, only including valid messages
			expected: []*wssession.ResponseMsg{
				{ID: "2", ConnID: "conn4", Type: "test", Message: "valid1"},
				{ID: "4", ConnID: "conn4", Type: "test", Message: "valid2"},
			},
		},
		{
			name:   "invalid json response",
			desc:   "Should return error if a message contains invalid JSON",
			connID: "conn5",
			keys: []string{
				"cache:conn5:1",
			},
			messages: map[string]string{
				"cache:conn5:1": `{invalid json`,
			},
			expectError: true,
		},
		{
			name:     "no keys found",
			desc:     "Should return empty slice (not nil) when no keys exist",
			connID:   "conn6",
			keys:     []string{},
			expected: []*wssession.ResponseMsg{},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mRClient := &mocks.MockRedisClient{}
			defer mRClient.AssertExpectations(t)

			// Setup Keys command expectation
			keysCmd := redis.NewStringSliceCmd(ctx)
			if tt.keysError != nil {
				keysCmd.SetErr(tt.keysError)
			} else {
				keysCmd.SetVal(tt.keys)
			}
			pattern := fmt.Sprintf("cache:%s:*", tt.connID)
			mRClient.EXPECT().Keys(mock.Anything, pattern).Return(keysCmd).Once()

			// The code will sort keys by timestamp and Get in that order
			// Sort keys by the numeric suffix which is our timestamp
			sortedKeys := make([]string, len(tt.keys))
			copy(sortedKeys, tt.keys)
			sort.Slice(sortedKeys, func(i, j int) bool {
				var t1, t2 int64
				fmt.Sscanf(sortedKeys[i], "cache:%*s:%d", &t1)
				fmt.Sscanf(sortedKeys[j], "cache:%*s:%d", &t2)
				return t1 < t2
			})

			// Setup Get command expectations in timestamp order
			for _, key := range sortedKeys {
				getCmd := redis.NewStringCmd(ctx)
				if err, ok := tt.getErrors[key]; ok {
					getCmd.SetErr(err)
				} else if val, ok := tt.messages[key]; ok {
					getCmd.SetVal(val)
				}
				mRClient.EXPECT().Get(mock.Anything, key).Return(getCmd).Once()
			}

			sut := &wssession.RedisCache{
				Client: mRClient,
			}

			msgs, err := sut.Items(tt.connID)
			if tt.expectError {
				assert.Error(err)
				return
			}

			assert.NoError(err)
			assert.Equal(len(tt.expected), len(msgs), "message count should match")
			for i, expected := range tt.expected {
				assert.Equal(expected.ID, msgs[i].ID)
				assert.Equal(expected.ConnID, msgs[i].ConnID)
				assert.Equal(expected.Type, msgs[i].Type)
				assert.Equal(expected.Message, msgs[i].Message)
			}
		})
	}
}