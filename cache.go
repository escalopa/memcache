package memcache

import (
	"hash/fnv"
	"time"
)

type Cache interface {
	Get(key string) (value interface{}, ok bool)
	Set(key string, value interface{}, ttl time.Duration)
	Delete(key string)
}

type memCache struct {
	nodes  int     // Count of sub caches
	caches []Cache // Sharded caches by key
}

// New returns a new instance of memeCache
// nodes: Number of sub caches
// capacity: Max number of key-value pairs in each sub cache before eviction
// newCacheImpl: Factory function to create a new instance of Cache, e.g. `NewLFU` or `NewLRU`
func New(nodes int, capacity int, newCacheImpl func(int) Cache) *memCache {
	if nodes <= 0 {
		panic("nodes must be greater than zero")
	}

	if capacity <= 0 {
		panic("capacity must be greater than zero")
	}

	caches := make([]Cache, nodes)
	for i := 0; i < nodes; i++ {
		caches[i] = newCacheImpl(capacity)
	}

	return &memCache{
		caches: caches,
		nodes:  nodes,
	}
}

// Get returns the value for the given key
func (mc *memCache) Get(key string) (value interface{}, ok bool) {
	id := mc.hashKey(key)
	return mc.caches[id].Get(key)
}

// Set sets the value for the given key
func (mc *memCache) Set(key string, value interface{}, ttl time.Duration) {
	id := mc.hashKey(key)
	mc.caches[id].Set(key, value, ttl)
}

// Delete deletes the value for the given key
func (mc *memCache) Delete(key string) {
	id := mc.hashKey(key)
	mc.caches[id].Delete(key)
}

// hashKey Calculates the cache ID where the key-value pair is stored
func (mc *memCache) hashKey(s string) int {
	h := fnv.New32()
	_, _ = h.Write([]byte(s))
	p := h.Sum32() % uint32(mc.nodes)
	return int(p)
}
