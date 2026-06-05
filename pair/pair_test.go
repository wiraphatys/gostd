package pair

import (
	"fmt"
	"testing"
)

func TestNewAndFields(t *testing.T) {
	p := New(1, "one")
	if p.First != 1 {
		t.Errorf("First = %d, want 1", p.First)
	}
	if p.Second != "one" {
		t.Errorf("Second = %q, want \"one\"", p.Second)
	}
}

func TestMakePair(t *testing.T) {
	p := MakePair(3.14, true)
	if p.First != 3.14 || p.Second != true {
		t.Errorf("MakePair = %+v, want {3.14 true}", p)
	}
}

func TestCompositeLiteral(t *testing.T) {
	p := Pair[string, int]{First: "age", Second: 30}
	if p.First != "age" || p.Second != 30 {
		t.Errorf("composite literal = %+v", p)
	}
}

func TestUnpack(t *testing.T) {
	p := New("key", 42)
	k, v := p.Unpack()
	if k != "key" || v != 42 {
		t.Errorf("Unpack = (%q, %d), want (\"key\", 42)", k, v)
	}
}

func TestSwap(t *testing.T) {
	a := New(1, "a")
	b := New(2, "b")
	a.Swap(&b)
	if a.First != 2 || a.Second != "b" {
		t.Errorf("after swap a = %+v, want {2 b}", a)
	}
	if b.First != 1 || b.Second != "a" {
		t.Errorf("after swap b = %+v, want {1 a}", b)
	}

	// Swapping with nil is a no-op.
	a.Swap(nil)
	if a.First != 2 || a.Second != "b" {
		t.Errorf("swap with nil mutated pair: %+v", a)
	}
}

func TestEqual(t *testing.T) {
	if !Equal(New(1, "x"), New(1, "x")) {
		t.Error("equal pairs reported unequal")
	}
	if Equal(New(1, "x"), New(1, "y")) {
		t.Error("pairs differing in Second reported equal")
	}
	if Equal(New(1, "x"), New(2, "x")) {
		t.Error("pairs differing in First reported equal")
	}
}

func TestLess(t *testing.T) {
	// First member decides when different.
	if !Less(New(1, 100), New(2, 0)) {
		t.Error("(1,100) should be < (2,0)")
	}
	// Second member breaks ties.
	if !Less(New(1, 1), New(1, 2)) {
		t.Error("(1,1) should be < (1,2)")
	}
	if Less(New(1, 2), New(1, 1)) {
		t.Error("(1,2) should not be < (1,1)")
	}
	// Equal pairs are not less.
	if Less(New(1, 1), New(1, 1)) {
		t.Error("(1,1) should not be < (1,1)")
	}
}

// Pairs of comparable types must be usable as map keys.
func TestPairAsMapKey(t *testing.T) {
	m := map[Pair[int, int]]string{}
	m[New(1, 2)] = "a"
	m[New(3, 4)] = "b"
	if m[New(1, 2)] != "a" {
		t.Error("pair map key lookup failed")
	}
	if len(m) != 2 {
		t.Errorf("map len = %d, want 2", len(m))
	}
}

func ExamplePair() {
	p := New("hello", 5)
	fmt.Printf("%s has length %d\n", p.First, p.Second)
	// Output: hello has length 5
}
