// Package rbtree implements a generic red-black tree. It is the shared,
// internal balanced-search-tree engine behind the ordered containers set,
// multiset, treemap and multimap.
//
// The implementation follows the algorithms in Cormen, Leiserson, Rivest and
// Stein, "Introduction to Algorithms" (3rd ed.): a single black sentinel node
// stands in for every nil pointer, and deletion relinks nodes by pointer
// (rather than copying keys) so that the surviving in-order successor node is
// returned correctly by EraseNode — exactly what std::map::erase needs to
// return the iterator following the erased element.
//
// All operations run in O(log n):
//
//	Insert / Find / LowerBound / UpperBound / EraseKey / EraseNode   O(log n)
//	Next / Prev (in-order step)                                      O(log n) worst, O(1) amortized
//	Min / Max                                                        O(log n)
//
// Externally, a missing node or the past-the-end position is reported as a Go
// nil *Node. Internally the black sentinel is used so the balancing code needs
// no nil checks.
package rbtree

type color bool

const (
	red   color = false
	black color = true
)

// Node is a single tree node. Key is immutable from the container's point of
// view (changing it would break the ordering); Value is freely mutable, which
// is what lets treemap return an assignable reference from its operator[]
// equivalent.
type Node[K, V any] struct {
	left, right, parent *Node[K, V]
	color               color
	Key                 K
	Value               V
}

// Tree is a red-black tree mapping keys of type K to values of type V.
type Tree[K, V any] struct {
	root  *Node[K, V]
	nilN  *Node[K, V] // sentinel: the single black "nil" node
	size  int
	less  func(a, b K) bool
	multi bool
}

// New creates an empty red-black tree ordered by less. If multi is true the
// tree permits multiple nodes with equivalent keys (multiset/multimap);
// otherwise keys are unique (set/map).
func New[K, V any](less func(a, b K) bool, multi bool) *Tree[K, V] {
	sentinel := &Node[K, V]{color: black}
	sentinel.left = sentinel
	sentinel.right = sentinel
	sentinel.parent = sentinel
	return &Tree[K, V]{
		root:  sentinel,
		nilN:  sentinel,
		less:  less,
		multi: multi,
	}
}

// Size returns the number of nodes in the tree.
func (t *Tree[K, V]) Size() int { return t.size }

// Clear removes all nodes. O(1) — the subtree becomes garbage.
func (t *Tree[K, V]) Clear() {
	t.root = t.nilN
	t.size = 0
}

// ext converts the internal sentinel to a Go nil for callers; real nodes pass
// through unchanged.
func (t *Tree[K, V]) ext(n *Node[K, V]) *Node[K, V] {
	if n == t.nilN {
		return nil
	}
	return n
}

func (t *Tree[K, V]) leftRotate(x *Node[K, V]) {
	y := x.right
	x.right = y.left
	if y.left != t.nilN {
		y.left.parent = x
	}
	y.parent = x.parent
	switch {
	case x.parent == t.nilN:
		t.root = y
	case x == x.parent.left:
		x.parent.left = y
	default:
		x.parent.right = y
	}
	y.left = x
	x.parent = y
}

func (t *Tree[K, V]) rightRotate(x *Node[K, V]) {
	y := x.left
	x.left = y.right
	if y.right != t.nilN {
		y.right.parent = x
	}
	y.parent = x.parent
	switch {
	case x.parent == t.nilN:
		t.root = y
	case x == x.parent.right:
		x.parent.right = y
	default:
		x.parent.left = y
	}
	y.right = x
	x.parent = y
}

// Insert inserts key/value. In unique mode, if an equivalent key already
// exists, the tree is left unchanged and that existing node is returned with
// inserted=false. In multi mode a fresh node is always created. The returned
// node is always the node now holding key (existing or new), so callers may
// read or mutate its Value.
func (t *Tree[K, V]) Insert(key K, value V) (node *Node[K, V], inserted bool) {
	return t.insert(key, value, false)
}

// InsertOrAssign behaves like Insert in unique mode, but if an equivalent key
// already exists its Value is overwritten with value. Returns the node and
// whether a new node was created.
func (t *Tree[K, V]) InsertOrAssign(key K, value V) (node *Node[K, V], inserted bool) {
	return t.insert(key, value, true)
}

