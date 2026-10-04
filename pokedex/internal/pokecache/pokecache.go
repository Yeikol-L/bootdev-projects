package pokecache

import (
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

type Cache struct {
	entries      map[string]cacheEntry
	duration     time.Duration
	entriesMutex sync.Mutex
}

func NewCache(duration time.Duration) *Cache {
	cache := Cache{
		entries:  map[string]cacheEntry{},
		duration: duration,
	}
	cache.reapLoop()
	return &cache
}

func (c *Cache) Add(key string, val []byte) {
	c.entriesMutex.Lock()
	c.entries[key] = cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
	c.entriesMutex.Unlock()
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.entriesMutex.Lock()
	entry, ok := c.entries[key]
	defer c.entriesMutex.Unlock()
	if !ok {
		return nil, ok
	}
	return entry.val, true
}

func (c *Cache) reapLoop() {
	ticker := time.NewTicker(c.duration)
	go func() {
		for tickTime := range ticker.C {
			c.entriesMutex.Lock()
			for i, v := range c.entries {
				if tickTime.Sub(v.createdAt) > c.duration {
					delete(c.entries, i)
				}
			}
			c.entriesMutex.Unlock()
		}
	}()
}
