// Package forwardlist implements a singly linked list that mirrors the API and
// behaviour of C++ std::forward_list.
//
// As in C++, mutation happens "after" a position: there is a before-begin
// sentinel, and InsertAfter/EraseAfter/SpliceAfter operate on the node
// following the supplied iterator. This keeps every node down to a single
// forward pointer.
//
//	PushFront / PopFront / Front          O(1)
//	InsertAfter / EraseAfter              O(1)
//	Sort                                  O(n log n), in-place merge sort
//	Reverse / Unique / RemoveIf / Merge   O(n)
//
// Unlike std::forward_list (which omits size() to save memory) this type tracks
// its length so Size runs in O(1) as a Go convenience.
package forwardlist

import "errors"

type node[T any] struct {
	next  *node[T]
	value T
}

// ForwardList is a singly linked list. Create one with New or NewFromSlice.
type ForwardList[T any] struct {
	before *node[T] // before-begin sentinel; before.next is the first element
	size   int
}

// Iterator is a forward-only cursor over a ForwardList.
type Iterator[T any] struct {
	node *node[T]
}

// Value returns the element the iterator refers to. Undefined for the End or
// BeforeBegin positions.
func (it Iterator[T]) Value() T { return it.node.value }

// SetValue overwrites the element the iterator refers to.
func (it Iterator[T]) SetValue(v T) { it.node.value = v }

// Next advances the iterator to the following element.
func (it *Iterator[T]) Next() { it.node = it.node.next }

// Equal reports whether two iterators refer to the same position.
func (it Iterator[T]) Equal(other Iterator[T]) bool { return it.node == other.node }

// New creates and returns a new empty forward list.
func New[T any]() *ForwardList[T] {
	return &ForwardList[T]{before: &node[T]{}}
}

// NewFromSlice creates a forward list initialized with a copy of items,
// preserving order.
func NewFromSlice[T any](items []T) *ForwardList[T] {
	fl := New[T]()
	prev := fl.before
	for _, v := range items {
		prev.next = &node[T]{value: v}
		prev = prev.next
		fl.size++
	}
	return fl
}

// Size returns the number of elements (O(1) here). Not present in C++.
func (fl *ForwardList[T]) Size() int { return fl.size }

// Empty reports whether the list has no elements. Equivalent to empty().
func (fl *ForwardList[T]) Empty() bool { return fl.size == 0 }

// BeforeBegin returns an iterator to the position before the first element,
// suitable for InsertAfter/EraseAfter at the front. Equivalent to
// before_begin().
func (fl *ForwardList[T]) BeforeBegin() Iterator[T] { return Iterator[T]{node: fl.before} }

// Begin returns an iterator to the first element (or End if empty).
func (fl *ForwardList[T]) Begin() Iterator[T] { return Iterator[T]{node: fl.before.next} }

// End returns the past-the-end iterator (nil).
func (fl *ForwardList[T]) End() Iterator[T] { return Iterator[T]{node: nil} }

// PushFront inserts value at the front in O(1). Equivalent to push_front().
func (fl *ForwardList[T]) PushFront(value T) {
	fl.before.next = &node[T]{value: value, next: fl.before.next}
	fl.size++
}

// EmplaceFront constructs an element in place at the front. Equivalent to
// PushFront in Go. Matches C++ emplace_front().
func (fl *ForwardList[T]) EmplaceFront(value T) { fl.PushFront(value) }

// PopFront removes the first element in O(1). Equivalent to pop_front().
func (fl *ForwardList[T]) PopFront() error {
	if fl.size == 0 {
		return errors.New("forwardlist: pop_front on empty list")
	}
	fl.before.next = fl.before.next.next
	fl.size--
	return nil
}

// Front returns the first element. Equivalent to C++ front().
func (fl *ForwardList[T]) Front() (T, error) {
	if fl.size == 0 {
		var zero T
		return zero, errors.New("forwardlist: front of empty list")
	}
	return fl.before.next.value, nil
}

// InsertAfter inserts value after the element pos refers to and returns an
// iterator to the new element. O(1). Equivalent to insert_after(pos, value).
func (fl *ForwardList[T]) InsertAfter(pos Iterator[T], value T) Iterator[T] {
	n := &node[T]{value: value, next: pos.node.next}
	pos.node.next = n
	fl.size++
	return Iterator[T]{node: n}
}

// EraseAfter removes the element following pos and returns an iterator to the
// element after the removed one (or End). O(1). Equivalent to erase_after(pos).
func (fl *ForwardList[T]) EraseAfter(pos Iterator[T]) Iterator[T] {
	target := pos.node.next
	if target == nil {
		return fl.End()
	}
	pos.node.next = target.next
	target.next = nil
	fl.size--
	return Iterator[T]{node: pos.node.next}
}

// EmplaceAfter constructs an element in place after pos. Equivalent to
// InsertAfter in Go. Matches C++ emplace_after().
func (fl *ForwardList[T]) EmplaceAfter(pos Iterator[T], value T) Iterator[T] {
	return fl.InsertAfter(pos, value)
}

// Clear removes all elements in O(1). Equivalent to C++ clear().
func (fl *ForwardList[T]) Clear() {
	fl.before.next = nil
	fl.size = 0
}

