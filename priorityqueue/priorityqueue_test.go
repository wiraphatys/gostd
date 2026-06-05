package priorityqueue

import (
	"math/rand"
	"sort"
	"testing"
)

// drain pops every element and returns them in pop order.
func drain[T any](pq *PriorityQueue[T]) []T {
	out := make([]T, 0, pq.Size())
	for !pq.Empty() {
		top, err := pq.Top()
		if err != nil {
			break
		}
		out = append(out, top)
		pq.Pop()
	}
	return out
}

func TestNewMaxHeapOrder(t *testing.T) {
	pq := New[int]()
	if !pq.Empty() {
		t.Error("new queue should be empty")
	}
	for _, x := range []int{3, 1, 4, 1, 5, 9, 2, 6} {
		pq.Push(x)
	}
	got := drain(pq)
	want := []int{9, 6, 5, 4, 3, 2, 1, 1}
	if !equalInts(got, want) {
		t.Errorf("max-heap pop order = %v, want %v", got, want)
	}
}

func TestNewMinHeapOrder(t *testing.T) {
	pq := NewMin[int]()
	for _, x := range []int{3, 1, 4, 1, 5, 9, 2, 6} {
		pq.Push(x)
	}
	got := drain(pq)
	want := []int{1, 1, 2, 3, 4, 5, 6, 9}
	if !equalInts(got, want) {
		t.Errorf("min-heap pop order = %v, want %v", got, want)
	}
}

func TestTopAndSize(t *testing.T) {
	pq := New[int]()
	if _, err := pq.Top(); err == nil {
		t.Error("Top on empty should error")
	}
	pq.Push(5)
	pq.Push(10)
	pq.Push(3)
	if pq.Size() != 3 {
		t.Errorf("Size = %d, want 3", pq.Size())
	}
	top, _ := pq.Top()
	if top != 10 {
		t.Errorf("Top = %d, want 10", top)
	}
	if pq.Size() != 3 {
		t.Error("Top must not change size")
	}
}

func TestPopEmpty(t *testing.T) {
	pq := New[int]()
	if err := pq.Pop(); err == nil {
		t.Error("Pop on empty should error")
	}
}

func TestEmplace(t *testing.T) {
	pq := New[int]()
	pq.Emplace(7)
	top, _ := pq.Top()
	if top != 7 {
		t.Errorf("Top = %d, want 7", top)
	}
}

func TestNewFromSliceHeapify(t *testing.T) {
	src := []int{5, 2, 8, 1, 9, 3, 7, 4, 6}
	pq := NewFromSlice(src)
	if pq.Size() != len(src) {
		t.Errorf("Size = %d, want %d", pq.Size(), len(src))
	}
	got := drain(pq)
	want := append([]int(nil), src...)
	sort.Sort(sort.Reverse(sort.IntSlice(want)))
	if !equalInts(got, want) {
		t.Errorf("heapify pop order = %v, want %v", got, want)
	}
	// Source must not be mutated.
	if src[0] != 5 {
		t.Error("NewFromSlice mutated the source")
	}
}

func TestNewFunc(t *testing.T) {
	type task struct {
		name     string
		priority int
	}
	// Highest priority value on top.
	pq := NewFunc(func(a, b task) bool { return a.priority < b.priority })
	pq.Push(task{"low", 1})
	pq.Push(task{"high", 10})
	pq.Push(task{"mid", 5})
	top, _ := pq.Top()
	if top.name != "high" {
		t.Errorf("Top = %q, want \"high\"", top.name)
	}
	order := []string{}
	for !pq.Empty() {
		top, _ := pq.Top()
		order = append(order, top.name)
		pq.Pop()
	}
	want := []string{"high", "mid", "low"}
	for i := range want {
		if order[i] != want[i] {
			t.Errorf("order = %v, want %v", order, want)
			break
		}
	}
}

func TestNewFuncNilPanics(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Error("NewFunc(nil) should panic")
		}
	}()
	NewFunc[int](nil)
}

func TestSwap(t *testing.T) {
	a := New[int]()
	a.Push(1)
	a.Push(2)
	b := NewMin[int]()
	b.Push(10)
	b.Push(20)
	b.Push(30)

	a.Swap(b)
	if a.Size() != 3 || b.Size() != 2 {
		t.Errorf("sizes after swap: a=%d b=%d", a.Size(), b.Size())
	}
	// a now holds b's min-heap comparator: top should be the minimum (10).
	topA, _ := a.Top()
	if topA != 10 {
		t.Errorf("a top after swap = %d, want 10 (min-heap semantics)", topA)
	}
	// b now holds a's max-heap comparator: top should be the maximum (2).
	topB, _ := b.Top()
	if topB != 2 {
		t.Errorf("b top after swap = %d, want 2 (max-heap semantics)", topB)
	}

	a.Swap(nil) // no-op
	if a.Size() != 3 {
		t.Error("swap with nil mutated queue")
	}
}

func TestClear(t *testing.T) {
	pq := New[int]()
	pq.Push(1)
	pq.Push(2)
	pq.Clear()
	if !pq.Empty() {
		t.Error("queue should be empty after Clear")
	}
	// Comparator must survive Clear.
	pq.Push(3)
	pq.Push(9)
	pq.Push(5)
	top, _ := pq.Top()
	if top != 9 {
		t.Errorf("Top after reuse = %d, want 9", top)
	}
}

// Property test: popping always yields descending order for a max-heap,
// regardless of insertion order.
func TestHeapInvariantRandom(t *testing.T) {
	rng := rand.New(rand.NewSource(42))
	for trial := 0; trial < 200; trial++ {
		n := rng.Intn(100)
		pq := New[int]()
		for i := 0; i < n; i++ {
			pq.Push(rng.Intn(1000))
		}
		got := drain(pq)
		for i := 1; i < len(got); i++ {
			if got[i-1] < got[i] {
				t.Fatalf("trial %d: not descending at %d: %v", trial, i, got)
			}
		}
		if len(got) != n {
			t.Fatalf("trial %d: drained %d, pushed %d", trial, len(got), n)
		}
	}
}

func equalInts(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func BenchmarkPush(b *testing.B) {
	pq := New[int]()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pq.Push(i)
	}
}

func BenchmarkPushPop(b *testing.B) {
	pq := New[int]()
	for i := 0; i < 1000; i++ {
		pq.Push(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pq.Push(i)
		pq.Pop()
	}
}

func BenchmarkHeapify(b *testing.B) {
	src := make([]int, 10000)
	for i := range src {
		src[i] = len(src) - i
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NewFromSlice(src)
	}
}
