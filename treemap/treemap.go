// Package treemap implements a sorted key-value associative container that
// mirrors the API and behaviour of C++ std::map. (It is named treemap rather
// than map because map is a reserved word in Go.)
//
// It is backed by the module's red-black tree, so entries are kept sorted by
// key and every core operation runs in O(log n):
//
//	Insert / Erase / Find / At / Ref / Contains          O(log n)
//	LowerBound / UpperBound / EqualRange                 O(log n)
//
// Because Go cannot return an assignable reference the way C++ operator[] does,
// the operator[] behaviour is split: Ref returns a *V that default-inserts and
// can be mutated in place, At does checked read access, Get is an ok-style
// lookup, and Set inserts or overwrites.
package treemap

import (
	"errors"

	"github.com/wiraphatys/gostd/constraints"
	"github.com/wiraphatys/gostd/internal/rbtree"
)

// Map is an ordered mapping from unique keys of type K to values of type V.
type Map[K, V any] struct {
	tree *rbtree.Tree[K, V]
}

// Iterator is a bidirectional cursor over a Map, visiting entries in ascending
// key order.
type Iterator[K, V any] struct {
	tree *rbtree.Tree[K, V]
	node *rbtree.Node[K, V]
}

// New creates an empty map ordered by the natural ordering of K.
func New[K constraints.Ordered, V any]() *Map[K, V] {
	return &Map[K, V]{tree: rbtree.New[K, V](constraints.Less[K], false)}
}

// NewFunc creates an empty map ordered by the given key comparator.
func NewFunc[K, V any](less func(a, b K) bool) *Map[K, V] {
	if less == nil {
		panic("treemap: nil comparator")
	}
	return &Map[K, V]{tree: rbtree.New[K, V](less, false)}
}

// Key returns the key the iterator refers to. Undefined at the End position.
func (it Iterator[K, V]) Key() K { return it.node.Key }

// Value returns the value the iterator refers to.
func (it Iterator[K, V]) Value() V { return it.node.Value }

// SetValue overwrites the value the iterator refers to (the key is unchanged).
func (it Iterator[K, V]) SetValue(v V) { it.node.Value = v }

// Next advances the iterator to the next entry by key order.
func (it *Iterator[K, V]) Next() { it.node = it.tree.Next(it.node) }

// Prev moves the iterator to the previous entry by key order.
func (it *Iterator[K, V]) Prev() { it.node = it.tree.Prev(it.node) }

// Equal reports whether two iterators refer to the same position.
func (it Iterator[K, V]) Equal(other Iterator[K, V]) bool { return it.node == other.node }

// IsEnd reports whether the iterator is the past-the-end iterator.
func (it Iterator[K, V]) IsEnd() bool { return it.node == nil }

func (m *Map[K, V]) iter(n *rbtree.Node[K, V]) Iterator[K, V] {
	return Iterator[K, V]{tree: m.tree, node: n}
}

// Size returns the number of entries. Equivalent to C++ size().
func (m *Map[K, V]) Size() int { return m.tree.Size() }

// Empty reports whether the map has no entries.
func (m *Map[K, V]) Empty() bool { return m.tree.Size() == 0 }

// Clear removes all entries.
func (m *Map[K, V]) Clear() { m.tree.Clear() }

// Insert adds key with value only if key is not already present, and reports
// whether it was inserted. Equivalent to std::map::insert.
func (m *Map[K, V]) Insert(key K, value V) bool {
	_, inserted := m.tree.Insert(key, value)
	return inserted
}

// InsertOrAssign adds key with value, overwriting any existing value, and
// reports whether a new entry was created. Equivalent to insert_or_assign.
func (m *Map[K, V]) InsertOrAssign(key K, value V) bool {
	_, inserted := m.tree.InsertOrAssign(key, value)
	return inserted
}

// Set is an alias for InsertOrAssign that ignores the inserted flag, for the
// common "just put this" use.
func (m *Map[K, V]) Set(key K, value V) { m.tree.InsertOrAssign(key, value) }

