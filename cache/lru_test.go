package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLRU(t *testing.T) {
	lru := NewLRU(100)

	key := "tom"
	entry := NewEntry(key, 12345, time.Now())
	lru.Put(key, *entry)
	expected := lru.Get(key)
	assert.Equal(t, *expected, *entry)
	assert.EqualValues(t, 1, lru.Size())

	lru.Delete(key)
	expected = lru.Get(key)
	assert.EqualValues(t, 0, lru.Size())

}

func TestLRUMaxSize(t *testing.T) {
	var (
		maxSize = int32(5)
		lru     = NewLRU(maxSize)
		values  = make(map[string]string)
	)
	values["key1"] = "value1"
	values["key2"] = "value2"
	values["key3"] = "value3"
	values["key4"] = "value4"
	values["key5"] = "value5"
	values["key6"] = "value6"
	values["key7"] = "value7"

	for k, v := range values {
		entry := NewEntry(k, v, time.Now())
		lru.Put(k, *entry)
	}
	assert.Equal(t, maxSize, lru.Size())

}
