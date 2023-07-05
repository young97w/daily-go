package ratelimit

import "time"

type TokenBucketLimiter struct {
	tokens chan struct{}
	close  chan struct{}
}

// NewTokenBucketLimiter interval 隔多久产生一个令牌
func NewTokenBucketLimiter(capacity int, interval time.Duration) *TokenBucketLimiter {
	ch := make(chan struct{}, capacity)
	closeCh := make(chan struct{})
	producer := time.NewTicker(interval)
	go func() {
		defer producer.Stop()
		for {
			select {
			case <-producer.C:
				select {
				case ch <- struct{}{}:
				default:
					// 令牌满了
				}
			case <-closeCh:
				return
			}
		}
	}()
	return &TokenBucketLimiter{
		tokens: ch,
		close:  closeCh,
	}
}
