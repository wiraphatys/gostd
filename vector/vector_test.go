package vector

import (
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	v := New[int]()
	if v == nil {
		t.Fatal("New returned nil")
	}
	if !v.Empty() {
		t.Error("new vector should be empty")
	}
	if v.Size() != 0 {
		t.Errorf("Size = %d, want 0", v.Size())
	}
}

func TestNewWithSize(t *testing.T) {
	v := NewWithSize[int](3)
	if v.Size() != 3 {
		t.Errorf("Size = %d, want 3", v.Size())
	}
	for i := 0; i < 3; i++ {
		if got := v.Get(i); got != 0 {
			t.Errorf("Get(%d) = %d, want zero value 0", i, got)
		}
	}
	// Negative clamps to zero.
	if NewWithSize[int](-5).Size() != 0 {
		t.Error("negative size should clamp to 0")
	}
}

func TestNewWithValue(t *testing.T) {
	v := NewWithValue(4, 7)
	if v.Size() != 4 {
		t.Errorf("Size = %d, want 4", v.Size())
	}
	for i := 0; i < 4; i++ {
		if got := v.Get(i); got != 7 {
			t.Errorf("Get(%d) = %d, want 7", i, got)
		}
	}
}

func TestNewFromSlice(t *testing.T) {
	src := []int{1, 2, 3}
	v := NewFromSlice(src)
	if !reflect.DeepEqual(v.ToSlice(), src) {
		t.Errorf("ToSlice = %v, want %v", v.ToSlice(), src)
	}
	// Mutating the source must not affect the vector (defensive copy).
	src[0] = 99
	if v.Get(0) == 99 {
		t.Error("vector aliased the source slice")
	}
}

func TestPushBackAndAccess(t *testing.T) {
	v := New[int]()
	for i := 1; i <= 5; i++ {
		v.PushBack(i * 10)
	}
	if v.Size() != 5 {
		t.Errorf("Size = %d, want 5", v.Size())
	}
	front, _ := v.Front()
	if front != 10 {
		t.Errorf("Front = %d, want 10", front)
	}
	back, _ := v.Back()
	if back != 50 {
		t.Errorf("Back = %d, want 50", back)
	}
	for i := 0; i < 5; i++ {
		if got := v.Get(i); got != (i+1)*10 {
			t.Errorf("Get(%d) = %d, want %d", i, got, (i+1)*10)
		}
	}
}

func TestAtBounds(t *testing.T) {
	v := NewFromSlice([]int{1, 2, 3})
	if _, err := v.At(-1); err == nil {
		t.Error("At(-1) should error")
	}
	if _, err := v.At(3); err == nil {
		t.Error("At(3) should error")
	}
	val, err := v.At(1)
	if err != nil || val != 2 {
		t.Errorf("At(1) = (%d, %v), want (2, nil)", val, err)
	}
}

func TestSet(t *testing.T) {
	v := NewWithSize[int](3)
	v.Set(1, 42)
	if v.Get(1) != 42 {
		t.Errorf("Get(1) = %d, want 42", v.Get(1))
	}
}

func TestFrontBackEmpty(t *testing.T) {
	v := New[int]()
	if _, err := v.Front(); err == nil {
		t.Error("Front on empty should error")
	}
	if _, err := v.Back(); err == nil {
		t.Error("Back on empty should error")
	}
}

func TestPopBack(t *testing.T) {
	v := New[int]()
	if err := v.PopBack(); err == nil {
		t.Error("PopBack on empty should error")
	}
	v.PushBack(1)
	v.PushBack(2)
	if err := v.PopBack(); err != nil {
		t.Errorf("PopBack errored: %v", err)
	}
	if v.Size() != 1 {
		t.Errorf("Size = %d, want 1", v.Size())
	}
	back, _ := v.Back()
	if back != 1 {
		t.Errorf("Back = %d, want 1", back)
	}
}

