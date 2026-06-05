package vector

import (
	"reflect"
	"testing"
)

func TestEmplaceBackAndAssignGrowth(t *testing.T) {
	v := New[int]()
	v.EmplaceBack(1)
	v.EmplaceBack(2)
	if !reflect.DeepEqual(v.ToSlice(), []int{1, 2}) {
		t.Errorf("EmplaceBack = %v, want [1 2]", v.ToSlice())
	}

	// Assign growing beyond capacity (forces reallocation).
	small := NewFromSlice([]int{1})
	small.Assign(10, 7)
	if small.Size() != 10 {
		t.Fatalf("Assign size = %d, want 10", small.Size())
	}
	for i := 0; i < 10; i++ {
		if small.Get(i) != 7 {
			t.Fatalf("Assign[%d] = %d, want 7", i, small.Get(i))
		}
	}
	// Negative count clamps to zero.
	small.Assign(-1, 0)
	if !small.Empty() {
		t.Error("Assign(-1) should empty the vector")
	}

	// AssignSlice growing beyond capacity.
	v2 := New[int]()
	v2.AssignSlice([]int{1, 2, 3, 4, 5})
	if !reflect.DeepEqual(v2.ToSlice(), []int{1, 2, 3, 4, 5}) {
		t.Errorf("AssignSlice = %v", v2.ToSlice())
	}

	// ShrinkToFit no-op when already tight.
	v2.ShrinkToFit()
	capBefore := v2.Capacity()
	v2.ShrinkToFit()
	if v2.Capacity() != capBefore {
		t.Error("redundant ShrinkToFit changed capacity")
	}

	// NewWithValue with negative count clamps to zero.
	if NewWithValue(-2, 9).Size() != 0 {
		t.Error("NewWithValue(-2) should be empty")
	}
}
