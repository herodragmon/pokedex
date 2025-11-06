package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	data     map[string]cacheEntry
	mu       sync.Mutex
	interval time.Duration
}

type cacheEntry struct {
	createdAt   time.Time 
	val         []byte
}

func NewCache(interval time.Duration) Cache {
	var c Cache
	c.data  = make(map[string]cacheEntry)

	c.interval = interval
	go c.reapLoop()
	return c
}

func (c Cache) reapLoop() {
	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()
	for range ticker.C {
		t := time.Now()
		c.mu.Lock()
		for key, entry := range c.data {
			if entry.createdAt + c.interval < t {
				delete (c.data, key)
			}
		}
		defer c.mu.Unlock()
	}
}

func (c Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	newEntry := cacheEntry{
							createdAt: time.Now()
							val:       val,
	}
	c.data[key] = newEntry
}

func (c Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	k, ok := c.data[key]
	if !ok {
		return nil, false
	}
	return c.data[key].val, true
}
