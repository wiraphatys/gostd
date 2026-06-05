// Package priorityqueue implements a heap-based priority queue that mirrors the
// API and behaviour of C++ std::priority_queue.
//
// It is backed by a binary heap stored in a contiguous slice, giving:
//
//	Push  O(log n)
//	Pop   O(log n)
//	Top   O(1)
//
// As in C++, the default ordering uses the natural less-than comparison, which
// yields a max-heap: Top returns the greatest element. Supply a custom
// comparator with NewFunc (or use NewMin) to obtain a min-heap or any other
// ordering. The comparator follows std::priority_queue semantics: less(a, b)
// must report whether a should appear *below* b in priority, so the element
// that is not less than any other ends up on top.
package priorityqueue

import (
	"errors"

	"github.com/wiraphatys/gostd/constraints"
)

// PriorityQueue is a container adapter that provides constant-time access to
// the largest (by the configured comparator) element.
type PriorityQueue[T any] struct {
	data []T
	less func(a, b T) bool
}

// New creates an empty max-priority-queue ordered by the natural ordering of T.
// Top will return the greatest element, matching C++ std::priority_queue<T>.
func New[T constraints.Ordered]() *PriorityQueue[T] {
	return &PriorityQueue[T]{less: constraints.Less[T]}
}

// NewMin creates an empty min-priority-queue ordered by the natural ordering of
// T. Top will return the smallest element, matching C++
// std::priority_queue<T, std::vector<T>, std::greater<T>>.
func NewMin[T constraints.Ordered]() *PriorityQueue[T] {
	return &PriorityQueue[T]{less: constraints.Greater[T]}
}

// NewFunc creates an empty priority queue ordered by the given comparator.
// less(a, b) must report whether a has lower priority than b; the element with
// the highest priority (not less than any other) is returned by Top.
func NewFunc[T any](less func(a, b T) bool) *PriorityQueue[T] {
	if less == nil {
		panic("priorityqueue: nil comparator")
	}
	return &PriorityQueue[T]{less: less}
}

// NewFromSlice builds a max-priority-queue from the elements of items in O(n)
// using Floyd's bottom-up heap construction.
func NewFromSlice[T constraints.Ordered](items []T) *PriorityQueue[T] {
	pq := &PriorityQueue[T]{
		data: make([]T, len(items)),
		less: constraints.Less[T],
	}
	copy(pq.data, items)
	pq.heapify()
	return pq
}

// NewFromSliceFunc builds a priority queue from items using the given
// comparator in O(n).
func NewFromSliceFunc[T any](items []T, less func(a, b T) bool) *PriorityQueue[T] {
	if less == nil {
		panic("priorityqueue: nil comparator")
	}
	pq := &PriorityQueue[T]{
		data: make([]T, len(items)),
		less: less,
	}
	copy(pq.data, items)
	pq.heapify()
	return pq
}

// heapify establishes the heap invariant over the whole backing slice in O(n).
func (pq *PriorityQueue[T]) heapify() {
	for i := len(pq.data)/2 - 1; i >= 0; i-- {
		pq.siftDown(i)
	}
}

// siftUp restores the heap property by moving the element at index i upward.
func (pq *PriorityQueue[T]) siftUp(i int) {
	for i > 0 {
		parent := (i - 1) / 2
		if !pq.less(pq.data[parent], pq.data[i]) {
			break
		}
		pq.data[parent], pq.data[i] = pq.data[i], pq.data[parent]
		i = parent
	}
}

// siftDown restores the heap property by moving the element at index i downward.
func (pq *PriorityQueue[T]) siftDown(i int) {
	n := len(pq.data)
	for {
		largest := i
		left := 2*i + 1
		right := 2*i + 2
		if left < n && pq.less(pq.data[largest], pq.data[left]) {
			largest = left
		}
		if right < n && pq.less(pq.data[largest], pq.data[right]) {
			largest = right
		}
		if largest == i {
			break
		}
		pq.data[i], pq.data[largest] = pq.data[largest], pq.data[i]
		i = largest
	}
}

// Push inserts value into the priority queue in O(log n).
// Equivalent to C++ push().
func (pq *PriorityQueue[T]) Push(value T) {
	pq.data = append(pq.data, value)
	pq.siftUp(len(pq.data) - 1)
}

// Emplace constructs an element in place. In Go there is no separate
// construction step, so it is equivalent to Push. Matches C++ emplace().
func (pq *PriorityQueue[T]) Emplace(value T) { pq.Push(value) }

// Pop removes the highest-priority element in O(log n).
// Equivalent to C++ pop(); like the original it does not return the value.
func (pq *PriorityQueue[T]) Pop() error {
	n := len(pq.data)
	if n == 0 {
		return errors.New("priorityqueue: pop on empty queue")
	}
	pq.data[0] = pq.data[n-1]
	var zero T
	pq.data[n-1] = zero
	pq.data = pq.data[:n-1]
	if len(pq.data) > 0 {
		pq.siftDown(0)
	}
	return nil
}

// Top returns the highest-priority element without removing it in O(1).
// Equivalent to C++ top().
func (pq *PriorityQueue[T]) Top() (T, error) {
	if len(pq.data) == 0 {
		var zero T
		return zero, errors.New("priorityqueue: top of empty queue")
	}
	return pq.data[0], nil
}

// Size returns the number of elements. Equivalent to C++ size().
func (pq *PriorityQueue[T]) Size() int { return len(pq.data) }

// Empty reports whether the queue is empty. Equivalent to C++ empty().
func (pq *PriorityQueue[T]) Empty() bool { return len(pq.data) == 0 }

// Swap exchanges the contents of the queue with another in O(1).
// Equivalent to C++ swap().
func (pq *PriorityQueue[T]) Swap(other *PriorityQueue[T]) {
	if other == nil {
		return
	}
	pq.data, other.data = other.data, pq.data
	pq.less, other.less = other.less, pq.less
}

// Clear removes all elements while keeping the configured comparator.
func (pq *PriorityQueue[T]) Clear() {
	var zero T
	for i := range pq.data {
		pq.data[i] = zero
	}
	pq.data = pq.data[:0]
}

// ToSlice returns a copy of the underlying heap array. The order is heap order,
// not sorted order. This is a Go convenience, not part of std::priority_queue.
func (pq *PriorityQueue[T]) ToSlice() []T {
	result := make([]T, len(pq.data))
	copy(result, pq.data)
	return result
}
