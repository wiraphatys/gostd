package treemap

import (
	"reflect"
	"testing"
)

func TestInsertGetAt(t *testing.T) {
	m := New[string, int]()
	if !m.Insert("a", 1) {
		t.Error("Insert(a) should return true")
	}
	if m.Insert("a", 2) {
		t.Error("duplicate Insert(a) should return false")
	}
	// Value preserved on duplicate Insert.
	v, _ := m.At("a")
	if v != 1 {
		t.Errorf("At(a) = %d, want 1 (duplicate insert must not overwrite)", v)
	}
	// Get comma-ok.
	if v, ok := m.Get("a"); !ok || v != 1 {
		t.Errorf("Get(a) = (%d,%v)", v, ok)
	}
	if _, ok := m.Get("zzz"); ok {
		t.Error("Get(zzz) should be not-ok")
	}
	// At on missing key errors.
	if _, err := m.At("zzz"); err == nil {
		t.Error("At(zzz) should error")
	}
}

func TestInsertOrAssignAndSet(t *testing.T) {
	m := New[string, int]()
	if !m.InsertOrAssign("k", 1) {
		t.Error("first InsertOrAssign should report inserted")
	}
	if m.InsertOrAssign("k", 2) {
		t.Error("second InsertOrAssign should report not inserted")
	}
	v, _ := m.At("k")
	if v != 2 {
		t.Errorf("At(k) = %d, want 2 after InsertOrAssign", v)
	}
	m.Set("k", 9)
	v, _ = m.At("k")
	if v != 9 {
		t.Errorf("At(k) = %d, want 9 after Set", v)
	}
}

func TestRefOperatorBracket(t *testing.T) {
	m := New[string, int]()
	// Ref on absent key default-inserts a zero value.
	p := m.Ref("counter")
	if *p != 0 {
		t.Errorf("Ref default = %d, want 0", *p)
	}
	*p = 5 // mutate through the pointer
	if v, _ := m.At("counter"); v != 5 {
		t.Errorf("At(counter) = %d, want 5 after mutation through Ref", v)
	}
	// Ref on existing key returns the live value.
	*m.Ref("counter")++
	if v, _ := m.At("counter"); v != 6 {
		t.Errorf("At(counter) = %d, want 6", v)
	}
	if m.Size() != 1 {
		t.Errorf("Size = %d, want 1", m.Size())
	}
}

func TestSortedKeysValues(t *testing.T) {
	m := New[int, string]()
	m.Set(3, "three")
	m.Set(1, "one")
	m.Set(2, "two")
	if !reflect.DeepEqual(m.Keys(), []int{1, 2, 3}) {
		t.Errorf("Keys = %v", m.Keys())
	}
	if !reflect.DeepEqual(m.Values(), []string{"one", "two", "three"}) {
		t.Errorf("Values = %v", m.Values())
	}
}

func TestErase(t *testing.T) {
	m := New[int, int]()
	for i := 1; i <= 5; i++ {
		m.Set(i, i*i)
	}
	if !m.Erase(3) {
		t.Error("Erase(3) should return true")
	}
	if m.Erase(3) {
		t.Error("Erase(3) again should return false")
	}
	if !reflect.DeepEqual(m.Keys(), []int{1, 2, 4, 5}) {
		t.Errorf("keys after erase = %v", m.Keys())
	}
}

func TestFindEraseIter(t *testing.T) {
	m := New[int, int]()
	for i := 1; i <= 5; i++ {
		m.Set(i, i)
	}
	it := m.Find(3)
	if it.IsEnd() || it.Key() != 3 {
		t.Errorf("Find(3) key = %v", it.Key())
	}
	next := m.EraseIter(it)
	if next.IsEnd() || next.Key() != 4 {
		t.Errorf("EraseIter next key = %v, want 4", next.Key())
	}
	if m.Contains(3) {
		t.Error("3 should be erased")
	}
}

func TestIterationAndSetValue(t *testing.T) {
	m := New[int, int]()
	for i := 1; i <= 4; i++ {
		m.Set(i, i)
	}
	// Double each value via iterator SetValue.
	for it := m.Begin(); !it.Equal(m.End()); it.Next() {
		it.SetValue(it.Value() * 2)
	}
	if !reflect.DeepEqual(m.Values(), []int{2, 4, 6, 8}) {
		t.Errorf("values after doubling = %v", m.Values())
	}
	// Reverse iteration over keys.
	var rk []int
	for it := m.RBegin(); !it.Equal(m.REnd()); it.Prev() {
		rk = append(rk, it.Key())
	}
	if !reflect.DeepEqual(rk, []int{4, 3, 2, 1}) {
		t.Errorf("reverse keys = %v", rk)
	}
}

func TestBounds(t *testing.T) {
	m := New[int, int]()
	for _, k := range []int{10, 20, 30} {
		m.Set(k, k)
	}
	if m.LowerBound(20).Key() != 20 {
		t.Error("LowerBound(20) wrong")
	}
	if m.UpperBound(20).Key() != 30 {
		t.Error("UpperBound(20) wrong")
	}
	if !m.LowerBound(99).IsEnd() {
		t.Error("LowerBound(99) should be End")
	}
}

func TestNewFunc(t *testing.T) {
	m := NewFunc[int, string](func(a, b int) bool { return a > b })
	m.Set(1, "a")
	m.Set(3, "c")
	m.Set(2, "b")
	if !reflect.DeepEqual(m.Keys(), []int{3, 2, 1}) {
		t.Errorf("descending keys = %v", m.Keys())
	}
}

func TestSwapClear(t *testing.T) {
	a := New[int, int]()
	a.Set(1, 1)
	b := New[int, int]()
	b.Set(2, 2)
	b.Set(3, 3)
	a.Swap(b)
	if a.Size() != 2 || b.Size() != 1 {
		t.Errorf("after swap a=%d b=%d", a.Size(), b.Size())
	}
	a.Clear()
	if !a.Empty() {
		t.Error("should be empty after Clear")
	}
}

func TestForEach(t *testing.T) {
	m := New[int, int]()
	for i := 1; i <= 3; i++ {
		m.Set(i, i*10)
	}
	sum := 0
	m.ForEach(func(k, v int) { sum += k + v })
	if sum != (1+2+3)+(10+20+30) {
		t.Errorf("ForEach sum = %d", sum)
	}
}

func BenchmarkSet(b *testing.B) {
	m := New[int, int]()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Set(i, i)
	}
}

func BenchmarkRefIncrement(b *testing.B) {
	m := New[int, int]()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		*m.Ref(i % 1000)++
	}
}
