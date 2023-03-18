package cache

import (
	"sync"
	"time"
)

type Clock interface {
	Now() time.Time
}

// implement cache interface
type MemoryCache struct {
	dataMap map[int]MemoryValue
	lock    *sync.Mutex
	clock   Clock
}

// expiration values in cache
type MemoryValue struct {
	SetTime    int64
	Expiration int64
}

// create new instance inmemory cache
func InitInMemoryCache(clock Clock) *MemoryCache {
	return &MemoryCache{
		dataMap: make(map[int]MemoryValue, 0),
		lock:    &sync.Mutex{},
		clock:   clock,
	}
}

// Set adds a value to the cache.
func (c *MemoryCache) Add(key int, expiration int64) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.dataMap[key] = MemoryValue{
		SetTime:    c.clock.Now().Unix(),
		Expiration: expiration,
	}

	return nil
}

func (c *MemoryCache) Get(key int) (bool, error) {
	c.lock.Lock()
	defer c.lock.Unlock()

	value, ok := c.dataMap[key]
	if ok && c.clock.Now().Unix()-value.SetTime > value.Expiration {
		return false, nil
	}

	return ok, nil
}

func (c *MemoryCache) Delete(key int) {
	c.lock.Lock()
	defer c.lock.Unlock()

	delete(c.dataMap, key)
}
