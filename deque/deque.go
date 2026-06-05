// Package deque implements a double-ended queue that mirrors the API and
// behaviour of C++ std::deque.
//
// It is backed by a growable circular buffer (ring buffer), which gives:
//
//	PushFront / PushBack   O(1) amortized
//	PopFront  / PopBack    O(1)
//	At (random access)     O(1)
//	Insert / Erase (middle) O(n) — shifts the smaller side
//
// The circular layout means front and back operations never shift elements,
// matching the constant-time guarantees of std::deque without the segmented
// storage of a typical libstdc++/libc++ implementation.
package deque

import "errors"

const minCapacity = 4

// Deque is a sequence container with fast insertion and removal at both ends.
type Deque[T any] struct {
	data []T // ring buffer; len(data) is the capacity
	head int // physical index of the logical front element
	size int // number of elements currently stored
}

// New creates and returns a new empty deque.
func New[T any]() *Deque[T] {
	return &Deque[T]{}
}

// NewFromSlice creates a deque initialized with a copy of the given slice,
// preserving order (items[0] becomes the front).
func NewFromSlice[T any](items []T) *Deque[T] {
	d := &Deque[T]{}
	if len(items) > 0 {
		d.data = make([]T, len(items))
		copy(d.data, items)
		d.size = len(items)
	}
	return d
}

// phys maps a logical index in [0, size) to its physical slot in the buffer.
func (d *Deque[T]) phys(i int) int {
	return (d.head + i) % len(d.data)
}

// ensureCap grows the buffer so it can hold at least minCap elements, copying
// the current elements into a freshly allocated, head-aligned buffer.
func (d *Deque[T]) ensureCap(minCap int) {
	if minCap <= len(d.data) {
		return
	}
	newCap := len(d.data) * 2
	if newCap < minCapacity {
		newCap = minCapacity
	}
	if newCap < minCap {
		newCap = minCap
	}
	newData := make([]T, newCap)
	for i := 0; i < d.size; i++ {
		newData[i] = d.data[d.phys(i)]
	}
	d.data = newData
	d.head = 0
}

// Size returns the number of elements. Equivalent to C++ size().
func (d *Deque[T]) Size() int { return d.size }

// Empty reports whether the deque has no elements. Equivalent to empty().
func (d *Deque[T]) Empty() bool { return d.size == 0 }

// Capacity returns the number of elements the deque can hold before growing.
// This is a Go convenience, not part of std::deque.
func (d *Deque[T]) Capacity() int { return len(d.data) }

// PushBack appends value to the back in amortized O(1). Equivalent to push_back.
func (d *Deque[T]) PushBack(value T) {
	d.ensureCap(d.size + 1)
	d.data[d.phys(d.size)] = value
	d.size++
}

// PushFront prepends value to the front in amortized O(1). Equivalent to
// push_front.
func (d *Deque[T]) PushFront(value T) {
	d.ensureCap(d.size + 1)
	d.head = (d.head - 1 + len(d.data)) % len(d.data)
	d.data[d.head] = value
	d.size++
}

// EmplaceBack constructs an element in place at the back. Equivalent to Push-
// Back in Go. Matches C++ emplace_back().
func (d *Deque[T]) EmplaceBack(value T) { d.PushBack(value) }

// EmplaceFront constructs an element in place at the front. Equivalent to
// PushFront in Go. Matches C++ emplace_front().
func (d *Deque[T]) EmplaceFront(value T) { d.PushFront(value) }

// PopFront removes the front element in O(1). Equivalent to pop_front().
func (d *Deque[T]) PopFront() error {
	if d.size == 0 {
		return errors.New("deque: pop_front on empty deque")
	}
	var zero T
	d.data[d.head] = zero
	d.head = (d.head + 1) % len(d.data)
	d.size--
	return nil
}

// PopBack removes the back element in O(1). Equivalent to pop_back().
func (d *Deque[T]) PopBack() error {
	if d.size == 0 {
		return errors.New("deque: pop_back on empty deque")
	}
	var zero T
	d.data[d.phys(d.size-1)] = zero
	d.size--
	return nil
}

// Front returns the front element. Equivalent to C++ front().
func (d *Deque[T]) Front() (T, error) {
	if d.size == 0 {
		var zero T
		return zero, errors.New("deque: front of empty deque")
	}
	return d.data[d.head], nil
}

// Back returns the back element. Equivalent to C++ back().
func (d *Deque[T]) Back() (T, error) {
	if d.size == 0 {
		var zero T
		return zero, errors.New("deque: back of empty deque")
	}
	return d.data[d.phys(d.size-1)], nil
}

// At returns the element at logical position i with bounds checking.
// Equivalent to C++ at().
func (d *Deque[T]) At(i int) (T, error) {
	if i < 0 || i >= d.size {
		var zero T
		return zero, errors.New("deque: index out of range")
	}
	return d.data[d.phys(i)], nil
}

