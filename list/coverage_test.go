package list

import (
	"reflect"
	"testing"
)

func TestEmplaceForEachAndGuards(t *testing.T) {
	l := New[int]()
	l.EmplaceBack(2)
	l.EmplaceFront(1)
	l.EmplaceBack(3) // [1 2 3]

	var fwd []int
	l.ForEach(func(v int) { fwd = append(fwd, v) })
	if !reflect.DeepEqual(fwd, []int{1, 2, 3}) {
		t.Errorf("ForEach = %v, want [1 2 3]", fwd)
	}
	var rev []int
	l.ForEachReverse(func(v int) { rev = append(rev, v) })
	if !reflect.DeepEqual(rev, []int{3, 2, 1}) {
		t.Errorf("ForEachReverse = %v, want [3 2 1]", rev)
	}

	// Swap with nil is a no-op.
	l.Swap(nil)
	if l.Size() != 3 {
		t.Error("Swap(nil) mutated the list")
	}

	// Splice/Merge with nil, self, or empty are no-ops.
	l.Splice(l.Begin(), nil)
	l.Splice(l.Begin(), l)
	l.Splice(l.Begin(), New[int]())
	l.Merge(nil, asc)
	l.Merge(l, asc)
	l.Merge(New[int](), asc)
	if l.Size() != 3 {
		t.Errorf("no-op splice/merge changed size to %d", l.Size())
	}
}
