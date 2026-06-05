package set

import (
	"reflect"
	"testing"
)

func TestInsertContainsSize(t *testing.T) {
	s := New[int]()
	if !s.Empty() {
		t.Error("new set should be empty")
	}
	if !s.Insert(5) {
		t.Error("first Insert(5) should return true")
	}
	if s.Insert(5) {
		t.Error("duplicate Insert(5) should return false")
	}
	s.Insert(3)
	s.Insert(8)
	if s.Size() != 3 {
		t.Errorf("Size = %d, want 3", s.Size())
	}
	if !s.Contains(3) || s.Contains(99) {
		t.Error("Contains wrong")
	}
	if s.Count(5) != 1 || s.Count(99) != 0 {
		t.Error("Count wrong")
	}
}

func TestSortedOrder(t *testing.T) {
	s := NewFromSlice([]int{5, 3, 8, 1, 9, 3, 5}) // duplicates collapse
	if !reflect.DeepEqual(s.ToSlice(), []int{1, 3, 5, 8, 9}) {
		t.Errorf("ToSlice = %v, want sorted unique", s.ToSlice())
	}
}

func TestErase(t *testing.T) {
	s := NewFromSlice([]int{1, 2, 3, 4, 5})
	if !s.Erase(3) {
		t.Error("Erase(3) should return true")
	}
	if s.Erase(3) {
		t.Error("Erase(3) again should return false")
	}
	if !reflect.DeepEqual(s.ToSlice(), []int{1, 2, 4, 5}) {
		t.Errorf("after Erase = %v", s.ToSlice())
	}
}

func TestFindAndEraseIter(t *testing.T) {
	s := NewFromSlice([]int{1, 2, 3, 4, 5})
	it := s.Find(3)
	if it.IsEnd() || it.Value() != 3 {
		t.Errorf("Find(3) = %v", it.Value())
	}
	next := s.EraseIter(it)
	if next.IsEnd() || next.Value() != 4 {
		t.Errorf("EraseIter returned %v, want 4", next.Value())
	}
	if s.Contains(3) {
		t.Error("3 should be erased")
	}
	// Find absent.
	if !s.Find(99).IsEnd() {
		t.Error("Find(99) should be End")
	}
}

func TestBounds(t *testing.T) {
	s := NewFromSlice([]int{10, 20, 30, 40})
	lb := s.LowerBound(20)
	if lb.Value() != 20 {
		t.Errorf("LowerBound(20) = %d, want 20", lb.Value())
	}
	ub := s.UpperBound(20)
	if ub.Value() != 30 {
		t.Errorf("UpperBound(20) = %d, want 30", ub.Value())
	}
	if !s.LowerBound(99).IsEnd() {
		t.Error("LowerBound(99) should be End")
	}
	lo, hi := s.EqualRange(25)
	if lo.Value() != 30 || hi.Value() != 30 {
		t.Errorf("EqualRange(25) = [%d,%d), want [30,30)", lo.Value(), hi.Value())
	}
}

func TestForwardReverseIteration(t *testing.T) {
	s := NewFromSlice([]int{3, 1, 4, 1, 5, 9, 2, 6})
	var fwd []int
	for it := s.Begin(); !it.Equal(s.End()); it.Next() {
		fwd = append(fwd, it.Value())
	}
	if !reflect.DeepEqual(fwd, []int{1, 2, 3, 4, 5, 6, 9}) {
		t.Errorf("forward = %v", fwd)
	}
	var rev []int
	for it := s.RBegin(); !it.Equal(s.REnd()); it.Prev() {
		rev = append(rev, it.Value())
	}
	if !reflect.DeepEqual(rev, []int{9, 6, 5, 4, 3, 2, 1}) {
		t.Errorf("reverse = %v", rev)
	}
}

func TestNewFunc(t *testing.T) {
	// Descending order via greater-than comparator.
	s := NewFunc(func(a, b int) bool { return a > b })
	for _, v := range []int{1, 5, 3, 2, 4} {
		s.Insert(v)
	}
	if !reflect.DeepEqual(s.ToSlice(), []int{5, 4, 3, 2, 1}) {
		t.Errorf("descending set = %v", s.ToSlice())
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

func TestSwapAndClear(t *testing.T) {
	a := NewFromSlice([]int{1, 2})
	b := NewFromSlice([]int{3, 4, 5})
	a.Swap(b)
	if !reflect.DeepEqual(a.ToSlice(), []int{3, 4, 5}) {
		t.Errorf("a after swap = %v", a.ToSlice())
	}
	a.Swap(nil)
	if a.Size() != 3 {
		t.Error("swap with nil mutated set")
	}
	a.Clear()
	if !a.Empty() {
		t.Error("set should be empty after Clear")
	}
}

func TestStringSet(t *testing.T) {
	s := New[string]()
	s.Insert("banana")
	s.Insert("apple")
	s.Insert("cherry")
	if !reflect.DeepEqual(s.ToSlice(), []string{"apple", "banana", "cherry"}) {
		t.Errorf("string set = %v", s.ToSlice())
	}
}

func TestEmptyIterationAndEraseIterEnd(t *testing.T) {
	s := New[int]()
	if !s.Begin().Equal(s.End()) {
		t.Error("empty Begin should equal End")
	}
	if !s.EraseIter(s.End()).IsEnd() {
		t.Error("EraseIter(End) should return End")
	}
}

func BenchmarkInsert(b *testing.B) {
	s := New[int]()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Insert(i)
	}
}

func BenchmarkContains(b *testing.B) {
	s := New[int]()
	for i := 0; i < 100000; i++ {
		s.Insert(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Contains(i % 100000)
	}
}
