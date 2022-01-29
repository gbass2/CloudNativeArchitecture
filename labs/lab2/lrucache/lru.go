package lru

type Cacher interface {
	Get(interface{}) (interface{}, error)
	Put(interface{}, interface{}) error
}

type lruCache struct {
	size      int
	remaining int
	cache     map[string]string
	queue     []string
}

func NewCache(size int) Cacher {
	return &lruCache{size: size, remaining: size, cache: make(map[string]string), queue: make([]string, 0)}
}

func (lru *lruCache) Get(key interface{}) (interface{}, error) {
	// Convert key and value to concrete type
	// iterate map O(n) to search for key.
	// if value exist return value
		// add key to queue
	// return error when key, value does not exist

	return "",nil
}

func (lru *lruCache) Put(key, val interface{}) error {
	// Convert key and value to concrete type
	k := key.(string)
	v := val.(string)

	b := false // boolean to remove from the front of the queue if a duplicate is not found
	// If key is in queue remove and add to map again
		// Leave map alone
	for _,k2 := range(lru.queue){
		if k == k2 {
			lru.queue = lru.queue[1:]
			b = true
		}
	}

	// Check if remaining value is 0
	// Remove the head queue if a duplicate has not been removed already
	if lru.remaining == 0 && !b {
		lru.queue = lru.queue[1:]
	}

	// Decrement remaining if not equal to zero
	if lru.remaining > 0 {
		lru.remaining--
	}

	// Add value to hashmap
	lru.cache[k]=v

	// add key to queue
	lru.queue = append(lru.queue, k)

	return nil
}

// Delete element from queue
func (lru *lruCache) qDel(ele string) {
	for i := 0; i < len(lru.queue); i++ {
		if lru.queue[i] == ele {
			oldlen := len(lru.queue)
			copy(lru.queue[i:], lru.queue[i+1:])
			lru.queue = lru.queue[:oldlen-1]
			break
		}
	}
}
