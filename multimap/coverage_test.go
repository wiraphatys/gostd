package multimap

import (
	"reflect"
	"testing"
)

func TestIteratorAndStateMethods(t *testing.T) {
	m := New[int, string]()
	if !m.Empty() {
		t.Error("new multimap should be Empty")
	}
	if !m.Begin().Equal(m.End()) {
		t.Error("empty Begin should equal End")
	}
	if !m.Begin().IsEnd() {
		t.Error("empty Begin should be IsEnd")
	}

	m.Insert(1, "a")
	m.Insert(2, "b")
	m.Insert(2, "c")

	if !m.Contains(2) || m.Contains(99) {
		t.Error("Contains wrong")
	}

	// Key/Value via Begin.
	it := m.Begin()
	if it.Key() != 1 || it.Value() != "a" {
		t.Errorf("Begin = (%d,%q), want (1,a)", it.Key(), it.Value())
	}

	// SetValue through an iterator mutates in place.
	it.SetValue("A")
	if v := m.Begin().Value(); v != "A" {
		t.Errorf("after SetValue, value = %q, want A", v)
	}

	// Reverse traversal via RBegin/REnd/Prev.
	var keys []int
	for rit := m.RBegin(); !rit.Equal(m.REnd()); rit.Prev() {
		keys = append(keys, rit.Key())
	}
	if !reflect.DeepEqual(keys, []int{2, 2, 1}) {
		t.Errorf("reverse keys = %v, want [2 2 1]", keys)
	}

	m.Clear()
	if !m.Empty() {
		t.Error("Clear should empty the multimap")
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
	m.Insert(1, 1)
	m.Insert(2, 2)
	if !m.EraseIter(m.End()).IsEnd() {
		t.Error("EraseIter(End) should return End")
	}
	if m.Size() != 2 {
		t.Error("EraseIter(End) mutated the multimap")
	}
	m.Swap(nil)
	if m.Size() != 2 {
		t.Error("Swap(nil) mutated the multimap")
	}
}
