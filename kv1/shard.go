package kv1

import (
	"time"
	"sync"

	"github.com/saintwish/kv/swiss"
)

type item[V any] struct {
	Object V
	Expire time.Time
}

//used internally
type shard[K comparable, V any] struct {
	Map *swiss.Map[K, item[V]]
	Expiration time.Duration
	sync.RWMutex //mutex
}

func newShard[K comparable, V any](ex time.Duration, size uint64, count uint64) *shard[K, V] {
	return &shard[K, V] {
		Map: swiss.NewMap[K, item[V]]( uint32(size/count) ),
		Expiration: ex,
	}
}

/*--------
	Raw get functions.
----------*/
func (m *shard[K, V]) get(key K) (val V) {
	m.RLock()

	val = m.Map.Get(key).Object

	m.RUnlock()

	return
}

func (m *shard[K, V]) getRenew(key K) (val V) {
	m.Lock()

	if ok,v := m.Map.GetHas(key); ok {
		v.Expire = time.Now().Add(m.Expiration)
		m.Map.Set(key, v)
		val = v.Object
	}

	m.Unlock()

	return
}

func (m *shard[K, V]) has(key K) (ok bool) {
	m.RLock()

	ok = m.Map.Has(key)

	m.RUnlock()

	return
}

func (m *shard[K, V]) getHas(key K) (val V, ok bool) {
	m.RLock()

	ok,v := m.Map.GetHas(key);
	val = v.Object

	m.RUnlock()

	return
}

func (m *shard[K, V]) getHasRenew(key K) (val V, ok bool) {
	m.Lock()

	if ok,v := m.Map.GetHas(key); ok {
		v.Expire = time.Now().Add(m.Expiration)
		m.Map.Set(key, v)
		val = v.Object
	}

	m.Unlock()

	return
}


/*--------
	Other functions
----------*/
func (m *shard[K, V]) set(key K, val V) {
	itm := item[V]{
		Object: val,
		Expire: time.Now().Add(m.Expiration),
	}

	m.Lock()

	m.Map.Set(key, itm)

	m.Unlock()
}

func (m *shard[K, V]) update(key K, val V) {
	itm := item[V]{
		Object: val,
		Expire: time.Now().Add(m.Expiration),
	}

	m.Lock()

	if ok := m.Map.Has(key); ok {
		m.Map.Set(key, itm)
	}

	m.Unlock()
}

func (m *shard[K, V]) delete(key K) (ok bool) {
	m.Lock()

	ok, _ = m.Map.Delete(key)

	m.Unlock()

	return
}

func (m *shard[K, V]) isExpired(key K) (ex bool) {
	m.RLock()

	if ok,v:= m.Map.GetHas(key); ok {
		ex = time.Now().After(v.Expire)
	}

	m.RUnlock()
	return
}

//Returns true if item is expired and thus evicted.
func (m *shard[K, V]) evictItem(key K, callback func(K,V)) (ex bool) {
	m.Lock()

	ex = false
	if ok,v := m.Map.GetHas(key); ok {
		if time.Now().After(v.Expire) {
			ex = true
			callback(key, v.Object)
			m.Map.Delete(key)
		}
	}

	m.Unlock()

	return
}

func (m *shard[K, V]) evictExpired(callback func(K,V)) {
	m.Lock()

	m.Map.Iter(func (key K, v item[V]) (stop bool) {
		if time.Now().After(v.Expire) {
			callback(key, v.Object)
			m.Map.Delete(key)
		}
		
		if stop {
			m.Unlock()
			return
		}

		return
	})

	m.Unlock()
}

func (m *shard[K, V]) renew(key K) {
	expire := time.Now().Add(m.Expiration)

	m.Lock()
	
	if ok,v := m.Map.GetHas(key); ok {
		v.Expire = expire
	}

	m.Unlock()
}

func (m *shard[K, V]) clear() {
	m.Map.Clear()
}

func (m *shard[K, V]) flush(callback func(K, V)) {
	m.Lock()

	m.Map.Iter(func(key K, val item[V]) (stop bool) {
		callback(key, val.Object)
		m.Map.Delete(key)
		
		if stop {
			return
		}

		return
	})

	m.Unlock()
}