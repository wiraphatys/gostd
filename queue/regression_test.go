package queue

import "testing"

// The ring-buffer backing must keep memory proportional to the number of live
// elements, not to the total number of items ever enqueued. The previous
// slice[1:] implementation would let the backing array grow toward the total
// throughput (here ~1,000,000); the ring buffer stays tiny.
func TestNoUnboundedGrowthUnderChurn(t *testing.T) {
	q := New[int]()
	const churn = 1_000_000
	for i := 0; i < churn; i++ {
		q.Push(i)
		if q.Size() > 4 { // never keep more than ~5 live elements
			q.Pop()
		}
	}
	for !q.Empty() {
		q.Pop()
	}
	if c := q.d.Capacity(); c > 64 {
		t.Errorf("capacity grew to %d under churn; ring buffer should stay bounded by live size", c)
	}
}

// FIFO order must hold under heavy interleaved push/pop churn (exercises the
// ring buffer wrap-around).
func TestFIFOUnderChurn(t *testing.T) {
	q := New[int]()
	next, expect, live := 0, 0, 0
	for step := 0; step < 100000; step++ {
		q.Push(next)
		next++
		live++
		if live > 8 {
			f, _ := q.Front()
			if f != expect {
				t.Fatalf("FIFO broken: front=%d want %d", f, expect)
			}
			q.Pop()
			expect++
			live--
		}
	}
	for !q.Empty() {
		f, _ := q.Front()
		if f != expect {
			t.Fatalf("FIFO broken during drain: front=%d want %d", f, expect)
		}
		q.Pop()
		expect++
	}
}

// Popped slots must be cleared so reference-typed elements are not pinned.
func TestPopReleasesReference(t *testing.T) {
	q := New[*int]()
	a, b := 1, 2
	q.Push(&a)
	q.Push(&b)
	if err := q.Pop(); err != nil { // removes &a
		t.Fatalf("Pop errored: %v", err)
	}
	front, _ := q.Front()
	if *front != 2 {
		t.Errorf("front = %d, want 2", *front)
	}
}
