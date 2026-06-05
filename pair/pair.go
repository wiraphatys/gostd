// Package pair provides a generic two-element tuple that mirrors C++ std::pair.
//
// A Pair couples together a pair of values, which may be of different types
// (T1 and T2). The individual values are accessed through the exported fields
// First and Second, exactly like std::pair::first and std::pair::second.
package pair

import "github.com/wiraphatys/gostd/constraints"

// Pair holds two heterogeneously typed values. It is the Go equivalent of
// C++ std::pair<T1, T2>.
//
// The fields are exported so a Pair can be created with a composite literal,
// used as a map key (when T1 and T2 are comparable), and compared with ==.
type Pair[T1, T2 any] struct {
	First  T1
	Second T2
}

// New constructs a Pair from the two given values.
// Equivalent to C++ std::make_pair, but the value type is explicit.
func New[T1, T2 any](first T1, second T2) Pair[T1, T2] {
	return Pair[T1, T2]{First: first, Second: second}
}

// MakePair is an alias for New, named to match C++ std::make_pair.
func MakePair[T1, T2 any](first T1, second T2) Pair[T1, T2] {
	return Pair[T1, T2]{First: first, Second: second}
}

// Unpack returns the two members of the pair, enabling Go-style
// multiple-assignment: a, b := p.Unpack().
func (p Pair[T1, T2]) Unpack() (T1, T2) {
	return p.First, p.Second
}

// Swap exchanges the contents of two pairs of the same type.
// Equivalent to C++ std::pair::swap.
func (p *Pair[T1, T2]) Swap(other *Pair[T1, T2]) {
	if other == nil {
		return
	}
	p.First, other.First = other.First, p.First
	p.Second, other.Second = other.Second, p.Second
}

// Equal reports whether two comparable pairs are equal, comparing First to
// First and Second to Second. Equivalent to C++ operator== for std::pair.
func Equal[T1, T2 comparable](a, b Pair[T1, T2]) bool {
	return a.First == b.First && a.Second == b.Second
}

// Less reports whether a is lexicographically ordered before b, comparing
// First members first and Second members only on a tie. Equivalent to C++
// operator< for std::pair.
func Less[T1, T2 constraints.Ordered](a, b Pair[T1, T2]) bool {
	if a.First != b.First {
		return a.First < b.First
	}
	return a.Second < b.Second
}
