package ratelimit

import (
	"container/list"
	"context"
	"errors"
	"google.golang.org/grpc"
	"sync"
	"time"
)

type SlideWindowLimiter struct {
	queue    *list.List
	interval int64
	rate     int
	mutex    sync.Mutex
}

func NewSlideWindowLimiter(interval time.Duration, rate int) *SlideWindowLimiter {
	return &SlideWindowLimiter{
		queue:    list.New(),
		interval: interval.Nanoseconds(),
		rate:     rate,
		mutex:    sync.Mutex{},
	}
}

func (l *SlideWindowLimiter) BuildServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		now := time.Now().UnixNano()
		boundary := now - l.interval
		l.mutex.Lock()
		if l.rate > l.queue.Len() {
			l.queue.PushBack(now)
			l.mutex.Unlock()
			resp, err = handler(ctx, req)
			return
		}
		// clear front items in list
		timestamp := l.queue.Front()
		for timestamp != nil && timestamp.Value.(int64) < boundary {
			l.queue.Remove(timestamp)
			timestamp = l.queue.Front()
		}
		l.mutex.Unlock()

		if l.rate <= l.queue.Len() {
			err = errors.New("到达瓶颈")
			return
		}
		resp, err = handler(ctx, req)
		now = time.Now().UnixNano()
		l.mutex.Lock()
		l.queue.PushBack(now)
		l.mutex.Unlock()
		return
	}
}
