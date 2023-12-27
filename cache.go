package cache

import "time"

type (
	//Iterator to traverse through elements in the Doubly linkedList
	Iterator interface {
		// Indicate if iterator has next item
		HasNext() bool
		//Next item in the iterator
		Next() any

		//Close iterator
		Close()
	}

	// Entry of a cache
	Entry interface {
		//Get key of the entry
		Key() any

		//Value of the entry
		Value() any

		//CreatedAt
		CreatedAt() time.Time
	}

	//General Cache interface
	Cache interface {
		Get(key any) any
		Put(key any, value any)
		Delete(key any)
		Size() int32
		Iterator() *Iterator
	}
)
