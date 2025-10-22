package main

import "sync"

type AtomicMap[T any] struct {
	internal map[string]T
	mutex    sync.Mutex
}

func NewAtomicMap[T any]() *AtomicMap[T] {
	return &AtomicMap[T]{
		internal: map[string]T{},
		mutex:    sync.Mutex{},
	}
}

func (am *AtomicMap[T]) GetOrDefault(key string, in T) (T, bool) {
	defer am.mutex.Unlock()
	am.mutex.Lock()

	out, found := am.internal[key]
	if !found {
		am.internal[key] = in
		return in, false
	}

	return out, true
}
