package main

import (
	"fmt"

	"github.com/escalopa/memcache"
)

func main() {
	nodes := 1_000
	capacity := 1_000_000
	mc := memcache.New(nodes, capacity, memcache.NewLRU) // Or use `memcache.NewLFU`

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
