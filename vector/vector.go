// Package vector implements a dynamic contiguous array that mirrors the API and
// behaviour of C++ std::vector.
//
// Like std::vector, elements are stored contiguously, random access is O(1),
// and PushBack runs in amortized O(1) time thanks to geometric capacity growth.
// Insertion or removal at the end is O(1) amortized; anywhere else is O(n)
// because the tail must be shifted, exactly as in C++.
package vector

import (
	"errors"
)

// Vector is a sequence container representing a dynamically sized array.
// The zero value is not ready for use; create one with New and friends.
type Vector[T any] struct {
	data []T
}

// New creates and returns a new empty vector.
func New[T any]() *Vector[T] {
	return &Vector[T]{data: make([]T, 0)}
}

// NewWithSize creates a vector containing n zero-valued elements.
// Equivalent to C++ std::vector<T>(n).
func NewWithSize[T any](n int) *Vector[T] {
	if n < 0 {
		n = 0
	}
	return &Vector[T]{data: make([]T, n)}
}

// NewWithValue creates a vector containing n copies of value.
// Equivalent to C++ std::vector<T>(n, value).
func NewWithValue[T any](n int, value T) *Vector[T] {
	if n < 0 {
		n = 0
	}
	data := make([]T, n)
	for i := range data {
		data[i] = value
	}
	return &Vector[T]{data: data}
}

// NewFromSlice creates a vector initialized with a copy of the given slice.
func NewFromSlice[T any](items []T) *Vector[T] {
	data := make([]T, len(items))
	copy(data, items)
	return &Vector[T]{data: data}
}

// Size returns the number of elements. Equivalent to C++ size().
func (v *Vector[T]) Size() int { return len(v.data) }

// Capacity returns the number of elements the vector can hold before
// reallocating. Equivalent to C++ capacity().
func (v *Vector[T]) Capacity() int { return cap(v.data) }

// Empty reports whether the vector contains no elements. Equivalent to empty().
func (v *Vector[T]) Empty() bool { return len(v.data) == 0 }

// At returns the element at position i with bounds checking.
// Equivalent to C++ at(): it reports an error instead of throwing
// std::out_of_range.
func (v *Vector[T]) At(i int) (T, error) {
	if i < 0 || i >= len(v.data) {
		var zero T
		return zero, errors.New("vector: index out of range")
	}
	return v.data[i], nil
}

// Get returns the element at position i without bounds checking.
// Equivalent to C++ operator[]: passing an out-of-range index panics, mirroring
// the undefined behaviour of the C++ original.
func (v *Vector[T]) Get(i int) T { return v.data[i] }

// Set assigns value to the element at position i without bounds checking.
// Equivalent to assigning through C++ operator[]. Panics if i is out of range.
func (v *Vector[T]) Set(i int, value T) { v.data[i] = value }

// Front returns the first element. Equivalent to C++ front().
func (v *Vector[T]) Front() (T, error) {
	if len(v.data) == 0 {
		var zero T
		return zero, errors.New("vector: front of empty vector")
	}
	return v.data[0], nil
}

// Back returns the last element. Equivalent to C++ back().
func (v *Vector[T]) Back() (T, error) {
	if len(v.data) == 0 {
		var zero T
		return zero, errors.New("vector: back of empty vector")
	}
	return v.data[len(v.data)-1], nil
}

// PushBack appends value to the end of the vector in amortized O(1) time.
// Equivalent to C++ push_back().
func (v *Vector[T]) PushBack(value T) {
	v.data = append(v.data, value)
}

// PopBack removes the last element. Equivalent to C++ pop_back().
// The vacated slot is zeroed so it does not keep references alive for the GC.
func (v *Vector[T]) PopBack() error {
	n := len(v.data)
	if n == 0 {
		return errors.New("vector: pop_back on empty vector")
	}
	var zero T
	v.data[n-1] = zero
	v.data = v.data[:n-1]
	return nil
}

// EmplaceBack constructs an element in place at the end of the vector.
// In Go there is no separate construction step, so it is equivalent to PushBack.
func (v *Vector[T]) EmplaceBack(value T) { v.PushBack(value) }

// Insert inserts value before position index, shifting subsequent elements.
// Runs in O(n). Equivalent to C++ insert(begin()+index, value).
func (v *Vector[T]) Insert(index int, value T) error {
	n := len(v.data)
	if index < 0 || index > n {
		return errors.New("vector: insert index out of range")
	}
	var zero T
	v.data = append(v.data, zero)
	copy(v.data[index+1:], v.data[index:n])
	v.data[index] = value
	return nil
}

