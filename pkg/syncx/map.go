package syncx

import (
	"maps"
	"sync"
)

// Map is a map that is safe for concurrent use (goroutines).
type Map[K comparable, V any] struct {
	sync.RWMutex
	d map[K]V
}

// NewMap returns a new empty [Map].
func NewMap[K comparable, V any]() *Map[K, V] {
	var m Map[K, V]
	m.d = map[K]V{}
	return &m
}

func (m *Map[K, V]) ShallowCopy() map[K]V {
	m.RLock()
	defer m.RUnlock()
	return maps.Clone(m.d)
}

// WalkWithLock iterates all entries with exclusive lock.
// WalkWithLock can read and update all entries.
func (m *Map[K, V]) WalkWithLock(f func(key K, value V) (newValue V, deleteThis bool)) {
	m.Lock()
	defer m.Unlock()

	for k, v := range m.d {
		newValue, deleteThis := f(k, v)
		if deleteThis {
			delete(m.d, k)
		} else {
			m.d[k] = newValue
		}
	}
}

// Tx reads a record and update it.
//
// `f` executes in a single transaction, where its key matches the key in `Tx` argument.
// The corresponding value is used, defaulting if absent, and `exist` indicates the presence of the value.
//
// If `f` returns `deleteThis` as true, the key is removed;
// if false, the key's value is updated to `newValue`.
// The return value `retValue` of `Tx` equals `newValue`,
// and `exist` matches `f`'s `exist` argument.
func (m *Map[K, V]) Tx(
	key K,
	f func(key K, value V, exist bool) (newValue V, deleteThis bool),
) (retValue V, exist bool) {
	m.Lock()
	defer m.Unlock()

	value, exist := m.d[key]
	newValue, deleteThis := f(key, value, exist)

	if deleteThis {
		delete(m.d, key)
	} else {
		m.d[key] = newValue
	}
	return newValue, exist
}

func (m *Map[K, V]) Del(key K) (V, bool) {
	return m.Tx(key, func(_ K, v V, _ bool) (V, bool) {
		return v, true
	})
}

func (m *Map[K, V]) Set(key K, value V) (V, bool) {
	return m.Tx(key, func(_ K, v V, _ bool) (V, bool) {
		return value, false
	})
}

func (m *Map[K, V]) Get(key K) (V, bool) {
	m.RLock()
	defer m.RUnlock()
	v, ok := m.d[key]
	return v, ok
}
