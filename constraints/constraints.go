// Package constraints defines a set of useful type constraints to be used with
// type parameters.
//
// It mirrors golang.org/x/exp/constraints (and the std cmp.Ordered introduced
// in Go 1.21), but is vendored here so that gostd can target Go 1.19 without any
// external dependencies. The ordered containers in this module (set, multiset,
// treemap, multimap, priorityqueue, ...) use constraints.Ordered to provide a
// default comparator that behaves like C++'s std::less<T>.
package constraints

// Signed is a constraint that permits any signed integer type.
// If future releases of Go add new predeclared signed integer types,
// this constraint will be modified to include them.
type Signed interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

// Unsigned is a constraint that permits any unsigned integer type.
// If future releases of Go add new predeclared unsigned integer types,
// this constraint will be modified to include them.
type Unsigned interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

// Integer is a constraint that permits any integer type.
// If future releases of Go add new predeclared integer types,
// this constraint will be modified to include them.
type Integer interface {
	Signed | Unsigned
}

// Float is a constraint that permits any floating-point type.
// If future releases of Go add new predeclared floating-point types,
// this constraint will be modified to include them.
type Float interface {
	~float32 | ~float64
}

// Complex is a constraint that permits any complex numeric type.
// If future releases of Go add new predeclared complex numeric types,
// this constraint will be modified to include them.
type Complex interface {
	~complex64 | ~complex128
}

// Number is a constraint that permits any real numeric type (no complex types).
type Number interface {
	Integer | Float
}

// Ordered is a constraint that permits any ordered type: any type that supports
// the operators < <= >= >. This is the Go equivalent of the set of types for
// which C++'s std::less<T> is defined by default.
//
// If future releases of Go add new ordered types, this constraint will be
// modified to include them.
type Ordered interface {
	Integer | Float | ~string
}

// Less reports whether a is ordered before b using the natural ordering of T.
// It provides the default "strict weak ordering" used by gostd's ordered
// containers, equivalent to C++'s std::less<T>{}(a, b).
func Less[T Ordered](a, b T) bool { return a < b }

// Greater reports whether a is ordered after b using the natural ordering of T.
// It is the default comparator equivalent to C++'s std::greater<T>{}(a, b),
// useful for building descending sets/maps or min-heaps.
func Greater[T Ordered](a, b T) bool { return a > b }

// Min returns the smaller of a and b. Equivalent to C++'s std::min.
func Min[T Ordered](a, b T) T {
	if a < b {
		return a
	}
	return b
}

// Max returns the larger of a and b. Equivalent to C++'s std::max.
func Max[T Ordered](a, b T) T {
	if a > b {
		return a
	}
	return b
}
