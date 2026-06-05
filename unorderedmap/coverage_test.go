package unorderedmap

import "testing"

func TestNewWithCapacityNegativeReserveNoopSwapNil(t *testing.T) {
	if NewWithCapacity[int, int](-5).Size() != 0 {
		t.Error("NewWithCapacity(-5) should be empty")
	}
	m := New[int, int]()
	m.Set(1, 1)
	m.Set(2, 2)
	// Reserve below current size is a no-op.
	m.Reserve(1)
	if m.Size() != 2 {
		t.Error("Reserve below size mutated the map")
	}
	// Count present/absent.
	if m.Count(1) != 1 || m.Count(99) != 0 {
		t.Error("Count wrong")
	}
	// Swap with nil is a no-op.
	m.Swap(nil)
	if m.Size() != 2 {
		t.Error("Swap(nil) mutated the map")
	}
}
