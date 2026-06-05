package bitset

import "testing"

func TestNewNegativeAndWordAlignedShifts(t *testing.T) {
	// New with a negative size clamps to zero.
	if New(-10).Size() != 0 {
		t.Error("New(-10) should have size 0")
	}

	// Word-aligned shifts (k a multiple of 64) hit the bitShift==0 fast path.
	b := New(130)
	b.Set(0)
	b.Set(5)
	b.ShiftLeft(64) // bits move up by exactly one word
	if !b.Test(64) || !b.Test(69) {
		t.Errorf("word-aligned ShiftLeft(64) failed: %s", b.ToString())
	}
	if b.Test(0) || b.Test(5) {
		t.Error("low bits should be zero after ShiftLeft(64)")
	}
	b.ShiftRight(64) // move them back down
	if !b.Test(0) || !b.Test(5) {
		t.Error("word-aligned ShiftRight(64) failed")
	}

	// Shift by zero is a no-op (returns receiver).
	before := b.ToString()
	b.ShiftLeft(0)
	b.ShiftRight(0)
	if b.ToString() != before {
		t.Error("shift by 0 should be a no-op")
	}
}

func TestExactWordMultipleSizeTrim(t *testing.T) {
	// n that is an exact multiple of 64 exercises trim's no-remainder path.
	b := New(128)
	b.SetAll()
	if b.Count() != 128 {
		t.Errorf("Count = %d, want 128", b.Count())
	}
	if !b.All() {
		t.Error("All should be true for a full 128-bit set")
	}
}
