// Package bitset implements a fixed-length sequence of bits that mirrors the
// API and behaviour of C++ std::bitset.
//
// Because Go has no const generic parameters, the length is chosen at runtime
// via New(n) rather than at compile time as in std::bitset<N>. Bits are packed
// into a []uint64, so per-bit operations are O(1) and whole-set operations
// (Count, And, Or, Xor, Not, shifts) are O(n/64).
//
// Following std::bitset, bit 0 is the least-significant bit, and ToString /
// NewFromString render the most-significant bit (index n-1) leftmost.
//
// Out-of-range bit positions are a programming error: the position-taking
// methods panic, which is the Go analogue of the std::out_of_range that
// std::bitset::set/test throw.
package bitset

import (
	"errors"
	"math/bits"
	"strings"
)

var errInvalidChar = errors.New("bitset: string must contain only '0' and '1'")

// BitSet is a fixed-length array of bits.
type BitSet struct {
	bits []uint64
	n    int
}

func numWords(n int) int { return (n + 63) >> 6 }

// New creates a bitset of n bits, all initialized to 0. A negative n is treated
// as 0.
func New(n int) *BitSet {
	if n < 0 {
		n = 0
	}
	return &BitSet{bits: make([]uint64, numWords(n)), n: n}
}

// NewFromUint64 creates a bitset of n bits initialized from the low bits of
// value (bit 0 of value becomes bit 0 of the set). Equivalent to the
// unsigned-long constructor of std::bitset.
func NewFromUint64(n int, value uint64) *BitSet {
	b := New(n)
	if len(b.bits) > 0 {
		b.bits[0] = value
		b.trim()
	}
	return b
}

// NewFromString creates a bitset from a string of '0' and '1' characters. As in
// std::bitset, the first character is the most-significant bit (index len-1).
// It returns an error if the string contains any other character.
func NewFromString(s string) (*BitSet, error) {
	b := New(len(s))
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '1':
			// s[0] is the most-significant bit -> index n-1.
			b.setBit(len(s) - 1 - i)
		case '0':
			// already 0
		default:
			return nil, errInvalidChar
		}
	}
	return b, nil
}

// trim clears any bits beyond n in the final word so that Count, comparisons
// and All stay correct.
func (b *BitSet) trim() {
	if b.n == 0 {
		return
	}
	if rem := uint(b.n & 63); rem != 0 {
		b.bits[len(b.bits)-1] &= (uint64(1) << rem) - 1
	}
}

func (b *BitSet) checkPos(pos int) {
	if pos < 0 || pos >= b.n {
		panic("bitset: position out of range")
	}
}

func (b *BitSet) setBit(pos int)   { b.bits[pos>>6] |= uint64(1) << uint(pos&63) }
func (b *BitSet) clearBit(pos int) { b.bits[pos>>6] &^= uint64(1) << uint(pos&63) }
func (b *BitSet) flipBit(pos int)  { b.bits[pos>>6] ^= uint64(1) << uint(pos&63) }
func (b *BitSet) testBit(pos int) bool {
	return b.bits[pos>>6]&(uint64(1)<<uint(pos&63)) != 0
}

// Size returns the number of bits in the set. Equivalent to C++ size().
func (b *BitSet) Size() int { return b.n }

// Set sets bit pos to 1. Panics if pos is out of range. Equivalent to set(pos).
func (b *BitSet) Set(pos int) *BitSet {
	b.checkPos(pos)
	b.setBit(pos)
	return b
}

// SetTo sets bit pos to the given value. Panics if pos is out of range.
// Equivalent to set(pos, value).
func (b *BitSet) SetTo(pos int, value bool) *BitSet {
	b.checkPos(pos)
	if value {
		b.setBit(pos)
	} else {
		b.clearBit(pos)
	}
	return b
}

// SetAll sets every bit to 1. Equivalent to set().
func (b *BitSet) SetAll() *BitSet {
	for i := range b.bits {
		b.bits[i] = ^uint64(0)
	}
	b.trim()
	return b
}

// Reset sets bit pos to 0. Panics if pos is out of range. Equivalent to
// reset(pos).
func (b *BitSet) Reset(pos int) *BitSet {
	b.checkPos(pos)
	b.clearBit(pos)
	return b
}

// ResetAll sets every bit to 0. Equivalent to reset().
func (b *BitSet) ResetAll() *BitSet {
	for i := range b.bits {
		b.bits[i] = 0
	}
	return b
}

// Flip toggles bit pos. Panics if pos is out of range. Equivalent to flip(pos).
func (b *BitSet) Flip(pos int) *BitSet {
	b.checkPos(pos)
	b.flipBit(pos)
	return b
}

// FlipAll toggles every bit. Equivalent to flip().
func (b *BitSet) FlipAll() *BitSet {
	for i := range b.bits {
		b.bits[i] = ^b.bits[i]
	}
	b.trim()
	return b
}

// Test reports whether bit pos is set. Panics if pos is out of range.
// Equivalent to test(pos).
func (b *BitSet) Test(pos int) bool {
	b.checkPos(pos)
	return b.testBit(pos)
}

// Count returns the number of bits set to 1. Equivalent to count().
func (b *BitSet) Count() int {
	c := 0
	for _, w := range b.bits {
		c += bits.OnesCount64(w)
	}
	return c
}

