package forwardlist

import (
	"reflect"
	"testing"
)

func TestEmplaceAndSetValue(t *testing.T) {
	fl := New[int]()
	fl.EmplaceFront(2)
	fl.EmplaceFront(1) // [1 2]
	it := fl.Begin()
	fl.EmplaceAfter(it, 99) // [1 99 2]
	if !reflect.DeepEqual(fl.ToSlice(), []int{1, 99, 2}) {
		t.Errorf("after emplace = %v, want [1 99 2]", fl.ToSlice())
	}
	// SetValue through the iterator.
	it.SetValue(100)
	if v, _ := fl.Front(); v != 100 {
		t.Errorf("front after SetValue = %d, want 100", v)
	}
}

func TestAssignSliceAndForEach(t *testing.T) {
	fl := NewFromSlice([]int{9, 9, 9})
	fl.AssignSlice([]int{1, 2, 3})
	if !reflect.DeepEqual(fl.ToSlice(), []int{1, 2, 3}) {
		t.Errorf("AssignSlice = %v, want [1 2 3]", fl.ToSlice())
	}
	sum := 0
	fl.ForEach(func(v int) { sum += v })
	if sum != 6 {
		t.Errorf("ForEach sum = %d, want 6", sum)
	}
}

func TestForwardListGuards(t *testing.T) {
	fl := NewFromSlice([]int{1, 2, 3})

	// Swap with nil is a no-op.
	fl.Swap(nil)
	if fl.Size() != 3 {
		t.Error("Swap(nil) mutated the list")
	}

	// SpliceAfter / Merge with nil, self, or empty are no-ops.
	fl.SpliceAfter(fl.Begin(), nil)
	fl.SpliceAfter(fl.Begin(), fl)
	fl.SpliceAfter(fl.Begin(), New[int]())
	fl.Merge(nil, asc)
	fl.Merge(fl, asc)
	fl.Merge(New[int](), asc)
	if fl.Size() != 3 {
		t.Errorf("no-op splice/merge changed size to %d", fl.Size())
	}

	// Unique on a single-element list returns 0.
	single := NewFromSlice([]int{42})
	if single.Unique(func(a, b int) bool { return a == b }) != 0 {
		t.Error("Unique on single element should remove 0")
	}

	// ResizeWithValue with a negative size clamps to zero.
	fl.ResizeWithValue(-1, 0)
	if !fl.Empty() {
		t.Error("negative ResizeWithValue should empty the list")
	}
}
