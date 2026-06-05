package treemap

import "testing"

func TestCountAndContains(t *testing.T) {
	m := New[int, int]()
	m.Set(1, 1)
	m.Set(2, 2)
	if m.Count(1) != 1 {
		t.Errorf("Count(1) = %d, want 1", m.Count(1))
	}
	if m.Count(99) != 0 {
		t.Errorf("Count(99) = %d, want 0", m.Count(99))
	}
	if !m.Contains(2) || m.Contains(99) {
		t.Error("Contains wrong")
	}
}

func TestEqualRange(t *testing.T) {
	m := New[int, int]()
	for _, k := range []int{10, 20, 30} {
		m.Set(k, k)
	}
	// Unique map: equal range of a present key spans exactly one element.
	lo, hi := m.EqualRange(20)
	count := 0
	for it := lo; !it.Equal(hi); it.Next() {
		if it.Key() != 20 {
			t.Errorf("EqualRange yielded key %d, want 20", it.Key())
		}
		count++
	}
	if count != 1 {
		t.Errorf("EqualRange(20) spanned %d entries, want 1", count)
	}
	// Absent key: empty range.
	lo, hi = m.EqualRange(25)
	if !lo.Equal(hi) {
		t.Error("EqualRange of absent key should be empty")
	}
}

func TestGuards(t *testing.T) {
	func() {
		defer func() {
			if recover() == nil {
				t.Error("NewFunc(nil) should panic")
			}
		}()
		NewFunc[int, int](nil)
	}()

	m := New[int, int]()
	m.Set(1, 1)
	m.Set(2, 2)
	if !m.EraseIter(m.End()).IsEnd() {
		t.Error("EraseIter(End) should return End")
	}
	if m.Size() != 2 {
		t.Error("EraseIter(End) mutated the map")
	}
	m.Swap(nil)
	if m.Size() != 2 {
		t.Error("Swap(nil) mutated the map")
	}
}
