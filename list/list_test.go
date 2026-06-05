package list

import (
	"math/rand"
	"reflect"
	"sort"
	"testing"
)

func asc(a, b int) bool { return a < b }

func collect[T any](l *List[T]) []T { return l.ToSlice() }

func TestPushPopEnds(t *testing.T) {
	l := New[int]()
	if !l.Empty() {
		t.Error("new list should be empty")
	}
	l.PushBack(2)
	l.PushBack(3)
	l.PushFront(1) // [1 2 3]
	if !reflect.DeepEqual(collect(l), []int{1, 2, 3}) {
		t.Errorf("got %v, want [1 2 3]", collect(l))
	}
	if l.Size() != 3 {
		t.Errorf("Size = %d, want 3", l.Size())
	}
	front, _ := l.Front()
	back, _ := l.Back()
	if front != 1 || back != 3 {
		t.Errorf("front=%d back=%d, want 1,3", front, back)
	}
	l.PopFront()
	l.PopBack()
	if !reflect.DeepEqual(collect(l), []int{2}) {
		t.Errorf("got %v, want [2]", collect(l))
	}
}

func TestPopEmpty(t *testing.T) {
	l := New[int]()
	if err := l.PopFront(); err == nil {
		t.Error("PopFront on empty should error")
	}
	if err := l.PopBack(); err == nil {
		t.Error("PopBack on empty should error")
	}
	if _, err := l.Front(); err == nil {
		t.Error("Front on empty should error")
	}
	if _, err := l.Back(); err == nil {
		t.Error("Back on empty should error")
	}
}

func TestForwardAndReverseIteration(t *testing.T) {
	l := NewFromSlice([]int{1, 2, 3, 4})
	var fwd []int
	for it := l.Begin(); !it.Equal(l.End()); it.Next() {
		fwd = append(fwd, it.Value())
	}
	if !reflect.DeepEqual(fwd, []int{1, 2, 3, 4}) {
		t.Errorf("forward = %v", fwd)
	}
	var rev []int
	for it := l.RBegin(); !it.Equal(l.REnd()); it.Prev() {
		rev = append(rev, it.Value())
	}
	if !reflect.DeepEqual(rev, []int{4, 3, 2, 1}) {
		t.Errorf("reverse = %v", rev)
	}
}

func TestInsertErase(t *testing.T) {
	l := NewFromSlice([]int{1, 2, 4})
	// Find the node with value 4 and insert 3 before it.
	it := l.Begin()
	for !it.Equal(l.End()) && it.Value() != 4 {
		it.Next()
	}
	l.Insert(it, 3)
	if !reflect.DeepEqual(collect(l), []int{1, 2, 3, 4}) {
		t.Errorf("after Insert = %v", collect(l))
	}
	// Erase the first element.
	next := l.Erase(l.Begin())
	if next.Value() != 2 {
		t.Errorf("Erase returned iterator to %d, want 2", next.Value())
	}
	if !reflect.DeepEqual(collect(l), []int{2, 3, 4}) {
		t.Errorf("after Erase = %v", collect(l))
	}
	// Erasing End is a no-op.
	l.Erase(l.End())
	if l.Size() != 3 {
		t.Error("erasing End changed size")
	}
}

func TestSetValueThroughIterator(t *testing.T) {
	l := NewFromSlice([]int{1, 2, 3})
	it := l.Begin()
	it.Next()
	it.SetValue(99)
	if !reflect.DeepEqual(collect(l), []int{1, 99, 3}) {
		t.Errorf("after SetValue = %v", collect(l))
	}
}

func TestReverse(t *testing.T) {
	l := NewFromSlice([]int{1, 2, 3, 4, 5})
	l.Reverse()
	if !reflect.DeepEqual(collect(l), []int{5, 4, 3, 2, 1}) {
		t.Errorf("after Reverse = %v", collect(l))
	}
	// Reverse preserves correct front/back wiring for further ops.
	l.PushBack(0)
	l.PushFront(6)
	if !reflect.DeepEqual(collect(l), []int{6, 5, 4, 3, 2, 1, 0}) {
		t.Errorf("after Reverse+push = %v", collect(l))
	}
	// Reverse of small lists is a no-op.
	single := NewFromSlice([]int{7})
	single.Reverse()
	if !reflect.DeepEqual(collect(single), []int{7}) {
		t.Errorf("reverse single = %v", collect(single))
	}
}

func TestSort(t *testing.T) {
	l := NewFromSlice([]int{5, 3, 1, 4, 2})
	l.Sort(asc)
	if !reflect.DeepEqual(collect(l), []int{1, 2, 3, 4, 5}) {
		t.Errorf("after Sort = %v", collect(l))
	}
	// Sorting keeps the list usable from both ends.
	front, _ := l.Front()
	back, _ := l.Back()
	if front != 1 || back != 5 {
		t.Errorf("front=%d back=%d after sort", front, back)
	}
}

