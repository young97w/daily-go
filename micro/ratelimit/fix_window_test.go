package ratelimit

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"testing"
	"time"
)

func TestNewFixWindowLimiter(t *testing.T) {
	testCases := []struct {
		name     string
		lt       *FixWindowLimiter
		ctx      context.Context
		handler  func(ctx context.Context, req interface{}) (interface{}, error)
		wantErr  error
		wantResp any
	}{
		{
			name:    "reach limit window",
			lt:      NewFixWindowLimiter(time.Second, 0),
			ctx:     context.Background(),
			wantErr: errors.New("触发瓶颈了"),
		},
		{
			name: "pass window",
			lt:   NewFixWindowLimiter(time.Second, 1),
			ctx:  context.Background(),
			handler: func(ctx context.Context, req interface{}) (interface{}, error) {
				return 1, nil
			},
			wantResp: 1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			interceptor := tc.lt.BuildServerInterceptor()
			resp, err := interceptor(tc.ctx, 1, &grpc.UnaryServerInfo{}, tc.handler)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantResp, resp)
		})
	}
}

func TestFixWindowLimiter_BuildServerInterceptor(t *testing.T) {
	limiter := NewFixWindowLimiter(time.Second, 1)
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return 1, nil
	}
	interceptor := limiter.BuildServerInterceptor()
	ctx := context.Background()
	resp, err := interceptor(ctx, 1, &grpc.UnaryServerInfo{}, handler)
	require.NoError(t, err)
	assert.Equal(t, 1, resp)
	// reach limit
	resp, err = interceptor(ctx, 1, &grpc.UnaryServerInfo{}, handler)
	assert.Equal(t, errors.New("触发瓶颈了"), err)
	assert.Equal(t, nil, resp)

	time.Sleep(time.Second)
	resp, err = interceptor(ctx, 1, &grpc.UnaryServerInfo{}, handler)
	require.NoError(t, err)
	assert.Equal(t, 1, resp)

}
