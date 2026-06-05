package treemap_test

import (
	"fmt"

	"github.com/wiraphatys/gostd/treemap"
)

func ExampleMap() {
	m := treemap.New[string, int]()
	m.Set("banana", 3)
	m.Set("apple", 5)
	m.Set("cherry", 7)

	// Entries are always visited in sorted key order.
	m.ForEach(func(k string, v int) {
		fmt.Printf("%s=%d\n", k, v)
	})
	// Output:
	// apple=5
	// banana=3
	// cherry=7
}

func ExampleMap_Ref() {
	// Ref works like C++ operator[]: it default-inserts and returns a pointer
	// you can write through. This makes counting idiomatic.
	counts := treemap.New[string, int]()
	for _, word := range []string{"a", "b", "a", "c", "a", "b"} {
		*counts.Ref(word)++
	}
	fmt.Println("a:", *counts.Ref("a"))
	fmt.Println("b:", *counts.Ref("b"))
	fmt.Println("c:", *counts.Ref("c"))
	// Output:
	// a: 3
	// b: 2
	// c: 1
}
