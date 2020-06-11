package repository

import (
	"sync"
	"time"
)

type Cache interface {
	Get(key string) []*News
	Set(key string, content []*News, duration time.Duration)
}

// Item is a cached reference
type Item struct {
	Content    []*News
	Expiration int64
}

// Expired returns true if the item has expired.
func (item *Item) Expired() bool {
	if item.Expiration == 0 {
		return false
	}
	return time.Now().UnixNano() > item.Expiration
}

//Cache mecanism for caching strings in memory
type cache struct {
	items map[string]Item
	mu    *sync.RWMutex
}

//NewCache creates a new in memory storage
func InitCache() Cache {
	return &cache{
		items: make(map[string]Item),
		mu:    &sync.RWMutex{},
	}
}

//Get a cached content by key
func (c *cache) Get(key string) []*News {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item := c.items[key]
	if item.Expired() {
		delete(c.items, key)
		return nil
	}
	return item.Content
}

//Set a cached content by key
func (c *cache) Set(key string, content []*News, duration time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[key] = Item{
		Content:    content,
		Expiration: time.Now().Add(duration).UnixNano(),
	}
}
