package bitset

import (
	"math/rand"
	"testing"
)

func TestBasicSetTestReset(t *testing.T) {
	b := New(10)
	if b.Size() != 10 {
		t.Errorf("Size = %d, want 10", b.Size())
	}
	if b.Any() {
		t.Error("new bitset should have no bits set")
	}
	b.Set(3)
	b.Set(7)
	if !b.Test(3) || !b.Test(7) {
		t.Error("Set/Test failed")
	}
	if b.Test(4) {
		t.Error("bit 4 should be unset")
	}
	if b.Count() != 2 {
		t.Errorf("Count = %d, want 2", b.Count())
	}
	b.Reset(3)
	if b.Test(3) {
		t.Error("Reset failed")
	}
	b.Flip(7) // 7 was set -> unset
	b.Flip(2) // 2 was unset -> set
	if b.Test(7) || !b.Test(2) {
		t.Error("Flip failed")
	}
	b.SetTo(5, true)
	b.SetTo(2, false)
	if !b.Test(5) || b.Test(2) {
		t.Error("SetTo failed")
	}
}

func TestOutOfRangePanics(t *testing.T) {
	b := New(8)
	for _, f := range []struct {
		name string
		fn   func()
	}{
		{"Set", func() { b.Set(8) }},
		{"Set-neg", func() { b.Set(-1) }},
		{"Test", func() { b.Test(8) }},
		{"Reset", func() { b.Reset(100) }},
		{"Flip", func() { b.Flip(8) }},
		{"SetTo", func() { b.SetTo(8, true) }},
	} {
		func() {
			defer func() {
				if recover() == nil {
					t.Errorf("%s out of range should panic", f.name)
				}
			}()
			f.fn()
		}()
	}
}

func TestAnyNoneAll(t *testing.T) {
	b := New(5)
	if !b.None() || b.Any() || b.All() {
		t.Error("empty bitset: None=true, Any=false, All=false expected")
	}
	b.SetAll()
	if !b.All() || !b.Any() || b.None() {
		t.Error("full bitset: All=true, Any=true, None=false expected")
	}
	if b.Count() != 5 {
		t.Errorf("Count after SetAll = %d, want 5", b.Count())
	}
	b.ResetAll()
	if !b.None() {
		t.Error("ResetAll should clear all bits")
	}
}

// SetAll must not leak bits beyond n in the last partial word.
func TestSetAllTrimAcrossWords(t *testing.T) {
	for _, n := range []int{1, 63, 64, 65, 100, 128, 129} {
		b := New(n)
		b.SetAll()
		if b.Count() != n {
			t.Errorf("n=%d: Count after SetAll = %d, want %d", n, b.Count(), n)
		}
		if !b.All() {
			t.Errorf("n=%d: All should be true after SetAll", n)
		}
		b.FlipAll()
		if b.Count() != 0 {
			t.Errorf("n=%d: Count after FlipAll(full) = %d, want 0", n, b.Count())
		}
	}
}

func TestBitwiseOps(t *testing.T) {
	a, _ := NewFromString("1100")
	c, _ := NewFromString("1010")

	andRes := a.Clone().And(c)
	if andRes.ToString() != "1000" {
		t.Errorf("AND = %s, want 1000", andRes.ToString())
	}
	orRes := a.Clone().Or(c)
	if orRes.ToString() != "1110" {
		t.Errorf("OR = %s, want 1110", orRes.ToString())
	}
	xorRes := a.Clone().Xor(c)
	if xorRes.ToString() != "0110" {
		t.Errorf("XOR = %s, want 0110", xorRes.ToString())
	}
	notRes := a.Clone().Not()
	if notRes.ToString() != "0011" {
		t.Errorf("NOT = %s, want 0011", notRes.ToString())
	}
	// Original unchanged (Clone independence).
	if a.ToString() != "1100" {
		t.Errorf("original mutated: %s", a.ToString())
	}
}

func TestSizeMismatchPanics(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Error("And with mismatched size should panic")
		}
	}()
	New(4).And(New(8))
}

func TestToStringMSBFirst(t *testing.T) {
	b := New(8)
	b.Set(0) // least-significant -> rightmost
	b.Set(7) // most-significant -> leftmost
	if got := b.ToString(); got != "10000001" {
		t.Errorf("ToString = %s, want 10000001", got)
	}
}

