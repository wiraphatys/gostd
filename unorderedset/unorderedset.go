// Package unorderedset implements a hash-based set of unique keys that mirrors
// the API and behaviour of C++ std::unordered_set.
//
// It is backed by Go's built-in map, which is a highly optimized open-hashing
// table, giving average O(1) Insert, Erase and Contains. As in C++, iteration
// order is unspecified. The element type must be comparable, which is Go's
// analogue of C++'s Hash + KeyEqual requirement.
//
// Bucket-level operations (bucket_count, load_factor, rehash to a specific
// bucket count, ...) are intentionally omitted because Go's runtime map does
// not expose its bucket structure; Reserve covers the practical need to
// pre-size the table.
package unorderedset

// Set is a collection of unique comparable keys with average O(1) operations.
type Set[T comparable] struct {
	m map[T]struct{}
}

// New creates and returns a new empty unordered set.
func New[T comparable]() *Set[T] {
	return &Set[T]{m: make(map[T]struct{})}
}

// NewWithCapacity creates an empty unordered set pre-sized for about n elements,
// reducing rehashing as it fills. Equivalent in spirit to reserve at
// construction.
func NewWithCapacity[T comparable](n int) *Set[T] {
	if n < 0 {
		n = 0
	}
	return &Set[T]{m: make(map[T]struct{}, n)}
}

// NewFromSlice creates an unordered set containing the unique elements of items.
func NewFromSlice[T comparable](items []T) *Set[T] {
	s := &Set[T]{m: make(map[T]struct{}, len(items))}
	for _, v := range items {
		s.m[v] = struct{}{}
	}
	return s
}

// Size returns the number of elements. Equivalent to C++ size().
func (s *Set[T]) Size() int { return len(s.m) }

// Empty reports whether the set has no elements.
func (s *Set[T]) Empty() bool { return len(s.m) == 0 }

// Insert adds value and reports whether it was newly inserted (false if it was
// already present). Equivalent to std::unordered_set::insert.
func (s *Set[T]) Insert(value T) bool {
	if _, ok := s.m[value]; ok {
		return false
	}
	s.m[value] = struct{}{}
	return true
}

// Erase removes value if present and reports whether something was removed.
// Equivalent to std::unordered_set::erase(key).
func (s *Set[T]) Erase(value T) bool {
	if _, ok := s.m[value]; !ok {
		return false
	}
	delete(s.m, value)
	return true
}

// Contains reports whether value is present (C++20 contains).
func (s *Set[T]) Contains(value T) bool {
	_, ok := s.m[value]
	return ok
}

// Count returns 1 if value is present, else 0.
func (s *Set[T]) Count(value T) int {
	if _, ok := s.m[value]; ok {
		return 1
	}
	return 0
}

// Clear removes all elements.
func (s *Set[T]) Clear() { s.m = make(map[T]struct{}) }

// Reserve ensures the table is sized for at least n elements. If the set is
// empty it pre-sizes a fresh table; otherwise it rehashes the existing
// elements into a table with the larger hint. Equivalent to reserve/rehash.
func (s *Set[T]) Reserve(n int) {
	if n <= len(s.m) {
		return
	}
	nm := make(map[T]struct{}, n)
	for k := range s.m {
		nm[k] = struct{}{}
	}
	s.m = nm
}

// Swap exchanges the contents of the set with another in O(1).
func (s *Set[T]) Swap(other *Set[T]) {
	if other == nil {
		return
	}
	s.m, other.m = other.m, s.m
}

// ToSlice returns all elements as a new slice in unspecified order.
func (s *Set[T]) ToSlice() []T {
	out := make([]T, 0, len(s.m))
	for k := range s.m {
		out = append(out, k)
	}
	return out
}

// ForEach calls fn for each element in unspecified order.
func (s *Set[T]) ForEach(fn func(value T)) {
	for k := range s.m {
		fn(k)
	}
}
