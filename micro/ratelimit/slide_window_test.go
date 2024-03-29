package ratelimit

import (
	"context"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"sync/atomic"
	"testing"
	"time"
)

func TestNewSlideWindowLimiter(t *testing.T) {
	testCases := []struct {
		name    string
		lt      *SlideWindowLimiter
		handler func(ctx context.Context, req interface{}) (interface{}, error)

		wantErr  error
		wantResp any
	}{
		{
			name: "pass limiter",
			lt:   NewSlideWindowLimiter(1, 20),
			handler: func(ctx context.Context, req interface{}) (interface{}, error) {
				return 1, nil
			},
			wantResp: 1,
		},
		{
			name:    "reach limit",
			lt:      NewSlideWindowLimiter(1, 0),
			wantErr: errors.New("到达瓶颈"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			interceptor := tc.lt.BuildServerInterceptor()
			resp, err := interceptor(context.Background(), 1, &grpc.UnaryServerInfo{}, tc.handler)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantResp, resp)
		})
	}
}

func TestSlideWindow(t *testing.T) {
	interceptor := NewSlideWindowLimiter(time.Millisecond*100, 100).BuildServerInterceptor()
	c := int64(0)
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return 1, nil
	}
	for i := 0; i < 150; i++ {
		go func() {
			_, err := interceptor(nil, nil, nil, handler)
			if err != nil {
				atomic.AddInt64(&c, 1)
			}
		}()
	}
	time.Sleep(time.Second)
	fmt.Println(c)
	c = 0

	for i := 0; i < 150; i++ {
		go func() {
			_, err := interceptor(nil, nil, nil, handler)
			if err != nil {
				atomic.AddInt64(&c, 1)
			}
		}()
	}
	time.Sleep(time.Second)
	fmt.Println(c)
}
