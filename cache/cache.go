package cache

import (
	"github.com/coocood/freecache"
)

const(
	DEFAULT_CACHE_SIZE = 10*1024*1024
)

type Cache struct {
	ExpiryTime int
	memoryCache *freecache.Cache
}

func NewCache(cacheSize int) *Cache {
	return &Cache{memoryCache:freecache.NewCache(cacheSize)}
}

func (c *Cache) Set(key []byte, value []byte) error {
	return c.memoryCache.Set(key, value, c.ExpiryTime)
}

func (c *Cache) Get(key []byte) ([]byte, error) {
	return c.memoryCache.Get(key)
}