// At returns the value mapped to key, or an error if key is absent.
// Equivalent to std::map::at (which throws std::out_of_range).
func (m *Map[K, V]) At(key K) (V, error) {
	n := m.tree.Find(key)
	if n == nil {
		var zero V
		return zero, errors.New("treemap: key not found")
	}
	return n.Value, nil
}

// Get returns the value mapped to key and whether it was present, in the usual
// Go comma-ok style.
func (m *Map[K, V]) Get(key K) (V, bool) {
	n := m.tree.Find(key)
	if n == nil {
		var zero V
		return zero, false
	}
	return n.Value, true
}

// Ref returns a pointer to the value mapped to key, inserting a zero value
// first if the key is absent. Writing through the pointer updates the map in
// place, mirroring C++ std::map::operator[].
func (m *Map[K, V]) Ref(key K) *V {
	var zero V
	n, _ := m.tree.Insert(key, zero)
	return &n.Value
}

// Erase removes key if present and reports whether something was removed.
// Equivalent to std::map::erase(key).
func (m *Map[K, V]) Erase(key K) bool { return m.tree.EraseKey(key) > 0 }

// EraseIter removes the entry the iterator refers to and returns an iterator to
// the next entry. Equivalent to std::map::erase(iterator).
func (m *Map[K, V]) EraseIter(it Iterator[K, V]) Iterator[K, V] {
	if it.node == nil {
		return m.End()
	}
	return m.iter(m.tree.EraseNode(it.node))
}

// Find returns an iterator to key, or End() if absent.
func (m *Map[K, V]) Find(key K) Iterator[K, V] { return m.iter(m.tree.Find(key)) }

// Contains reports whether key is present (C++20 std::map::contains).
func (m *Map[K, V]) Contains(key K) bool { return m.tree.Contains(key) }

// Count returns 1 if key is present, else 0.
func (m *Map[K, V]) Count(key K) int {
	if m.tree.Contains(key) {
		return 1
	}
	return 0
}

// LowerBound returns an iterator to the first entry whose key is not less than
// key, or End().
func (m *Map[K, V]) LowerBound(key K) Iterator[K, V] { return m.iter(m.tree.LowerBound(key)) }

// UpperBound returns an iterator to the first entry whose key is greater than
// key, or End().
func (m *Map[K, V]) UpperBound(key K) Iterator[K, V] { return m.iter(m.tree.UpperBound(key)) }

// EqualRange returns the range [LowerBound, UpperBound) for key.
func (m *Map[K, V]) EqualRange(key K) (Iterator[K, V], Iterator[K, V]) {
	return m.LowerBound(key), m.UpperBound(key)
}

// Begin returns an iterator to the entry with the smallest key (or End()).
func (m *Map[K, V]) Begin() Iterator[K, V] { return m.iter(m.tree.Min()) }

// End returns the past-the-end iterator.
func (m *Map[K, V]) End() Iterator[K, V] { return m.iter(nil) }

// RBegin returns an iterator to the entry with the largest key, for reverse
// traversal with Prev.
func (m *Map[K, V]) RBegin() Iterator[K, V] { return m.iter(m.tree.Max()) }

// REnd returns the position before the first entry.
func (m *Map[K, V]) REnd() Iterator[K, V] { return m.iter(nil) }

// Swap exchanges the contents of the map with another in O(1).
func (m *Map[K, V]) Swap(other *Map[K, V]) {
	if other == nil {
		return
	}
	m.tree, other.tree = other.tree, m.tree
}

// Keys returns all keys in ascending order.
func (m *Map[K, V]) Keys() []K {
	out := make([]K, 0, m.tree.Size())
	for n := m.tree.Min(); n != nil; n = m.tree.Next(n) {
		out = append(out, n.Key)
	}
	return out
}

// Values returns all values ordered by their keys.
func (m *Map[K, V]) Values() []V {
	out := make([]V, 0, m.tree.Size())
	for n := m.tree.Min(); n != nil; n = m.tree.Next(n) {
		out = append(out, n.Value)
	}
	return out
}

// ForEach calls fn for each entry in ascending key order.
func (m *Map[K, V]) ForEach(fn func(key K, value V)) {
	for n := m.tree.Min(); n != nil; n = m.tree.Next(n) {
		fn(n.Key, n.Value)
	}
}
