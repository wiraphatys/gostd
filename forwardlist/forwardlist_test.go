package forwardlist

import (
	"math/rand"
	"reflect"
	"sort"
	"testing"
)

func asc(a, b int) bool { return a < b }

func TestPushPopFront(t *testing.T) {
	fl := New[int]()
	if !fl.Empty() {
		t.Error("new list should be empty")
	}
	fl.PushFront(3)
	fl.PushFront(2)
	fl.PushFront(1) // [1 2 3]
	if !reflect.DeepEqual(fl.ToSlice(), []int{1, 2, 3}) {
		t.Errorf("got %v, want [1 2 3]", fl.ToSlice())
	}
	front, _ := fl.Front()
	if front != 1 {
		t.Errorf("Front = %d, want 1", front)
	}
	fl.PopFront()
	if !reflect.DeepEqual(fl.ToSlice(), []int{2, 3}) {
		t.Errorf("after PopFront = %v", fl.ToSlice())
	}
}

func TestPopFrontEmpty(t *testing.T) {
	fl := New[int]()
	if err := fl.PopFront(); err == nil {
		t.Error("PopFront on empty should error")
	}
	if _, err := fl.Front(); err == nil {
		t.Error("Front on empty should error")
	}
}

func TestBeforeBeginInsertAfter(t *testing.T) {
	fl := New[int]()
	// Insert after before_begin == push at the front.
	fl.InsertAfter(fl.BeforeBegin(), 1)
	it := fl.Begin()
	fl.InsertAfter(it, 2) // after 1 -> [1 2]
	fl.InsertAfter(it, 3) // after 1 -> [1 3 2]
	if !reflect.DeepEqual(fl.ToSlice(), []int{1, 3, 2}) {
		t.Errorf("got %v, want [1 3 2]", fl.ToSlice())
	}
}

func TestEraseAfter(t *testing.T) {
	fl := NewFromSlice([]int{1, 2, 3, 4})
	// Erase the element after the first one (removes 2).
	next := fl.EraseAfter(fl.Begin())
	if next.Value() != 3 {
		t.Errorf("EraseAfter returned %d, want 3", next.Value())
	}
	if !reflect.DeepEqual(fl.ToSlice(), []int{1, 3, 4}) {
		t.Errorf("after EraseAfter = %v", fl.ToSlice())
	}
	// Erase after the last element is a no-op returning End.
	it := fl.Begin()
	it.Next()
	it.Next() // points at 4 (last)
	res := fl.EraseAfter(it)
	if !res.Equal(fl.End()) {
		t.Error("EraseAfter at last should return End")
	}
	if fl.Size() != 3 {
		t.Errorf("size = %d, want 3", fl.Size())
	}
}

func TestIteration(t *testing.T) {
	fl := NewFromSlice([]int{10, 20, 30})
	var got []int
	for it := fl.Begin(); !it.Equal(fl.End()); it.Next() {
		got = append(got, it.Value())
	}
	if !reflect.DeepEqual(got, []int{10, 20, 30}) {
		t.Errorf("iteration = %v", got)
	}
}

func TestReverse(t *testing.T) {
	fl := NewFromSlice([]int{1, 2, 3, 4})
	fl.Reverse()
	if !reflect.DeepEqual(fl.ToSlice(), []int{4, 3, 2, 1}) {
		t.Errorf("after Reverse = %v", fl.ToSlice())
	}
	// Front must still work and push must prepend correctly.
	fl.PushFront(5)
	if front, _ := fl.Front(); front != 5 {
		t.Errorf("front after reverse = %d, want 5", front)
	}
}

func TestSort(t *testing.T) {
	fl := NewFromSlice([]int{4, 1, 3, 5, 2})
	fl.Sort(asc)
	if !reflect.DeepEqual(fl.ToSlice(), []int{1, 2, 3, 4, 5}) {
		t.Errorf("after Sort = %v", fl.ToSlice())
	}
}

func TestMerge(t *testing.T) {
	a := NewFromSlice([]int{1, 3, 5})
	b := NewFromSlice([]int{2, 4})
	a.Merge(b, asc)
	if !reflect.DeepEqual(a.ToSlice(), []int{1, 2, 3, 4, 5}) {
		t.Errorf("after Merge = %v", a.ToSlice())
	}
	if !b.Empty() {
		t.Error("merged-from list should be empty")
	}
}

