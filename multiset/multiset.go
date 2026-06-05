// Package multiset implements a sorted associative container that mirrors the
// API and behaviour of C++ std::multiset. Unlike a set it permits multiple
// equivalent keys.
//
// It is backed by the module's red-black tree, keeping elements sorted with all
// core operations in O(log n) (Erase of a key is O((1+count)·log n)). Equivalent
// keys are kept adjacent so EqualRange, Count and ordered traversal behave just
// like std::multiset.
package multiset

import (
	"github.com/wiraphatys/gostd/constraints"
	"github.com/wiraphatys/gostd/internal/rbtree"
)

// MultiSet is a collection of keys, sorted by the configured comparator, that
// allows duplicates.
type MultiSet[T any] struct {
	tree *rbtree.Tree[T, struct{}]
}

// Iterator is a bidirectional cursor over a MultiSet, yielding keys in sorted
// order (equivalent keys appear consecutively).
type Iterator[T any] struct {
	tree *rbtree.Tree[T, struct{}]
	node *rbtree.Node[T, struct{}]
}

// New creates an empty multiset ordered by the natural ordering of T.
func New[T constraints.Ordered]() *MultiSet[T] {
	return &MultiSet[T]{tree: rbtree.New[T, struct{}](constraints.Less[T], true)}
}

// NewFunc creates an empty multiset ordered by the given comparator.
func NewFunc[T any](less func(a, b T) bool) *MultiSet[T] {
	if less == nil {
		panic("multiset: nil comparator")
	}
	return &MultiSet[T]{tree: rbtree.New[T, struct{}](less, true)}
}

// NewFromSlice creates a multiset containing all elements of items, keeping
// duplicates.
func NewFromSlice[T constraints.Ordered](items []T) *MultiSet[T] {
	m := New[T]()
	for _, v := range items {
		m.tree.Insert(v, struct{}{})
	}
	return m
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

func (m *MultiSet[T]) iter(n *rbtree.Node[T, struct{}]) Iterator[T] {
	return Iterator[T]{tree: m.tree, node: n}
}

// Size returns the number of elements (counting duplicates). C++ size().
func (m *MultiSet[T]) Size() int { return m.tree.Size() }

// Empty reports whether the multiset has no elements.
func (m *MultiSet[T]) Empty() bool { return m.tree.Size() == 0 }

// Clear removes all elements.
func (m *MultiSet[T]) Clear() { m.tree.Clear() }

// Insert adds value; a multiset always stores it, so this never fails.
// Equivalent to std::multiset::insert.
func (m *MultiSet[T]) Insert(value T) { m.tree.Insert(value, struct{}{}) }

// Erase removes every element equivalent to value and returns how many were
// removed. Equivalent to std::multiset::erase(key).
func (m *MultiSet[T]) Erase(value T) int { return m.tree.EraseKey(value) }

// EraseOne removes a single element equivalent to value (if any) and reports
// whether one was removed.
func (m *MultiSet[T]) EraseOne(value T) bool {
	n := m.tree.Find(value)
	if n == nil {
		return false
	}
	m.tree.EraseNode(n)
	return true
}

// EraseIter removes the element the iterator refers to and returns an iterator
// to the following element. Equivalent to std::multiset::erase(iterator).
func (m *MultiSet[T]) EraseIter(it Iterator[T]) Iterator[T] {
	if it.node == nil {
		return m.End()
	}
	return m.iter(m.tree.EraseNode(it.node))
}

// Find returns an iterator to some element equivalent to value, or End().
func (m *MultiSet[T]) Find(value T) Iterator[T] { return m.iter(m.tree.Find(value)) }

// Contains reports whether at least one element equivalent to value is present.
func (m *MultiSet[T]) Contains(value T) bool { return m.tree.Contains(value) }

// Count returns the number of elements equivalent to value.
func (m *MultiSet[T]) Count(value T) int { return m.tree.Count(value) }

// LowerBound returns an iterator to the first element not less than value.
func (m *MultiSet[T]) LowerBound(value T) Iterator[T] { return m.iter(m.tree.LowerBound(value)) }

// UpperBound returns an iterator to the first element greater than value.
func (m *MultiSet[T]) UpperBound(value T) Iterator[T] { return m.iter(m.tree.UpperBound(value)) }

// EqualRange returns the range [LowerBound, UpperBound) of elements equivalent
// to value. Equivalent to std::multiset::equal_range.
func (m *MultiSet[T]) EqualRange(value T) (Iterator[T], Iterator[T]) {
	return m.LowerBound(value), m.UpperBound(value)
}

// Begin returns an iterator to the smallest element (or End() if empty).
func (m *MultiSet[T]) Begin() Iterator[T] { return m.iter(m.tree.Min()) }

// End returns the past-the-end iterator.
func (m *MultiSet[T]) End() Iterator[T] { return m.iter(nil) }

// RBegin returns an iterator to the largest element, for reverse traversal.
func (m *MultiSet[T]) RBegin() Iterator[T] { return m.iter(m.tree.Max()) }

// REnd returns the position before the first element.
func (m *MultiSet[T]) REnd() Iterator[T] { return m.iter(nil) }

// Swap exchanges the contents of the multiset with another in O(1).
func (m *MultiSet[T]) Swap(other *MultiSet[T]) {
	if other == nil {
		return
	}
	m.tree, other.tree = other.tree, m.tree
}

// ToSlice returns all elements as a new slice in sorted order (with duplicates).
func (m *MultiSet[T]) ToSlice() []T {
	out := make([]T, 0, m.tree.Size())
	for n := m.tree.Min(); n != nil; n = m.tree.Next(n) {
		out = append(out, n.Key)
	}
	return out
}

// ForEach calls fn for each element in sorted order.
func (m *MultiSet[T]) ForEach(fn func(value T)) {
	for n := m.tree.Min(); n != nil; n = m.tree.Next(n) {
		fn(n.Key)
	}
}
