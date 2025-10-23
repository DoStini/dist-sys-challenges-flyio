package atomicmap

import "sync"

type AtomicMap[K comparable, T any] struct {
	internal map[K]T
	mutex    sync.Mutex
}

func NewAtomicMap[K comparable, T any]() *AtomicMap[K, T] {
	return &AtomicMap[K, T]{
		internal: map[K]T{},
		mutex:    sync.Mutex{},
	}
}

func (am *AtomicMap[K, T]) Delete(key K) {
	defer am.mutex.Unlock()
	am.mutex.Lock()

	delete(am.internal, key)
}

func (am *AtomicMap[K, T]) Get(key K) (T, bool) {
	defer am.mutex.Unlock()
	am.mutex.Lock()

	v, found := am.internal[key]

	return v, found
}

func (am *AtomicMap[K, T]) Set(key K, val T) {
	defer am.mutex.Unlock()
	am.mutex.Lock()

	am.internal[key] = val
}

func (am *AtomicMap[K, T]) GetOrDefault(key K, in T) (T, bool) {
	defer am.mutex.Unlock()
	am.mutex.Lock()

	out, found := am.internal[key]
	if !found {
		am.internal[key] = in
		return in, false
	}

	return out, true
}
