package rbtree

import (
	"math/rand"
	"sort"
	"testing"
)

func intLess(a, b int) bool { return a < b }

func newUnique() *Tree[int, int] { return New[int, int](intLess, false) }
func newMulti() *Tree[int, int]  { return New[int, int](intLess, true) }

// validate asserts every red-black and BST invariant holds, failing tb on any
// violation. It returns nothing; correctness is signalled through tb.
func (t *Tree[K, V]) validate(tb testing.TB) {
	tb.Helper()
	if t.root.color != black {
		tb.Fatal("invariant: root is not black")
	}
	if t.nilN.color != black {
		tb.Fatal("invariant: sentinel is not black")
	}
	if t.root != t.nilN && t.root.parent != t.nilN {
		tb.Fatal("invariant: root.parent is not the sentinel")
	}
	t.checkSubtree(tb, t.root)

	// Global in-order ordering and size agreement.
	count := 0
	var prev *Node[K, V]
	havePrev := false
	for n := t.Min(); n != nil; n = t.Next(n) {
		if havePrev {
			if t.less(n.Key, prev.Key) {
				tb.Fatalf("invariant: in-order not sorted (%v before %v)", prev.Key, n.Key)
			}
			if !t.multi && !t.less(prev.Key, n.Key) {
				tb.Fatalf("invariant: duplicate key %v in unique tree", n.Key)
			}
		}
		prev = n
		havePrev = true
		count++
	}
	if count != t.size {
		tb.Fatalf("invariant: in-order count %d != size %d", count, t.size)
	}
}

// checkSubtree returns the black-height of the subtree (counting the sentinel
// leaf as black) and verifies the red and child-ordering properties.
func (t *Tree[K, V]) checkSubtree(tb testing.TB, n *Node[K, V]) int {
	if n == t.nilN {
		return 1
	}
	if n.color == red {
		if n.left.color != black || n.right.color != black {
			tb.Fatalf("invariant: red node %v has a red child", n.Key)
		}
	}
	if n.left != t.nilN {
		if n.left.parent != n {
			tb.Fatalf("invariant: left child of %v has wrong parent", n.Key)
		}
		if t.less(n.Key, n.left.Key) {
			tb.Fatalf("invariant: left child %v > parent %v", n.left.Key, n.Key)
		}
	}
	if n.right != t.nilN {
		if n.right.parent != n {
			tb.Fatalf("invariant: right child of %v has wrong parent", n.Key)
		}
		if t.less(n.right.Key, n.Key) {
			tb.Fatalf("invariant: right child %v < parent %v", n.right.Key, n.Key)
		}
	}
	lh := t.checkSubtree(tb, n.left)
	rh := t.checkSubtree(tb, n.right)
	if lh != rh {
		tb.Fatalf("invariant: black-height mismatch at %v (%d vs %d)", n.Key, lh, rh)
	}
	if n.color == black {
		return lh + 1
	}
	return lh
}

func (t *Tree[K, V]) keysInOrder() []K {
	var out []K
	for n := t.Min(); n != nil; n = t.Next(n) {
		out = append(out, n.Key)
	}
	return out
}

func TestEmpty(t *testing.T) {
	tr := newUnique()
	tr.validate(t)
	if tr.Size() != 0 {
		t.Error("empty tree size != 0")
	}
	if tr.Find(1) != nil {
		t.Error("Find on empty should be nil")
	}
	if tr.Min() != nil || tr.Max() != nil {
		t.Error("Min/Max on empty should be nil")
	}
	if tr.EraseKey(5) != 0 {
		t.Error("EraseKey on empty should be 0")
	}
}

