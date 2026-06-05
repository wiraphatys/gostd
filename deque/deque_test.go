package deque

import (
	"math/rand"
	"reflect"
	"testing"
)

func TestNewEmpty(t *testing.T) {
	d := New[int]()
	if !d.Empty() || d.Size() != 0 {
		t.Error("new deque should be empty with size 0")
	}
	if err := d.PopFront(); err == nil {
		t.Error("PopFront on empty should error")
	}
	if err := d.PopBack(); err == nil {
		t.Error("PopBack on empty should error")
	}
	if _, err := d.Front(); err == nil {
		t.Error("Front on empty should error")
	}
	if _, err := d.Back(); err == nil {
		t.Error("Back on empty should error")
	}
}

func TestNewFromSlice(t *testing.T) {
	d := NewFromSlice([]int{1, 2, 3})
	if !reflect.DeepEqual(d.ToSlice(), []int{1, 2, 3}) {
		t.Errorf("ToSlice = %v, want [1 2 3]", d.ToSlice())
	}
	front, _ := d.Front()
	back, _ := d.Back()
	if front != 1 || back != 3 {
		t.Errorf("front=%d back=%d, want 1 and 3", front, back)
	}
}

func TestPushBackPopFront(t *testing.T) {
	d := New[int]()
	for i := 1; i <= 5; i++ {
		d.PushBack(i)
	}
	// FIFO order via PopFront.
	for i := 1; i <= 5; i++ {
		front, _ := d.Front()
		if front != i {
			t.Errorf("Front = %d, want %d", front, i)
		}
		d.PopFront()
	}
	if !d.Empty() {
		t.Error("deque should be empty")
	}
}

func TestPushFrontPopBack(t *testing.T) {
	d := New[int]()
	for i := 1; i <= 5; i++ {
		d.PushFront(i) // 5,4,3,2,1
	}
	if !reflect.DeepEqual(d.ToSlice(), []int{5, 4, 3, 2, 1}) {
		t.Errorf("ToSlice = %v, want [5 4 3 2 1]", d.ToSlice())
	}
	for i := 1; i <= 5; i++ {
		back, _ := d.Back()
		if back != i {
			t.Errorf("Back = %d, want %d", back, i)
		}
		d.PopBack()
	}
}

func TestMixedEnds(t *testing.T) {
	d := New[int]()
	d.PushBack(1)  // [1]
	d.PushFront(2) // [2 1]
	d.PushBack(3)  // [2 1 3]
	d.PushFront(4) // [4 2 1 3]
	if !reflect.DeepEqual(d.ToSlice(), []int{4, 2, 1, 3}) {
		t.Errorf("ToSlice = %v, want [4 2 1 3]", d.ToSlice())
	}
}

func TestAtAndSet(t *testing.T) {
	d := NewFromSlice([]int{10, 20, 30})
	for i, want := range []int{10, 20, 30} {
		got, err := d.At(i)
		if err != nil || got != want {
			t.Errorf("At(%d) = (%d, %v), want (%d, nil)", i, got, err, want)
		}
	}
	if _, err := d.At(-1); err == nil {
		t.Error("At(-1) should error")
	}
	if _, err := d.At(3); err == nil {
		t.Error("At(3) should error")
	}
	d.Set(1, 99)
	if d.Get(1) != 99 {
		t.Errorf("Get(1) = %d, want 99", d.Get(1))
	}
}

func TestInsert(t *testing.T) {
	d := NewFromSlice([]int{1, 2, 4, 5})
	if err := d.Insert(2, 3); err != nil { // insert before index 2
		t.Fatalf("Insert errored: %v", err)
	}
	if !reflect.DeepEqual(d.ToSlice(), []int{1, 2, 3, 4, 5}) {
		t.Errorf("after Insert = %v, want [1 2 3 4 5]", d.ToSlice())
	}
	// Insert at both extremes.
	d.Insert(0, 0)
	d.Insert(d.Size(), 6)
	if !reflect.DeepEqual(d.ToSlice(), []int{0, 1, 2, 3, 4, 5, 6}) {
		t.Errorf("after extreme inserts = %v", d.ToSlice())
	}
	if err := d.Insert(-1, 0); err == nil {
		t.Error("Insert(-1) should error")
	}
	if err := d.Insert(d.Size()+1, 0); err == nil {
		t.Error("Insert beyond size should error")
	}
}

func TestErase(t *testing.T) {
	d := NewFromSlice([]int{1, 2, 3, 4, 5})
	if err := d.Erase(2); err != nil { // remove the 3
		t.Fatalf("Erase errored: %v", err)
	}
	if !reflect.DeepEqual(d.ToSlice(), []int{1, 2, 4, 5}) {
		t.Errorf("after Erase = %v, want [1 2 4 5]", d.ToSlice())
	}
	d.Erase(0)            // front
	d.Erase(d.Size() - 1) // back
	if !reflect.DeepEqual(d.ToSlice(), []int{2, 4}) {
		t.Errorf("after extreme erases = %v, want [2 4]", d.ToSlice())
	}
	if err := d.Erase(-1); err == nil {
		t.Error("Erase(-1) should error")
	}
	if err := d.Erase(d.Size()); err == nil {
		t.Error("Erase past end should error")
	}
}

func TestResize(t *testing.T) {
	d := NewFromSlice([]int{1, 2, 3})
	d.Resize(5)
	if !reflect.DeepEqual(d.ToSlice(), []int{1, 2, 3, 0, 0}) {
		t.Errorf("grow Resize = %v", d.ToSlice())
	}
	d.Resize(2)
	if !reflect.DeepEqual(d.ToSlice(), []int{1, 2}) {
		t.Errorf("shrink Resize = %v", d.ToSlice())
	}
	d.ResizeWithValue(4, 7)
	if !reflect.DeepEqual(d.ToSlice(), []int{1, 2, 7, 7}) {
		t.Errorf("ResizeWithValue = %v", d.ToSlice())
	}
}

