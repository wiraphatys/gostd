package stack_test

import (
	"fmt"

	"github.com/wiraphatys/gostd/stack"
)

func ExampleStack() {
	// Create a new stack
	s := stack.New[int]()

	// Push elements to the stack
	s.Push(10)
	s.Push(20)
	s.Push(30)

	// Check if stack is empty
	fmt.Printf("Is empty: %v\n", s.Empty())

	// Get the size of the stack
	fmt.Printf("Size: %d\n", s.Size())

	// Access top element
	top, _ := s.Top()
	fmt.Printf("Top: %d\n", top)

	// Pop elements from the stack (LIFO)
	for !s.Empty() {
		top, _ := s.Top()
		fmt.Printf("Popping: %d\n", top)
		s.Pop()
	}

	// Output:
	// Is empty: false
	// Size: 3
	// Top: 30
	// Popping: 30
	// Popping: 20
	// Popping: 10
}

func ExampleStack_stringStack() {
	// Create a stack of strings
	s := stack.New[string]()

	// Add some strings
	s.Push("first")
	s.Push("second")
	s.Push("third")

	// Process all elements (LIFO order)
	for !s.Empty() {
		item, _ := s.Top()
		fmt.Printf("Processing: %s\n", item)
		s.Pop()
	}

	// Output:
	// Processing: third
	// Processing: second
	// Processing: first
}

func ExampleNewFromSlice() {
	// Create a stack from an existing slice
	items := []int{1, 2, 3, 4, 5}
	s := stack.NewFromSlice(items)

	fmt.Printf("Stack size: %d\n", s.Size())

	// The last element of the slice becomes the top
	top, _ := s.Top()
	fmt.Printf("Top: %d\n", top)

	// Pop a few elements to show LIFO behavior
	s.Pop()
	s.Pop()
	top, _ = s.Top()
	fmt.Printf("Top after 2 pops: %d\n", top)

	// Output:
	// Stack size: 5
	// Top: 5
	// Top after 2 pops: 3
}

func ExampleStack_Swap() {
	// Create two stacks
	s1 := stack.New[int]()
	s2 := stack.New[int]()

	s1.Push(1)
	s1.Push(2)

	s2.Push(10)
	s2.Push(20)
	s2.Push(30)

	fmt.Printf("Before swap - S1 size: %d, S2 size: %d\n", s1.Size(), s2.Size())

	// Get tops before swap
	top1, _ := s1.Top()
	top2, _ := s2.Top()
	fmt.Printf("Before swap - S1 top: %d, S2 top: %d\n", top1, top2)

	// Swap the stacks
	s1.Swap(s2)

	fmt.Printf("After swap - S1 size: %d, S2 size: %d\n", s1.Size(), s2.Size())

	// Get tops after swap
	newTop1, _ := s1.Top()
	newTop2, _ := s2.Top()
	fmt.Printf("After swap - S1 top: %d, S2 top: %d\n", newTop1, newTop2)

	// Output:
	// Before swap - S1 size: 2, S2 size: 3
	// Before swap - S1 top: 2, S2 top: 30
	// After swap - S1 size: 3, S2 size: 2
	// After swap - S1 top: 30, S2 top: 2
}
