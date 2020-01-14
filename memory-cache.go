package dfuse

import (
	"time"

	gocache "github.com/patrickmn/go-cache"
)

type CacheInter interface {
	Set(key string, value interface{}, expirationTime time.Duration)
	Get(key string) (interface{}, bool)
	Add(key string, value interface{}, expirationTime time.Duration) error
	Delete(key string)
	Flush()
}

type memoryCache struct {
	cache *gocache.Cache
}

func NewMemoryCache(expireTime, cleanupInterval time.Duration) *memoryCache {
	cache := gocache.New(expireTime, cleanupInterval)
	cache.Flush()
	return &memoryCache{cache: cache}
}

// 设置缓存
func (c *memoryCache) Set(key string, value interface{}, expirationTime time.Duration) {
	c.cache.Set(key, value, expirationTime)
}

// 获取缓存
func (c *memoryCache) Get(key string) (interface{}, bool) {
	return c.cache.Get(key)
}

// 添加缓存当key已存在会报错
func (c *memoryCache) Add(key string, value interface{}, expirationTime time.Duration) error {
	err := c.cache.Add(key, value, expirationTime)
	if err != nil {
		return err
	}
	return nil
}

// 删除缓存
func (c *memoryCache) Delete(key string) {
	c.cache.Delete(key)
}

// 刷新缓存
func (c *memoryCache) Flush() {
	c.cache.Flush()
}
