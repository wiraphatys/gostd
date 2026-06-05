package queue

import (
	"errors"

	"github.com/wiraphatys/gostd/deque"
)

// Queue represents a FIFO (first-in-first-out) data structure
// that mimics the behavior of C++ std::queue.
//
// Like C++ std::queue (which is a container adapter over std::deque by
// default), this Queue is backed by a growable circular ring buffer. That
// gives O(1) Push, Pop, Front and Back with memory proportional to the number
// of live elements — not to the total number of items ever enqueued — and it
// clears popped slots immediately so they do not keep references alive.
type Queue[T any] struct {
	d *deque.Deque[T]
}

// New creates and returns a new empty queue
func New[T any]() *Queue[T] {
	return &Queue[T]{d: deque.New[T]()}
}

// NewFromSlice creates a new queue initialized with elements from a slice
func NewFromSlice[T any](items []T) *Queue[T] {
	return &Queue[T]{d: deque.NewFromSlice(items)}
}

// Push adds an element to the back of the queue (enqueue operation)
// Equivalent to C++ std::queue::push. Runs in amortized O(1).
func (q *Queue[T]) Push(value T) {
	q.d.PushBack(value)
}

// Pop removes the front element from the queue (dequeue operation)
// Equivalent to C++ std::queue::pop. Runs in O(1).
// Note: In C++, pop() doesn't return the value, it just removes it
func (q *Queue[T]) Pop() error {
	if q.d.Empty() {
		return errors.New("queue is empty")
	}
	return q.d.PopFront()
}

// Front returns a reference to the front element of the queue
// Equivalent to C++ std::queue::front
func (q *Queue[T]) Front() (T, error) {
	if q.d.Empty() {
		var zero T
		return zero, errors.New("queue is empty")
	}
	return q.d.Front()
}

// Back returns a reference to the back element of the queue
// Equivalent to C++ std::queue::back
func (q *Queue[T]) Back() (T, error) {
	if q.d.Empty() {
		var zero T
		return zero, errors.New("queue is empty")
	}
	return q.d.Back()
}

// Empty checks whether the queue is empty
// Equivalent to C++ std::queue::empty
func (q *Queue[T]) Empty() bool {
	return q.d.Empty()
}

// Size returns the number of elements in the queue
// Equivalent to C++ std::queue::size
func (q *Queue[T]) Size() int {
	return q.d.Size()
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
	q.d.Swap(other.d)
}

// Clear removes all elements from the queue
// This is a convenience method not present in C++ std::queue
func (q *Queue[T]) Clear() {
	q.d.Clear()
}

// ToSlice returns a copy of the queue as a slice
// This is a convenience method for Go users
// The slice is ordered from front to back (index 0 = front)
func (q *Queue[T]) ToSlice() []T {
	return q.d.ToSlice()
}
