package cache

import (
	"container/list"
	"sync"
	"time"
)

type (
	lru struct {
		mut          sync.RWMutex
		valuesKeyMap map[any]*list.Element
		valueList    *list.List
		maxSize      int32
		curSize      int32
		ttl          time.Duration
	}

	iterator struct {
		lru       *lru
		nextItem  *list.Element
		createdAt time.Time
	}

	entry struct {
		key       any
		value     any
		createdAt time.Time
	}
)

func NewEntry(key any, value any, createdAt time.Time) *entry {
	return &entry{
		key:       key,
		value:     value,
		createdAt: createdAt,
	}
}

func (e *entry) Key() any {
	return e.key
}

func (e *entry) Value() any {
	return e.value
}

func (e *entry) CreatedAt() time.Time {
	return e.createdAt
}

func newIterator(lru *lru, createdAt time.Time) *iterator {
	return &iterator{
		lru:       lru,
		createdAt: createdAt,
	}
}

func (itr *iterator) HasNext() bool {
	return itr.nextItem != nil
}

func (itr *iterator) Next() *list.Element {
	next := itr.nextItem
	itr.prepareNext()
	return next
}

func (itr *iterator) Close() {
	itr.lru.mut.Unlock()
}

func (itr *iterator) prepareNext() {
	itr.nextItem = itr.lru.valueList.Front().Next()
}

func NewLRU(maxSize int32) *lru {
	return &lru{
		valuesKeyMap: make(map[any]*list.Element),
		valueList:    list.New(),
		maxSize:      maxSize,
		curSize:      0,
	}
}

func (c *lru) Get(key any) *entry {
	if c.maxSize == 0 {
		return nil
	}
	c.mut.Lock()
	defer c.mut.Unlock()

	element, ok := c.valuesKeyMap[key]
	if !ok {
		return nil
	}
	val := element.Value.(entry)
	return &val
}

func (c *lru) Put(key any, value entry) {
	c.mut.Lock()
	defer c.mut.Unlock()
	if c.curSize == c.maxSize {
		//Evict least used
		lastEl := c.valueList.Back()
		val := lastEl.Value.(entry)
		delete(c.valuesKeyMap, val.key)
		c.valueList.Remove(lastEl)
		c.curSize--
	}
	element := c.valueList.PushFront(value)
	c.valuesKeyMap[key] = element
	c.curSize++
}

func (c *lru) Delete(key any) {
	c.mut.Lock()
	defer c.mut.Unlock()

	element, ok := c.valuesKeyMap[key]
	if !ok {
		return
	}
	delete(c.valuesKeyMap, key)
	c.valueList.Remove(element)
	c.curSize--
}

// Creates iterator for LRU. Make sure to defer call iterator.Close()
func (c *lru) Iterator() *iterator {
	c.mut.Lock()
	itr := newIterator(c, time.Now())
	itr.prepareNext()
	return itr
}

func (c *lru) Size() int32 {
	return c.curSize
}
