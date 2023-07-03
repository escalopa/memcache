package memcache

import (
	"strconv"
	"testing"
	"time"
)

func TestNewLFU(t *testing.T) {
	tests := []struct {
		name     string
		capacity int
		expPaic  bool
	}{
		{
			name:     "success",
			capacity: 1,
			expPaic:  false,
		},
		{
			name:     "capacity is zero",
			capacity: 0,
			expPaic:  true,
		},
		{
			name:     "capacity is negative",
			capacity: -1,
			expPaic:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.expPaic {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("expected panic")
					}
				}()
			}
			NewLFU(tt.capacity)
		})
	}
}

func TestLFU(t *testing.T) {
	mc := NewLFU(5)

	for i := 0; i < 6; i++ {
		val := strconv.Itoa(i)
		mc.Set("key"+val, "value"+val, 1*time.Hour)
	}

	shouldEvict(t, mc, "key0")

	// key1 should be promoted to the head of the list
	// key2 should be evicted to make room for `key7`
	mc.Set("key1", "value1", 0)
	mc.Set("key7", "value7", 0)
	shouldEvict(t, mc, "key2")

	// key4 should be evicted because it was the removed on purpose using `Delete`
	mc.Delete("key4")
	shouldEvict(t, mc, "key4")

	// key8 won't evict key3 because because the list is not full yet (Since key4 was deleted)
	// key9 will evict key3 because the list is full again
	mc.Set("key8", "value8", 0)
	mc.Set("key9", "value9", 0)
	shouldEvict(t, mc, "key3")

	// key10 will evict key5 because it was the least recently used
	mc.Set("key10", "value10", 0)
	shouldEvict(t, mc, "key5")
}
