package rbtree

import "testing"

func TestClearContainsNilGuards(t *testing.T) {
	tr := newUnique()
	for i := 0; i < 5; i++ {
		tr.Insert(i, i)
	}
	if !tr.Contains(3) || tr.Contains(99) {
		t.Error("Contains wrong")
	}
	tr.Clear()
	if tr.Size() != 0 {
		t.Errorf("Size after Clear = %d, want 0", tr.Size())
	}
	if tr.Contains(3) {
		t.Error("Contains should be false after Clear")
	}

	// Nil/sentinel guards on the public stepping API.
	if tr.Next(nil) != nil {
		t.Error("Next(nil) should be nil")
	}
	if tr.Prev(nil) != nil {
		t.Error("Prev(nil) should be nil")
	}
	if tr.EraseNode(nil) != nil {
		t.Error("EraseNode(nil) should be nil")
	}
}

func TestCountUniqueZeroAndMultiBreak(t *testing.T) {
	// Unique-mode Count returns 0 for an absent key (the early-return path).
	u := newUnique()
	u.Insert(1, 1)
	if u.Count(2) != 0 {
		t.Error("unique Count(absent) should be 0")
	}
	if u.Count(1) != 1 {
		t.Error("unique Count(present) should be 1")
	}

	// Multi-mode Count must stop at the first key greater than the target
	// (the loop-break branch), not run to the end.
	m := newMulti()
	for _, k := range []int{1, 2, 2, 2, 3, 4} {
		m.Insert(k, 0)
	}
	if got := m.Count(2); got != 3 {
		t.Errorf("multi Count(2) = %d, want 3", got)
	}
	if got := m.Count(5); got != 0 {
		t.Errorf("multi Count(5) = %d, want 0", got)
	}
}
