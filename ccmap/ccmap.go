package ccmap

import (
	"fmt"
	"sync"
	"encoding/json"
)

type Cache[K comparable, V any] struct {
	Map map[K]V //cached items
	OnEvicted func(K, V) //function that's called when cached item is deleted automatically

	sync.RWMutex //mutex
}

func New[K comparable, V any]() *Cache[K, V] {
	return &Cache[K, V] {
		Map: make(map[K]V, 0),
	}
}

func (c *Cache[K, V]) SetOnEvicted(f func(K, V)) {
	c.OnEvicted = f
}

func (c *Cache[K, V]) Get(key K) (data V) {
	c.RLock()

	data,_ = c.Map[key]

	c.RUnlock()
	return
}

func (c *Cache[K, V]) GetHas(key K) (data V, ok bool) {
	c.RLock()

	data, ok = c.Map[key]

	c.RUnlock()
	return
}

func (c *Cache[K, V]) Has(key K) (ok bool) {
	c.RLock()

	_, ok = c.Map[key];

	c.RUnlock()
	return
}

func (c *Cache[K, V]) Set(key K, val V) {
	c.Lock()

	c.Map[key] = val

	c.Unlock()
}

func (c *Cache[K, V]) Add(key K, val V) error {
	if c.Has(key) {
		return fmt.Errorf("ccmap: Data already exists with given key %T", key)
	}
	
	c.Set(key, val)

	return nil
}

func (c *Cache[K, V]) Update(key K, val V) error {
	if !c.Has(key) {
		return fmt.Errorf("ccmap: Data doesn't exists with given key %T", key)
	}

	c.Set(key, val)

	return nil
}

func (c *Cache[K, V]) Delete(key K) {
	if c.Has(key) {
		delete(c.Map, key)
	}
}

func (c *Cache[K, V]) Flush() {
	for k,v := range c.Map {
		c.OnEvicted(k, v)
		delete(c.Map, k)
	}
}

func (c *Cache[K, V]) LoadFromJSON(b []byte) (err error) {
	c.Lock()

	err = json.Unmarshal(b, &c.Map)

	c.Unlock()
	return
}