func TestFromStringRoundTrip(t *testing.T) {
	for _, s := range []string{"0", "1", "1010", "1111000011110000", "1"} {
		b, err := NewFromString(s)
		if err != nil {
			t.Fatalf("NewFromString(%q) errored: %v", s, err)
		}
		if got := b.ToString(); got != s {
			t.Errorf("round trip %q -> %q", s, got)
		}
	}
	if _, err := NewFromString("10201"); err == nil {
		t.Error("invalid char should error")
	}
}

func TestUint64Conversions(t *testing.T) {
	b := NewFromUint64(16, 0xABCD)
	if v, ok := b.ToUint64(); !ok || v != 0xABCD {
		t.Errorf("ToUint64 = (%x, %v), want (abcd, true)", v, ok)
	}
	// A value with bits above 64 does not fit.
	big := New(100)
	big.Set(70)
	if _, ok := big.ToUint64(); ok {
		t.Error("ToUint64 should report not-ok when a high bit is set")
	}
	// NewFromUint64 truncates to n bits.
	small := NewFromUint64(4, 0xFF)
	if small.Count() != 4 {
		t.Errorf("NewFromUint64(4, 0xFF) Count = %d, want 4", small.Count())
	}
}

func TestEqualAndClone(t *testing.T) {
	a, _ := NewFromString("10110")
	if !a.Equal(a.Clone()) {
		t.Error("clone should be equal")
	}
	b, _ := NewFromString("10111")
	if a.Equal(b) {
		t.Error("different bits should not be equal")
	}
	if a.Equal(New(4)) {
		t.Error("different sizes should not be equal")
	}
	if a.Equal(nil) {
		t.Error("Equal(nil) should be false")
	}
	// Mutating the clone must not affect the original.
	cl := a.Clone()
	cl.Set(0)
	if a.Test(0) == cl.Test(0) && !a.Test(0) {
		// a.Test(0) was 0 ("10110" -> bit0 = 0); cl.Set(0) makes it 1
		t.Error("clone is not independent")
	}
}

// Randomized verification of shifts against a []bool reference model.
func TestShiftsAgainstModel(t *testing.T) {
	rng := rand.New(rand.NewSource(5))
	for trial := 0; trial < 500; trial++ {
		n := 1 + rng.Intn(200)
		b := New(n)
		model := make([]bool, n)
		// Random initial bits.
		for i := 0; i < n; i++ {
			if rng.Intn(2) == 0 {
				b.Set(i)
				model[i] = true
			}
		}
		k := rng.Intn(n + 5)
		left := rng.Intn(2) == 0

		if left {
			b.ShiftLeft(k)
			model = shiftLeftModel(model, k)
		} else {
			b.ShiftRight(k)
			model = shiftRightModel(model, k)
		}
		for i := 0; i < n; i++ {
			if b.Test(i) != model[i] {
				t.Fatalf("trial %d (left=%v k=%d n=%d): bit %d = %v, want %v",
					trial, left, k, n, i, b.Test(i), model[i])
			}
		}
	}
}

func shiftLeftModel(m []bool, k int) []bool {
	n := len(m)
	out := make([]bool, n)
	for i := 0; i < n; i++ {
		if i-k >= 0 {
			out[i] = m[i-k]
		}
	}
	return out
}

func shiftRightModel(m []bool, k int) []bool {
	n := len(m)
	out := make([]bool, n)
	for i := 0; i < n; i++ {
		if i+k < n {
			out[i] = m[i+k]
		}
	}
	return out
}

func TestZeroLength(t *testing.T) {
	b := New(0)
	if b.Size() != 0 || b.Any() || !b.None() || !b.All() {
		t.Error("zero-length bitset edge cases wrong")
	}
	if b.ToString() != "" {
		t.Errorf("ToString = %q, want empty", b.ToString())
	}
	if v, ok := b.ToUint64(); !ok || v != 0 {
		t.Errorf("ToUint64 = (%d,%v), want (0,true)", v, ok)
	}
}

func BenchmarkSet(b *testing.B) {
	bs := New(1024)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bs.Set(i % 1024)
	}
}

func BenchmarkCount(b *testing.B) {
	bs := New(100000)
	bs.SetAll()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bs.Count()
	}
}

func BenchmarkAnd(b *testing.B) {
	x := New(100000)
	x.SetAll()
	y := New(100000)
	y.SetAll()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		x.And(y)
	}
}
