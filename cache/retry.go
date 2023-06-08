package cache

import "time"

type RetryStrategy interface {
	Next() (time.Duration, bool)
}

type FXRetry struct {
	Interval time.Duration
	MaxCnt   int
	cnt      int
}

func (f *FXRetry) Next() (time.Duration, bool) {
	if f.cnt >= f.MaxCnt {
		return 0, false
	}
	return f.Interval, true
}
