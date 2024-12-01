package wssession

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"golang.org/x/net/context"
)

type RedisClient interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	Get(ctx context.Context, key string) *redis.StringCmd
	Keys(ctx context.Context, pattern string) *redis.StringSliceCmd
}

type RedisCache struct {
	Client      RedisClient
	TTL         time.Duration
	TimestampFn func() time.Time // Defaults to time.Now, typically you won't want to modify this unless writing a test
}

var ctx = context.Background()

func (r *RedisCache) ttl() time.Duration {
	if r.TTL == 0 {
		return time.Minute
	}
	return r.TTL
}

func (r *RedisCache) timeNow() time.Time {
	if r.TimestampFn == nil {
		return time.Now()
	}
	return r.TimestampFn()
}

func (r *RedisCache) Add(connID string, msg ResponseMsg) error {
	if connID == "" {
		return fmt.Errorf("connection ID cannot be empty")
	}
	msgStr, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}
	// Use the current timestamp as part of the key
	timestamp := r.timeNow().UnixNano()
	key := fmt.Sprintf("cache:%s:%d", connID, timestamp)

	// Add the message with a TTL
	resp := r.Client.Set(ctx, key, string(msgStr), r.ttl())
	if resp.Err() != nil {
		return fmt.Errorf("failed to set message: %w", resp.Err())
	}
	return nil
}

func (r *RedisCache) Items(connID string) ([]*ResponseMsg, error) {
	pattern := fmt.Sprintf("cache:%s:*", connID)
	keys, err := r.Client.Keys(ctx, pattern).Result()
	if err != nil {
		return nil, err
	}

	sortedKeys := r.sortKeysByTimestamp(keys)

	// Use capacity of sortedKeys as it's our upper bound
	resp := make([]*ResponseMsg, 0, len(sortedKeys))
	for _, key := range sortedKeys {
		message, err := r.Client.Get(ctx, key).Result()
		if err != nil {
			if err == redis.Nil {
				continue
			}
			return nil, err
		}
		var m ResponseMsg
		err = json.Unmarshal([]byte(message), &m)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal message: %w", err)
		}
		resp = append(resp, &m)
	}
	return resp, nil
}

func (r *RedisCache) sortKeysByTimestamp(keys []string) []string {
	type keyWithTimestamp struct {
		key       string
		timestamp int64
	}
	keyTimestamps := make([]keyWithTimestamp, len(keys))

	for i, key := range keys {
		// Split by : and take the last part as timestamp
		parts := strings.Split(key, ":")
		if len(parts) >= 3 {
			timestamp, _ := strconv.ParseInt(parts[len(parts)-1], 10, 64)
			keyTimestamps[i] = keyWithTimestamp{key, timestamp}
		}
	}

	sort.Slice(keyTimestamps, func(i, j int) bool {
		return keyTimestamps[i].timestamp < keyTimestamps[j].timestamp
	})

	sortedKeys := make([]string, len(keys))
	for i, kt := range keyTimestamps {
		sortedKeys[i] = kt.key
	}
	return sortedKeys
}
