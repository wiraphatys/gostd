package stack

import (
	"errors"
)

// Stack represents a LIFO (last-in-first-out) data structure
// that mimics the behavior of C++ std::stack
type Stack[T any] struct {
	items []T
}

// New creates and returns a new empty stack
func New[T any]() *Stack[T] {
	return &Stack[T]{
		items: make([]T, 0),
	}
}

// NewFromSlice creates a new stack initialized with elements from a slice
// Elements are pushed in order, so the last element of the slice becomes the top
func NewFromSlice[T any](items []T) *Stack[T] {
	s := &Stack[T]{
		items: make([]T, len(items)),
	}
	copy(s.items, items)
	return s
}

// Push adds an element to the top of the stack
// Equivalent to C++ std::stack::push
func (s *Stack[T]) Push(value T) {
	s.items = append(s.items, value)
}

// Pop removes the top element from the stack
// Equivalent to C++ std::stack::pop
// Note: In C++, pop() doesn't return the value, it just removes it
func (s *Stack[T]) Pop() error {
	if s.Empty() {
		return errors.New("stack is empty")
	}
	s.items = s.items[:len(s.items)-1]
	return nil
}

// Top returns a reference to the top element of the stack
// Equivalent to C++ std::stack::top
func (s *Stack[T]) Top() (T, error) {
	var zero T
	if s.Empty() {
		return zero, errors.New("stack is empty")
	}
	return s.items[len(s.items)-1], nil
}

// Empty checks whether the stack is empty
// Equivalent to C++ std::stack::empty
func (s *Stack[T]) Empty() bool {
	return len(s.items) == 0
}

// Size returns the number of elements in the stack
// Equivalent to C++ std::stack::size
func (s *Stack[T]) Size() int {
	return len(s.items)
}

// Emplace constructs an element in-place at the top of the stack
// Equivalent to C++ std::stack::emplace (C++11)
// In Go, this is effectively the same as Push since we don't have constructors
func (s *Stack[T]) Emplace(value T) {
	s.Push(value)
}

// Swap exchanges the contents of the stack with another stack
// Equivalent to C++ std::stack::swap (C++11)
func (s *Stack[T]) Swap(other *Stack[T]) {
	if other == nil {
		return
	}
	s.items, other.items = other.items, s.items
}

// Clear removes all elements from the stack
// This is a convenience method not present in C++ std::stack
func (s *Stack[T]) Clear() {
	s.items = s.items[:0]
}

// ToSlice returns a copy of the stack as a slice
// This is a convenience method for Go users
// The slice is ordered from bottom to top (index 0 = bottom, last index = top)
func (s *Stack[T]) ToSlice() []T {
	result := make([]T, len(s.items))
	copy(result, s.items)
	return result
}