//go:build e2e

package cache

import (
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestRedisCache_Set(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	c := NewRedisCache(rdb)
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	err := c.Set(ctx, "key1", "value1", time.Minute)
	require.NoError(t, err)

	val, err := c.Get(ctx, "key1")
	require.NoError(t, err)
	assert.Equal(t, "value1", val)
}

func TestRedisCache_e2e_SetV1(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	testCases := []struct {
		name string
		//before func(t *testing.T)
		after func(t *testing.T)

		key        string
		value      string
		expiration time.Duration

		wantErr error
	}{
		{
			name: "set value",
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				defer cancel()
				res, err := rdb.Get(ctx, "key1").Result()
				require.NoError(t, err)
				assert.Equal(t, "value1", res)
				_, err = rdb.Del(ctx, "key1").Result()
				require.NoError(t, err)
			},
			key:        "key1",
			value:      "value1",
			expiration: time.Minute,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c := NewRedisCache(rdb)
			//tc.before()
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
			defer cancel()
			err := c.Set(ctx, tc.key, tc.value, tc.expiration)
			require.NoError(t, err)
			tc.after(t)
		})
	}

}
