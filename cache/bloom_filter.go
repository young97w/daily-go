package cache

import (
	"context"
	"golang.org/x/sync/singleflight"
)

type BloomFilter interface {
	HasKey(ctx context.Context, key string) bool
}

type BloomFilterCache struct {
	ReadThroughCache
}

func NewBloomFilterCache(cache Cache, bf BloomFilter, loadFunc func(ctx context.Context, key string) (any, error)) *BloomFilterCache {
	return &BloomFilterCache{ReadThroughCache{
		Cache: cache,
		LoadFunc: func(ctx context.Context, key string) (any, error) {
			if !bf.HasKey(ctx, key) {
				return nil, errKeyNotFound
			}
			return loadFunc(ctx, key)
		},
		g: singleflight.Group{},
	}}
}
