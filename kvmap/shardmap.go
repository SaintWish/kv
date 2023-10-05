package kvmap

import (
	"sync"
)

//used internally
type shardMap[K comparable, V any] struct {
	Map map[K]V
	sync.RWMutex //mutex
}

func newShardMap[K comparable, V any]() *shardMap[K, V] {
	return &shardMap[K, V] {
		Map: make(map[K]V, 0),
	}
}

func (m *shardMap[K, V]) set(key K, val V) {
	m.Lock()

	m.Map[key] = val

	m.Unlock()
}

func (m *shardMap[K, V]) get(key K) (val V) {
	m.RLock()

	val = m.Map[key]

	m.RUnlock()

	return
}