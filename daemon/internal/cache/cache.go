package cache

import (
	"sync"
	"time"

	"daemon/internal/protocol"
)

type CacheEntry struct {
	Data        *protocol.ServerStatus
	ExpiresAt   time.Time
	CreatedAt   time.Time
	AccessCount int
}

type Cache struct {
	entries     map[string]*CacheEntry
	mutex       sync.RWMutex
	ttl         time.Duration
	maxEntries  int
	stats       CacheStats
	stopCleanup chan bool
}

type CacheStats struct {
	Hits        int64
	Misses      int64
	Evictions   int64
	TotalSize   int64
	LastCleanup time.Time
}

func NewCache(ttl time.Duration) *Cache {
	cache := &Cache{
		entries:     make(map[string]*CacheEntry),
		ttl:         ttl,
		maxEntries:  1000,
		stopCleanup: make(chan bool),
	}

	go cache.startCleanup()

	return cache
}

func NewCacheWithConfig(ttl time.Duration, maxEntries int) *Cache {
	cache := &Cache{
		entries:     make(map[string]*CacheEntry),
		ttl:         ttl,
		maxEntries:  maxEntries,
		stopCleanup: make(chan bool),
	}

	go cache.startCleanup()

	return cache
}

func (c *Cache) Get(key string) (*protocol.ServerStatus, bool) {
	c.mutex.RLock()
	entry, exists := c.entries[key]
	c.mutex.RUnlock()

	if !exists {
		c.stats.Misses++
		return nil, false
	}

	if time.Now().After(entry.ExpiresAt) {
		c.mutex.Lock()
		delete(c.entries, key)
		c.mutex.Unlock()
		c.stats.Misses++
		return nil, false
	}

	entry.AccessCount++
	c.stats.Hits++
	return entry.Data, true
}

func (c *Cache) Set(key string, data *protocol.ServerStatus) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if len(c.entries) >= c.maxEntries && len(c.entries) > 0 {
		c.evictOldest()
	}

	c.entries[key] = &CacheEntry{
		Data:        data,
		ExpiresAt:   time.Now().Add(c.ttl),
		CreatedAt:   time.Now(),
		AccessCount: 1,
	}

	c.stats.TotalSize = int64(len(c.entries))
}

func (c *Cache) Delete(key string) {
	c.mutex.Lock()
	delete(c.entries, key)
	c.mutex.Unlock()
}

func (c *Cache) Clear() {
	c.mutex.Lock()
	c.entries = make(map[string]*CacheEntry)
	c.mutex.Unlock()
}

func (c *Cache) GetStats() CacheStats {
	c.mutex.RLock()
	stats := c.stats
	stats.TotalSize = int64(len(c.entries))
	c.mutex.RUnlock()
	return stats
}

func (c *Cache) GetKeys() []string {
	c.mutex.RLock()
	keys := make([]string, 0, len(c.entries))
	for key := range c.entries {
		keys = append(keys, key)
	}
	c.mutex.RUnlock()
	return keys
}

func (c *Cache) evictOldest() {
	var oldestKey string
	var oldestTime time.Time
	var oldestAccessCount int

	first := true
	for key, entry := range c.entries {
		if first || entry.CreatedAt.Before(oldestTime) ||
			(entry.CreatedAt.Equal(oldestTime) && entry.AccessCount < oldestAccessCount) {
			oldestKey = key
			oldestTime = entry.CreatedAt
			oldestAccessCount = entry.AccessCount
			first = false
		}
	}

	if oldestKey != "" {
		delete(c.entries, oldestKey)
		c.stats.Evictions++
	}
}

func (c *Cache) cleanup() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	now := time.Now()
	removed := 0

	for key, entry := range c.entries {
		if now.After(entry.ExpiresAt) {
			delete(c.entries, key)
			removed++
		}
	}

	c.stats.Evictions += int64(removed)
	c.stats.LastCleanup = now
	c.stats.TotalSize = int64(len(c.entries))
}

func (c *Cache) startCleanup() {
	ticker := time.NewTicker(c.ttl / 2)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.cleanup()
		case <-c.stopCleanup:
			return
		}
	}
}

func (c *Cache) StopCleanup() {
	close(c.stopCleanup)
}

func (c *Cache) Size() int {
	c.mutex.RLock()
	size := len(c.entries)
	c.mutex.RUnlock()
	return size
}
