// Package unorderedmap implements a hash-based key-value container that mirrors
// the API and behaviour of C++ std::unordered_map.
//
// It is backed by Go's built-in map, giving average O(1) Insert, Erase, At and
// Get. As in C++, iteration order is unspecified, and the key type must be
// comparable (Go's analogue of Hash + KeyEqual).
//
// Note on operator[]: C++ operator[] returns an assignable reference, but Go map
// values are not addressable, so there is no pointer-returning Ref here (unlike
// treemap, whose tree nodes are addressable). Use Set to insert/overwrite and
// GetOr for the read-with-default pattern. If you need stable, mutable
// references to values, use treemap.
package unorderedmap

import "errors"

var errKeyNotFound = errors.New("unorderedmap: key not found")

// Map is a hash map from comparable keys of type K to values of type V.
type Map[K comparable, V any] struct {
	m map[K]V
}

// New creates and returns a new empty unordered map.
func New[K comparable, V any]() *Map[K, V] {
	return &Map[K, V]{m: make(map[K]V)}
}

// NewWithCapacity creates an empty unordered map pre-sized for about n entries.
func NewWithCapacity[K comparable, V any](n int) *Map[K, V] {
	if n < 0 {
		n = 0
	}
	return &Map[K, V]{m: make(map[K]V, n)}
}

// Size returns the number of entries. Equivalent to C++ size().
func (m *Map[K, V]) Size() int { return len(m.m) }

// Empty reports whether the map has no entries.
func (m *Map[K, V]) Empty() bool { return len(m.m) == 0 }

// Insert adds key with value only if key is absent, and reports whether it was
// inserted. Equivalent to std::unordered_map::insert.
func (m *Map[K, V]) Insert(key K, value V) bool {
	if _, ok := m.m[key]; ok {
		return false
	}
	m.m[key] = value
	return true
}

// InsertOrAssign adds key with value, overwriting any existing value, and
// reports whether a new entry was created. Equivalent to insert_or_assign.
func (m *Map[K, V]) InsertOrAssign(key K, value V) bool {
	_, existed := m.m[key]
	m.m[key] = value
	return !existed
}

// Set inserts or overwrites key with value, ignoring the inserted flag.
func (m *Map[K, V]) Set(key K, value V) { m.m[key] = value }

// At returns the value mapped to key, or an error if key is absent.
// Equivalent to std::unordered_map::at.
func (m *Map[K, V]) At(key K) (V, error) {
	v, ok := m.m[key]
	if !ok {
		var zero V
		return zero, errKeyNotFound
	}
	return v, nil
}

// Get returns the value mapped to key and whether it was present (comma-ok).
func (m *Map[K, V]) Get(key K) (V, bool) {
	v, ok := m.m[key]
	return v, ok
}

// GetOr returns the value mapped to key, or def if key is absent. The map is
// not modified. This is the idiomatic replacement for the read side of C++
// operator[].
func (m *Map[K, V]) GetOr(key K, def V) V {
	if v, ok := m.m[key]; ok {
		return v
	}
	return def
}

// Erase removes key if present and reports whether something was removed.
func (m *Map[K, V]) Erase(key K) bool {
	if _, ok := m.m[key]; !ok {
		return false
	}
	delete(m.m, key)
	return true
}

// Contains reports whether key is present (C++20 contains).
func (m *Map[K, V]) Contains(key K) bool {
	_, ok := m.m[key]
	return ok
}

// Count returns 1 if key is present, else 0.
func (m *Map[K, V]) Count(key K) int {
	if _, ok := m.m[key]; ok {
		return 1
	}
	return 0
}

// Clear removes all entries.
func (m *Map[K, V]) Clear() { m.m = make(map[K]V) }

// Reserve ensures the table is sized for at least n entries, rehashing existing
// entries into a larger table when needed.
func (m *Map[K, V]) Reserve(n int) {
	if n <= len(m.m) {
		return
	}
	nm := make(map[K]V, n)
	for k, v := range m.m {
		nm[k] = v
	}
	m.m = nm
}

// Swap exchanges the contents of the map with another in O(1).
func (m *Map[K, V]) Swap(other *Map[K, V]) {
	if other == nil {
		return
	}
	m.m, other.m = other.m, m.m
}

// Keys returns all keys in unspecified order.
func (m *Map[K, V]) Keys() []K {
	out := make([]K, 0, len(m.m))
	for k := range m.m {
		out = append(out, k)
	}
	return out
}

// Values returns all values in unspecified order.
func (m *Map[K, V]) Values() []V {
	out := make([]V, 0, len(m.m))
	for _, v := range m.m {
		out = append(out, v)
	}
	return out
}

// ForEach calls fn for each entry in unspecified order.
func (m *Map[K, V]) ForEach(fn func(key K, value V)) {
	for k, v := range m.m {
		fn(k, v)
	}
}