// Assign replaces the contents with n copies of value.
func (fl *ForwardList[T]) Assign(n int, value T) {
	fl.Clear()
	prev := fl.before
	for i := 0; i < n; i++ {
		prev.next = &node[T]{value: value}
		prev = prev.next
		fl.size++
	}
}

// AssignSlice replaces the contents with a copy of items.
func (fl *ForwardList[T]) AssignSlice(items []T) {
	fl.Clear()
	prev := fl.before
	for _, v := range items {
		prev.next = &node[T]{value: v}
		prev = prev.next
		fl.size++
	}
}

// Reverse reverses the order of the elements in O(n). Equivalent to reverse().
func (fl *ForwardList[T]) Reverse() {
	var prev *node[T]
	cur := fl.before.next
	for cur != nil {
		next := cur.next
		cur.next = prev
		prev = cur
		cur = next
	}
	fl.before.next = prev
}

// Swap exchanges the contents of the list with another in O(1).
func (fl *ForwardList[T]) Swap(other *ForwardList[T]) {
	if other == nil {
		return
	}
	fl.before, other.before = other.before, fl.before
	fl.size, other.size = other.size, fl.size
}

// RemoveIf removes every element for which pred returns true and returns the
// number removed. O(n). Equivalent to remove_if(pred).
func (fl *ForwardList[T]) RemoveIf(pred func(value T) bool) int {
	removed := 0
	prev := fl.before
	for prev.next != nil {
		if pred(prev.next.value) {
			doomed := prev.next
			prev.next = doomed.next
			doomed.next = nil
			fl.size--
			removed++
		} else {
			prev = prev.next
		}
	}
	return removed
}

// Unique removes consecutive duplicates judged equal by eq, keeping the first
// of each run, and returns the number removed. O(n). Equivalent to unique(pred).
func (fl *ForwardList[T]) Unique(eq func(a, b T) bool) int {
	if fl.size < 2 {
		return 0
	}
	removed := 0
	cur := fl.before.next
	for cur.next != nil {
		if eq(cur.value, cur.next.value) {
			cur.next = cur.next.next
			fl.size--
			removed++
		} else {
			cur = cur.next
		}
	}
	return removed
}

// SpliceAfter moves all elements of other into this list after pos, leaving
// other empty. other must not be the same list. O(distance in other) because a
// singly linked list cannot find its tail in constant time. Equivalent to the
// whole-list overload of splice_after.
func (fl *ForwardList[T]) SpliceAfter(pos Iterator[T], other *ForwardList[T]) {
	if other == nil || other == fl || other.size == 0 {
		return
	}
	first := other.before.next
	last := first
	for last.next != nil {
		last = last.next
	}
	last.next = pos.node.next
	pos.node.next = first
	fl.size += other.size
	other.before.next = nil
	other.size = 0
}

// ResizeWithValue changes the number of elements to n, padding any growth with
// copies of value. Equivalent to resize(n, value).
func (fl *ForwardList[T]) ResizeWithValue(n int, value T) {
	if n < 0 {
		n = 0
	}
	switch {
	case n < fl.size:
		prev := fl.before
		for i := 0; i < n; i++ {
			prev = prev.next
		}
		prev.next = nil
		fl.size = n
	case n > fl.size:
		prev := fl.before
		for prev.next != nil {
			prev = prev.next
		}
		for fl.size < n {
			prev.next = &node[T]{value: value}
			prev = prev.next
			fl.size++
		}
	}
}

// Resize changes the number of elements to n, zero-filling any growth.
func (fl *ForwardList[T]) Resize(n int) {
	var zero T
	fl.ResizeWithValue(n, zero)
}

// Sort sorts the list in ascending order according to less using a stable
// merge sort in O(n log n). Equivalent to forward_list::sort(comp).
func (fl *ForwardList[T]) Sort(less func(a, b T) bool) {
	if fl.size < 2 {
		return
	}
	fl.before.next = mergeSortChain(fl.before.next, less)
}

// Merge merges the sorted list other into this sorted list, preserving order
// and leaving other empty. Both lists must already be sorted. O(n+m).
// Equivalent to forward_list::merge(other, comp).
func (fl *ForwardList[T]) Merge(other *ForwardList[T], less func(a, b T) bool) {
	if other == nil || other == fl || other.size == 0 {
		return
	}
	fl.before.next = mergeChains(fl.before.next, other.before.next, less)
	fl.size += other.size
	other.before.next = nil
	other.size = 0
}

func mergeSortChain[T any](head *node[T], less func(a, b T) bool) *node[T] {
	if head == nil || head.next == nil {
		return head
	}
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

func mergeChains[T any](a, b *node[T], less func(a, b T) bool) *node[T] {
	var dummy node[T]
	tail := &dummy
	for a != nil && b != nil {
		if less(b.value, a.value) {
			tail.next = b
			b = b.next
		} else {
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
func (fl *ForwardList[T]) ToSlice() []T {
	result := make([]T, 0, fl.size)
	for n := fl.before.next; n != nil; n = n.next {
		result = append(result, n.value)
	}
	return result
}

// ForEach calls fn for each element in front-to-back order.
func (fl *ForwardList[T]) ForEach(fn func(value T)) {
	for n := fl.before.next; n != nil; n = n.next {
		fn(n.value)
	}
}
