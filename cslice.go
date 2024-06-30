package cslice

import (
	"sync"
)

type CSlice[T any] struct {
	sync.RWMutex
	items []T
}

func New[T any](size ...int) CSlice[T] {
	switch len(size) {
	case 0:
		return CSlice[T]{}
	case 1:
		return CSlice[T]{items: make([]T, size[0])}
	default:
		return CSlice[T]{items: make([]T, size[0], size[1])}
	}
}

func From[T any](items []T) CSlice[T] {
	return CSlice[T]{items: items}
}

func (slice *CSlice[T]) PushFront(items ...T) {
	slice.Lock()
	defer slice.Unlock()
	slice.items = append(items, slice.items...)
}

func (slice *CSlice[T]) PopFront() (T, bool) {
	slice.Lock()
	defer slice.Unlock()
	var item T
	if len(slice.items) == 0 {
		return item, false
	}
	item, slice.items = slice.items[0], slice.items[1:]
	return item, true
}

func (slice *CSlice[T]) Front() (T, bool) {
	slice.RLock()
	defer slice.RUnlock()
	if len(slice.items) == 0 {
		var zero T
		return zero, false
	}
	return slice.items[0], true
}

func (slice *CSlice[T]) PushBack(items ...T) {
	slice.Lock()
	defer slice.Unlock()
	slice.items = append(slice.items, items...)
}

func (slice *CSlice[T]) PopBack() (T, bool) {
	slice.Lock()
	defer slice.Unlock()
	var item T
	if len(slice.items) == 0 {
		return item, false
	}
	item, slice.items = slice.items[len(slice.items)-1], slice.items[:len(slice.items)-1]
	return item, true
}

func (slice *CSlice[T]) Back() (T, bool) {
	slice.RLock()
	defer slice.RUnlock()
	if len(slice.items) == 0 {
		var zero T
		return zero, false
	}
	return slice.items[len(slice.items)-1], true
}

func (slice *CSlice[T]) Append(items ...T) {
	slice.Lock()
	defer slice.Unlock()
	slice.items = append(slice.items, items...)
}

func (slice *CSlice[T]) Prepend(items ...T) {
	slice.Lock()
	defer slice.Unlock()
	slice.items = append(items, slice.items...)
}

func (slice *CSlice[T]) Copy(items []T) int {
	slice.Lock()
	defer slice.Unlock()
	if len(slice.items) != len(items) {
		slice.items = make([]T, len(items))
	}
	return copy(slice.items, items)
}

func (slice *CSlice[T]) Get(index int) (T, bool) {
	slice.RLock()
	defer slice.RUnlock()
	if index < 0 || index >= len(slice.items) {
		var zero T
		return zero, false
	}
	return slice.items[index], true
}

func (slice *CSlice[T]) Set(index int, value T) {
	slice.RLock()
	defer slice.RUnlock()
	if index < 0 {
		return
	}
	if index >= len(slice.items) {
		slice.items = append(slice.items, make([]T, index-len(slice.items)+1)...)
	}
	slice.items[index] = value
}

func (slice *CSlice[T]) Overwrite(items []T) {
	slice.RLock()
	defer slice.RUnlock()
	slice.items = items
}

func (slice *CSlice[T]) Remove(index int) (T, bool) {
	slice.Lock()
	defer slice.Unlock()
	if index < 0 || index >= len(slice.items) {
		var zero T
		return zero, false
	}
	item := slice.items[index]
	slice.items = append(slice.items[:index], slice.items[index+1:]...)
	return item, true
}

func (slice *CSlice[T]) Delete(fn func(T) bool) bool {
	slice.Lock()
	defer slice.Unlock()
	index := -1
	for i, item := range slice.items {
		if fn(item) {
			index = i
			break
		}
	}
	if index == -1 {
		return false
	}
	slice.items = append(slice.items[:index], slice.items[index+1:]...)
	return true
}

func (slice *CSlice[T]) Len() int {
	slice.RLock()
	defer slice.RUnlock()
	return len(slice.items)
}

func (slice *CSlice[T]) Empty() bool {
	slice.RLock()
	defer slice.RUnlock()
	return len(slice.items) == 0
}

func (slice *CSlice[T]) Range(f func(index int, value T) bool) {
	len := slice.Len()
	for i := 0; i < len; i++ {
		value, ok := slice.Get(i)
		if !ok {
			break
		}
		if !f(i, value) {
			break
		}
	}
}

func (slice *CSlice[T]) Iter() chan T {
	ch := make(chan T)
	go func() {
		slice.Range(func(index int, value T) bool {
			ch <- value
			return true
		})
		close(ch)
	}()
	return ch
}

func (slice *CSlice[T]) Clear() {
	slice.RWMutex.Lock()
	defer slice.RWMutex.Unlock()
	slice.items = nil
}

func (slice *CSlice[T]) Slice(params ...int) []T {
	slice.RWMutex.Lock()
	defer slice.RWMutex.Unlock()
	var start, end int
	switch len(params) {
	case 0:
		start = 0
		end = len(slice.items)
	case 1:
		start = params[0]
		end = len(slice.items)
	default:
		start = params[0]
		end = params[1]
		if end < 0 {
			end = len(slice.items) + 1 + end
		}
	}
	if start > end {
		start, end = end, start
	}
	clone := make([]T, end-start)
	copy(clone, slice.items[start:end])
	return clone
}
