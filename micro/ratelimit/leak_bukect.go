package ratelimit

import (
	"context"
	"errors"
	"google.golang.org/grpc"
	"time"
)

type LeakBucketLimiter struct {
	producer *time.Ticker
}

func NewLeakBucketLimiter(interval time.Duration) *LeakBucketLimiter {
	return &LeakBucketLimiter{producer: time.NewTicker(interval)}
}

func (l *LeakBucketLimiter) BuildServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		select {
		case <-ctx.Done():
			err = ctx.Err()
		case <-l.producer.C:
			resp, err = handler(ctx, req)
		default:
			err = errors.New("触发瓶颈了")
		}
		return
	}
}

func (l *LeakBucketLimiter) Close() error {
	l.producer.Stop()
	return nil
}
