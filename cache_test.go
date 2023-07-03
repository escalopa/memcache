package memcache

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name         string
		nodes        int
		capacity     int
		newCacheImpl func(int) Cache
		expPanic     bool
	}{
		{
			name:         "nodes is zero",
			nodes:        0,
			capacity:     1,
			newCacheImpl: WithLRU,
			expPanic:     true,
		},
		{
			name:         "capacity is zero",
			nodes:        1,
			capacity:     0,
			newCacheImpl: WithLRU,
			expPanic:     true,
		},
		{
			name:         "success",
			nodes:        1,
			capacity:     1,
			newCacheImpl: WithLRU,
			expPanic:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.expPanic {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("expected panic")
					}
				}()
			}

			New(tt.nodes, tt.capacity, tt.newCacheImpl)
		})
	}
}

func TestHashKey(t *testing.T) {
	tests := []struct {
		name  string
		nodes int

		key string
	}{
		{
			name:  "success",
			nodes: 10,
			key:   "foo",
		},
		{
			name:  "success",
			nodes: 100,
			key:   "bar",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := New(tt.nodes, 10, WithLRU)
			got := mc.hashKey(tt.key)
			if got >= tt.nodes {
				t.Errorf("expected less than %d, got %d", tt.nodes, got)
			}
		})
	}
}

func TestMemCacheGet(t *testing.T) {
	tests := []struct {
		name         string
		nodes        int
		capacity     int
		newCacheImpl func(int) Cache

		key string
	}{
		{
			name:         "success",
			nodes:        1,
			capacity:     1,
			newCacheImpl: WithLRU,

			key: "foo",
		},
		{
			name:         "success",
			nodes:        1,
			capacity:     1,
			newCacheImpl: WithLFU,

			key: "foo",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := New(tt.nodes, tt.capacity, tt.newCacheImpl)
			mc.Set(tt.key, "bar", 0)
			got, ok := mc.Get(tt.key)
			require.True(t, ok)
			require.Equal(t, "bar", got)
		})
	}
}

func TestMemCacheSet(t *testing.T) {
	tests := []struct {
		name         string
		nodes        int
		capacity     int
		newCacheImpl func(int) Cache

		key   string
		value interface{}
	}{
		{
			name:         "success",
			nodes:        1,
			capacity:     1,
			newCacheImpl: WithLRU,

			key:   "foo",
			value: "bar",
		},
		{
			name:         "success",
			nodes:        1,
			capacity:     1,
			newCacheImpl: WithLFU,

			key:   "foo",
			value: "bar",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := New(tt.nodes, tt.capacity, tt.newCacheImpl)
			mc.Set(tt.key, tt.value, 0)

			got, ok := mc.Get(tt.key)
			require.True(t, ok)
			require.Equal(t, tt.value, got)
		})
	}
}

func TestMemCacheDelete(t *testing.T) {
	tests := []struct {
		name         string
		nodes        int
		capacity     int
		newCacheImpl func(int) Cache

		key string
	}{
		{
			name:         "success LRU",
			nodes:        1,
			capacity:     1,
			newCacheImpl: WithLRU,

			key: "foo",
		},
		{
			name:         "success LFU",
			nodes:        1,
			capacity:     1,
			newCacheImpl: WithLFU,

			key: "foo",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := New(tt.nodes, tt.capacity, tt.newCacheImpl)
			mc.Set(tt.key, "bar", 0)

			got, ok := mc.Get(tt.key)
			require.True(t, ok)
			require.Equal(t, "bar", got)

			mc.Delete(tt.key)
			got, ok = mc.Get(tt.key)
			require.False(t, ok)
			require.Nil(t, got)
		})
	}
}

func shouldEvict(t *testing.T, mc Cache, key string) {
	if _, ok := mc.Get(key); ok {
		t.Errorf("%s should be evicted", key)
	}
}
