// Package multimap implements a sorted key-value associative container that
// mirrors the API and behaviour of C++ std::multimap. Unlike a map it permits
// multiple entries with equivalent keys.
//
// It is backed by the module's red-black tree, keeping entries sorted by key
// with core operations in O(log n) (Erase of a key is O((1+count)·log n)).
// Entries with equivalent keys are stored adjacently so EqualRange, Count and
// ordered traversal behave like std::multimap. As in C++, there is no operator[]
// or At because a key may map to many values.
package multimap

import (
	"github.com/wiraphatys/gostd/constraints"
	"github.com/wiraphatys/gostd/internal/rbtree"
)

// MultiMap is an ordered mapping from keys of type K to values of type V that
// allows multiple entries per key.
type MultiMap[K, V any] struct {
	tree *rbtree.Tree[K, V]
}

// Iterator is a bidirectional cursor over a MultiMap, visiting entries in
// ascending key order (equivalent keys appear consecutively).
type Iterator[K, V any] struct {
	tree *rbtree.Tree[K, V]
	node *rbtree.Node[K, V]
}

// New creates an empty multimap ordered by the natural ordering of K.
func New[K constraints.Ordered, V any]() *MultiMap[K, V] {
	return &MultiMap[K, V]{tree: rbtree.New[K, V](constraints.Less[K], true)}
}

// NewFunc creates an empty multimap ordered by the given key comparator.
func NewFunc[K, V any](less func(a, b K) bool) *MultiMap[K, V] {
	if less == nil {
		panic("multimap: nil comparator")
	}
	return &MultiMap[K, V]{tree: rbtree.New[K, V](less, true)}
}

// Key returns the key the iterator refers to. Undefined at the End position.
func (it Iterator[K, V]) Key() K { return it.node.Key }

// Value returns the value the iterator refers to.
func (it Iterator[K, V]) Value() V { return it.node.Value }

// SetValue overwrites the value the iterator refers to.
func (it Iterator[K, V]) SetValue(v V) { it.node.Value = v }

// Next advances the iterator to the next entry by key order.
func (it *Iterator[K, V]) Next() { it.node = it.tree.Next(it.node) }

// Prev moves the iterator to the previous entry by key order.
func (it *Iterator[K, V]) Prev() { it.node = it.tree.Prev(it.node) }

// Equal reports whether two iterators refer to the same position.
func (it Iterator[K, V]) Equal(other Iterator[K, V]) bool { return it.node == other.node }

// IsEnd reports whether the iterator is the past-the-end iterator.
func (it Iterator[K, V]) IsEnd() bool { return it.node == nil }

func (m *MultiMap[K, V]) iter(n *rbtree.Node[K, V]) Iterator[K, V] {
	return Iterator[K, V]{tree: m.tree, node: n}
}

// Size returns the number of entries (counting duplicate keys).
func (m *MultiMap[K, V]) Size() int { return m.tree.Size() }

// Empty reports whether the multimap has no entries.
func (m *MultiMap[K, V]) Empty() bool { return m.tree.Size() == 0 }

// Clear removes all entries.
func (m *MultiMap[K, V]) Clear() { m.tree.Clear() }

// Insert adds an entry mapping key to value; a multimap always stores it.
// Equivalent to std::multimap::insert.
func (m *MultiMap[K, V]) Insert(key K, value V) { m.tree.Insert(key, value) }

// Erase removes every entry with a key equivalent to key and returns how many
// were removed. Equivalent to std::multimap::erase(key).
func (m *MultiMap[K, V]) Erase(key K) int { return m.tree.EraseKey(key) }

// EraseIter removes the entry the iterator refers to and returns an iterator to
// the next entry. Equivalent to std::multimap::erase(iterator).
func (m *MultiMap[K, V]) EraseIter(it Iterator[K, V]) Iterator[K, V] {
	if it.node == nil {
		return m.End()
	}
	return m.iter(m.tree.EraseNode(it.node))
}

// Find returns an iterator to some entry with the given key, or End().
func (m *MultiMap[K, V]) Find(key K) Iterator[K, V] { return m.iter(m.tree.Find(key)) }

// Contains reports whether at least one entry has the given key.
func (m *MultiMap[K, V]) Contains(key K) bool { return m.tree.Contains(key) }

// Count returns the number of entries with the given key.
func (m *MultiMap[K, V]) Count(key K) int { return m.tree.Count(key) }

// LowerBound returns an iterator to the first entry whose key is not less than
// key, or End().
func (m *MultiMap[K, V]) LowerBound(key K) Iterator[K, V] { return m.iter(m.tree.LowerBound(key)) }

// UpperBound returns an iterator to the first entry whose key is greater than
// key, or End().
func (m *MultiMap[K, V]) UpperBound(key K) Iterator[K, V] { return m.iter(m.tree.UpperBound(key)) }

// EqualRange returns the range [LowerBound, UpperBound) of entries with the
// given key — the canonical way to iterate every value mapped to a key.
func (m *MultiMap[K, V]) EqualRange(key K) (Iterator[K, V], Iterator[K, V]) {
	return m.LowerBound(key), m.UpperBound(key)
}

// Begin returns an iterator to the entry with the smallest key (or End()).
func (m *MultiMap[K, V]) Begin() Iterator[K, V] { return m.iter(m.tree.Min()) }

// End returns the past-the-end iterator.
func (m *MultiMap[K, V]) End() Iterator[K, V] { return m.iter(nil) }

// RBegin returns an iterator to the entry with the largest key, for reverse
// traversal.
func (m *MultiMap[K, V]) RBegin() Iterator[K, V] { return m.iter(m.tree.Max()) }

// REnd returns the position before the first entry.
func (m *MultiMap[K, V]) REnd() Iterator[K, V] { return m.iter(nil) }

// Swap exchanges the contents of the multimap with another in O(1).
func (m *MultiMap[K, V]) Swap(other *MultiMap[K, V]) {
	if other == nil {
		return
	}
	m.tree, other.tree = other.tree, m.tree
}

// ForEach calls fn for each entry in ascending key order.
func (m *MultiMap[K, V]) ForEach(fn func(key K, value V)) {
	for n := m.tree.Min(); n != nil; n = m.tree.Next(n) {
		fn(n.Key, n.Value)
	}
}
