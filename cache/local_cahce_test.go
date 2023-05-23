package cache

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestBuildInMapCache_Get(t *testing.T) {
	testCases := []struct {
		name    string
		key     string
		cache   func() *BuildInMapCache
		wantVal any
		wantErr error
	}{
		{
			name: "not found",
			key:  "key",
			cache: func() *BuildInMapCache {
				return NewBuildInMapCache(1)
			},
			wantErr: errKeyNotFound,
		},
		{
			name: "key1",
			key:  "key1",
			cache: func() *BuildInMapCache {
				res := NewBuildInMapCache(1)
				res.Set(context.Background(), "key1", 111, time.Second*2)
				return res
			},
			wantVal: 111,
		},
		{
			name: "expired",
			key:  "key1",
			cache: func() *BuildInMapCache {
				res := NewBuildInMapCache(1)
				res.Set(context.Background(), "key1", 111, time.Second*1)
				time.Sleep(2 * time.Second)
				return res
			},
			wantErr: errKeyNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			val, err := tc.cache().Get(context.Background(), tc.key)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantVal, val)
		})
	}
}

func TestBuildInMapCache_Set(t *testing.T) {
	ctn := 0
	c := NewBuildInMapCache(1, WithEvictedCallBack(func(key string, val any) {
		ctn++
	}))
	err := c.Set(context.Background(), "key1", "val1", 1)
	require.NoError(t, err)
	time.Sleep(2 * time.Second)
	assert.Equal(t, 1, ctn)
	fmt.Println(c.data)
}
