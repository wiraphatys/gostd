// Package set implements a sorted associative container of unique keys that
// mirrors the API and behaviour of C++ std::set.
//
// It is backed by the module's red-black tree, so the elements are always kept
// in sorted order and every core operation runs in O(log n):
//
//	Insert / Erase / Find / Count / Contains            O(log n)
//	LowerBound / UpperBound / EqualRange                O(log n)
//	Begin / End / RBegin / REnd                         O(log n) to reach, O(1) amortized to step
//
// By default elements are ordered by the natural less-than ordering (like
// std::less); supply a custom comparator with NewFunc for any other ordering.
package set

import (
	"github.com/wiraphatys/gostd/constraints"
	"github.com/wiraphatys/gostd/internal/rbtree"
)

// Set is a collection of unique keys, sorted by the configured comparator.
type Set[T any] struct {
	tree *rbtree.Tree[T, struct{}]
}

// Iterator is a bidirectional cursor over a Set, yielding keys in sorted order.
type Iterator[T any] struct {
	tree *rbtree.Tree[T, struct{}]
	node *rbtree.Node[T, struct{}]
}

// New creates an empty set ordered by the natural ordering of T.
func New[T constraints.Ordered]() *Set[T] {
	return &Set[T]{tree: rbtree.New[T, struct{}](constraints.Less[T], false)}
}

// NewFunc creates an empty set ordered by the given comparator. less(a, b) must
// report whether a should be ordered before b (a strict weak ordering).
func NewFunc[T any](less func(a, b T) bool) *Set[T] {
	if less == nil {
		panic("set: nil comparator")
	}
	return &Set[T]{tree: rbtree.New[T, struct{}](less, false)}
}

// NewFromSlice creates a set containing the unique elements of items.
func NewFromSlice[T constraints.Ordered](items []T) *Set[T] {
	s := New[T]()
	for _, v := range items {
		s.tree.Insert(v, struct{}{})
	}
	return s
}

// Value returns the key the iterator refers to. Undefined at the End position.
func (it Iterator[T]) Value() T { return it.node.Key }

// Next advances the iterator to the next element in sorted order.
func (it *Iterator[T]) Next() { it.node = it.tree.Next(it.node) }

// Prev moves the iterator to the previous element in sorted order.
func (it *Iterator[T]) Prev() { it.node = it.tree.Prev(it.node) }

// Equal reports whether two iterators refer to the same position.
func (it Iterator[T]) Equal(other Iterator[T]) bool { return it.node == other.node }

// IsEnd reports whether the iterator is the past-the-end iterator.
func (it Iterator[T]) IsEnd() bool { return it.node == nil }

func (s *Set[T]) iter(n *rbtree.Node[T, struct{}]) Iterator[T] {
	return Iterator[T]{tree: s.tree, node: n}
}

// Size returns the number of elements. Equivalent to C++ size().
func (s *Set[T]) Size() int { return s.tree.Size() }

// Empty reports whether the set has no elements. Equivalent to empty().
func (s *Set[T]) Empty() bool { return s.tree.Size() == 0 }

// Clear removes all elements. Equivalent to C++ clear().
func (s *Set[T]) Clear() { s.tree.Clear() }

// Insert adds value and reports whether it was newly inserted (false if an
// equivalent element was already present). Equivalent to std::set::insert,
// whose returned pair's bool this mirrors.
func (s *Set[T]) Insert(value T) bool {
	_, inserted := s.tree.Insert(value, struct{}{})
	return inserted
}

// Erase removes value if present and reports whether something was removed.
// Equivalent to std::set::erase(key) (whose count return is 0 or 1).
func (s *Set[T]) Erase(value T) bool {
	return s.tree.EraseKey(value) > 0
}

// EraseIter removes the element the iterator refers to and returns an iterator
// to the following element. Equivalent to std::set::erase(iterator).
func (s *Set[T]) EraseIter(it Iterator[T]) Iterator[T] {
	if it.node == nil {
		return s.End()
	}
	return s.iter(s.tree.EraseNode(it.node))
}

// Find returns an iterator to value, or End() if it is not present.
// Equivalent to std::set::find.
func (s *Set[T]) Find(value T) Iterator[T] { return s.iter(s.tree.Find(value)) }

// Contains reports whether value is present (C++20 std::set::contains).
func (s *Set[T]) Contains(value T) bool { return s.tree.Contains(value) }

// Count returns 1 if value is present, else 0. Equivalent to std::set::count.
func (s *Set[T]) Count(value T) int {
	if s.tree.Contains(value) {
		return 1
	}
	return 0
}

// LowerBound returns an iterator to the first element not less than value, or
// End(). Equivalent to std::set::lower_bound.
func (s *Set[T]) LowerBound(value T) Iterator[T] { return s.iter(s.tree.LowerBound(value)) }

// UpperBound returns an iterator to the first element greater than value, or
// End(). Equivalent to std::set::upper_bound.
func (s *Set[T]) UpperBound(value T) Iterator[T] { return s.iter(s.tree.UpperBound(value)) }

// EqualRange returns the range [LowerBound, UpperBound) of elements equivalent
// to value. Equivalent to std::set::equal_range.
func (s *Set[T]) EqualRange(value T) (Iterator[T], Iterator[T]) {
	return s.LowerBound(value), s.UpperBound(value)
}

// Begin returns an iterator to the smallest element (or End() if empty).
func (s *Set[T]) Begin() Iterator[T] { return s.iter(s.tree.Min()) }

// End returns the past-the-end iterator.
func (s *Set[T]) End() Iterator[T] { return s.iter(nil) }

// RBegin returns an iterator to the largest element, for reverse traversal with
// Prev (or REnd if empty).
func (s *Set[T]) RBegin() Iterator[T] { return s.iter(s.tree.Max()) }

// REnd returns the position before the first element, the stop condition for
// reverse traversal.
func (s *Set[T]) REnd() Iterator[T] { return s.iter(nil) }

// Swap exchanges the contents of the set with another in O(1).
func (s *Set[T]) Swap(other *Set[T]) {
	if other == nil {
		return
	}
	s.tree, other.tree = other.tree, s.tree
}

// ToSlice returns all elements as a new slice in sorted order.
func (s *Set[T]) ToSlice() []T {
	out := make([]T, 0, s.tree.Size())
	for n := s.tree.Min(); n != nil; n = s.tree.Next(n) {
		out = append(out, n.Key)
	}
	return out
}

// ForEach calls fn for each element in sorted order.
func (s *Set[T]) ForEach(fn func(value T)) {
	for n := s.tree.Min(); n != nil; n = s.tree.Next(n) {
		fn(n.Key)
	}
}
