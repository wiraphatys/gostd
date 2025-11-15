package stack

import (
	"testing"
)

func TestNew(t *testing.T) {
	s := New[int]()

	if s == nil {
		t.Error("New() should not return nil")
	}

	if !s.Empty() {
		t.Error("New stack should be empty")
	}

	if s.Size() != 0 {
		t.Errorf("New stack size should be 0, got %d", s.Size())
	}
}

func TestNewFromSlice(t *testing.T) {
	items := []int{1, 2, 3, 4, 5}
	s := NewFromSlice(items)

	if s.Size() != len(items) {
		t.Errorf("Stack size should be %d, got %d", len(items), s.Size())
	}

	top, err := s.Top()
	if err != nil {
		t.Errorf("Top() should not return error: %v", err)
	}
	if top != 5 {
		t.Errorf("Top should be 5, got %d", top)
	}

	// Test that elements are in correct order (LIFO)
	expected := []int{5, 4, 3, 2, 1}
	for i, exp := range expected {
		top, err := s.Top()
		if err != nil {
			t.Errorf("Top() should not return error at iteration %d: %v", i, err)
		}
		if top != exp {
			t.Errorf("Top should be %d at iteration %d, got %d", exp, i, top)
		}
		s.Pop()
	}
}

func TestPush(t *testing.T) {
	s := New[int]()

	s.Push(10)
	if s.Empty() {
		t.Error("Stack should not be empty after push")
	}
	if s.Size() != 1 {
		t.Errorf("Stack size should be 1, got %d", s.Size())
	}

	s.Push(20)
	s.Push(30)
	if s.Size() != 3 {
		t.Errorf("Stack size should be 3, got %d", s.Size())
	}

	// Check that the last pushed element is on top
	top, err := s.Top()
	if err != nil {
		t.Errorf("Top() should not return error: %v", err)
	}
	if top != 30 {
		t.Errorf("Top should be 30, got %d", top)
	}
}

func TestPop(t *testing.T) {
	s := New[int]()

	// Test pop on empty stack
	err := s.Pop()
	if err == nil {
		t.Error("Pop on empty stack should return error")
	}

	// Test normal pop
	s.Push(10)
	s.Push(20)
	s.Push(30)

	err = s.Pop()
	if err != nil {
		t.Errorf("Pop should not return error: %v", err)
	}

	if s.Size() != 2 {
		t.Errorf("Stack size should be 2 after pop, got %d", s.Size())
	}

	top, err := s.Top()
	if err != nil {
		t.Errorf("Top() should not return error: %v", err)
	}
	if top != 20 {
		t.Errorf("Top should be 20 after pop, got %d", top)
	}
}

func TestTop(t *testing.T) {
	s := New[int]()

	// Test top on empty stack
	_, err := s.Top()
	if err == nil {
		t.Error("Top on empty stack should return error")
	}

	// Test normal top
	s.Push(10)
	s.Push(20)

	top, err := s.Top()
	if err != nil {
		t.Errorf("Top() should not return error: %v", err)
	}
	if top != 20 {
		t.Errorf("Top should be 20, got %d", top)
	}

	// Size should remain the same after Top()
	if s.Size() != 2 {
		t.Errorf("Stack size should remain 2, got %d", s.Size())
	}
}

func TestEmpty(t *testing.T) {
	s := New[int]()

	if !s.Empty() {
		t.Error("New stack should be empty")
	}

	s.Push(10)
	if s.Empty() {
		t.Error("Stack with elements should not be empty")
	}

	s.Pop()
	if !s.Empty() {
		t.Error("Stack should be empty after popping all elements")
	}
}

func TestSize(t *testing.T) {
	s := New[int]()

	if s.Size() != 0 {
		t.Errorf("Empty stack size should be 0, got %d", s.Size())
	}

	for i := 1; i <= 5; i++ {
		s.Push(i * 10)
		if s.Size() != i {
			t.Errorf("Stack size should be %d, got %d", i, s.Size())
		}
	}

	for i := 4; i >= 0; i-- {
		s.Pop()
		if s.Size() != i {
			t.Errorf("Stack size should be %d, got %d", i, s.Size())
		}
	}
}

func TestEmplace(t *testing.T) {
	s := New[int]()

	s.Emplace(100)
	if s.Size() != 1 {
		t.Errorf("Stack size should be 1 after emplace, got %d", s.Size())
	}

	top, err := s.Top()
	if err != nil {
		t.Errorf("Top() should not return error: %v", err)
	}
	if top != 100 {
		t.Errorf("Top should be 100, got %d", top)
	}
}

