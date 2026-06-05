package forwardlist_test

import (
	"fmt"

	"github.com/wiraphatys/gostd/forwardlist"
)

func ExampleForwardList() {
	fl := forwardlist.New[int]()
	fl.PushFront(3)
	fl.PushFront(2)
	fl.PushFront(1) // [1 2 3]

	for it := fl.Begin(); !it.Equal(fl.End()); it.Next() {
		fmt.Print(it.Value(), " ")
	}
	fmt.Println()
	// Output: 1 2 3
}

func ExampleForwardList_InsertAfter() {
	fl := forwardlist.NewFromSlice([]int{1, 3})
	it := fl.Begin()      // points at 1
	fl.InsertAfter(it, 2) // [1 2 3]
	fmt.Println(fl.ToSlice())
	// Output: [1 2 3]
}