func TestSortStability(t *testing.T) {
	type kv struct{ k, ord int }
	src := []kv{{1, 0}, {2, 1}, {1, 2}, {2, 3}, {1, 4}}
	l := NewFromSlice(src)
	l.Sort(func(a, b kv) bool { return a.k < b.k })
	got := collect(l)
	// Same-key elements must keep their original relative order.
	want := []kv{{1, 0}, {1, 2}, {1, 4}, {2, 1}, {2, 3}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("stable sort = %v, want %v", got, want)
	}
}

func TestMerge(t *testing.T) {
	a := NewFromSlice([]int{1, 3, 5})
	b := NewFromSlice([]int{2, 4, 6})
	a.Merge(b, asc)
	if !reflect.DeepEqual(collect(a), []int{1, 2, 3, 4, 5, 6}) {
		t.Errorf("after Merge = %v", collect(a))
	}
	if !b.Empty() {
		t.Error("merged-from list should be empty")
	}
	// Merging into an empty list.
	empty := New[int]()
	c := NewFromSlice([]int{7, 8})
	empty.Merge(c, asc)
	if !reflect.DeepEqual(collect(empty), []int{7, 8}) {
		t.Errorf("merge into empty = %v", collect(empty))
	}
}

func TestSplice(t *testing.T) {
	a := NewFromSlice([]int{1, 4})
	b := NewFromSlice([]int{2, 3})
	// Splice b before the element 4 (second position).
	it := a.Begin()
	it.Next()
	a.Splice(it, b)
	if !reflect.DeepEqual(collect(a), []int{1, 2, 3, 4}) {
		t.Errorf("after Splice = %v", collect(a))
	}
	if !b.Empty() {
		t.Error("spliced-from list should be empty")
	}
	if a.Size() != 4 {
		t.Errorf("size = %d, want 4", a.Size())
	}
}

func TestUnique(t *testing.T) {
	l := NewFromSlice([]int{1, 1, 2, 3, 3, 3, 2, 2})
	removed := l.Unique(func(a, b int) bool { return a == b })
	if !reflect.DeepEqual(collect(l), []int{1, 2, 3, 2}) {
		t.Errorf("after Unique = %v, want [1 2 3 2]", collect(l))
	}
	if removed != 4 {
		t.Errorf("removed = %d, want 4", removed)
	}
}

func TestRemoveIf(t *testing.T) {
	l := NewFromSlice([]int{1, 2, 3, 4, 5, 6})
	removed := l.RemoveIf(func(v int) bool { return v%2 == 0 })
	if !reflect.DeepEqual(collect(l), []int{1, 3, 5}) {
		t.Errorf("after RemoveIf = %v", collect(l))
	}
	if removed != 3 {
		t.Errorf("removed = %d, want 3", removed)
	}
	// Remove all.
	l.RemoveIf(func(v int) bool { return true })
	if !l.Empty() {
		t.Error("list should be empty after removing all")
	}
}

func TestSwapAndClear(t *testing.T) {
	a := NewFromSlice([]int{1, 2})
	b := NewFromSlice([]int{3, 4, 5})
	a.Swap(b)
	if !reflect.DeepEqual(collect(a), []int{3, 4, 5}) {
		t.Errorf("a after swap = %v", collect(a))
	}
	if !reflect.DeepEqual(collect(b), []int{1, 2}) {
		t.Errorf("b after swap = %v", collect(b))
	}
	a.Clear()
	if !a.Empty() {
		t.Error("list should be empty after Clear")
	}
	a.PushBack(9) // usable after clear
	if v, _ := a.Front(); v != 9 {
		t.Errorf("front after clear+push = %d", v)
	}
}

func TestAssign(t *testing.T) {
	l := New[int]()
	l.Assign(3, 7)
	if !reflect.DeepEqual(collect(l), []int{7, 7, 7}) {
		t.Errorf("Assign = %v", collect(l))
	}
	l.AssignSlice([]int{1, 2})
	if !reflect.DeepEqual(collect(l), []int{1, 2}) {
		t.Errorf("AssignSlice = %v", collect(l))
	}
}

// Randomized check that Sort matches the standard library's stable sort.
func TestSortRandom(t *testing.T) {
	rng := rand.New(rand.NewSource(99))
	for trial := 0; trial < 300; trial++ {
		n := rng.Intn(200)
		vals := make([]int, n)
		for i := range vals {
			vals[i] = rng.Intn(50)
		}
		l := NewFromSlice(vals)
		l.Sort(asc)

		want := append([]int(nil), vals...)
		sort.Ints(want)
		got := collect(l)
		if len(got) != len(want) {
			t.Fatalf("trial %d: len %d != %d", trial, len(got), len(want))
		}
		for i := range want {
			if got[i] != want[i] {
				t.Fatalf("trial %d: sorted mismatch at %d: %v vs %v", trial, i, got, want)
			}
		}
	}
}

func BenchmarkPushBack(b *testing.B) {
	l := New[int]()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l.PushBack(i)
	}
}

func BenchmarkSort(b *testing.B) {
	rng := rand.New(rand.NewSource(1))
	vals := make([]int, 10000)
	for i := range vals {
		vals[i] = rng.Int()
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		l := NewFromSlice(vals)
		b.StartTimer()
		l.Sort(asc)
	}
}
