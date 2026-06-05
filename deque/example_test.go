package deque_test

import (
	"fmt"

	"github.com/wiraphatys/gostd/deque"
)

func ExampleDeque() {
	d := deque.New[int]()

	d.PushBack(1)  // [1]
	d.PushBack(2)  // [1 2]
	d.PushFront(0) // [0 1 2]

	front, _ := d.Front()
	back, _ := d.Back()
	fmt.Printf("front=%d back=%d size=%d\n", front, back, d.Size())

	d.PopFront()
	fmt.Println(d.ToSlice())

	// Output:
	// front=0 back=2 size=3
	// [1 2]
}

func ExampleDeque_asQueue() {
	// A deque works as a FIFO queue: push at the back, pop from the front.
	d := deque.New[string]()
	d.PushBack("a")
	d.PushBack("b")
	d.PushBack("c")
	for !d.Empty() {
		front, _ := d.Front()
		fmt.Print(front)
		d.PopFront()
	}
	fmt.Println()
	// Output: abc
}
