package kv1s

// A key, value resizeable cache using a modified swiss map.

import (
	"sync"

	"github.com/saintwish/kv/swiss"
)

//used internally
type shard[K comparable, V any] struct {
	Map *swiss.Map[K, V]
	sync.RWMutex //mutex
}

func newShard[K comparable, V any](size uint64, count uint64) *shard[K, V] {
	return &shard[K, V] {
		Map: swiss.NewMap[K, V]( uint32(size/count) ),
	}
}

func (m *shard[K, V]) has(key K) bool {
	m.RLock()

	ok := m.Map.Has(key)

	m.RUnlock()

	return ok
}

func (m *shard[K, V]) getHas(key K) (V, bool) {
	m.RLock()

	val, ok := m.Map.GetHas(key)

	m.RUnlock()

	return val, ok
}

func (m *shard[K, V]) getHasRenew(key K) (V, bool) {
	return m.getHas(key)
}

func (m *shard[K, V]) get(key K) V {
	m.RLock()

	val := m.Map.Get(key)

	m.RUnlock()

	return val
}

func (m *shard[K, V]) getRenew(key K) V {
	return m.get(key)
}

/*--------
	Other functions
----------*/
func (m *shard[K, V]) set(key K, val V, callback func(K, V)) {
	m.Lock()

	m.Map.Set(key, val)

	m.Unlock()
}

func (m *shard[K, V]) update(key K, val V) {
	m.Lock()

	if val, ok := m.Map.GetHas(key); ok {
		m.Map.Set(key, val)
	}

	m.Unlock()
}

func (m *shard[K, V]) delete(key K) bool {
	m.Lock()

	ok,_ := m.Map.Delete(key)

	m.Unlock()

	return ok
}

func (m *shard[K, V]) deleteCallback(key K, callback func(K, V)) bool {
	m.Lock()

	ok, val := m.Map.Delete(key)
	callback(key, val)

	m.Unlock()

	return ok
}

func (m *shard[K, V]) clear() {
	m.Map.Clear()
}

func (m *shard[K, V]) flush(callback func(K, V)) {
	m.Lock()

	m.Map.Iter(func(key K, val V) (stop bool) {
		callback(key, val)
		m.Map.Delete(key)
		
		if stop {
			return
		}

		return
	})

	m.Unlock()
}