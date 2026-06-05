package unorderedset

import "testing"

func TestNewWithCapacityNegativeAndReserveNoop(t *testing.T) {
	if NewWithCapacity[int](-5).Size() != 0 {
		t.Error("NewWithCapacity(-5) should be empty")
	}
	s := NewFromSlice([]int{1, 2, 3})
	// Reserving below current size is a no-op.
	s.Reserve(1)
	if s.Size() != 3 {
		t.Error("Reserve below size mutated the set")
	}
	if s.Count(2) != 1 {
		t.Error("Reserve no-op lost an element")
	}
}
