package cache

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLRU(t *testing.T) {
	lru := NewLRU(100)

	key := "tom"
	value := 123
	lru.Put(key, value)
	expected := lru.Get(key).(int)
	assert.Equal(t, expected, value)
	assert.EqualValues(t, 1, lru.Size())

	lru.Delete(key)
	assert.EqualValues(t, 0, lru.Size())

}

func TestLRUMaxSize(t *testing.T) {
	var (
		maxSize = int32(5)
		lru     = NewLRU(maxSize)
		values  = map[string]string{
			"key1": "value1",
			"key2": "value2",
			"key3": "value3",
			"key4": "value4",
			"key5": "value5",
			"key6": "value6",
			"key7": "value7",
		}
	)

	for k, v := range values {
		lru.Put(k, v)
	}
	assert.Equal(t, maxSize, lru.Size())

}
func TestLRUIterator(t *testing.T) {
	var (
		maxSize = int32(100)
		lru     = NewLRU(maxSize)
		values  = map[string]string{
			"key1": "value1",
			"key2": "value2",
			"key3": "value3",
			"key4": "value4",
			"key5": "value5",
			"key6": "value6",
			"key7": "value7",
		}
	)

	var entries []string
	for k, v := range values {
		entries = append(entries, v)
		lru.Put(k, v)
	}
	itr := lru.Iterator()
	defer itr.Close()
	cur := itr.Next().Value.(entry)
	assert.Equal(t, entries[6], cur.value)

	cur = itr.Next().Value.(entry)
	assert.Equal(t, entries[5], cur.value)

	cur = itr.Next().Value.(entry)
	assert.Equal(t, entries[4], cur.value)

	cur = itr.Next().Value.(entry)
	assert.Equal(t, entries[3], cur.value)

	cur = itr.Next().Value.(entry)
	assert.Equal(t, entries[2], cur.value)

	cur = itr.Next().Value.(entry)
	assert.Equal(t, entries[1], cur.value)

	cur = itr.Next().Value.(entry)
	assert.Equal(t, entries[0], cur.value)

	assert.Panics(t, func() { itr.Next() })

}
