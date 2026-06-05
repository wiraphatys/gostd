package unorderedset

import (
	"sort"
	"testing"
)

func sortedSlice(s *Set[int]) []int {
	out := s.ToSlice()
	sort.Ints(out)
	return out
}

func TestInsertContains(t *testing.T) {
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
	if !s.Contains(5) || s.Contains(99) {
		t.Error("Contains wrong")
	}
	if s.Count(5) != 1 || s.Count(99) != 0 {
		t.Error("Count wrong")
	}
}

func TestErase(t *testing.T) {
	s := NewFromSlice([]int{1, 2, 3})
	if !s.Erase(2) {
		t.Error("Erase(2) should return true")
	}
	if s.Erase(2) {
		t.Error("Erase(2) again should return false")
	}
	if s.Contains(2) {
		t.Error("2 should be gone")
	}
	if got := sortedSlice(s); len(got) != 2 || got[0] != 1 || got[1] != 3 {
		t.Errorf("after erase = %v", got)
	}
}

func TestNewFromSliceDedup(t *testing.T) {
	s := NewFromSlice([]int{1, 1, 2, 2, 3})
	if s.Size() != 3 {
		t.Errorf("Size = %d, want 3 (dedup)", s.Size())
	}
}

func TestClearReserveSwap(t *testing.T) {
	s := NewWithCapacity[int](100)
	for i := 0; i < 50; i++ {
		s.Insert(i)
	}
	s.Reserve(1000) // rehash, contents preserved
	if s.Size() != 50 {
		t.Errorf("Size after Reserve = %d, want 50", s.Size())
	}
	if !s.Contains(25) {
		t.Error("Reserve lost an element")
	}

	other := NewFromSlice([]int{100, 200})
	s.Swap(other)
	if s.Size() != 2 || other.Size() != 50 {
		t.Errorf("after swap s=%d other=%d", s.Size(), other.Size())
	}
	s.Swap(nil) // no-op
	if s.Size() != 2 {
		t.Error("swap with nil mutated set")
	}

	s.Clear()
	if !s.Empty() {
		t.Error("should be empty after Clear")
	}
}

func TestForEach(t *testing.T) {
	s := NewFromSlice([]int{1, 2, 3, 4})
	sum := 0
	s.ForEach(func(v int) { sum += v })
	if sum != 10 {
		t.Errorf("ForEach sum = %d, want 10", sum)
	}
}

func TestStringSet(t *testing.T) {
	s := New[string]()
	s.Insert("a")
	s.Insert("b")
	s.Insert("a")
	if s.Size() != 2 {
		t.Errorf("Size = %d, want 2", s.Size())
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
