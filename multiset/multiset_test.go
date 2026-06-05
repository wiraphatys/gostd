package multiset

import (
	"reflect"
	"testing"
)

func TestInsertDuplicatesAndCount(t *testing.T) {
	m := New[int]()
	m.Insert(5)
	m.Insert(5)
	m.Insert(5)
	m.Insert(3)
	if m.Size() != 4 {
		t.Errorf("Size = %d, want 4", m.Size())
	}
	if m.Count(5) != 3 {
		t.Errorf("Count(5) = %d, want 3", m.Count(5))
	}
	if m.Count(3) != 1 {
		t.Errorf("Count(3) = %d, want 1", m.Count(3))
	}
	if m.Count(99) != 0 {
		t.Errorf("Count(99) = %d, want 0", m.Count(99))
	}
}

func TestSortedWithDuplicates(t *testing.T) {
	m := NewFromSlice([]int{3, 1, 2, 3, 1, 3})
	if !reflect.DeepEqual(m.ToSlice(), []int{1, 1, 2, 3, 3, 3}) {
		t.Errorf("ToSlice = %v", m.ToSlice())
	}
}

func TestEraseAll(t *testing.T) {
	m := NewFromSlice([]int{5, 5, 5, 3, 8})
	if removed := m.Erase(5); removed != 3 {
		t.Errorf("Erase(5) = %d, want 3", removed)
	}
	if m.Count(5) != 0 {
		t.Error("all 5s should be gone")
	}
	if !reflect.DeepEqual(m.ToSlice(), []int{3, 8}) {
		t.Errorf("after Erase = %v", m.ToSlice())
	}
}

func TestEraseOne(t *testing.T) {
	m := NewFromSlice([]int{5, 5, 5})
	if !m.EraseOne(5) {
		t.Error("EraseOne(5) should return true")
	}
	if m.Count(5) != 2 {
		t.Errorf("Count(5) = %d, want 2 after EraseOne", m.Count(5))
	}
	if m.EraseOne(99) {
		t.Error("EraseOne(99) should return false")
	}
}

func TestEqualRange(t *testing.T) {
	m := NewFromSlice([]int{1, 2, 2, 2, 3})
	lo, hi := m.EqualRange(2)
	var got []int
	for it := lo; !it.Equal(hi); it.Next() {
		got = append(got, it.Value())
	}
	if !reflect.DeepEqual(got, []int{2, 2, 2}) {
		t.Errorf("EqualRange(2) collected %v, want [2 2 2]", got)
	}
}

func TestEraseIter(t *testing.T) {
	m := NewFromSlice([]int{1, 2, 2, 3})
	it := m.Find(2)
	if it.IsEnd() {
		t.Fatal("Find(2) should not be End")
	}
	next := m.EraseIter(it)
	if next.IsEnd() {
		t.Error("EraseIter should return a valid iterator")
	}
	if m.Count(2) != 1 {
		t.Errorf("Count(2) = %d, want 1", m.Count(2))
	}
}

func TestReverseIteration(t *testing.T) {
	m := NewFromSlice([]int{1, 2, 2, 3})
	var rev []int
	for it := m.RBegin(); !it.Equal(m.REnd()); it.Prev() {
		rev = append(rev, it.Value())
	}
	if !reflect.DeepEqual(rev, []int{3, 2, 2, 1}) {
		t.Errorf("reverse = %v", rev)
	}
}

func TestNewFunc(t *testing.T) {
	m := NewFunc(func(a, b int) bool { return a > b })
	for _, v := range []int{1, 2, 2, 3} {
		m.Insert(v)
	}
	if !reflect.DeepEqual(m.ToSlice(), []int{3, 2, 2, 1}) {
		t.Errorf("descending multiset = %v", m.ToSlice())
	}
}

func TestSwap(t *testing.T) {
	a := NewFromSlice([]int{1, 1})
	b := NewFromSlice([]int{2, 2, 2})
	a.Swap(b)
	if a.Size() != 3 || b.Size() != 2 {
		t.Errorf("after swap a=%d b=%d", a.Size(), b.Size())
	}
}

func BenchmarkInsert(b *testing.B) {
	m := New[int]()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Insert(i % 1000)
	}
}
