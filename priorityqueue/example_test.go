package priorityqueue_test

import (
	"fmt"

	"github.com/wiraphatys/gostd/priorityqueue"
)

func ExamplePriorityQueue() {
	// Default is a max-heap: the largest element is always on top.
	pq := priorityqueue.New[int]()
	pq.Push(3)
	pq.Push(1)
	pq.Push(4)
	pq.Push(1)
	pq.Push(5)

	for !pq.Empty() {
		top, _ := pq.Top()
		fmt.Print(top, " ")
		pq.Pop()
	}
	fmt.Println()
	// Output: 5 4 3 1 1
}

func ExampleNewMin() {
	// A min-heap returns the smallest element first.
	pq := priorityqueue.NewMin[int]()
	for _, x := range []int{3, 1, 4, 1, 5} {
		pq.Push(x)
	}
	for !pq.Empty() {
		top, _ := pq.Top()
		fmt.Print(top, " ")
		pq.Pop()
	}
	fmt.Println()
	// Output: 1 1 3 4 5
}

func ExampleNewFunc() {
	type job struct {
		name     string
		priority int
	}
	// Order jobs so the highest priority value comes out first.
	pq := priorityqueue.NewFunc(func(a, b job) bool { return a.priority < b.priority })
	pq.Push(job{"email", 2})
	pq.Push(job{"alarm", 9})
	pq.Push(job{"log", 1})

	top, _ := pq.Top()
	fmt.Println("most urgent:", top.name)
	// Output: most urgent: alarm
}
