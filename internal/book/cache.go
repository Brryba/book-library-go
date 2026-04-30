package book

import (
	"sync"
	"time"
)

type Cache struct {
	mu      sync.RWMutex
	items   map[string]cacheEntry
	ttl     time.Duration
	addCh   chan addRequest
	closeCh chan struct{}
}

type cacheEntry struct {
	value     []Book
	expiresAt time.Time
}

type addRequest struct {
	key   string
	value []Book
}

func NewCache(ttl time.Duration) *Cache {
	c := &Cache{
		items:   make(map[string]cacheEntry),
		ttl:     ttl,
		addCh:   make(chan addRequest, 100),
		closeCh: make(chan struct{}),
	}
	go c.run()
	return c
}

func (c *Cache) Get(key string) ([]Book, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	e, ok := c.items[key]
	if !ok || time.Now().After(e.expiresAt) {
		return nil, false
	}
	return e.value, true
}

func (c *Cache) Add(key string, value []Book) {
	c.addCh <- addRequest{key: key, value: value}
}

func (c *Cache) Invalidate(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, key)
}

func (c *Cache) Close() {
	close(c.closeCh)
}

func (c *Cache) run() {
	ticker := time.NewTicker(c.ttl)
	defer ticker.Stop()

	for {
		select {
		case req := <-c.addCh:
			c.mu.Lock()
			c.items[req.key] = cacheEntry{
				value:     req.value,
				expiresAt: time.Now().Add(c.ttl),
			}
			c.mu.Unlock()

		case <-ticker.C:
			c.mu.Lock()
			now := time.Now()
			for k, e := range c.items {
				if now.After(e.expiresAt) {
					delete(c.items, k)
				}
			}
			c.mu.Unlock()

		case <-c.closeCh:
			return
		}
	}
}
