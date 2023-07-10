package ratelimit

import (
	"context"
	"errors"
	"fmt"
	__gen "geektime/micro/proto/.gen"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"testing"
	"time"
)

func TestTokenBucketLimiter_BuildServerInterceptor(t *testing.T) {
	testCases := []struct {
		name    string
		ctx     context.Context
		tl      *TokenBucketLimiter
		handler func(ctx context.Context, req interface{}) (interface{}, error)

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
			tl: func() *TokenBucketLimiter {
				ch := make(chan struct{})
				token := make(chan struct{}, 1)
				return &TokenBucketLimiter{
					tokens: token,
					close:  ch,
				}
			}(),
			wantErr: context.Canceled,
		},
		{
			name: "close channel closed",
			ctx: func() context.Context {
				return context.Background()
			}(),
			tl: func() *TokenBucketLimiter {
				ch := make(chan struct{}, 1)
				close(ch)
				token := make(chan struct{}, 1)
				return &TokenBucketLimiter{
					tokens: token,
					close:  ch,
				}
			}(),
			wantErr: errors.New("缺乏保护，拒绝请求"),
		},
		{
			name: "get token",
			ctx: func() context.Context {
				return context.Background()
			}(),
			tl: func() *TokenBucketLimiter {
				ch := make(chan struct{})
				token := make(chan struct{}, 1)
				token <- struct{}{}
				return &TokenBucketLimiter{
					tokens: token,
					close:  ch,
				}
			}(),
			wantResp: "hello",
			handler: func(ctx context.Context, req interface{}) (interface{}, error) {
				return "hello", nil
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			interceptor := tc.tl.BuildServerInterceptor()

			resp, err := interceptor(tc.ctx, &__gen.GetByIdReq{}, &grpc.UnaryServerInfo{}, tc.handler)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantResp, resp)
		})
	}
}

func TestTokenBucketLimiter_Tokens(t *testing.T) {
	limiter := NewTokenBucketLimiter(1, time.Second*20)
	time.Sleep(time.Second)
	defer limiter.Close()

	//limiter.tokens <- struct{}{}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return 1, nil
	}
	interceptor := limiter.BuildServerInterceptor()
	resp, err := interceptor(context.Background(), 1, &grpc.UnaryServerInfo{}, handler)
	assert.Equal(t, err, nil)
	assert.Equal(t, resp, 1)

	//触发限流
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	resp, err = interceptor(ctx, 1, &grpc.UnaryServerInfo{}, handler)
	assert.Equal(t, errors.New("到达瓶颈"), err)
	assert.Nil(t, resp)
}

func TestGetToken(t *testing.T) {
	limiter := NewTokenBucketLimiter(1, time.Second*2)
	fmt.Println("initialisation at ", time.Now().Format("2006-01-02 04:05"))
	go func() {
		for {
			select {
			case <-limiter.tokens:
				fmt.Println("get token at ", time.Now().Format("2006-01-02 04:05"))
			}
		}
	}()

	time.Sleep(time.Second * 10)
}
