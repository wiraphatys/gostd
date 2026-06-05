package multiset

import (
	"reflect"
	"testing"
)

func TestStateAndTraversalMethods(t *testing.T) {
	m := New[int]()
	if !m.Empty() {
		t.Error("new multiset should be Empty")
	}
	if !m.Begin().Equal(m.End()) {
		t.Error("empty Begin should equal End")
	}

	for _, v := range []int{3, 1, 2, 1} {
		m.Insert(v)
	}
	if !m.Contains(1) || m.Contains(99) {
		t.Error("Contains wrong")
	}

	// Begin points at the smallest element.
	if m.Begin().Value() != 1 {
		t.Errorf("Begin = %d, want 1", m.Begin().Value())
	}

	// ForEach visits in sorted order with duplicates.
	var got []int
	m.ForEach(func(v int) { got = append(got, v) })
	if !reflect.DeepEqual(got, []int{1, 1, 2, 3}) {
		t.Errorf("ForEach = %v, want [1 1 2 3]", got)
	}

	m.Clear()
	if !m.Empty() {
		t.Error("Clear should empty the multiset")
	}
}

func TestGuards(t *testing.T) {
	// NewFunc(nil) panics.
	func() {
		defer func() {
			if recover() == nil {
				t.Error("NewFunc(nil) should panic")
			}
		}()
		NewFunc[int](nil)
	}()

	m := NewFromSlice([]int{1, 2})
	// EraseIter(End) returns End and does nothing.
	if !m.EraseIter(m.End()).IsEnd() {
		t.Error("EraseIter(End) should return End")
	}
	if m.Size() != 2 {
		t.Error("EraseIter(End) mutated the multiset")
	}
	// Swap with nil is a no-op.
	m.Swap(nil)
	if m.Size() != 2 {
		t.Error("Swap(nil) mutated the multiset")
	}
}
