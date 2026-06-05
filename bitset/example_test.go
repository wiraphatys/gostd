package bitset_test

import (
	"fmt"

	"github.com/wiraphatys/gostd/bitset"
)

func ExampleBitSet() {
	b := bitset.New(8)
	b.Set(0)
	b.Set(2)
	b.Set(7)

	fmt.Println("string:", b.ToString())
	fmt.Println("count:", b.Count())
	fmt.Println("test 2:", b.Test(2))
	// Output:
	// string: 10000101
	// count: 3
	// test 2: true
}

func ExampleBitSet_bitwise() {
	a, _ := bitset.NewFromString("1100")
	b, _ := bitset.NewFromString("1010")

	fmt.Println("a AND b:", a.Clone().And(b).ToString())
	fmt.Println("a OR  b:", a.Clone().Or(b).ToString())
	fmt.Println("a XOR b:", a.Clone().Xor(b).ToString())
	// Output:
	// a AND b: 1000
	// a OR  b: 1110
	// a XOR b: 0110
}
