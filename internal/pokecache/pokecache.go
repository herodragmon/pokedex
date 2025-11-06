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
    createdAt time.Time
    val       []byte
}

func NewCache(interval time.Duration) Cache {
    c := Cache{
        data:     make(map[string]cacheEntry),
        interval: interval,
    }
    go c.reapLoop()
    return c
}

func (c Cache) reapLoop() {
    ticker := time.NewTicker(c.interval)
    defer ticker.Stop()
    for range ticker.C {
        now := time.Now()
        c.mu.Lock()
        for key, entry := range c.data {
            if entry.createdAt.Add(c.interval).Before(now) {
                delete(c.data, key)
            }
        }
        c.mu.Unlock()
    }
}

func (c Cache) Add(key string, val []byte) {
    c.mu.Lock()
    c.data[key] = cacheEntry{
        createdAt: time.Now(),
        val:       val,
    }
    c.mu.Unlock()
}

func (c Cache) Get(key string) ([]byte, bool) {
    c.mu.Lock()
    entry, ok := c.data[key]
    c.mu.Unlock()
    if !ok {
        return nil, false
    }
    return entry.val, true
}

