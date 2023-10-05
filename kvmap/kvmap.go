package kvmap

import (
	"fmt"
	
	"github.com/dolthub/maphash"
)

type Cache[K comparable, V any] struct {
	shards []*shardMap[K, V]
	shardCount uint64
	hash maphash.Hasher[K]

	OnEvicted func(K, V) //function that's called when cached item is deleted by the system
}

func New[K comparable, V any](sc uint64) *Cache[K, V] {
	cache := Cache[K, V] {}
	cache.shards = make([]*shardMap[K, V], sc)
	cache.hash = maphash.NewHasher[K]()
	cache.shardCount = sc

	for i := 0; i < int(sc); i++ {
		cache.shards[i] = newShardMap[K, V]()
	}

	return &cache
}

func (c *Cache[K, V]) getShardIndex(key K) uint64 {
	sum := c.hash.Hash(key)

	return sum % c.shardCount
}

func (c *Cache[K, V]) getShard(key K) *shardMap[K, V] {
	sum := c.hash.Hash(key)
	fmt.Println(sum)
	return c.shards[sum%c.shardCount]
}

func (c *Cache[K, V]) SetOnEvicted(f func(K, V)) {
	c.OnEvicted = f
}

func (c *Cache[K, V]) Set(key K, val V) {
	shard := c.getShard(key)
	shard.set(key, val)
}

func (c *Cache[K, V]) Get(key K) V {
	shard := c.getShard(key)
	return shard.get(key)
}

func (c *Cache[K, V]) Has(key K) bool {
	shard := c.getShard(key)
	shard.RLock()
	defer shard.RUnlock()

	if _, ok := shard.Map[key]; ok {
		return true
	}

	return false
}

func (c *Cache[K, V]) Add(key K, val V) (err error) {
	shard := c.getShard(key)
	shard.Lock()

	if _, ok := shard.Map[key]; ok {
		err = fmt.Errorf("kvmap: Data already exists with given key %T", key)
		shard.Unlock()
		return
	}
	
	c.Set(key, val)

	shard.Unlock()
	return nil
}

func (c *Cache[K, V]) Update(key K, val V) (err error) {
	shard := c.getShard(key)
	shard.Lock()

	if _, ok := shard.Map[key]; !ok {
		return fmt.Errorf("kvmap: Data doesn't exists with given key %T", key)
		shard.Unlock()
		return
	}

	c.Set(key, val)

	shard.Unlock()
	return nil
}

func (c *Cache[K, V]) Delete(key K) {
	if c.Has(key) {
		shard := c.getShard(key)
		delete(shard.Map, key)
	}
}

func (c *Cache[K, V]) Flush() {
	for i := 0; i < len(c.shards); i++ {
		shard := c.shards[i]

		for k,v := range shard.Map {
			c.OnEvicted(k, v)
			delete(shard.Map, k)
		}
	}
}