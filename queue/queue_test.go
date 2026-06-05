package queue

import (
	"testing"
)

func TestNew(t *testing.T) {
	q := New[int]()

	if q == nil {
		t.Error("New() should not return nil")
	}

	if !q.Empty() {
		t.Error("New queue should be empty")
	}

	if q.Size() != 0 {
		t.Errorf("New queue size should be 0, got %d", q.Size())
	}
}

func TestNewFromSlice(t *testing.T) {
	items := []int{1, 2, 3, 4, 5}
	q := NewFromSlice(items)

	if q.Size() != len(items) {
		t.Errorf("Queue size should be %d, got %d", len(items), q.Size())
	}

	front, err := q.Front()
	if err != nil {
		t.Errorf("Front() should not return error: %v", err)
	}
	if front != 1 {
		t.Errorf("Front should be 1, got %d", front)
	}

	back, err := q.Back()
	if err != nil {
		t.Errorf("Back() should not return error: %v", err)
	}
	if back != 5 {
		t.Errorf("Back should be 5, got %d", back)
	}
}

func TestPush(t *testing.T) {
	q := New[int]()

	q.Push(10)
	if q.Empty() {
		t.Error("Queue should not be empty after push")
	}
	if q.Size() != 1 {
		t.Errorf("Queue size should be 1, got %d", q.Size())
	}

	q.Push(20)
	q.Push(30)
	if q.Size() != 3 {
		t.Errorf("Queue size should be 3, got %d", q.Size())
	}
}

func TestPop(t *testing.T) {
	q := New[int]()

	// Test pop on empty queue
	err := q.Pop()
	if err == nil {
		t.Error("Pop on empty queue should return error")
	}

	// Test normal pop
	q.Push(10)
	q.Push(20)
	q.Push(30)

	err = q.Pop()
	if err != nil {
		t.Errorf("Pop should not return error: %v", err)
	}

	if q.Size() != 2 {
		t.Errorf("Queue size should be 2 after pop, got %d", q.Size())
	}

	front, err := q.Front()
	if err != nil {
		t.Errorf("Front() should not return error: %v", err)
	}
	if front != 20 {
		t.Errorf("Front should be 20 after pop, got %d", front)
	}
}

func TestFront(t *testing.T) {
	q := New[int]()

	// Test front on empty queue
	_, err := q.Front()
	if err == nil {
		t.Error("Front on empty queue should return error")
	}

	// Test normal front
	q.Push(10)
	q.Push(20)

	front, err := q.Front()
	if err != nil {
		t.Errorf("Front() should not return error: %v", err)
	}
	if front != 10 {
		t.Errorf("Front should be 10, got %d", front)
	}

	// Size should remain the same after Front()
	if q.Size() != 2 {
		t.Errorf("Queue size should remain 2, got %d", q.Size())
	}
}

func TestBack(t *testing.T) {
	q := New[int]()

	// Test back on empty queue
	_, err := q.Back()
	if err == nil {
		t.Error("Back on empty queue should return error")
	}

	// Test normal back
	q.Push(10)
	q.Push(20)
	q.Push(30)

	back, err := q.Back()
	if err != nil {
		t.Errorf("Back() should not return error: %v", err)
	}
	if back != 30 {
		t.Errorf("Back should be 30, got %d", back)
	}

	// Size should remain the same after Back()
	if q.Size() != 3 {
		t.Errorf("Queue size should remain 3, got %d", q.Size())
	}
}

func TestEmpty(t *testing.T) {
	q := New[int]()

	if !q.Empty() {
		t.Error("New queue should be empty")
	}

	q.Push(10)
	if q.Empty() {
		t.Error("Queue with elements should not be empty")
	}

	q.Pop()
	if !q.Empty() {
		t.Error("Queue should be empty after popping all elements")
	}
}

func TestSize(t *testing.T) {
	q := New[int]()

	if q.Size() != 0 {
		t.Errorf("Empty queue size should be 0, got %d", q.Size())
	}

	for i := 1; i <= 5; i++ {
		q.Push(i * 10)
		if q.Size() != i {
			t.Errorf("Queue size should be %d, got %d", i, q.Size())
		}
	}

	for i := 4; i >= 0; i-- {
		q.Pop()
		if q.Size() != i {
			t.Errorf("Queue size should be %d, got %d", i, q.Size())
		}
	}
}

