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

func TestClient_TryLock_e2e(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{Addr: "localhost:6379"})

	testCases := []struct {
		name       string
		key        string
		expiration time.Duration

		before func(t *testing.T)
		after  func(t *testing.T)

		wantErr  error
		wantLock *Lock
	}{
		{
			name:       "locked",
			key:        "key1",
			expiration: time.Second * 3,
			before:     func(t *testing.T) {},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				defer cancel()
				res, err := rdb.GetDel(ctx, "key1").Result()
				require.NoError(t, err)
				assert.NotEmpty(t, res)
			},
			wantLock: &Lock{
				key:        "key1",
				expiration: time.Second * 3,
			},
		},
		{
			name:       "hold by others",
			key:        "key1",
			expiration: time.Second * 3,
			before: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				defer cancel()
				res, err := rdb.Set(ctx, "key1", "value1", time.Second*3).Result()
				require.NoError(t, err)
				assert.NotEmpty(t, res)
			},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				defer cancel()
				res, err := rdb.GetDel(ctx, "key1").Result()
				require.NoError(t, err)
				assert.NotEmpty(t, res)
			},
			wantErr: ErrFailedToPreemptLock,
		},
	}

	client := NewClient(rdb)
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
			defer cancel()

			tc.before(t)
			lock, err := client.TryLock(ctx, tc.key, tc.expiration)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantLock.key, lock.key)
			assert.Equal(t, tc.wantLock.expiration, lock.expiration)
			assert.NotEmpty(t, lock.value)
			assert.NotNil(t, lock.client)
			tc.after(t)
		})
	}
}

func TestLock_Unlock_e2e(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{Addr: "localhost:6379"})

	testCases := []struct {
		name string
		lock *Lock

		before func(t *testing.T)
		after  func(t *testing.T)

		wantErr error
	}{
		{
			name: "unlocked",
			lock: &Lock{
				client:     rdb,
				key:        "key1",
				value:      "value1",
				expiration: time.Second * 3,
			},
			before: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				defer cancel()
				res, err := rdb.Set(ctx, "key1", "value1", time.Second*3).Result()
				require.NoError(t, err)
				assert.Equal(t, "OK", res)
			},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				defer cancel()
				res, err := rdb.Exists(ctx, "key1").Result()
				require.NoError(t, err)
				assert.Equal(t, int64(0), res)
			},
		},
		{
			name: "key not hold",
			lock: &Lock{
				client:     rdb,
				key:        "key1",
				value:      "value1",
				expiration: time.Second * 3,
			},
			before:  func(t *testing.T) {},
			after:   func(t *testing.T) {},
			wantErr: ErrLockNotHold,
		},
		{
			name: "key hold by others",
			lock: &Lock{
				client:     rdb,
				key:        "key1",
				value:      "value1",
				expiration: time.Second * 3,
			},
			before: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				defer cancel()
				res, err := rdb.Set(ctx, "key1", "123", time.Second*3).Result()
				require.NoError(t, err)
				assert.Equal(t, "OK", res)
			},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				defer cancel()
				res, err := rdb.GetDel(ctx, "key1").Result()
				require.NoError(t, err)
				assert.NotEmpty(t, res)
			},
			wantErr: ErrLockNotHold,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.before(t)
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
			defer cancel()

			err := tc.lock.Unlock(ctx)
			assert.Equal(t, tc.wantErr, err)
			tc.after(t)
		})
	}
}

func TestLock_Refresh_e2e(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{Addr: "localhost:6379"})

	testCases := []struct {
		name string
		lock *Lock

		before func(t *testing.T)
		after  func(t *testing.T)

		wantErr error
	}{
		{
			name: "refreshed",
			lock: &Lock{
				client:     rdb,
				key:        "key1",
				value:      "value1",
				expiration: time.Minute,
			},
			before: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				defer cancel()
				res, err := rdb.Set(ctx, "key1", "value1", time.Second*3).Result()
				require.NoError(t, err)
				assert.Equal(t, "OK", res)
			},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				defer cancel()
				duration, err := rdb.TTL(ctx, "key1").Result()
				require.NoError(t, err)
				assert.True(t, duration > time.Second*3)

				res, err := rdb.GetDel(ctx, "key1").Result()
				require.NoError(t, err)
				assert.NotEmpty(t, res)
			},
		},
		{
			name: "key not hold",
			lock: &Lock{
				client:     rdb,
				key:        "key1",
				value:      "value1",
				expiration: time.Minute,
			},
			before:  func(t *testing.T) {},
			after:   func(t *testing.T) {},
			wantErr: ErrLockNotHold,
		},
		{
			name: "key hold by others",
			lock: &Lock{
				client:     rdb,
				key:        "key1",
				value:      "value1",
				expiration: time.Minute,
			},
			before: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				defer cancel()
				res, err := rdb.Set(ctx, "key1", "123", time.Second*3).Result()
				require.NoError(t, err)
				assert.Equal(t, "OK", res)
			},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				defer cancel()
				res, err := rdb.GetDel(ctx, "key1").Result()
				require.NoError(t, err)
				assert.NotEmpty(t, res)
			},
			wantErr: ErrLockNotHold,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.before(t)
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
			defer cancel()

			err := tc.lock.Refresh(ctx)
			assert.Equal(t, tc.wantErr, err)
			tc.after(t)
		})
	}
}
