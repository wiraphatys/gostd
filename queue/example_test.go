package queue_test

import (
	"fmt"

	"github.com/wiraphatys/gostd/queue"
)

func ExampleQueue() {
	// Create a new queue
	q := queue.New[int]()

	// Push elements to the queue
	q.Push(10)
	q.Push(20)
	q.Push(30)

	// Check if queue is empty
	fmt.Printf("Is empty: %v\n", q.Empty())

	// Get the size of the queue
	fmt.Printf("Size: %d\n", q.Size())

	// Access front and back elements
	front, _ := q.Front()
	back, _ := q.Back()
	fmt.Printf("Front: %d, Back: %d\n", front, back)

	// Pop elements from the queue (FIFO)
	for !q.Empty() {
		front, _ := q.Front()
		fmt.Printf("Popping: %d\n", front)
		q.Pop()
	}

	// Output:
	// Is empty: false
	// Size: 3
	// Front: 10, Back: 30
	// Popping: 10
	// Popping: 20
	// Popping: 30
}

func ExampleQueue_stringQueue() {
	// Create a queue of strings
	q := queue.New[string]()

	// Add some strings
	q.Push("first")
	q.Push("second")
	q.Push("third")

	// Process all elements
	for !q.Empty() {
		item, _ := q.Front()
		fmt.Printf("Processing: %s\n", item)
		q.Pop()
	}

	// Output:
	// Processing: first
	// Processing: second
	// Processing: third
}

func ExampleNewFromSlice() {
	// Create a queue from an existing slice
	items := []int{1, 2, 3, 4, 5}
	q := queue.NewFromSlice(items)

	fmt.Printf("Queue size: %d\n", q.Size())

	front, _ := q.Front()
	back, _ := q.Back()
	fmt.Printf("Front: %d, Back: %d\n", front, back)

	// Output:
	// Queue size: 5
	// Front: 1, Back: 5
}

func ExampleQueue_Swap() {
	// Create two queues
	q1 := queue.New[int]()
	q2 := queue.New[int]()

	q1.Push(1)
	q1.Push(2)

	q2.Push(10)
	q2.Push(20)
	q2.Push(30)

	fmt.Printf("Before swap - Q1 size: %d, Q2 size: %d\n", q1.Size(), q2.Size())

	// Swap the queues
	q1.Swap(q2)

	fmt.Printf("After swap - Q1 size: %d, Q2 size: %d\n", q1.Size(), q2.Size())

	front1, _ := q1.Front()
	front2, _ := q2.Front()
	fmt.Printf("Q1 front: %d, Q2 front: %d\n", front1, front2)

	// Output:
	// Before swap - Q1 size: 2, Q2 size: 3
	// After swap - Q1 size: 3, Q2 size: 2
	// Q1 front: 10, Q2 front: 1
}
