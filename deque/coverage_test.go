package deque

import (
	"reflect"
	"testing"
)

func TestEmplaceAndEdgeCases(t *testing.T) {
	d := New[int]()
	d.EmplaceBack(2)
	d.EmplaceFront(1) // [1 2]
	if !reflect.DeepEqual(d.ToSlice(), []int{1, 2}) {
		t.Errorf("emplace = %v, want [1 2]", d.ToSlice())
	}

	// ResizeWithValue with a negative size clamps to zero.
	d.ResizeWithValue(-3, 0)
	if !d.Empty() {
		t.Error("negative ResizeWithValue should empty the deque")
	}

	// ShrinkToFit on an already-tight buffer is a no-op.
	d.PushBack(1)
	d.ShrinkToFit()
	capBefore := d.Capacity()
	d.ShrinkToFit()
	if d.Capacity() != capBefore {
		t.Error("redundant ShrinkToFit should not change capacity")
	}
}
