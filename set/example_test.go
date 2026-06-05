package set_test

import (
	"fmt"

	"github.com/wiraphatys/gostd/set"
)

func ExampleSet() {
	s := set.New[int]()
	s.Insert(3)
	s.Insert(1)
	s.Insert(4)
	s.Insert(1) // duplicate, ignored

	fmt.Println("size:", s.Size())
	fmt.Println("contains 4:", s.Contains(4))
	fmt.Println("sorted:", s.ToSlice())
	// Output:
	// size: 3
	// contains 4: true
	// sorted: [1 3 4]
}

func ExampleSet_iterate() {
	s := set.NewFromSlice([]int{5, 2, 8, 2, 1})
	for it := s.Begin(); !it.Equal(s.End()); it.Next() {
		fmt.Print(it.Value(), " ")
	}
	fmt.Println()
	// Output: 1 2 5 8
}

func ExampleNewFunc() {
	// A set ordered from largest to smallest.
	s := set.NewFunc(func(a, b int) bool { return a > b })
	s.Insert(1)
	s.Insert(3)
	s.Insert(2)
	fmt.Println(s.ToSlice())
	// Output: [3 2 1]
}
