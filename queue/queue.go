package queue

import (
	"errors"
)

// Queue represents a FIFO (first-in-first-out) data structure
// that mimics the behavior of C++ std::queue
type Queue[T any] struct {
	items []T
}

// New creates and returns a new empty queue
func New[T any]() *Queue[T] {
	return &Queue[T]{
		items: make([]T, 0),
	}
}

// NewFromSlice creates a new queue initialized with elements from a slice
func NewFromSlice[T any](items []T) *Queue[T] {
	q := &Queue[T]{
		items: make([]T, len(items)),
	}
	copy(q.items, items)
	return q
}

// Push adds an element to the back of the queue (enqueue operation)
// Equivalent to C++ std::queue::push
func (q *Queue[T]) Push(value T) {
	q.items = append(q.items, value)
}

// Pop removes the front element from the queue (dequeue operation)
// Equivalent to C++ std::queue::pop
// Note: In C++, pop() doesn't return the value, it just removes it
func (q *Queue[T]) Pop() error {
	if q.Empty() {
		return errors.New("queue is empty")
	}
	q.items = q.items[1:]
	return nil
}

// Front returns a reference to the front element of the queue
// Equivalent to C++ std::queue::front
func (q *Queue[T]) Front() (T, error) {
	var zero T
	if q.Empty() {
		return zero, errors.New("queue is empty")
	}
	return q.items[0], nil
}

// Back returns a reference to the back element of the queue
// Equivalent to C++ std::queue::back
func (q *Queue[T]) Back() (T, error) {
	var zero T
	if q.Empty() {
		return zero, errors.New("queue is empty")
	}
	return q.items[len(q.items)-1], nil
}

// Empty checks whether the queue is empty
// Equivalent to C++ std::queue::empty
func (q *Queue[T]) Empty() bool {
	return len(q.items) == 0
}

// Size returns the number of elements in the queue
// Equivalent to C++ std::queue::size
func (q *Queue[T]) Size() int {
	return len(q.items)
}

// Emplace constructs an element in-place at the back of the queue
// Equivalent to C++ std::queue::emplace (C++11)
// In Go, this is effectively the same as Push since we don't have constructors
func (q *Queue[T]) Emplace(value T) {
	q.Push(value)
}

// Swap exchanges the contents of the queue with another queue
// Equivalent to C++ std::queue::swap (C++11)
func (q *Queue[T]) Swap(other *Queue[T]) {
	if other == nil {
		return
	}
	q.items, other.items = other.items, q.items
}

// Clear removes all elements from the queue
// This is a convenience method not present in C++ std::queue
func (q *Queue[T]) Clear() {
	q.items = q.items[:0]
}

// ToSlice returns a copy of the queue as a slice
// This is a convenience method for Go users
func (q *Queue[T]) ToSlice() []T {
	result := make([]T, len(q.items))
	copy(result, q.items)
	return result
}