// Get returns the element at position i without bounds checking.
// Equivalent to C++ operator[]; out-of-range access panics.
func (d *Deque[T]) Get(i int) T { return d.data[d.phys(i)] }

// Set assigns value at position i without bounds checking.
// Equivalent to assigning through C++ operator[]; out-of-range access panics.
func (d *Deque[T]) Set(i int, value T) { d.data[d.phys(i)] = value }

// Insert inserts value before logical position index, shifting the smaller of
// the two sides for efficiency. Runs in O(min(index, size-index)). Equivalent
// to C++ insert(begin()+index, value).
func (d *Deque[T]) Insert(index int, value T) error {
	if index < 0 || index > d.size {
		return errors.New("deque: insert index out of range")
	}
	if index == 0 {
		d.PushFront(value)
		return nil
	}
	if index == d.size {
		d.PushBack(value)
		return nil
	}
	d.ensureCap(d.size + 1)
	cap := len(d.data)
	if index < d.size-index {
		// Shift the front part [0, index) one slot toward the front.
		newHead := (d.head - 1 + cap) % cap
		for k := 0; k < index; k++ {
			d.data[(newHead+k)%cap] = d.data[(d.head+k)%cap]
		}
		d.head = newHead
		d.size++
		d.data[d.phys(index)] = value
	} else {
		// Shift the back part [index, size) one slot toward the back.
		d.size++
		for k := d.size - 1; k > index; k-- {
			d.data[d.phys(k)] = d.data[d.phys(k-1)]
		}
		d.data[d.phys(index)] = value
	}
	return nil
}

// Erase removes the element at logical position index, shifting the smaller of
// the two sides. Runs in O(min(index, size-index)). Equivalent to C++
// erase(begin()+index).
func (d *Deque[T]) Erase(index int) error {
	if index < 0 || index >= d.size {
		return errors.New("deque: erase index out of range")
	}
	if index == 0 {
		return d.PopFront()
	}
	if index == d.size-1 {
		return d.PopBack()
	}
	var zero T
	if index < d.size-1-index {
		// Shift the front part [0, index) one slot toward the back.
		for k := index; k > 0; k-- {
			d.data[d.phys(k)] = d.data[d.phys(k-1)]
		}
		d.data[d.head] = zero
		d.head = (d.head + 1) % len(d.data)
	} else {
		// Shift the back part (index, size) one slot toward the front.
		for k := index; k < d.size-1; k++ {
			d.data[d.phys(k)] = d.data[d.phys(k+1)]
		}
		d.data[d.phys(d.size-1)] = zero
	}
	d.size--
	return nil
}

// Reserve grows capacity to hold at least n elements without changing size.
// This is a Go convenience, not part of std::deque.
func (d *Deque[T]) Reserve(n int) { d.ensureCap(n) }

// Resize changes the number of elements to n, zero-filling any growth.
// Equivalent to C++ resize(n).
func (d *Deque[T]) Resize(n int) {
	var zero T
	d.ResizeWithValue(n, zero)
}

// ResizeWithValue changes the number of elements to n, padding any growth with
// copies of value. Equivalent to C++ resize(n, value).
func (d *Deque[T]) ResizeWithValue(n int, value T) {
	if n < 0 {
		n = 0
	}
	for d.size > n {
		d.PopBack()
	}
	if n > d.size {
		d.ensureCap(n)
		for d.size < n {
			d.data[d.phys(d.size)] = value
			d.size++
		}
	}
}

// ShrinkToFit reallocates so capacity equals size, releasing unused memory.
// Equivalent to C++ shrink_to_fit().
func (d *Deque[T]) ShrinkToFit() {
	if len(d.data) == d.size {
		return
	}
	newData := make([]T, d.size)
	for i := 0; i < d.size; i++ {
		newData[i] = d.data[d.phys(i)]
	}
	d.data = newData
	d.head = 0
}

// Clear removes all elements while retaining capacity. Equivalent to clear().
func (d *Deque[T]) Clear() {
	var zero T
	for i := 0; i < d.size; i++ {
		d.data[d.phys(i)] = zero
	}
	d.head = 0
	d.size = 0
}

// Swap exchanges the contents of the deque with another in O(1).
// Equivalent to C++ swap().
func (d *Deque[T]) Swap(other *Deque[T]) {
	if other == nil {
		return
	}
	d.data, other.data = other.data, d.data
	d.head, other.head = other.head, d.head
	d.size, other.size = other.size, d.size
}

// ToSlice returns the elements as a new slice in front-to-back order.
func (d *Deque[T]) ToSlice() []T {
	result := make([]T, d.size)
	for i := 0; i < d.size; i++ {
		result[i] = d.data[d.phys(i)]
	}
	return result
}

// ForEach calls fn for each element in front-to-back order.
func (d *Deque[T]) ForEach(fn func(index int, value T)) {
	for i := 0; i < d.size; i++ {
		fn(i, d.data[d.phys(i)])
	}
}
