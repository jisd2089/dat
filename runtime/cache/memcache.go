package cache

import (
	"fmt"
	"sync"
	"time"
)

var __global_memcache__ *MemCache

type MemCache struct {
	dict map[string]interface{}
	mu   sync.RWMutex //读写锁
}

func GetMemCacheInstance() *MemCache {
	if __global_memcache__ == nil {
		__global_memcache__ = new(MemCache)
		__global_memcache__.dict = make(map[string]interface{})
	}

	return __global_memcache__
}

func (c *MemCache) SetMemCache(cache_key string, obj interface{}, timeout int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	timeout_key := fmt.Sprintf("%s_timeout", cache_key)
	_timeout := int(time.Now().Unix()) + timeout

	c.dict[cache_key] = obj
	c.dict[timeout_key] = _timeout

}

func (c *MemCache) GetMemCache(cache_key string) interface{} {
	c.mu.RLock()
	defer c.mu.RUnlock()
	data, ok := c.dict[cache_key]
	if ok {
		return data
	}

	return nil
}
