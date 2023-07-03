# memcache 📚

[![Go Report Card](https://goreportcard.com/badge/github.com/escalopa/memcache)](https://goreportcard.com/report/github.com/escalopa/memcache) [![codecov](https://codecov.io/gh/escalopa/memcache/branch/main/graph/badge.svg?token=GYCQFM7WUM)](https://codecov.io/gh/escalopa/memcache) [![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fescalopa%2Fmemcache.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2Fescalopa%2Fmemcache?ref=badge_shield)

[![DeepSource](https://app.deepsource.com/gh/escalopa/memcache.svg/?label=active+issues&show_trend=true&token=le3CGl9jnv3HKOckiiT5r1pE)](https://app.deepsource.com/gh/escalopa/memcache/?ref=repository-badge)

DHT cache, Built with consitent hashing &amp; LRU/LFU cache

## Example

### Memcache

```go
package main

import (
	"fmt"

	"github.com/escalopa/memcache"
)

func main() {
	nodes := 1_000
	capacity := 1_000_000
	mc := memcache.New(nodes, capacity, memcache.NewLRU) // Or use `memcache.NewLFU`

	// If no nodes are needed, use:
	// mc := meme.NewLRU(capacity)
	// Or
	// mc := meme.NewLFU(capacity)

	var value interface{}
	var ok bool

	value, ok = mc.Get("foo")
	fmt.Println(value, ok)
	// Output: <nil> false

	mc.Set("foo", "bar", 0)
	value, ok = mc.Get("foo")
	fmt.Println(value, ok)
	// Output: bar true

	mc.Delete("foo")
	value, ok = mc.Get("foo")
	fmt.Println(value, ok)
	// Output: <nil> false

	mc.Set("foo", "bar", 0)
	value, ok = mc.Get("foo")
	fmt.Println(value, ok)
	// Output: bar true
}

```

### LFU

```go
package main

import (
  "fmt"

  "github.com/escalopa/memcache"
)


## About 

### DHT

Distributed Hash Table (DHT) is a distributed system that provides a lookup service similar to a hash table. It assigns keys to nodes in the network using a hash function. This allows the nodes to efficiently retrieve the value associated with a given key.

**Notice: This project implements DHT on a single node. It is not distributed.**

### LRU

Least Recently Used (LRU) is a common caching strategy. It defines that the least recently used items are discarded first.

### LFU

Least Frequently Used (LFU) is a caching strategy whereby the least frequently used items are discarded first.

## References

- [Consistent Hashing](https://en.wikipedia.org/wiki/Consistent_hashing)
- [Least Recently Used (LRU)](https://en.wikipedia.org/wiki/Cache_replacement_policies#Least_recently_used_(LRU))
- [Least Frequently Used (LFU)](https://en.wikipedia.org/wiki/Cache_replacement_policies#Least-frequently_used_(LFU))