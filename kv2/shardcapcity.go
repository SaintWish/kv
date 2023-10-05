package kv2

import (
	"sync"

	"github.com/saintwish/kv/swiss"
	"github.com/saintwish/kv/stack"
)

type item[V any] struct {
	Object V
	Index int
}

//used internally
type shardCapacity[K comparable, V any] struct {
	Map *swiss.Map[K, item[V]]
	Stack *stack.Stack[K]
	sync.RWMutex //mutex
}

func newShardCapacity[K comparable, V any](size uint64, count uint64) *shardCapacity[K, V] {
	return &shardCapacity[K, V] {
		Map: swiss.NewMap[K, item[V]]( uint32(size/count) ),
		Stack: stack.New[K](int(size/count)),
	}
}

func (m *shardCapacity[K, V]) has(key K) bool {
	m.RLock()

	ok := m.Map.Has(key)

	m.RUnlock()

	return ok
}

func (m *shardCapacity[K, V]) getHasRenew(key K) (V, bool) {
	m.Lock()

	val, ok := m.Map.GetHas(key)
	val.Index = m.Stack.MoveToBack(val.Index)
	m.Map.Set(key, val)

	m.Unlock()

	return val.Object, ok
}

func (m *shardCapacity[K, V]) getHas(key K) (V, bool) {
	m.RLock()

	val, ok := m.Map.GetHas(key)

	m.RUnlock()

	return val.Object, ok
}

func (m *shardCapacity[K, V]) getRenew(key K) V {
	m.Lock()

	val := m.Map.Get(key)
	val.Index = m.Stack.MoveToBack(val.Index)
	m.Map.Set(key, val)

	m.Unlock()

	return val.Object
}

func (m *shardCapacity[K, V]) get(key K) V {
	m.RLock()

	val := m.Map.Get(key)

	m.RUnlock()

	return val.Object
}

/*--------
	Other functions
----------*/
func (m *shardCapacity[K, V]) set(key K, val V, callback func(K, V)) {
	itm := item[V]{
		Object: val,
	}

	m.Lock()

	
	if m.Map.Capacity() > 0 {
		itm.Index = m.Stack.Push(key)
		m.Map.Set(key, itm)
	}

	if m.Map.Capacity() == 0 {
		_, oldKey := m.Stack.Pop()
		_,v := m.Map.Delete(oldKey)
		callback(oldKey, v.Object)

		itm.Index = m.Stack.Push(key)
		m.Map.Set(key, itm)
	}

	m.Unlock()
}

func (m *shardCapacity[K, V]) update(key K, val V) {
	m.Lock()

	if v, ok := m.Map.GetHas(key); ok {
		v.Object = val
		v.Index = m.Stack.MoveToBack(v.Index)
		m.Map.Set(key, v)
	}

	m.Unlock()
}

func (m *shardCapacity[K, V]) delete(key K) bool {
	m.Lock()

	ok, val := m.Map.Delete(key)
	m.Stack.Remove(val.Index)

	m.Unlock()

	return ok
}

func (m *shardCapacity[K, V]) deleteCallback(key K, callback func(K, V)) bool {
	m.Lock()

	ok, val := m.Map.Delete(key)
	m.Stack.Remove(val.Index)
	callback(key, val.Object)

	m.Unlock()

	return ok
}

func (m *shardCapacity[K, V]) clear() {
	m.Map.Clear()
}

func (m *shardCapacity[K, V]) flush(callback func(K, V)) {
	m.Lock()

	m.Stack.Clear()
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