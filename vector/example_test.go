package vector_test

import (
	"fmt"

	"github.com/wiraphatys/gostd/vector"
)

func ExampleVector() {
	v := vector.New[int]()

	v.PushBack(10)
	v.PushBack(20)
	v.PushBack(30)

	fmt.Printf("Size: %d\n", v.Size())
	front, _ := v.Front()
	back, _ := v.Back()
	fmt.Printf("Front: %d, Back: %d\n", front, back)

	v.Set(1, 99)
	fmt.Printf("Element at 1: %d\n", v.Get(1))

	v.ForEach(func(i, val int) {
		fmt.Printf("[%d] = %d\n", i, val)
	})

	// Output:
	// Size: 3
	// Front: 10, Back: 30
	// Element at 1: 99
	// [0] = 10
	// [1] = 99
	// [2] = 30
}

func ExampleVector_Insert() {
	v := vector.NewFromSlice([]int{1, 2, 4})
	v.Insert(2, 3) // insert 3 before index 2
	fmt.Println(v.ToSlice())
	// Output: [1 2 3 4]
}

func ExampleVector_Erase() {
	v := vector.NewFromSlice([]int{1, 2, 3, 4, 5})
	v.EraseRange(1, 4) // remove indices 1,2,3
	fmt.Println(v.ToSlice())
	// Output: [1 5]
}
