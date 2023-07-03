package memcache

import (
	"container/list"
	"sync"
	"time"
)

type lruObj struct {
	key   string
	value interface{}

	exp *time.Time
}

// newLRUObj returns a new instance of lruObj with the given key, value and ttl
func newLRUObj(key string, value interface{}, ttl time.Duration) *lruObj {
	var t *time.Time

	if ttl > 0 {
		temp := time.Now().Add(ttl)
		t = &temp
	}

	return &lruObj{
		key:   key,
		value: value,
		exp:   t,
	}
}

// lruCache is a thread-safe fixed size cache that evicts the least recently used item
type lruCache struct {
	c int
	l *list.List
	m map[string]*list.Element

	mu sync.Mutex
}

// NewLRU returns a new instance of lruCache with the given capacity
func NewLRU(capacity int) *lruCache {
	if capacity <= 0 {
		panic("NewLRU: capacity must be greater than 0")
	}

	return &lruCache{
		c: capacity,
		l: list.New(),
		m: make(map[string]*list.Element),
	}
}

// Get returns the value associated with the key and true if the key exists
// If the key exists, it will be promoted to the front of the list
func (c *lruCache) Get(key string) (value interface{}, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if elem, ok := c.m[key]; ok {
		// TODO: Check if the key has expired(ttl) and delete it if it has
		c.l.MoveToFront(elem)
		return elem.Value.(*lruObj).value, true
	}

	return nil, false
}

// Set adds a key-value pair to the cache
// If the key already exists, its value will be updated and the key will be promoted to the front of the list
// If the key doesn't exist, it will be added to the cache and the least recently used key will be evicted if the cache is full
func (c *lruCache) Set(key string, value interface{}, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if elem, ok := c.m[key]; ok {
		elem.Value.(*lruObj).value = value
		c.l.MoveToFront(elem)
	} else {
		// remove the least recently used key if the cache is full
		if c.l.Len() == c.c {
			elem := c.l.Back()
			e := elem.Value.(*lruObj)
			c.delete(e.key)
		}

		// add the new key to the front of the list
		elem := c.l.PushFront(newLRUObj(key, value, ttl))
		c.m[key] = elem
	}
}

// Delete removes the key from the cache if it exists
func (c *lruCache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.delete(key)
}

// delete removes the key from the cache if it exists
func (c *lruCache) delete(key string) {
	if elem, ok := c.m[key]; ok {
		delete(c.m, key)
		c.l.Remove(elem)
	}
}
