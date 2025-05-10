package pokecache

import (
	"time"
	"sync"
)

type cacheEntry struct {
	createdAt		time.Time
	data				[]byte
}

type Cache struct {
	entries		map[string]cacheEntry
	mu	     	*sync.RWMutex
}

func NewCache(expireTime time.Duration) *Cache {
	cache := &Cache {
		entries: 	make(map[string]cacheEntry),
		mu: 			&sync.RWMutex{},
	}

	go func () {
		ticker := time.NewTicker(expireTime)
		defer ticker.Stop()

		for range ticker.C {
			cache.reapLoop(expireTime)
		}
	}()

	return cache
}

func (c *Cache) reapLoop(expireTime time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for key, entry := range c.entries {
		now := time.Now()
		if now.Sub(entry.createdAt) > expireTime {
			delete(c.entries, key)
		}	
	}
}

func (c *Cache) Add(key string, data []byte) {
	entry := cacheEntry {
		createdAt:	time.Now(),
		data:	data,
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.entries[key] = entry
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if entry, ok := c.entries[key]; !ok {
		return nil, false
	} else {
		return entry.data, true
	}
}