func TestInsertFindUnique(t *testing.T) {
	tr := newUnique()
	for _, k := range []int{5, 3, 8, 1, 4, 7, 9, 2, 6} {
		n, inserted := tr.Insert(k, k*10)
		if !inserted {
			t.Errorf("Insert(%d) reported not inserted", k)
		}
		if n.Key != k || n.Value != k*10 {
			t.Errorf("Insert returned wrong node: %+v", n)
		}
		tr.validate(t)
	}
	// Duplicate insert is rejected, value preserved.
	n, inserted := tr.Insert(5, 999)
	if inserted {
		t.Error("duplicate Insert reported inserted")
	}
	if n.Value != 50 {
		t.Errorf("duplicate Insert overwrote value: %d", n.Value)
	}
	// InsertOrAssign overwrites.
	n, inserted = tr.InsertOrAssign(5, 555)
	if inserted || n.Value != 555 {
		t.Errorf("InsertOrAssign failed: inserted=%v value=%d", inserted, n.Value)
	}
	if got := tr.keysInOrder(); !sortedEqual(got, []int{1, 2, 3, 4, 5, 6, 7, 8, 9}) {
		t.Errorf("in-order keys = %v", got)
	}
}

func TestBounds(t *testing.T) {
	tr := newUnique()
	for _, k := range []int{10, 20, 30, 40, 50} {
		tr.Insert(k, 0)
	}
	cases := []struct {
		key            int
		wantLB, wantUB int // -1 means nil
	}{
		{5, 10, 10},
		{10, 10, 20},
		{25, 30, 30},
		{50, 50, -1},
		{55, -1, -1},
	}
	for _, c := range cases {
		lb := tr.LowerBound(c.key)
		ub := tr.UpperBound(c.key)
		if !nodeKeyEq(lb, c.wantLB) {
			t.Errorf("LowerBound(%d) = %v, want %d", c.key, nodeKey(lb), c.wantLB)
		}
		if !nodeKeyEq(ub, c.wantUB) {
			t.Errorf("UpperBound(%d) = %v, want %d", c.key, nodeKey(ub), c.wantUB)
		}
	}
}

func TestMultiInsertCountErase(t *testing.T) {
	tr := newMulti()
	for _, k := range []int{5, 5, 5, 3, 3, 8} {
		_, inserted := tr.Insert(k, 0)
		if !inserted {
			t.Error("multi Insert should always insert")
		}
		tr.validate(t)
	}
	if tr.Size() != 6 {
		t.Errorf("size = %d, want 6", tr.Size())
	}
	if c := tr.Count(5); c != 3 {
		t.Errorf("Count(5) = %d, want 3", c)
	}
	if c := tr.Count(3); c != 2 {
		t.Errorf("Count(3) = %d, want 2", c)
	}
	if c := tr.Count(99); c != 0 {
		t.Errorf("Count(99) = %d, want 0", c)
	}
	// EraseKey removes all equivalent.
	if removed := tr.EraseKey(5); removed != 3 {
		t.Errorf("EraseKey(5) = %d, want 3", removed)
	}
	tr.validate(t)
	if tr.Count(5) != 0 {
		t.Error("5 should be gone")
	}
	if tr.Size() != 3 {
		t.Errorf("size = %d, want 3", tr.Size())
	}
}

func TestEraseNodeReturnsSuccessor(t *testing.T) {
	tr := newUnique()
	for _, k := range []int{1, 2, 3, 4, 5} {
		tr.Insert(k, 0)
	}
	// Erase 3; successor should be 4.
	n := tr.Find(3)
	succ := tr.EraseNode(n)
	if succ == nil || succ.Key != 4 {
		t.Errorf("successor after erasing 3 = %v, want 4", nodeKey(succ))
	}
	tr.validate(t)
	// Erase the max; successor should be nil.
	max := tr.Find(5)
	if s := tr.EraseNode(max); s != nil {
		t.Errorf("successor after erasing max = %v, want nil", nodeKey(s))
	}
	tr.validate(t)
}

