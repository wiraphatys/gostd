// Package list implements a doubly linked list that mirrors the API and
// behaviour of C++ std::list.
//
// It uses a single sentinel node whose next points at the front and whose prev
// points at the back, so the list is a circular ring. This gives:
//
//	PushFront / PushBack / PopFront / PopBack   O(1)
//	Insert / Erase given an iterator            O(1)
//	Splice (whole list)                         O(1)
//	Sort                                        O(n log n), in-place merge sort
//	Reverse / Unique / RemoveIf / Merge         O(n)
//
// Iterators returned by Begin/End/etc. remain valid across insertions and
// across erasures of *other* elements, exactly like std::list iterators. Using
// an iterator to an erased element is undefined behaviour, as in C++.
package list

import "errors"

// node is an internal list element.
type node[T any] struct {
	prev, next *node[T]
	value      T
}

// List is a doubly linked list. Create one with New or NewFromSlice; the zero
// value is not ready for use.
type List[T any] struct {
	sentinel *node[T]
	size     int
}

// Iterator is a bidirectional cursor over a List. It is the Go analogue of a
// std::list iterator.
type Iterator[T any] struct {
	node *node[T]
}

// Value returns the element the iterator refers to. Calling it on the past-the-
// end iterator (End) is undefined behaviour.
func (it Iterator[T]) Value() T { return it.node.value }

// SetValue overwrites the element the iterator refers to.
func (it Iterator[T]) SetValue(v T) { it.node.value = v }

// Next advances the iterator to the following element.
func (it *Iterator[T]) Next() { it.node = it.node.next }

// Prev moves the iterator to the preceding element.
func (it *Iterator[T]) Prev() { it.node = it.node.prev }

// Equal reports whether two iterators refer to the same position.
func (it Iterator[T]) Equal(other Iterator[T]) bool { return it.node == other.node }

// New creates and returns a new empty list.
func New[T any]() *List[T] {
	l := &List[T]{}
	l.init()
	return l
}

// NewFromSlice creates a list initialized with a copy of the given slice,
// preserving order.
func NewFromSlice[T any](items []T) *List[T] {
	l := New[T]()
	for _, v := range items {
		l.PushBack(v)
	}
	return l
}

func (l *List[T]) init() {
	s := &node[T]{}
	s.next = s
	s.prev = s
	l.sentinel = s
	l.size = 0
}

// linkBefore inserts n immediately before at.
func linkBefore[T any](n, at *node[T]) {
	n.prev = at.prev
	n.next = at
	at.prev.next = n
	at.prev = n
}

// unlink removes n from its list.
func unlink[T any](n *node[T]) {
	n.prev.next = n.next
	n.next.prev = n.prev
}

// Size returns the number of elements. Equivalent to C++ size().
func (l *List[T]) Size() int { return l.size }

// Empty reports whether the list has no elements. Equivalent to empty().
func (l *List[T]) Empty() bool { return l.size == 0 }

// Begin returns an iterator to the first element (or End if empty).
func (l *List[T]) Begin() Iterator[T] { return Iterator[T]{node: l.sentinel.next} }

// End returns an iterator to the past-the-end position (the sentinel).
func (l *List[T]) End() Iterator[T] { return Iterator[T]{node: l.sentinel} }

// RBegin returns an iterator to the last element, for reverse traversal with
// Prev (or End/the sentinel if the list is empty).
func (l *List[T]) RBegin() Iterator[T] { return Iterator[T]{node: l.sentinel.prev} }

// REnd returns the position before the first element (the sentinel), used as
// the stop condition for reverse traversal.
func (l *List[T]) REnd() Iterator[T] { return Iterator[T]{node: l.sentinel} }

// PushBack appends value to the end in O(1). Equivalent to push_back().
func (l *List[T]) PushBack(value T) {
	linkBefore(&node[T]{value: value}, l.sentinel)
	l.size++
}

// PushFront prepends value to the front in O(1). Equivalent to push_front().
func (l *List[T]) PushFront(value T) {
	linkBefore(&node[T]{value: value}, l.sentinel.next)
	l.size++
}

// EmplaceBack constructs an element in place at the back. Equivalent to
// PushBack in Go. Matches C++ emplace_back().
func (l *List[T]) EmplaceBack(value T) { l.PushBack(value) }

