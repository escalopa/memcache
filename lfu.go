package memcache

import (
	"container/list"
	"sync"
	"time"
)

type lfuObj struct {
	key   string
	value interface{}
	freq  int

	exp *time.Time
}

// newLFUObj returns a new instance of lfuObj with the given key, value and ttl
func newLFUObj(key string, value interface{}, ttl time.Duration) *lfuObj {
	var t *time.Time

	if ttl > 0 {
		temp := time.Now().Add(ttl)
		t = &temp
	}

	return &lfuObj{
		key:   key,
		value: value,
		freq:  1,
		exp:   t,
	}
}

// lfuCache is a thread-safe fixed size cache that evicts the least frequently used item
// It uses a map to store key-value pairs and a doubly linked list to keep track of the frequency of each key
// The least frequently used key is always at the back of the list and evicted first if the cache is full when a new key is added
type lfuCache struct {
	c int

	freq map[int]*list.List
	elem map[string]*list.Element
	min  int

	mu sync.Mutex
}

// NewLFU returns a new instance of lfuCache with the given capacity
func NewLFU(capacity int) Cache {
	if capacity <= 0 {
		panic("NewLFU: capacity must be greater than 0")
	}

	return &lfuCache{
		c:    capacity,
		freq: make(map[int]*list.List),
		elem: make(map[string]*list.Element),
		min:  1,
	}
}

// Get returns the value associated with the key and true if the key exists
// If the key exists, it will be promoted and its frequency will be increased by 1 otherwise it will return nil, false
func (c *lfuCache) Get(key string) (value interface{}, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if elem, ok := c.elem[key]; ok {
		obj := elem.Value.(*lfuObj)
		// TODO: Check if the key has expired(ttl) and delete it if it has
		c.fix(elem)
		return obj.value, true
	}

	return nil, false
}

// Set adds a key-value pair to the cache
// If the key already exists, its value will be updated and the key will be promoted
// If the key doesn't exist, it will be added to the cache and the least frequently used key will be evicted if the cache is full
func (c *lfuCache) Set(key string, value interface{}, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if elem, ok := c.elem[key]; ok {
		obj := elem.Value.(*lfuObj)
		obj.value = value
		c.fix(elem)
	} else {
		// If the cache is full, the least frequently used key will be evicted
		if len(c.elem) == c.c {
			e := c.freq[c.min].Back().Value.(*lfuObj)
			delete(c.elem, e.key)
			c.freq[c.min].Remove(c.freq[c.min].Back())
		}

		// Add the key to the front of the list with frequency 1
		if c.freq[1] == nil {
			c.freq[1] = list.New()
		}
		c.elem[key] = c.freq[1].PushFront(newLFUObj(key, value, ttl))

		// The minimum frequency will always be 1
		c.min = 1
	}
}

// Delete removes the key-value pair from the cache
func (c *lfuCache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if elem, ok := c.elem[key]; ok {
		obj := elem.Value.(*lfuObj)
		delete(c.elem, key)
		c.freq[obj.freq].Remove(elem)

		// If the key was the only key with the minimum frequency and it was deleted,
		// the minimum frequency will be increased until a key with a higher frequency is found
		for c.min <= obj.freq && c.freq[c.min].Len() == 0 {
			c.min++
		}
	}
}

// fix promotes the given key and increases its frequency by 1
func (c *lfuCache) fix(elem *list.Element) {
	// Remove the key from the list with the `obj.freq` frequency
	obj := elem.Value.(*lfuObj)
	c.freq[obj.freq].Remove(elem)

	// Move the key to the front of the list with the `obj.freq+1` frequency
	obj.freq++
	if c.freq[obj.freq] == nil {
		c.freq[obj.freq] = list.New()
	}
	c.elem[obj.key] = c.freq[obj.freq].PushFront(obj)

	// If the key was the only key with the minimum frequency and it was promoted,
	// the minimum frequency will be increased by 1
	if c.min == obj.freq-1 && c.freq[obj.freq-1].Len() == 0 {
		c.min++
	}
}
