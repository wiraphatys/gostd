package unorderedmap

import (
	"sort"
	"testing"
)

func TestInsertAtGet(t *testing.T) {
	m := New[string, int]()
	if !m.Insert("a", 1) {
		t.Error("Insert(a) should return true")
	}
	if m.Insert("a", 2) {
		t.Error("duplicate Insert(a) should return false")
	}
	v, err := m.At("a")
	if err != nil || v != 1 {
		t.Errorf("At(a) = (%d, %v), want (1, nil)", v, err)
	}
	if _, err := m.At("missing"); err == nil {
		t.Error("At(missing) should error")
	}
	if v, ok := m.Get("a"); !ok || v != 1 {
		t.Errorf("Get(a) = (%d, %v)", v, ok)
	}
	if _, ok := m.Get("missing"); ok {
		t.Error("Get(missing) should be not-ok")
	}
}

func TestInsertOrAssignSet(t *testing.T) {
	m := New[string, int]()
	if !m.InsertOrAssign("k", 1) {
		t.Error("first InsertOrAssign should report inserted")
	}
	if m.InsertOrAssign("k", 2) {
		t.Error("second InsertOrAssign should report not inserted")
	}
	if v, _ := m.At("k"); v != 2 {
		t.Errorf("value = %d, want 2", v)
	}
	m.Set("k", 9)
	if v, _ := m.At("k"); v != 9 {
		t.Errorf("value = %d, want 9", v)
	}
}

func TestGetOr(t *testing.T) {
	m := New[string, int]()
	m.Set("a", 5)
	if got := m.GetOr("a", -1); got != 5 {
		t.Errorf("GetOr(a) = %d, want 5", got)
	}
	if got := m.GetOr("missing", -1); got != -1 {
		t.Errorf("GetOr(missing) = %d, want -1", got)
	}
	// GetOr must not insert.
	if m.Contains("missing") {
		t.Error("GetOr should not insert the key")
	}
	// The read-with-default + Set counting idiom.
	counts := New[string, int]()
	for _, w := range []string{"x", "y", "x", "x"} {
		counts.Set(w, counts.GetOr(w, 0)+1)
	}
	if v, _ := counts.At("x"); v != 3 {
		t.Errorf("count x = %d, want 3", v)
	}
}

func TestEraseContainsCount(t *testing.T) {
	m := New[int, int]()
	for i := 0; i < 5; i++ {
		m.Set(i, i*i)
	}
	if !m.Erase(2) {
		t.Error("Erase(2) should return true")
	}
	if m.Erase(2) {
		t.Error("Erase(2) again should return false")
	}
	if m.Contains(2) || m.Count(2) != 0 {
		t.Error("2 should be gone")
	}
	if m.Size() != 4 {
		t.Errorf("Size = %d, want 4", m.Size())
	}
}

func TestKeysValues(t *testing.T) {
	m := New[int, int]()
	for i := 1; i <= 3; i++ {
		m.Set(i, i*10)
	}
	keys := m.Keys()
	sort.Ints(keys)
	if len(keys) != 3 || keys[0] != 1 || keys[2] != 3 {
		t.Errorf("Keys = %v", keys)
	}
	vals := m.Values()
	sort.Ints(vals)
	if len(vals) != 3 || vals[0] != 10 || vals[2] != 30 {
		t.Errorf("Values = %v", vals)
	}
}

func TestReserveSwapClear(t *testing.T) {
	m := NewWithCapacity[int, int](10)
	for i := 0; i < 20; i++ {
		m.Set(i, i)
	}
	m.Reserve(1000)
	if m.Size() != 20 {
		t.Errorf("Size after Reserve = %d, want 20", m.Size())
	}
	if v, _ := m.At(10); v != 10 {
		t.Error("Reserve lost a value")
	}

	other := New[int, int]()
	other.Set(99, 99)
	m.Swap(other)
	if m.Size() != 1 || other.Size() != 20 {
		t.Errorf("after swap m=%d other=%d", m.Size(), other.Size())
	}

	m.Clear()
	if !m.Empty() {
		t.Error("should be empty after Clear")
	}
}

func TestForEach(t *testing.T) {
	m := New[int, int]()
	for i := 1; i <= 4; i++ {
		m.Set(i, i)
	}
	sumK, sumV := 0, 0
	m.ForEach(func(k, v int) {
		sumK += k
		sumV += v
	})
	if sumK != 10 || sumV != 10 {
		t.Errorf("ForEach sums = (%d,%d), want (10,10)", sumK, sumV)
	}
}

func BenchmarkSet(b *testing.B) {
	m := New[int, int]()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Set(i, i)
	}
}