func TestInsert(t *testing.T) {
	v := NewFromSlice([]int{1, 2, 4})
	if err := v.Insert(2, 3); err != nil {
		t.Fatalf("Insert errored: %v", err)
	}
	want := []int{1, 2, 3, 4}
	if !reflect.DeepEqual(v.ToSlice(), want) {
		t.Errorf("after Insert = %v, want %v", v.ToSlice(), want)
	}
	// Insert at front and back.
	v.Insert(0, 0)
	v.Insert(v.Size(), 5)
	if !reflect.DeepEqual(v.ToSlice(), []int{0, 1, 2, 3, 4, 5}) {
		t.Errorf("after boundary inserts = %v", v.ToSlice())
	}
	// Out of range.
	if err := v.Insert(-1, 0); err == nil {
		t.Error("Insert(-1) should error")
	}
	if err := v.Insert(v.Size()+1, 0); err == nil {
		t.Error("Insert beyond size should error")
	}
}

func TestInsertRange(t *testing.T) {
	v := NewFromSlice([]int{1, 5})
	if err := v.InsertRange(1, []int{2, 3, 4}); err != nil {
		t.Fatalf("InsertRange errored: %v", err)
	}
	want := []int{1, 2, 3, 4, 5}
	if !reflect.DeepEqual(v.ToSlice(), want) {
		t.Errorf("after InsertRange = %v, want %v", v.ToSlice(), want)
	}
	// Empty range is a no-op.
	if err := v.InsertRange(0, nil); err != nil {
		t.Errorf("InsertRange(nil) errored: %v", err)
	}
	if v.Size() != 5 {
		t.Errorf("Size after empty insert = %d, want 5", v.Size())
	}
	if err := v.InsertRange(99, []int{1}); err == nil {
		t.Error("InsertRange out of range should error")
	}
}

func TestErase(t *testing.T) {
	v := NewFromSlice([]int{1, 2, 3, 4, 5})
	if err := v.Erase(2); err != nil {
		t.Fatalf("Erase errored: %v", err)
	}
	want := []int{1, 2, 4, 5}
	if !reflect.DeepEqual(v.ToSlice(), want) {
		t.Errorf("after Erase = %v, want %v", v.ToSlice(), want)
	}
	if err := v.Erase(-1); err == nil {
		t.Error("Erase(-1) should error")
	}
	if err := v.Erase(v.Size()); err == nil {
		t.Error("Erase past end should error")
	}
}

func TestEraseRange(t *testing.T) {
	v := NewFromSlice([]int{1, 2, 3, 4, 5})
	if err := v.EraseRange(1, 4); err != nil {
		t.Fatalf("EraseRange errored: %v", err)
	}
	want := []int{1, 5}
	if !reflect.DeepEqual(v.ToSlice(), want) {
		t.Errorf("after EraseRange = %v, want %v", v.ToSlice(), want)
	}
	// Empty range no-op.
	if err := v.EraseRange(1, 1); err != nil {
		t.Errorf("empty EraseRange errored: %v", err)
	}
	if err := v.EraseRange(0, 99); err == nil {
		t.Error("EraseRange out of range should error")
	}
	if err := v.EraseRange(2, 1); err == nil {
		t.Error("EraseRange with first > last should error")
	}
}

func TestReserveAndCapacity(t *testing.T) {
	v := New[int]()
	v.Reserve(100)
	if v.Capacity() < 100 {
		t.Errorf("Capacity = %d, want >= 100", v.Capacity())
	}
	if v.Size() != 0 {
		t.Error("Reserve must not change size")
	}
	// Reserving smaller is a no-op.
	capBefore := v.Capacity()
	v.Reserve(10)
	if v.Capacity() != capBefore {
		t.Error("Reserve smaller should not shrink capacity")
	}
	// Pushing within reserved capacity should not reallocate.
	for i := 0; i < 100; i++ {
		v.PushBack(i)
	}
	if v.Capacity() != capBefore {
		t.Errorf("Capacity changed after pushing within reservation: %d", v.Capacity())
	}
}

func TestResize(t *testing.T) {
	v := NewFromSlice([]int{1, 2, 3})
	v.Resize(5)
	if !reflect.DeepEqual(v.ToSlice(), []int{1, 2, 3, 0, 0}) {
		t.Errorf("grow Resize = %v", v.ToSlice())
	}
	v.Resize(2)
	if !reflect.DeepEqual(v.ToSlice(), []int{1, 2}) {
		t.Errorf("shrink Resize = %v", v.ToSlice())
	}
	v.ResizeWithValue(4, 9)
	if !reflect.DeepEqual(v.ToSlice(), []int{1, 2, 9, 9}) {
		t.Errorf("ResizeWithValue = %v", v.ToSlice())
	}
	v.Resize(-3)
	if v.Size() != 0 {
		t.Error("negative resize should clamp to 0")
	}
}