func (t *Tree[K, V]) insert(key K, value V, assign bool) (*Node[K, V], bool) {
	y := t.nilN
	x := t.root
	for x != t.nilN {
		y = x
		if t.less(key, x.Key) {
			x = x.left
		} else if !t.multi && !t.less(x.Key, key) {
			// Equivalent key in a unique tree.
			if assign {
				x.Value = value
			}
			return x, false
		} else {
			x = x.right
		}
	}

	z := &Node[K, V]{
		Key:    key,
		Value:  value,
		color:  red,
		left:   t.nilN,
		right:  t.nilN,
		parent: y,
	}
	switch {
	case y == t.nilN:
		t.root = z
	case t.less(key, y.Key):
		y.left = z
	default:
		y.right = z
	}
	t.size++
	t.insertFixup(z)
	return z, true
}

func (t *Tree[K, V]) insertFixup(z *Node[K, V]) {
	for z.parent.color == red {
		if z.parent == z.parent.parent.left {
			y := z.parent.parent.right
			if y.color == red {
				z.parent.color = black
				y.color = black
				z.parent.parent.color = red
				z = z.parent.parent
			} else {
				if z == z.parent.right {
					z = z.parent
					t.leftRotate(z)
				}
				z.parent.color = black
				z.parent.parent.color = red
				t.rightRotate(z.parent.parent)
			}
		} else {
			y := z.parent.parent.left
			if y.color == red {
				z.parent.color = black
				y.color = black
				z.parent.parent.color = red
				z = z.parent.parent
			} else {
				if z == z.parent.left {
					z = z.parent
					t.rightRotate(z)
				}
				z.parent.color = black
				z.parent.parent.color = red
				t.leftRotate(z.parent.parent)
			}
		}
	}
	t.root.color = black
}

// transplant replaces the subtree rooted at u with the subtree rooted at v.
func (t *Tree[K, V]) transplant(u, v *Node[K, V]) {
	switch {
	case u.parent == t.nilN:
		t.root = v
	case u == u.parent.left:
		u.parent.left = v
	default:
		u.parent.right = v
	}
	v.parent = u.parent
}

func (t *Tree[K, V]) minimum(x *Node[K, V]) *Node[K, V] {
	for x.left != t.nilN {
		x = x.left
	}
	return x
}

func (t *Tree[K, V]) maximum(x *Node[K, V]) *Node[K, V] {
	for x.right != t.nilN {
		x = x.right
	}
	return x
}

// successor returns the in-order successor of x, or the sentinel if none.
func (t *Tree[K, V]) successor(x *Node[K, V]) *Node[K, V] {
	if x.right != t.nilN {
		return t.minimum(x.right)
	}
	y := x.parent
	for y != t.nilN && x == y.right {
		x = y
		y = y.parent
	}
	return y
}

// predecessor returns the in-order predecessor of x, or the sentinel if none.
func (t *Tree[K, V]) predecessor(x *Node[K, V]) *Node[K, V] {
	if x.left != t.nilN {
		return t.maximum(x.left)
	}
	y := x.parent
	for y != t.nilN && x == y.left {
		x = y
		y = y.parent
	}
	return y
}

// EraseNode removes node z from the tree and returns its in-order successor
// (or nil if z was the last element). z must belong to this tree and not be
// nil. The pointer-relinking delete guarantees the returned successor node is
// the element that now follows the gap left by z, matching std::map::erase.
func (t *Tree[K, V]) EraseNode(z *Node[K, V]) *Node[K, V] {
	if z == nil || z == t.nilN {
		return nil
	}
	succ := t.successor(z)

	y := z
	yOriginalColor := y.color
	var x *Node[K, V]
	switch {
	case z.left == t.nilN:
		x = z.right
		t.transplant(z, z.right)
	case z.right == t.nilN:
		x = z.left
		t.transplant(z, z.left)
	default:
		y = t.minimum(z.right)
		yOriginalColor = y.color
		x = y.right
		if y.parent == z {
			x.parent = y
		} else {
			t.transplant(y, y.right)
			y.right = z.right
			y.right.parent = y
		}
		t.transplant(z, y)
		y.left = z.left
		y.left.parent = y
		y.color = z.color
	}
	if yOriginalColor == black {
		t.deleteFixup(x)
	}
	t.size--

	// Detach z to help the garbage collector and to surface misuse early.
	z.left, z.right, z.parent = nil, nil, nil
	return t.ext(succ)
}

