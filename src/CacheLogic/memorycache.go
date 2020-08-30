package CacheLogic

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

type cacheItem struct {
	Value      string
	Expiration int64
}

type Cache struct {
	*cache
}

type cache struct {
	items   map[string]cacheItem
	mu      sync.RWMutex
	cleanup *cleanup
}

type cleanup struct {
	Interval time.Duration
	stop     chan bool
}

func (item cacheItem) Expired() bool {

	return time.Now().UnixNano() > item.Expiration
}

func (c *cache) Insert(k string, x *string, m int) {

	var e int64
	d := time.Duration(m) * time.Minute
	if d > 0 {
		e = time.Now().Add(d).UnixNano()
	}
	c.mu.Lock()
	c.items[k] = cacheItem{
		Value:      *x,
		Expiration: e,
	}

	c.mu.Unlock()
}

// Get an item from the cache. Returns the item or nil, and a bool indicating
// whether the key was found.
func (c *cache) Get(k string) string {
	c.mu.RLock()
	// "Inlining" of get and Expired
	item, found := c.items[k]
	if !found {
		c.mu.RUnlock()
		panic("Key not found")
	}
	if item.Expiration > 0 {
		if time.Now().UnixNano() > item.Expiration {
			c.mu.RUnlock()
			panic("Key not found")
		}
	}
	c.mu.RUnlock()
	return item.Value
}

func (c *cache) KeyExists(k string) bool {
	c.mu.RLock()
	// "Inlining" of get and Expired
	item, found := c.items[k]
	if !found {
		c.mu.RUnlock()
		return false
	}
	if item.Expiration > 0 {
		if time.Now().UnixNano() > item.Expiration {
			c.mu.RUnlock()
			return false
		}
	}
	c.mu.RUnlock()
	return true
}

func (c *cache) delete(k string) {

	delete(c.items, k)

}

type keyAndValue struct {
	key   string
	value string
}

func (c *cache) DeleteExpired() {
	fmt.Println("Cleanup job")
	now := time.Now().UnixNano()
	c.mu.Lock()
	for k, v := range c.items {

		if v.Expiration > 0 && now > v.Expiration {
			c.delete(k)
			fmt.Println("Key expired :", k)
		}
	}
	c.mu.Unlock()

}

func (j *cleanup) Run(c *cache) {
	ticker := time.NewTicker(j.Interval)
	for {
		select {
		case <-ticker.C:
			c.DeleteExpired()
		case <-j.stop:
			ticker.Stop()
			return
		}
	}
}

func stopCleanup(c *Cache) {
	c.cleanup.stop <- true
}

func runCleanup(c *cache, ci time.Duration) {
	j := &cleanup{
		Interval: ci,
		stop:     make(chan bool),
	}
	c.cleanup = j
	go j.Run(c)
}

func newCache(m map[string]cacheItem) *cache {

	c := &cache{

		items: m,
	}
	return c
}

func newCacheWithCleanup(ci time.Duration, m map[string]cacheItem) *Cache {
	c := newCache(m)

	C := &Cache{c}
	if ci > 0 {
		runCleanup(c, ci)
		runtime.SetFinalizer(C, stopCleanup)
	}
	return C
}

func NewMemoryCache(cleanupInterval time.Duration) *Cache {
	items := make(map[string]cacheItem)
	return newCacheWithCleanup(cleanupInterval, items)
}
