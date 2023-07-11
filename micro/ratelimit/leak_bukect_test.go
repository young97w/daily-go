package ratelimit

import (
	"context"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"testing"
	"time"
)

func TestNewLeakBucketLimiter(t *testing.T) {
	testCases := []struct {
		name        string
		ctx         context.Context
		lt          *LeakBucketLimiter
		handler     func(ctx context.Context, req interface{}) (interface{}, error)
		latencyFunc func()

		wantErr  error
		wantResp any
	}{
		{
			name: "context cancel",
			ctx: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				defer cancel()
				return ctx
			}(),
			lt:      NewLeakBucketLimiter(time.Second),
			wantErr: context.Canceled,
		},
		{
			name: "reach limit",
			ctx: func() context.Context {
				ctx, _ := context.WithCancel(context.Background())
				return ctx
			}(),
			lt:      NewLeakBucketLimiter(time.Second),
			wantErr: errors.New("触发瓶颈了"),
		},
		{
			name: "get token",
			ctx: func() context.Context {
				ctx, _ := context.WithCancel(context.Background())
				return ctx
			}(),
			lt: NewLeakBucketLimiter(time.Millisecond * 100),
			latencyFunc: func() {
				time.Sleep(time.Millisecond * 110)
			},
			handler: func(ctx context.Context, req interface{}) (interface{}, error) {
				return 1, nil
			},
			wantResp: 1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			interceptor := tc.lt.BuildServerInterceptor()
			if tc.latencyFunc != nil {
				tc.latencyFunc()
			}
			resp, err := interceptor(tc.ctx, 1, &grpc.UnaryServerInfo{}, tc.handler)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantResp, resp)
		})
	}
}

func TestProducer(t *testing.T) {
	l := NewLeakBucketLimiter(time.Millisecond * 10)
	cnt := 0
	fmt.Println("initialisation at ", time.Now().Format("2006-01-02 04:05"))
	go func() {
		for {
			select {
			case <-l.producer.C:
				fmt.Println("get token at ", time.Now().Format("2006-01-02 04:05"))
			}
		}
	}()

	time.Sleep(time.Millisecond * 100)
	fmt.Println(cnt)
}
