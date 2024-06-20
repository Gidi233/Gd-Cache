package Cache

import (
	"sync"

	"github.com/Gidi233/Gd-Cache/lru"
)

type cache struct {
	mu         sync.Mutex
	lru        *lru.Cache
	cacheBytes int64
}

func newCache(cacheBytes int64) *cache {
	return &cache{
		cacheBytes: cacheBytes,
	}
}

func (c *cache) add(key string, value ByteView) {
	c.mu.Lock()
	defer c.mu.Unlock()
	// 延迟初始化: 在第一次用到 lru 时才初始化
	if c.lru == nil {
		c.lru = lru.New(c.cacheBytes, nil)
	}
	c.lru.Add(key, value)
}

func (c *cache) get(key string) (value ByteView, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil {
		return
	}

	if v, ok := c.lru.Get(key); ok {
		return v.(ByteView), ok
	}

	return
}

// func (c *cache) delete(key string) (ok bool) {
// 	c.mu.Lock()
// 	defer c.mu.Unlock()
// 	if c.lru == nil {
// 		return
// 	}

// 	if ok := c.lru.Delete(key); ok {
// 		return true
// 	}

// 	return false
// }
