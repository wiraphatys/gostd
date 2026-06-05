package stack

import "testing"

// Pop must zero the vacated backing slot so reference-typed elements can be
// garbage collected rather than being pinned by the underlying array.
func TestPopZeroesSlotForGC(t *testing.T) {
	s := New[*int]()
	a, b := 1, 2
	s.Push(&a)
	s.Push(&b)
	s.Pop() // removes &b

	// The popped element now lives beyond len but within cap; it must be nil.
	backing := s.items[:cap(s.items)]
	if backing[1] != nil {
		t.Error("Pop did not zero the removed slot (reference would leak)")
	}
}

// Clear must zero every slot, not just reset the length.
func TestClearZeroesSlotsForGC(t *testing.T) {
	s := New[*int]()
	a, b, c := 1, 2, 3
	s.Push(&a)
	s.Push(&b)
	s.Push(&c)
	s.Clear()

	backing := s.items[:cap(s.items)]
	for i := 0; i < 3; i++ {
		if backing[i] != nil {
			t.Errorf("Clear did not zero slot %d (reference would leak)", i)
		}
	}
}
