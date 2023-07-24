package cache

import (
	"fmt"
	"log"
	"sync"
	"time"
)

// Cache: object that holds the hashmap
// to store key value pairs items to be cached
type Cache struct {
	lock sync.RWMutex
	// In this code, `map` is used to create a hashmap data structure. A map is an unordered
	// collection of key-value pairs, where each key is unique. It allows efficient lookup,
	// insertion, and deletion of elements based on the key. In this case, the `data` field of the
	// `Cache` struct is a map that stores key-value pairs of type `string` and `[]byte`.
	data map[string][]byte
} 

func New() *Cache {
	return &Cache {
		data: make(map[string][]byte),
	}
}

func (c* Cache) Has(key []byte) bool {
	_, ok := c.data[string(key)]
	return ok
}

func (c* Cache) Get(key []byte) ([]byte, error) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	var keyStr string = string(key)
	val, ok := c.data[keyStr]
	if ok {
		return val, nil
	}
	return nil, fmt.Errorf("key (%s) doesn't exist in cache", keyStr)
}

func (c* Cache) Set(key, value []byte, ttl time.Duration) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.data[string(key)] = value
	fmt.Printf("SET %s to %s\n", string(key), string(value))
	// This portion activates as the TTL to delete the entries after a certain period of time
	// alternatively we can use tinker = time.NewTicker(ttl)
	// within the go func method we call <-tinker.C
	go func() {
		<-time.After(ttl)
		delete(c.data, string(key))
	}()
	return nil 
}

func (c* Cache) Delete(key []byte) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	delete(c.data, string(key))
	return nil
}