// EmplaceFront constructs an element in place at the front. Equivalent to
// PushFront in Go. Matches C++ emplace_front().
func (l *List[T]) EmplaceFront(value T) { l.PushFront(value) }

// PopFront removes the first element in O(1). Equivalent to pop_front().
func (l *List[T]) PopFront() error {
	if l.size == 0 {
		return errors.New("list: pop_front on empty list")
	}
	n := l.sentinel.next
	unlink(n)
	n.prev, n.next = nil, nil
	l.size--
	return nil
}

// PopBack removes the last element in O(1). Equivalent to pop_back().
func (l *List[T]) PopBack() error {
	if l.size == 0 {
		return errors.New("list: pop_back on empty list")
	}
	n := l.sentinel.prev
	unlink(n)
	n.prev, n.next = nil, nil
	l.size--
	return nil
}

// Front returns the first element. Equivalent to C++ front().
func (l *List[T]) Front() (T, error) {
	if l.size == 0 {
		var zero T
		return zero, errors.New("list: front of empty list")
	}
	return l.sentinel.next.value, nil
}

// Back returns the last element. Equivalent to C++ back().
func (l *List[T]) Back() (T, error) {
	if l.size == 0 {
		var zero T
		return zero, errors.New("list: back of empty list")
	}
	return l.sentinel.prev.value, nil
}

// Insert inserts value before the position pos refers to and returns an
// iterator to the new element. O(1). Equivalent to C++ insert(pos, value).
func (l *List[T]) Insert(pos Iterator[T], value T) Iterator[T] {
	n := &node[T]{value: value}
	linkBefore(n, pos.node)
	l.size++
	return Iterator[T]{node: n}
}

// Erase removes the element pos refers to and returns an iterator to the next
// element. O(1). Erasing End is a no-op. Equivalent to C++ erase(pos).
func (l *List[T]) Erase(pos Iterator[T]) Iterator[T] {
	if pos.node == l.sentinel {
		return pos
	}
	next := pos.node.next
	unlink(pos.node)
	pos.node.prev, pos.node.next = nil, nil
	l.size--
	return Iterator[T]{node: next}
}

// Clear removes all elements in O(1) (the nodes become garbage).
// Equivalent to C++ clear().
func (l *List[T]) Clear() { l.init() }

// Assign replaces the contents with n copies of value.
// Equivalent to C++ assign(n, value).
func (l *List[T]) Assign(n int, value T) {
	l.init()
	for i := 0; i < n; i++ {
		l.PushBack(value)
	}
}

// AssignSlice replaces the contents with a copy of the given slice.
func (l *List[T]) AssignSlice(items []T) {
	l.init()
	for _, v := range items {
		l.PushBack(v)
	}
}

// Reverse reverses the order of the elements in O(n). Equivalent to reverse().
func (l *List[T]) Reverse() {
	if l.size < 2 {
		return
	}
	cur := l.sentinel
	for {
		cur.prev, cur.next = cur.next, cur.prev
		cur = cur.prev // prev now holds the old next
		if cur == l.sentinel {
			break
		}
	}
}

// Swap exchanges the contents of the list with another in O(1).
// Equivalent to C++ swap().
func (l *List[T]) Swap(other *List[T]) {
	if other == nil {
		return
	}
	l.sentinel, other.sentinel = other.sentinel, l.sentinel
	l.size, other.size = other.size, l.size
}

// RemoveIf removes every element for which pred returns true and returns the
// number removed. O(n). Equivalent to C++ remove_if(pred).
func (l *List[T]) RemoveIf(pred func(value T) bool) int {
	removed := 0
	cur := l.sentinel.next
	for cur != l.sentinel {
		next := cur.next
		if pred(cur.value) {
			unlink(cur)
			cur.prev, cur.next = nil, nil
			l.size--
			removed++
		}
		cur = next
	}
	return removed
}

// Unique removes consecutive duplicate elements judged equal by eq, keeping the
// first of each run, and returns the number removed. O(n). Equivalent to C++
// unique(pred). Call Sort first to remove all duplicates.
func (l *List[T]) Unique(eq func(a, b T) bool) int {
	if l.size < 2 {
		return 0
	}
	removed := 0
	cur := l.sentinel.next
	for cur.next != l.sentinel {
		if eq(cur.value, cur.next.value) {
			dup := cur.next
			unlink(dup)
			dup.prev, dup.next = nil, nil
			l.size--
			removed++
		} else {
			cur = cur.next
		}
	}
	return removed
}

