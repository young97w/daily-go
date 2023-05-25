package cache

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/sync/singleflight"
	"log"
	"time"
)

var (
	ErrFailedToRefreshCache = errors.New("刷新缓存失败")
)

type ReadThroughCache struct {
	Cache
	LoadFunc   func(ctx context.Context, key string) (any, error)
	Expiration time.Duration
	g          singleflight.Group
}

func (r *ReadThroughCache) Get(ctx context.Context, key string) (any, error) {
	val, err := r.Cache.Get(ctx, key)
	if err == errKeyNotFound {
		val, err = r.LoadFunc(ctx, key)
		if err == nil {
			err2 := r.Cache.Set(ctx, key, val, r.Expiration)
			if err2 != nil {
				return val, fmt.Errorf("%w, 原因：%s", ErrFailedToRefreshCache, err2.Error())
			}
		}
	}
	return val, err
}

// GetV1 异步模式，无论有没有找到直接返回
func (r *ReadThroughCache) GetV1(ctx context.Context, key string) (any, error) {
	val, err := r.Cache.Get(ctx, key)
	if err == errKeyNotFound {
		go func() {
			val, err = r.LoadFunc(ctx, key)
			if err == nil {
				err2 := r.Cache.Set(ctx, key, val, r.Expiration)
				if err2 != nil {
					log.Fatalln(fmt.Errorf("%w, 原因：%s", ErrFailedToRefreshCache, err2.Error()))
				}
			}
		}()
	}
	return val, err
}

// GetV2 异步模式，如果cache miss ， 则从db捞数据再返回
func (r *ReadThroughCache) GetV2(ctx context.Context, key string) (any, error) {
	val, err := r.Cache.Get(ctx, key)
	if err == errKeyNotFound {
		val, err = r.LoadFunc(ctx, key)
		go func() {
			if err == nil {
				err2 := r.Cache.Set(ctx, key, val, r.Expiration)
				if err2 != nil {
					log.Fatalln(fmt.Errorf("%w, 原因：%s", ErrFailedToRefreshCache, err2.Error()))
				}
			}
		}()
	}
	return val, err
}

// GetV3 如果cache miss ， 则从db捞数据再返回，使用singleflight
func (r *ReadThroughCache) GetV3(ctx context.Context, key string) (any, error) {
	val, err := r.Cache.Get(ctx, key)
	if err == errKeyNotFound {
		val, err, _ = r.g.Do(key, func() (interface{}, error) {
			val, err = r.LoadFunc(ctx, key)
			if err == nil {
				err2 := r.Cache.Set(ctx, key, val, r.Expiration)
				if err2 != nil {
					return val, fmt.Errorf("%w, 原因：%s", ErrFailedToRefreshCache, err2.Error())
				}
			}
			return val, err
		})

	}
	return val, err
}