func (t *Tree[K, V]) deleteFixup(x *Node[K, V]) {
	for x != t.root && x.color == black {
		if x == x.parent.left {
			w := x.parent.right
			if w.color == red {
				w.color = black
				x.parent.color = red
				t.leftRotate(x.parent)
				w = x.parent.right
			}
			if w.left.color == black && w.right.color == black {
				w.color = red
				x = x.parent
			} else {
				if w.right.color == black {
					w.left.color = black
					w.color = red
					t.rightRotate(w)
					w = x.parent.right
				}
				w.color = x.parent.color
				x.parent.color = black
				w.right.color = black
				t.leftRotate(x.parent)
				x = t.root
			}
		} else {
			w := x.parent.left
			if w.color == red {
				w.color = black
				x.parent.color = red
				t.rightRotate(x.parent)
				w = x.parent.left
			}
			if w.right.color == black && w.left.color == black {
				w.color = red
				x = x.parent
			} else {
				if w.left.color == black {
					w.right.color = black
					w.color = red
					t.leftRotate(w)
					w = x.parent.left
				}
				w.color = x.parent.color
				x.parent.color = black
				w.left.color = black
				t.rightRotate(x.parent)
				x = t.root
			}
		}
	}
	x.color = black
}

// findInternal returns the node holding key, or the sentinel if absent. In
// multi mode it returns the first equivalent node encountered on the path.
func (t *Tree[K, V]) findInternal(key K) *Node[K, V] {
	x := t.root
	for x != t.nilN {
		if t.less(key, x.Key) {
			x = x.left
		} else if t.less(x.Key, key) {
			x = x.right
		} else {
			return x
		}
	}
	return t.nilN
}

// Find returns the node holding key, or nil if absent.
func (t *Tree[K, V]) Find(key K) *Node[K, V] { return t.ext(t.findInternal(key)) }

// Contains reports whether key is present.
func (t *Tree[K, V]) Contains(key K) bool { return t.findInternal(key) != t.nilN }

// LowerBound returns the first node whose key is not less than key (i.e. >=),
// or nil if every key is less than key. Equivalent to C++ lower_bound.
func (t *Tree[K, V]) LowerBound(key K) *Node[K, V] {
	result := t.nilN
	x := t.root
	for x != t.nilN {
		if !t.less(x.Key, key) { // x.Key >= key
			result = x
			x = x.left
		} else {
			x = x.right
		}
	}
	return t.ext(result)
}

// UpperBound returns the first node whose key is greater than key, or nil if
// none. Equivalent to C++ upper_bound.
func (t *Tree[K, V]) UpperBound(key K) *Node[K, V] {
	result := t.nilN
	x := t.root
	for x != t.nilN {
		if t.less(key, x.Key) { // x.Key > key
			result = x
			x = x.left
		} else {
			x = x.right
		}
	}
	return t.ext(result)
}

// Count returns the number of nodes whose key is equivalent to key.
func (t *Tree[K, V]) Count(key K) int {
	if !t.multi {
		if t.findInternal(key) != t.nilN {
			return 1
		}
		return 0
	}
	count := 0
	lo := t.LowerBound(key)
	for n := lo; n != nil; n = t.Next(n) {
		if t.less(key, n.Key) {
			break
		}
		count++
	}
	return count
}

// EraseKey removes every node whose key is equivalent to key and returns the
// number removed (0 or 1 in unique mode). O((1+count)·log n).
func (t *Tree[K, V]) EraseKey(key K) int {
	count := 0
	for {
		n := t.findInternal(key)
		if n == t.nilN {
			break
		}
		t.EraseNode(n)
		count++
	}
	return count
}

// Min returns the node with the smallest key, or nil if empty.
func (t *Tree[K, V]) Min() *Node[K, V] {
	if t.root == t.nilN {
		return nil
	}
	return t.minimum(t.root)
}

// Max returns the node with the largest key, or nil if empty.
func (t *Tree[K, V]) Max() *Node[K, V] {
	if t.root == t.nilN {
		return nil
	}
	return t.maximum(t.root)
}

// Next returns the in-order successor of n, or nil if n is the last node.
func (t *Tree[K, V]) Next(n *Node[K, V]) *Node[K, V] {
	if n == nil || n == t.nilN {
		return nil
	}
	return t.ext(t.successor(n))
}

// Prev returns the in-order predecessor of n, or nil if n is the first node.
func (t *Tree[K, V]) Prev(n *Node[K, V]) *Node[K, V] {
	if n == nil || n == t.nilN {
		return nil
	}
	return t.ext(t.predecessor(n))
}