func TestResizeReusesCapacityWithoutStaleData(t *testing.T) {
	v := New[int]()
	v.Reserve(10)
	v.PushBack(1)
	v.PushBack(2)
	v.PopBack() // leaves a stale 2 in the backing array beyond len
	v.Resize(3) // must zero-fill, not resurrect the stale 2
	if !reflect.DeepEqual(v.ToSlice(), []int{1, 0, 0}) {
		t.Errorf("Resize resurfaced stale data: %v", v.ToSlice())
	}
}

func TestShrinkToFit(t *testing.T) {
	v := New[int]()
	v.Reserve(100)
	v.PushBack(1)
	v.PushBack(2)
	v.ShrinkToFit()
	if v.Capacity() != 2 {
		t.Errorf("Capacity = %d, want 2 after ShrinkToFit", v.Capacity())
	}
	if !reflect.DeepEqual(v.ToSlice(), []int{1, 2}) {
		t.Errorf("ShrinkToFit altered contents: %v", v.ToSlice())
	}
}

func TestAssign(t *testing.T) {
	v := NewFromSlice([]int{1, 2, 3})
	v.Assign(2, 8)
	if !reflect.DeepEqual(v.ToSlice(), []int{8, 8}) {
		t.Errorf("Assign = %v, want [8 8]", v.ToSlice())
	}
	v.AssignSlice([]int{4, 5, 6, 7})
	if !reflect.DeepEqual(v.ToSlice(), []int{4, 5, 6, 7}) {
		t.Errorf("AssignSlice = %v", v.ToSlice())
	}
}

func TestClear(t *testing.T) {
	v := NewFromSlice([]int{1, 2, 3})
	capBefore := v.Capacity()
	v.Clear()
	if !v.Empty() {
		t.Error("vector should be empty after Clear")
	}
	if v.Capacity() != capBefore {
		t.Error("Clear should retain capacity")
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
	a.Swap(nil) // no-op
	if a.Size() != 3 {
		t.Error("swap with nil mutated vector")
	}
}

func TestData(t *testing.T) {
	v := NewFromSlice([]int{1, 2, 3})
	d := v.Data()
	d[0] = 100 // mutating Data must mutate the vector
	if v.Get(0) != 100 {
		t.Error("Data should return the live backing slice")
	}
}

func TestToSliceIsCopy(t *testing.T) {
	v := NewFromSlice([]int{1, 2, 3})
	s := v.ToSlice()
	s[0] = 100
	if v.Get(0) == 100 {
		t.Error("ToSlice must return an independent copy")
	}
}

func TestForEach(t *testing.T) {
	v := NewFromSlice([]int{2, 4, 6})
	sum := 0
	indices := []int{}
	v.ForEach(func(i, val int) {
		indices = append(indices, i)
		sum += val
	})
	if sum != 12 {
		t.Errorf("sum = %d, want 12", sum)
	}
	if !reflect.DeepEqual(indices, []int{0, 1, 2}) {
		t.Errorf("indices = %v, want [0 1 2]", indices)
	}
}

func TestStringVector(t *testing.T) {
	v := New[string]()
	v.PushBack("a")
	v.PushBack("b")
	if back, _ := v.Back(); back != "b" {
		t.Errorf("Back = %q, want \"b\"", back)
	}
}

func BenchmarkPushBack(b *testing.B) {
	v := New[int]()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		v.PushBack(i)
	}
}

func BenchmarkPushBackReserved(b *testing.B) {
	v := New[int]()
	v.Reserve(b.N)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		v.PushBack(i)
	}
}

func BenchmarkGet(b *testing.B) {
	v := New[int]()
	for i := 0; i < 1000; i++ {
		v.PushBack(i)
	}
	b.ResetTimer()
	sum := 0
	for i := 0; i < b.N; i++ {
		sum += v.Get(i % 1000)
	}
	_ = sum
}