func TestSwap(t *testing.T) {
	s1 := New[int]()
	s2 := New[int]()

	// Test swap with nil
	s1.Push(10)
	s1.Swap(nil)
	if s1.Size() != 1 {
		t.Error("Swap with nil should not change stack")
	}

	// Test normal swap
	s1.Push(20)
	s2.Push(100)
	s2.Push(200)
	s2.Push(300)

	s1Size := s1.Size()
	s2Size := s2.Size()

	s1Top, _ := s1.Top()
	s2Top, _ := s2.Top()

	s1.Swap(s2)

	if s1.Size() != s2Size {
		t.Errorf("After swap, s1 size should be %d, got %d", s2Size, s1.Size())
	}
	if s2.Size() != s1Size {
		t.Errorf("After swap, s2 size should be %d, got %d", s1Size, s2.Size())
	}

	newS1Top, _ := s1.Top()
	newS2Top, _ := s2.Top()

	if newS1Top != s2Top {
		t.Errorf("After swap, s1 top should be %d, got %d", s2Top, newS1Top)
	}
	if newS2Top != s1Top {
		t.Errorf("After swap, s2 top should be %d, got %d", s1Top, newS2Top)
	}
}

func TestClear(t *testing.T) {
	s := New[int]()
	s.Push(10)
	s.Push(20)
	s.Push(30)

	if s.Size() != 3 {
		t.Errorf("Stack size should be 3 before clear, got %d", s.Size())
	}

	s.Clear()

	if !s.Empty() {
		t.Error("Stack should be empty after clear")
	}
	if s.Size() != 0 {
		t.Errorf("Stack size should be 0 after clear, got %d", s.Size())
	}
}

func TestToSlice(t *testing.T) {
	s := New[int]()
	items := []int{1, 2, 3, 4, 5}

	for _, item := range items {
		s.Push(item)
	}

	slice := s.ToSlice()

	if len(slice) != len(items) {
		t.Errorf("Slice length should be %d, got %d", len(items), len(slice))
	}

	// Slice should be ordered from bottom to top
	for i, item := range items {
		if slice[i] != item {
			t.Errorf("Slice[%d] should be %d, got %d", i, item, slice[i])
		}
	}

	// Modifying slice should not affect stack
	slice[0] = 999
	bottom := s.ToSlice()[0]
	if bottom == 999 {
		t.Error("Modifying slice should not affect original stack")
	}
}

func TestLIFOBehavior(t *testing.T) {
	s := New[int]()

	// Test LIFO (Last In, First Out) behavior
	items := []int{1, 2, 3, 4, 5}

	for _, item := range items {
		s.Push(item)
	}

	// Should pop in reverse order
	expected := []int{5, 4, 3, 2, 1}
	for i, exp := range expected {
		top, err := s.Top()
		if err != nil {
			t.Errorf("Top() should not return error at iteration %d: %v", i, err)
		}
		if top != exp {
			t.Errorf("Expected %d at iteration %d, got %d", exp, i, top)
		}

		err = s.Pop()
		if err != nil {
			t.Errorf("Pop() should not return error at iteration %d: %v", i, err)
		}
	}

	if !s.Empty() {
		t.Error("Stack should be empty after popping all elements")
	}
}

func TestStringStack(t *testing.T) {
	s := New[string]()

	s.Push("first")
	s.Push("second")
	s.Push("third")

	top, err := s.Top()
	if err != nil {
		t.Errorf("Top() should not return error: %v", err)
	}
	if top != "third" {
		t.Errorf("Top should be 'third', got '%s'", top)
	}

	// Test LIFO order
	s.Pop()
	top, err = s.Top()
	if err != nil {
		t.Errorf("Top() should not return error: %v", err)
	}
	if top != "second" {
		t.Errorf("Top should be 'second', got '%s'", top)
	}
}

func BenchmarkPush(b *testing.B) {
	s := New[int]()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		s.Push(i)
	}
}

func BenchmarkPop(b *testing.B) {
	s := New[int]()
	for i := 0; i < b.N; i++ {
		s.Push(i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Pop()
	}
}

func BenchmarkTop(b *testing.B) {
	s := New[int]()
	s.Push(42)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Top()
	}
}
