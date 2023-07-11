package ratelimit

import (
	"context"
	"errors"
	"google.golang.org/grpc"
	"sync/atomic"
	"time"
)

type FixWindowLimiter struct {
	timestamp int64
	interval  int64
	rate      int64
	cnt       int64
}

func NewFixWindowLimiter(interval time.Duration, rate int64) *FixWindowLimiter {
	return &FixWindowLimiter{
		timestamp: time.Now().UnixNano(),
		interval:  interval.Nanoseconds(),
		rate:      rate,
	}
}

func (f *FixWindowLimiter) BuildServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		// 对比timestamp和now，如果到了下一个窗口期，重置
		// 否则cnt++
		timestamp := atomic.LoadInt64(&f.timestamp)
		current := time.Now().UnixNano()
		cnt := atomic.LoadInt64(&f.cnt)
		if timestamp+f.interval < current {
			if atomic.CompareAndSwapInt64(&f.timestamp, timestamp, current) {
				atomic.CompareAndSwapInt64(&f.cnt, cnt, 0)
			}
		}
		cnt = atomic.AddInt64(&f.cnt, 1)
		if cnt > f.rate {
			err = errors.New("触发瓶颈了")
			return
		}
		resp, err = handler(ctx, req)
		return
	}
}
