package cache

import (
	"context"
	"fmt"
	"geektime/cache/mocks"
	"github.com/golang/mock/gomock"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestRedisCache_SetU(t *testing.T) {
	testCases := []struct {
		name       string
		mock       func(ctrl *gomock.Controller) redis.Cmdable
		key        string
		value      string
		expiration time.Duration
		wantErr    error
	}{
		{
			name: "set ok",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := mocks.NewMockCmdable(ctrl)

				status := redis.NewStatusCmd(context.Background())
				status.SetVal("OK")

				cmd.EXPECT().Set(context.Background(), "key1", "value1", time.Second).Return(status)
				return cmd
			},
			expiration: time.Second,
			key:        "key1",
			value:      "value1",
		},
		{
			name: "time out",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := mocks.NewMockCmdable(ctrl)
				status := redis.NewStatusCmd(context.Background())
				status.SetErr(context.DeadlineExceeded)
				cmd.EXPECT().
					Set(context.Background(), "key1", "value1", time.Second).
					Return(status)
				return cmd
			},
			key:        "key1",
			value:      "value1",
			expiration: time.Second,
			wantErr:    context.DeadlineExceeded,
		},
		{
			name: "unexpected msg",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := mocks.NewMockCmdable(ctrl)
				status := redis.NewStatusCmd(context.Background())
				status.SetVal("NO OK")
				cmd.EXPECT().
					Set(context.Background(), "key1", "value1", time.Second).
					Return(status)
				return cmd
			},
			key:        "key1",
			value:      "value1",
			expiration: time.Second,
			wantErr:    fmt.Errorf("%w, 返回信息 %s", errFailedToSetCache, "NO OK"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			c := NewRedisCache(tc.mock(ctrl))
			err := c.Set(context.Background(), tc.key, tc.value, tc.expiration)
			assert.Equal(t, tc.wantErr, err)
		})
	}
}

func TestRedisCache_Get(t *testing.T) {
	testCases := []struct {
		name    string
		mock    func(ctrl *gomock.Controller) redis.Cmdable
		key     string
		wantErr error
		wantVal string
	}{
		{
			name: "get value",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := mocks.NewMockCmdable(ctrl)

				str := redis.NewStringCmd(context.Background())
				str.SetVal("value1")

				cmd.EXPECT().Get(context.Background(), "key1").Return(str)
				return cmd
			},
			key:     "key1",
			wantVal: "value1",
		},
		{
			name: "timeout",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := mocks.NewMockCmdable(ctrl)

				str := redis.NewStringCmd(context.Background())
				str.SetErr(context.DeadlineExceeded)

				cmd.EXPECT().Get(context.Background(), "key1").Return(str)
				return cmd
			},
			key:     "key1",
			wantErr: context.DeadlineExceeded,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			c := NewRedisCache(tc.mock(ctrl))

			val, err := c.Get(context.Background(), tc.key)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantVal, val)
		})
	}
}
