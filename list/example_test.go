package list_test

import (
	"fmt"

	"github.com/wiraphatys/gostd/list"
)

func ExampleList() {
	l := list.New[int]()
	l.PushBack(2)
	l.PushBack(3)
	l.PushFront(1) // [1 2 3]

	for it := l.Begin(); !it.Equal(l.End()); it.Next() {
		fmt.Print(it.Value(), " ")
	}
	fmt.Println()
	// Output: 1 2 3
}

func ExampleList_Sort() {
	l := list.NewFromSlice([]int{5, 3, 1, 4, 2})
	l.Sort(func(a, b int) bool { return a < b })
	fmt.Println(l.ToSlice())
	// Output: [1 2 3 4 5]
}

func ExampleList_Splice() {
	a := list.NewFromSlice([]int{1, 4})
	b := list.NewFromSlice([]int{2, 3})

	it := a.Begin()
	it.Next() // points at 4
	a.Splice(it, b)

	fmt.Println(a.ToSlice())
	fmt.Println("b empty:", b.Empty())
	// Output:
	// [1 2 3 4]
	// b empty: true
}