func TestShrinkToFit(t *testing.T) {
	d := New[int]()
	d.Reserve(100)
	d.PushBack(1)
	d.PushBack(2)
	if d.Capacity() < 100 {
		t.Errorf("Capacity = %d, want >= 100 after Reserve", d.Capacity())
	}
	d.ShrinkToFit()
	if d.Capacity() != 2 {
		t.Errorf("Capacity = %d, want 2 after ShrinkToFit", d.Capacity())
	}
	if !reflect.DeepEqual(d.ToSlice(), []int{1, 2}) {
		t.Errorf("ShrinkToFit altered contents: %v", d.ToSlice())
	}
}

func TestShrinkToFitWrapped(t *testing.T) {
	// Force a wrapped layout (head != 0), then shrink.
	d := New[int]()
	d.PushFront(2)
	d.PushFront(1) // [1 2] but stored wrapped
	d.PushBack(3)  // [1 2 3]
	d.ShrinkToFit()
	if !reflect.DeepEqual(d.ToSlice(), []int{1, 2, 3}) {
		t.Errorf("wrapped ShrinkToFit = %v, want [1 2 3]", d.ToSlice())
	}
}

func TestClear(t *testing.T) {
	d := NewFromSlice([]int{1, 2, 3})
	d.Clear()
	if !d.Empty() {
		t.Error("deque should be empty after Clear")
	}
	// Reuse after clear.
	d.PushBack(9)
	front, _ := d.Front()
	if front != 9 {
		t.Errorf("Front after reuse = %d, want 9", front)
	}
}

func TestSwap(t *testing.T) {
	a := NewFromSlice([]int{1, 2})
	b := NewFromSlice([]int{3, 4, 5})
	a.Swap(b)
	if !reflect.DeepEqual(a.ToSlice(), []int{3, 4, 5}) {
		t.Errorf("a after swap = %v", a.ToSlice())
	}
	if !reflect.DeepEqual(b.ToSlice(), []int{1, 2}) {
		t.Errorf("b after swap = %v", b.ToSlice())
	}
	a.Swap(nil)
	if a.Size() != 3 {
		t.Error("swap with nil mutated deque")
	}
}

func TestForEach(t *testing.T) {
	d := NewFromSlice([]int{5, 6, 7})
	var got []int
	d.ForEach(func(_ int, v int) { got = append(got, v) })
	if !reflect.DeepEqual(got, []int{5, 6, 7}) {
		t.Errorf("ForEach order = %v", got)
	}
}

func TestWrapAroundGrowth(t *testing.T) {
	// Exercise the ring buffer wrapping and growing under interleaved ops.
	d := New[int]()
	for i := 0; i < 1000; i++ {
		d.PushBack(i)
		d.PushFront(-i)
	}
	if d.Size() != 2000 {
		t.Fatalf("Size = %d, want 2000", d.Size())
	}
	front, _ := d.Front()
	back, _ := d.Back()
	if front != -999 || back != 999 {
		t.Errorf("front=%d back=%d, want -999 and 999", front, back)
	}
}

// Differential test: run random operations against a reference []int model and
// assert the deque always matches it. This is the strongest correctness check
// for the circular-buffer index arithmetic.
func TestRandomizedAgainstModel(t *testing.T) {
	rng := rand.New(rand.NewSource(7))
	d := New[int]()
	var model []int

	check := func(step int) {
		if d.Size() != len(model) {
			t.Fatalf("step %d: size %d != model %d", step, d.Size(), len(model))
		}
		got := d.ToSlice()
		for i := range model {
			if got[i] != model[i] {
				t.Fatalf("step %d: deque %v != model %v", step, got, model)
			}
		}
	}

	for step := 0; step < 20000; step++ {
		switch rng.Intn(8) {
		case 0:
			v := rng.Intn(1000)
			d.PushBack(v)
			model = append(model, v)
		case 1:
			v := rng.Intn(1000)
			d.PushFront(v)
			model = append([]int{v}, model...)
		case 2:
			if len(model) > 0 {
				d.PopBack()
				model = model[:len(model)-1]
			}
		case 3:
			if len(model) > 0 {
				d.PopFront()
				model = model[1:]
			}
		case 4:
			if len(model) > 0 {
				i := rng.Intn(len(model))
				v := rng.Intn(1000)
				d.Set(i, v)
				model[i] = v
			}
		case 5:
			i := rng.Intn(len(model) + 1)
			v := rng.Intn(1000)
			d.Insert(i, v)
			model = append(model, 0)
			copy(model[i+1:], model[i:])
			model[i] = v
		case 6:
			if len(model) > 0 {
				i := rng.Intn(len(model))
				d.Erase(i)
				model = append(model[:i], model[i+1:]...)
			}
		case 7:
			if len(model) > 0 {
				i := rng.Intn(len(model))
				got, _ := d.At(i)
				if got != model[i] {
					t.Fatalf("step %d: At(%d)=%d, model=%d", step, i, got, model[i])
				}
			}
		}
		if step%100 == 0 {
			check(step)
		}
	}
	check(20000)
}

func BenchmarkPushBack(b *testing.B) {
	d := New[int]()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		d.PushBack(i)
	}
}

func BenchmarkPushFront(b *testing.B) {
	d := New[int]()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		d.PushFront(i)
	}
}

func BenchmarkQueueChurn(b *testing.B) {
	d := New[int]()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		d.PushBack(i)
		if i%2 == 0 {
			d.PopFront()
		}
	}
}