func TestEmplace(t *testing.T) {
	q := New[int]()

	q.Emplace(100)
	if q.Size() != 1 {
		t.Errorf("Queue size should be 1 after emplace, got %d", q.Size())
	}

	front, err := q.Front()
	if err != nil {
		t.Errorf("Front() should not return error: %v", err)
	}
	if front != 100 {
		t.Errorf("Front should be 100, got %d", front)
	}
}

func TestSwap(t *testing.T) {
	q1 := New[int]()
	q2 := New[int]()

	// Test swap with nil
	q1.Push(10)
	q1.Swap(nil)
	if q1.Size() != 1 {
		t.Error("Swap with nil should not change queue")
	}

	// Test normal swap
	q1.Push(20)
	q2.Push(100)
	q2.Push(200)
	q2.Push(300)

	q1Size := q1.Size()
	q2Size := q2.Size()

	q1Front, _ := q1.Front()
	q2Front, _ := q2.Front()

	q1.Swap(q2)

	if q1.Size() != q2Size {
		t.Errorf("After swap, q1 size should be %d, got %d", q2Size, q1.Size())
	}
	if q2.Size() != q1Size {
		t.Errorf("After swap, q2 size should be %d, got %d", q1Size, q2.Size())
	}

	newQ1Front, _ := q1.Front()
	newQ2Front, _ := q2.Front()

	if newQ1Front != q2Front {
		t.Errorf("After swap, q1 front should be %d, got %d", q2Front, newQ1Front)
	}
	if newQ2Front != q1Front {
		t.Errorf("After swap, q2 front should be %d, got %d", q1Front, newQ2Front)
	}
}

func TestClear(t *testing.T) {
	q := New[int]()
	q.Push(10)
	q.Push(20)
	q.Push(30)

	if q.Size() != 3 {
		t.Errorf("Queue size should be 3 before clear, got %d", q.Size())
	}

	q.Clear()

	if !q.Empty() {
		t.Error("Queue should be empty after clear")
	}
	if q.Size() != 0 {
		t.Errorf("Queue size should be 0 after clear, got %d", q.Size())
	}
}

func TestToSlice(t *testing.T) {
	q := New[int]()
	items := []int{1, 2, 3, 4, 5}

	for _, item := range items {
		q.Push(item)
	}

	slice := q.ToSlice()

	if len(slice) != len(items) {
		t.Errorf("Slice length should be %d, got %d", len(items), len(slice))
	}

	for i, item := range items {
		if slice[i] != item {
			t.Errorf("Slice[%d] should be %d, got %d", i, item, slice[i])
		}
	}

	// Modifying slice should not affect queue
	slice[0] = 999
	front, _ := q.Front()
	if front == 999 {
		t.Error("Modifying slice should not affect original queue")
	}
}

func TestFIFOBehavior(t *testing.T) {
	q := New[int]()

	// Test FIFO (First In, First Out) behavior
	items := []int{1, 2, 3, 4, 5}

	for _, item := range items {
		q.Push(item)
	}

	for i, expected := range items {
		front, err := q.Front()
		if err != nil {
			t.Errorf("Front() should not return error at iteration %d: %v", i, err)
		}
		if front != expected {
			t.Errorf("Expected %d at iteration %d, got %d", expected, i, front)
		}

		err = q.Pop()
		if err != nil {
			t.Errorf("Pop() should not return error at iteration %d: %v", i, err)
		}
	}

	if !q.Empty() {
		t.Error("Queue should be empty after popping all elements")
	}
}

func TestStringQueue(t *testing.T) {
	q := New[string]()

	q.Push("hello")
	q.Push("world")

	front, err := q.Front()
	if err != nil {
		t.Errorf("Front() should not return error: %v", err)
	}
	if front != "hello" {
		t.Errorf("Front should be 'hello', got '%s'", front)
	}

	back, err := q.Back()
	if err != nil {
		t.Errorf("Back() should not return error: %v", err)
	}
	if back != "world" {
		t.Errorf("Back should be 'world', got '%s'", back)
	}
}

func BenchmarkPush(b *testing.B) {
	q := New[int]()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		q.Push(i)
	}
}

func BenchmarkPop(b *testing.B) {
	q := New[int]()
	for i := 0; i < b.N; i++ {
		q.Push(i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		q.Pop()
	}
}

func BenchmarkFront(b *testing.B) {
	q := New[int]()
	q.Push(42)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		q.Front()
	}
}