func TestUnique(t *testing.T) {
	fl := NewFromSlice([]int{1, 1, 1, 2, 2, 3, 1})
	removed := fl.Unique(func(a, b int) bool { return a == b })
	if !reflect.DeepEqual(fl.ToSlice(), []int{1, 2, 3, 1}) {
		t.Errorf("after Unique = %v", fl.ToSlice())
	}
	if removed != 3 {
		t.Errorf("removed = %d, want 3", removed)
	}
}

func TestRemoveIf(t *testing.T) {
	fl := NewFromSlice([]int{1, 2, 3, 4, 5})
	removed := fl.RemoveIf(func(v int) bool { return v > 3 })
	if !reflect.DeepEqual(fl.ToSlice(), []int{1, 2, 3}) {
		t.Errorf("after RemoveIf = %v", fl.ToSlice())
	}
	if removed != 2 {
		t.Errorf("removed = %d, want 2", removed)
	}
	// Removing the head element too.
	fl.RemoveIf(func(v int) bool { return v == 1 })
	if !reflect.DeepEqual(fl.ToSlice(), []int{2, 3}) {
		t.Errorf("after head RemoveIf = %v", fl.ToSlice())
	}
}

func TestSpliceAfter(t *testing.T) {
	a := NewFromSlice([]int{1, 4})
	b := NewFromSlice([]int{2, 3})
	a.SpliceAfter(a.Begin(), b) // after 1 -> [1 2 3 4]
	if !reflect.DeepEqual(a.ToSlice(), []int{1, 2, 3, 4}) {
		t.Errorf("after SpliceAfter = %v", a.ToSlice())
	}
	if !b.Empty() {
		t.Error("spliced-from list should be empty")
	}
}

func TestResize(t *testing.T) {
	fl := NewFromSlice([]int{1, 2, 3})
	fl.Resize(5)
	if !reflect.DeepEqual(fl.ToSlice(), []int{1, 2, 3, 0, 0}) {
		t.Errorf("grow Resize = %v", fl.ToSlice())
	}
	fl.Resize(2)
	if !reflect.DeepEqual(fl.ToSlice(), []int{1, 2}) {
		t.Errorf("shrink Resize = %v", fl.ToSlice())
	}
	fl.ResizeWithValue(4, 9)
	if !reflect.DeepEqual(fl.ToSlice(), []int{1, 2, 9, 9}) {
		t.Errorf("ResizeWithValue = %v", fl.ToSlice())
	}
	fl.Resize(0)
	if !fl.Empty() {
		t.Error("Resize(0) should empty the list")
	}
}

func TestSwapClearAssign(t *testing.T) {
	a := NewFromSlice([]int{1, 2})
	b := NewFromSlice([]int{3, 4, 5})
	a.Swap(b)
	if !reflect.DeepEqual(a.ToSlice(), []int{3, 4, 5}) {
		t.Errorf("a after swap = %v", a.ToSlice())
	}
	a.Clear()
	if !a.Empty() {
		t.Error("should be empty after Clear")
	}
	a.Assign(2, 8)
	if !reflect.DeepEqual(a.ToSlice(), []int{8, 8}) {
		t.Errorf("Assign = %v", a.ToSlice())
	}
}

func TestSortRandom(t *testing.T) {
	rng := rand.New(rand.NewSource(13))
	for trial := 0; trial < 300; trial++ {
		n := rng.Intn(200)
		vals := make([]int, n)
		for i := range vals {
			vals[i] = rng.Intn(50)
		}
		fl := NewFromSlice(vals)
		fl.Sort(asc)
		want := append([]int(nil), vals...)
		sort.Ints(want)
		got := fl.ToSlice()
		if len(got) != len(want) {
			t.Fatalf("trial %d: len %d != %d", trial, len(got), len(want))
		}
		for i := range want {
			if got[i] != want[i] {
				t.Fatalf("trial %d: mismatch at %d", trial, i)
			}
		}
	}
}

func BenchmarkPushFront(b *testing.B) {
	fl := New[int]()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fl.PushFront(i)
	}
}