// Any reports whether at least one bit is set. Equivalent to any().
func (b *BitSet) Any() bool {
	for _, w := range b.bits {
		if w != 0 {
			return true
		}
	}
	return false
}

// None reports whether no bit is set. Equivalent to none().
func (b *BitSet) None() bool { return !b.Any() }

// All reports whether every bit is set. Equivalent to all().
func (b *BitSet) All() bool { return b.Count() == b.n }

func (b *BitSet) sameSize(other *BitSet) {
	if b.n != other.n {
		panic("bitset: size mismatch")
	}
}

// And replaces this bitset with the bitwise AND of itself and other (which must
// have the same size) and returns it. Equivalent to operator&=.
func (b *BitSet) And(other *BitSet) *BitSet {
	b.sameSize(other)
	for i := range b.bits {
		b.bits[i] &= other.bits[i]
	}
	return b
}

// Or replaces this bitset with the bitwise OR. Equivalent to operator|=.
func (b *BitSet) Or(other *BitSet) *BitSet {
	b.sameSize(other)
	for i := range b.bits {
		b.bits[i] |= other.bits[i]
	}
	return b
}

// Xor replaces this bitset with the bitwise XOR. Equivalent to operator^=.
func (b *BitSet) Xor(other *BitSet) *BitSet {
	b.sameSize(other)
	for i := range b.bits {
		b.bits[i] ^= other.bits[i]
	}
	return b
}

// Not flips every bit in place and returns the bitset. Equivalent to operator~
// applied in place.
func (b *BitSet) Not() *BitSet { return b.FlipAll() }

// ShiftLeft shifts all bits toward higher indices by k positions (like <<=),
// shifting in zeros. A negative k shifts right. Equivalent to operator<<=.
func (b *BitSet) ShiftLeft(k int) *BitSet {
	if k < 0 {
		return b.ShiftRight(-k)
	}
	if k == 0 {
		return b
	}
	if k >= b.n {
		return b.ResetAll()
	}
	wordShift := k >> 6
	bitShift := uint(k & 63)
	nw := len(b.bits)
	if bitShift == 0 {
		for i := nw - 1; i >= 0; i-- {
			if i-wordShift >= 0 {
				b.bits[i] = b.bits[i-wordShift]
			} else {
				b.bits[i] = 0
			}
		}
	} else {
		for i := nw - 1; i >= 0; i-- {
			var v uint64
			if src := i - wordShift; src >= 0 {
				v = b.bits[src] << bitShift
				if src-1 >= 0 {
					v |= b.bits[src-1] >> (64 - bitShift)
				}
			}
			b.bits[i] = v
		}
	}
	b.trim()
	return b
}

// ShiftRight shifts all bits toward lower indices by k positions (like >>=),
// shifting in zeros. A negative k shifts left. Equivalent to operator>>=.
func (b *BitSet) ShiftRight(k int) *BitSet {
	if k < 0 {
		return b.ShiftLeft(-k)
	}
	if k == 0 {
		return b
	}
	if k >= b.n {
		return b.ResetAll()
	}
	wordShift := k >> 6
	bitShift := uint(k & 63)
	nw := len(b.bits)
	if bitShift == 0 {
		for i := 0; i < nw; i++ {
			if i+wordShift < nw {
				b.bits[i] = b.bits[i+wordShift]
			} else {
				b.bits[i] = 0
			}
		}
	} else {
		for i := 0; i < nw; i++ {
			var v uint64
			if src := i + wordShift; src < nw {
				v = b.bits[src] >> bitShift
				if src+1 < nw {
					v |= b.bits[src+1] << (64 - bitShift)
				}
			}
			b.bits[i] = v
		}
	}
	b.trim()
	return b
}

// Equal reports whether two bitsets have the same size and the same bits.
func (b *BitSet) Equal(other *BitSet) bool {
	if other == nil || b.n != other.n {
		return false
	}
	for i := range b.bits {
		if b.bits[i] != other.bits[i] {
			return false
		}
	}
	return true
}

// Clone returns an independent copy of the bitset.
func (b *BitSet) Clone() *BitSet {
	cp := &BitSet{bits: make([]uint64, len(b.bits)), n: b.n}
	copy(cp.bits, b.bits)
	return cp
}

// ToString renders the bitset as a string of '0'/'1', most-significant bit
// (index n-1) first, matching std::bitset::to_string.
func (b *BitSet) ToString() string {
	var sb strings.Builder
	sb.Grow(b.n)
	for i := b.n - 1; i >= 0; i-- {
		if b.testBit(i) {
			sb.WriteByte('1')
		} else {
			sb.WriteByte('0')
		}
	}
	return sb.String()
}

// ToUint64 returns the low 64 bits as a uint64 and reports whether the value
// fits (i.e. no set bit at index >= 64). Equivalent to to_ullong, which throws
// std::overflow_error when the value does not fit.
func (b *BitSet) ToUint64() (uint64, bool) {
	for i := 1; i < len(b.bits); i++ {
		if b.bits[i] != 0 {
			return 0, false
		}
	}
	if len(b.bits) == 0 {
		return 0, true
	}
	return b.bits[0], true
}
