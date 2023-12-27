package kv1

import (
	"fmt"
	"time"
	
	"github.com/dolthub/maphash"
)

type Cache[K comparable, V any] struct {
	shards []*shard[K, V]
	shardCount uint64
	hash maphash.Hasher[K]

	OnEvicted func(K, V) //function that's called when cached item is deleted by the system
}

func New[K comparable, V any](ex time.Duration, sz uint64, sc uint64) *Cache[K, V] {
	if sc > sz {
		panic("kv1: shard count must be smaller than cache size!")
	}

	cache := Cache[K, V] {}
	cache.shards = make([]*shard[K, V], sc)
	cache.hash = maphash.NewHasher[K]()
	cache.shardCount = sc

	for i := 0; i < int(sc); i++ {
		cache.shards[i] = newShard[K, V](ex, sz, sc)
	}

	return &cache
}

func (c *Cache[K, V]) getShardIndex(key K) uint64 {
	sum := c.hash.Hash(key)

	return sum % c.shardCount
}

func (c *Cache[K, V]) getShard(key K) *shard[K,V] {
	sum := c.hash.Hash(key)
	return c.shards[sum%c.shardCount]
}

func (c *Cache[K, V]) SetOnEvicted(f func(K, V)) {
	c.OnEvicted = f
}

func (c *Cache[K, V]) Get(key K) V {
	shard := c.getShard(key)
	return shard.get(key)
}

func (c *Cache[K, V]) GetRenew(key K) V {
	shard := c.getShard(key)
	return shard.getRenew(key)
}

func (c *Cache[K, V]) GetHas(key K) (V, bool) {
	shard := c.getShard(key)
	return shard.getHas(key)
}

func (c *Cache[K, V]) GetHasRenew(key K) (V, bool) {
	shard := c.getShard(key)
	return shard.getHasRenew(key)
}

func (c *Cache[K, V]) Has(key K) bool {
	shard := c.getShard(key)
	return shard.has(key)
}

// Sets the key with value, will overwrite if key exists
func (c *Cache[K, V]) Set(key K, val V) {
	shard := c.getShard(key)
	shard.set(key, val)
}

// Adds key with value to map, will error if key already exists.
func (c *Cache[K, V]) Add(key K, val V) error {
	shard := c.getShard(key)
	if shard.has(key) {
		return fmt.Errorf("kv1: Data already exists with given key %T", key)
	}

	shard.set(key, val)
	return nil
}

// Updates given key, errors if key doesn't already exists.
func (c *Cache[K, V]) Update(key K, val V) error {
	shard := c.getShard(key)
	if !shard.has(key) {
		return fmt.Errorf("kv1: Data doesn't exists with given key %T", key)
	}

	shard.update(key, val)
	return nil
}

// Will Set or Update said key depending if exists or not.
func (c *Cache[K, V]) SetOrUpdate(key K, val V) {
	shard := c.getShard(key)
	if shard.has(key) {
		shard.update(key, val)
	}else{
		shard.set(key, val)
	}
}

func (c *Cache[K, V]) Renew(key K) {
	shard := c.getShard(key)
	shard.renew(key)
}

func (c *Cache[K, V]) Delete(key K) bool {
	shard := c.getShard(key)
	return shard.delete(key)
}

func (c *Cache[K, V]) ShardCount() uint64 {
	return c.shardCount
}

func (c *Cache[K, V]) GetShardSize(key K) int {
	shard := c.getShard(key)
	return shard.Map.Count()
}

func (c *Cache[K, V]) GetShardMaxSize(key K) int {
	shard := c.getShard(key)
	return shard.Map.MaxCapacity()
}

func (c *Cache[K, V]) GetShardCapacity(key K) int {
	shard := c.getShard(key)
	return shard.Map.Capacity()
}

// Gets the current amount of elements in the cache.
func (c *Cache[K, V]) Count() (count int) {
	for i := 0; i < len(c.shards); i++ {
		shard := c.shards[i]
		count = count + shard.Map.Count()
	}
	return
}

// Clears the cache with calling OnEviction callback
func (c *Cache[K, V]) Flush() {
	for i := 0; i < len(c.shards); i++ {
		shard := c.shards[i]
		shard.flush(c.OnEvicted)
	}
}

// Clears the cache without calling OnEviction callback
func (c *Cache[K, V]) Clear() {
	for i := 0; i < len(c.shards); i++ {
		shard := c.shards[i]
		shard.clear()
	}
}

func (c *Cache[K, V]) ForEach(f func(key K, val item[V])) {
	for i := 0; i < len(c.shards); i++ {
		shard := c.shards[i]
		shard.Lock()
		defer shard.Unlock()

		shard.Map.Iter(func(key K, val item[V]) (stop bool) {
			f(key, val)

			if stop {
				return
			}

			return
		})
	}
}

func (c *Cache[K,V]) DeleteExpired() {
	for i := 0; i < len(c.shards); i++ {
		shard := c.shards[i]
		shard.evictExpired(c.OnEvicted)
	}
}

func (c *Cache[K,V]) IsExpired(key K) bool {
	shard := c.getShard(key)
	return shard.isExpired(key)
}