// InsertRange inserts all of values before position index, preserving order.
// Runs in O(n + len(values)). Equivalent to the range overload of C++ insert.
func (v *Vector[T]) InsertRange(index int, values []T) error {
	n := len(v.data)
	if index < 0 || index > n {
		return errors.New("vector: insert index out of range")
	}
	m := len(values)
	if m == 0 {
		return nil
	}
	v.data = append(v.data, make([]T, m)...)
	copy(v.data[index+m:], v.data[index:n])
	copy(v.data[index:index+m], values)
	return nil
}

// Erase removes the element at position index, shifting subsequent elements.
// Runs in O(n). Equivalent to C++ erase(begin()+index).
func (v *Vector[T]) Erase(index int) error {
	n := len(v.data)
	if index < 0 || index >= n {
		return errors.New("vector: erase index out of range")
	}
	copy(v.data[index:], v.data[index+1:])
	var zero T
	v.data[n-1] = zero
	v.data = v.data[:n-1]
	return nil
}

// EraseRange removes the elements in the half-open range [first, last).
// Runs in O(n). Equivalent to C++ erase(begin()+first, begin()+last).
func (v *Vector[T]) EraseRange(first, last int) error {
	n := len(v.data)
	if first < 0 || last > n || first > last {
		return errors.New("vector: erase range out of range")
	}
	count := last - first
	if count == 0 {
		return nil
	}
	copy(v.data[first:], v.data[last:])
	var zero T
	for i := n - count; i < n; i++ {
		v.data[i] = zero
	}
	v.data = v.data[:n-count]
	return nil
}

// Reserve increases the capacity to at least n elements. It never shrinks the
// vector and never changes its size. Equivalent to C++ reserve().
func (v *Vector[T]) Reserve(n int) {
	if n <= cap(v.data) {
		return
	}
	newData := make([]T, len(v.data), n)
	copy(newData, v.data)
	v.data = newData
}

// Resize changes the number of elements to n. If the vector grows, the new
// elements are zero-valued. Equivalent to C++ resize(n).
func (v *Vector[T]) Resize(n int) {
	var zero T
	v.ResizeWithValue(n, zero)
}

// ResizeWithValue changes the number of elements to n, padding any growth with
// copies of value. Equivalent to C++ resize(n, value).
func (v *Vector[T]) ResizeWithValue(n int, value T) {
	if n < 0 {
		n = 0
	}
	cur := len(v.data)
	switch {
	case n < cur:
		var zero T
		for i := n; i < cur; i++ {
			v.data[i] = zero
		}
		v.data = v.data[:n]
	case n > cur:
		if n <= cap(v.data) {
			v.data = v.data[:n]
		} else {
			newData := make([]T, n)
			copy(newData, v.data)
			v.data = newData
		}
		for i := cur; i < n; i++ {
			v.data[i] = value
		}
	}
}

// ShrinkToFit reduces capacity to match the current size, releasing any unused
// memory. Equivalent to C++ shrink_to_fit().
func (v *Vector[T]) ShrinkToFit() {
	if cap(v.data) == len(v.data) {
		return
	}
	newData := make([]T, len(v.data))
	copy(newData, v.data)
	v.data = newData
}

// Assign replaces the contents with n copies of value.
// Equivalent to C++ assign(n, value).
func (v *Vector[T]) Assign(n int, value T) {
	if n < 0 {
		n = 0
	}
	if n <= cap(v.data) {
		v.data = v.data[:n]
	} else {
		v.data = make([]T, n)
	}
	for i := 0; i < n; i++ {
		v.data[i] = value
	}
}

// AssignSlice replaces the contents with a copy of the given slice.
// Equivalent to the range overload of C++ assign.
func (v *Vector[T]) AssignSlice(items []T) {
	if len(items) <= cap(v.data) {
		v.data = v.data[:len(items)]
	} else {
		v.data = make([]T, len(items))
	}
	copy(v.data, items)
}

// Clear removes all elements but retains capacity, matching C++ clear().
func (v *Vector[T]) Clear() {
	var zero T
	for i := range v.data {
		v.data[i] = zero
	}
	v.data = v.data[:0]
}

// Swap exchanges the contents of the vector with another vector in O(1).
// Equivalent to C++ swap().
func (v *Vector[T]) Swap(other *Vector[T]) {
	if other == nil {
		return
	}
	v.data, other.data = other.data, v.data
}

// Data returns the underlying slice. Mutating it mutates the vector, mirroring
// the pointer returned by C++ data(). The slice is valid until the next
// operation that reallocates the vector.
func (v *Vector[T]) Data() []T { return v.data }

// ToSlice returns a copy of the vector's elements as a new slice.
func (v *Vector[T]) ToSlice() []T {
	result := make([]T, len(v.data))
	copy(result, v.data)
	return result
}

// ForEach calls fn for each element in order, passing its index and value.
// It is a Go-friendly substitute for iterating with begin()/end().
func (v *Vector[T]) ForEach(fn func(index int, value T)) {
	for i, item := range v.data {
		fn(i, item)
	}
}
