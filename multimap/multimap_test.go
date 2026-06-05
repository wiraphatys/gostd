package multimap

import (
	"reflect"
	"testing"
)

func TestInsertCount(t *testing.T) {
	m := New[string, int]()
	m.Insert("a", 1)
	m.Insert("a", 2)
	m.Insert("a", 3)
	m.Insert("b", 10)
	if m.Size() != 4 {
		t.Errorf("Size = %d, want 4", m.Size())
	}
	if m.Count("a") != 3 {
		t.Errorf("Count(a) = %d, want 3", m.Count("a"))
	}
	if m.Count("b") != 1 {
		t.Errorf("Count(b) = %d, want 1", m.Count("b"))
	}
	if m.Count("z") != 0 {
		t.Errorf("Count(z) = %d, want 0", m.Count("z"))
	}
}

func TestEqualRangeCollectsAllValues(t *testing.T) {
	m := New[string, int]()
	m.Insert("k", 1)
	m.Insert("k", 2)
	m.Insert("k", 3)
	m.Insert("other", 99)

	var got []int
	lo, hi := m.EqualRange("k")
	for it := lo; !it.Equal(hi); it.Next() {
		got = append(got, it.Value())
	}
	if len(got) != 3 {
		t.Fatalf("collected %v, want 3 values", got)
	}
	sum := 0
	for _, v := range got {
		sum += v
	}
	if sum != 6 {
		t.Errorf("values sum = %d, want 6", sum)
	}
}

func TestEraseAllAndOne(t *testing.T) {
	m := New[int, int]()
	m.Insert(1, 10)
	m.Insert(1, 11)
	m.Insert(1, 12)
	m.Insert(2, 20)

	if removed := m.Erase(1); removed != 3 {
		t.Errorf("Erase(1) = %d, want 3", removed)
	}
	if m.Count(1) != 0 {
		t.Error("all key-1 entries should be gone")
	}
	if m.Size() != 1 {
		t.Errorf("Size = %d, want 1", m.Size())
	}

	// Erase a single entry via iterator.
	m.Insert(3, 30)
	m.Insert(3, 31)
	it := m.Find(3)
	m.EraseIter(it)
	if m.Count(3) != 1 {
		t.Errorf("Count(3) = %d, want 1 after single erase", m.Count(3))
	}
}

func TestOrderedTraversal(t *testing.T) {
	m := New[int, string]()
	m.Insert(3, "c")
	m.Insert(1, "a")
	m.Insert(2, "b")
	m.Insert(1, "a2")

	var keys []int
	m.ForEach(func(k int, _ string) { keys = append(keys, k) })
	if !reflect.DeepEqual(keys, []int{1, 1, 2, 3}) {
		t.Errorf("keys in order = %v", keys)
	}
}

func TestNewFunc(t *testing.T) {
	m := NewFunc[int, int](func(a, b int) bool { return a > b })
	m.Insert(1, 1)
	m.Insert(2, 2)
	m.Insert(2, 22)
	m.Insert(3, 3)
	var keys []int
	m.ForEach(func(k, _ int) { keys = append(keys, k) })
	if !reflect.DeepEqual(keys, []int{3, 2, 2, 1}) {
		t.Errorf("descending keys = %v", keys)
	}
}

func TestSwap(t *testing.T) {
	a := New[int, int]()
	a.Insert(1, 1)
	b := New[int, int]()
	b.Insert(2, 2)
	b.Insert(2, 3)
	a.Swap(b)
	if a.Size() != 2 || b.Size() != 1 {
		t.Errorf("after swap a=%d b=%d", a.Size(), b.Size())
	}
}

func BenchmarkInsert(b *testing.B) {
	m := New[int, int]()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Insert(i%1000, i)
	}
}