// Splice moves all elements of other into this list before pos, in O(1),
// leaving other empty. other must not be the same list. Equivalent to the
// whole-list overload of C++ splice.
func (l *List[T]) Splice(pos Iterator[T], other *List[T]) {
	if other == nil || other == l || other.size == 0 {
		return
	}
	first := other.sentinel.next
	last := other.sentinel.prev
	at := pos.node
	before := at.prev

	before.next = first
	first.prev = before
	last.next = at
	at.prev = last

	l.size += other.size
	other.init()
}

// Sort sorts the list in ascending order according to less, using a stable,
// in-place bottom-up-style merge sort in O(n log n) time and O(log n) stack.
// Equivalent to C++ list::sort(comp).
func (l *List[T]) Sort(less func(a, b T) bool) {
	if l.size < 2 {
		return
	}
	head := l.breakRing()
	head = mergeSortChain(head, less)
	l.rebuildRing(head, l.size)
}

// Merge merges the sorted list other into this sorted list, preserving order
// and leaving other empty. Both lists must already be sorted by less. O(n+m).
// Equivalent to C++ list::merge(other, comp).
func (l *List[T]) Merge(other *List[T], less func(a, b T) bool) {
	if other == nil || other == l || other.size == 0 {
		return
	}
	if l.size == 0 {
		l.Swap(other)
		return
	}
	total := l.size + other.size
	a := l.breakRing()
	b := other.breakRing()
	l.rebuildRing(mergeChains(a, b, less), total)
	other.init()
}

// breakRing returns the elements as a nil-terminated forward (next) chain
// without changing size. Returns nil when empty.
func (l *List[T]) breakRing() *node[T] {
	if l.size == 0 {
		return nil
	}
	head := l.sentinel.next
	l.sentinel.prev.next = nil
	return head
}

// rebuildRing relinks the circular doubly linked ring from a nil-terminated
// forward chain and sets size to count.
func (l *List[T]) rebuildRing(head *node[T], count int) {
	s := l.sentinel
	prev := s
	for n := head; n != nil; n = n.next {
		n.prev = prev
		prev.next = n
		prev = n
	}
	prev.next = s
	s.prev = prev
	l.size = count
}

// mergeSortChain sorts a nil-terminated forward chain and returns the new head.
func mergeSortChain[T any](head *node[T], less func(a, b T) bool) *node[T] {
	if head == nil || head.next == nil {
		return head
	}
	// Split into two halves with the slow/fast pointer technique.
	slow, fast := head, head.next
	for fast != nil && fast.next != nil {
		slow = slow.next
		fast = fast.next.next
	}
	mid := slow.next
	slow.next = nil
	left := mergeSortChain(head, less)
	right := mergeSortChain(mid, less)
	return mergeChains(left, right, less)
}

// mergeChains stably merges two sorted nil-terminated forward chains.
func mergeChains[T any](a, b *node[T], less func(a, b T) bool) *node[T] {
	var dummy node[T]
	tail := &dummy
	for a != nil && b != nil {
		if less(b.value, a.value) { // b strictly less: take b
			tail.next = b
			b = b.next
		} else { // a <= b: take a (keeps the sort stable)
			tail.next = a
			a = a.next
		}
		tail = tail.next
	}
	if a != nil {
		tail.next = a
	} else {
		tail.next = b
	}
	return dummy.next
}

// ToSlice returns the elements as a new slice in front-to-back order.
func (l *List[T]) ToSlice() []T {
	result := make([]T, 0, l.size)
	for n := l.sentinel.next; n != l.sentinel; n = n.next {
		result = append(result, n.value)
	}
	return result
}

// ForEach calls fn for each element in front-to-back order.
func (l *List[T]) ForEach(fn func(value T)) {
	for n := l.sentinel.next; n != l.sentinel; n = n.next {
		fn(n.value)
	}
}

// ForEachReverse calls fn for each element in back-to-front order.
func (l *List[T]) ForEachReverse(fn func(value T)) {
	for n := l.sentinel.prev; n != l.sentinel; n = n.prev {
		fn(n.value)
	}
}
