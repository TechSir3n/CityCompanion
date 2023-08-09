package cache

import (
	"github.com/patrickmn/go-cache"
	"time"
)

type Cache interface {
	Set(queryStr string, x interface{}, d time.Duration)

	Get(queryStr string) (interface{}, bool)
}

type CacheImpl struct {
	cache *cache.Cache
}

const (
	defaultExpiration = time.Minute * 5
	cleanupTime       = time.Minute * 10
)

func NewCache() *CacheImpl {
	cache := cache.New(defaultExpiration, cleanupTime)
	return &CacheImpl{
		cache: cache,
	}
}

func (c *CacheImpl) Set(queryStr string, x interface{}, d time.Duration) {
	c.cache.Set(queryStr, x, d)
}

func (c *CacheImpl) Get(queryStr string) (interface{}, bool) {
	return c.cache.Get(queryStr)
}
