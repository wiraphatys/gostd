package priorityqueue

import "testing"

func TestNewFromSliceFunc(t *testing.T) {
	type item struct{ p int }
	// Highest p on top.
	pq := NewFromSliceFunc([]item{{3}, {1}, {5}, {2}}, func(a, b item) bool { return a.p < b.p })
	if pq.Size() != 4 {
		t.Fatalf("Size = %d, want 4", pq.Size())
	}
	var order []int
	for !pq.Empty() {
		top, _ := pq.Top()
		order = append(order, top.p)
		pq.Pop()
	}
	want := []int{5, 3, 2, 1}
	for i := range want {
		if order[i] != want[i] {
			t.Fatalf("order = %v, want %v", order, want)
		}
	}
}

func TestNewFromSliceFuncNilPanics(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Error("NewFromSliceFunc with nil comparator should panic")
		}
	}()
	NewFromSliceFunc([]int{1}, nil)
}

func TestToSlice(t *testing.T) {
	pq := New[int]()
	pq.Push(1)
	pq.Push(2)
	pq.Push(3)
	s := pq.ToSlice()
	if len(s) != 3 {
		t.Fatalf("ToSlice len = %d, want 3", len(s))
	}
	// ToSlice returns heap order; the max must be at the root (index 0).
	if s[0] != 3 {
		t.Errorf("ToSlice[0] = %d, want 3 (heap root)", s[0])
	}
	// Mutating the returned slice must not affect the queue.
	s[0] = 999
	top, _ := pq.Top()
	if top != 3 {
		t.Error("ToSlice should return an independent copy")
	}
}