func TestNextPrevTraversal(t *testing.T) {
	tr := newUnique()
	for _, k := range []int{3, 1, 2, 5, 4} {
		tr.Insert(k, 0)
	}
	// Forward.
	var fwd []int
	for n := tr.Min(); n != nil; n = tr.Next(n) {
		fwd = append(fwd, n.Key)
	}
	if !sortedEqual(fwd, []int{1, 2, 3, 4, 5}) {
		t.Errorf("forward = %v", fwd)
	}
	// Backward.
	var bwd []int
	for n := tr.Max(); n != nil; n = tr.Prev(n) {
		bwd = append(bwd, n.Key)
	}
	if !sortedEqual(bwd, []int{5, 4, 3, 2, 1}) {
		t.Errorf("backward = %v", bwd)
	}
}

// Differential fuzz: unique tree vs a Go map reference, validating invariants
// throughout. This is the decisive correctness test for the balancing code.
func TestFuzzUnique(t *testing.T) {
	rng := rand.New(rand.NewSource(2024))
	tr := newUnique()
	ref := map[int]bool{}

	for step := 0; step < 30000; step++ {
		k := rng.Intn(500)
		if rng.Intn(2) == 0 {
			_, inserted := tr.Insert(k, k)
			if inserted == ref[k] {
				t.Fatalf("step %d: Insert(%d) inserted=%v but ref present=%v", step, k, inserted, ref[k])
			}
			ref[k] = true
		} else {
			removed := tr.EraseKey(k)
			wantRemoved := 0
			if ref[k] {
				wantRemoved = 1
			}
			if removed != wantRemoved {
				t.Fatalf("step %d: EraseKey(%d) = %d, want %d", step, k, removed, wantRemoved)
			}
			delete(ref, k)
		}

		if tr.Size() != len(ref) {
			t.Fatalf("step %d: size %d != ref %d", step, tr.Size(), len(ref))
		}
		if step%500 == 0 {
			tr.validate(t)
			// Compare full ordering.
			want := make([]int, 0, len(ref))
			for k := range ref {
				want = append(want, k)
			}
			sort.Ints(want)
			if got := tr.keysInOrder(); !sortedEqual(got, want) {
				t.Fatalf("step %d: keys mismatch", step)
			}
		}
	}
	tr.validate(t)
}

// Differential fuzz for multi mode with EraseNode removing a single occurrence.
func TestFuzzMulti(t *testing.T) {
	rng := rand.New(rand.NewSource(7))
	tr := newMulti()
	counts := map[int]int{}
	total := 0

	for step := 0; step < 20000; step++ {
		k := rng.Intn(100)
		switch rng.Intn(3) {
		case 0, 1:
			tr.Insert(k, 0)
			counts[k]++
			total++
		case 2:
			n := tr.Find(k)
			if n != nil {
				tr.EraseNode(n)
				counts[k]--
				if counts[k] == 0 {
					delete(counts, k)
				}
				total--
			}
		}
		if tr.Size() != total {
			t.Fatalf("step %d: size %d != total %d", step, tr.Size(), total)
		}
		if step%500 == 0 {
			tr.validate(t)
			for k, want := range counts {
				if got := tr.Count(k); got != want {
					t.Fatalf("step %d: Count(%d) = %d, want %d", step, k, got, want)
				}
			}
		}
	}
	tr.validate(t)
}

// --- helpers ---

func nodeKey(n *Node[int, int]) interface{} {
	if n == nil {
		return nil
	}
	return n.Key
}

func nodeKeyEq(n *Node[int, int], want int) bool {
	if want == -1 {
		return n == nil
	}
	return n != nil && n.Key == want
}

func sortedEqual(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func BenchmarkInsert(b *testing.B) {
	rng := rand.New(rand.NewSource(1))
	tr := newUnique()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tr.Insert(rng.Int(), 0)
	}
}

func BenchmarkFind(b *testing.B) {
	tr := newUnique()
	for i := 0; i < 100000; i++ {
		tr.Insert(i, 0)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tr.Find(i % 100000)
	}
}
