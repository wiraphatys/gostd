# GoSTD — Go Standard Data Structures

A Go library providing **C++ STL-compatible data structures** with familiar
method names and semantics, implemented with **optimal algorithms and data
structures** (not naive ones). It targets developers who know the C++ STL and
want the same containers — `vector`, `deque`, `list`, `set`, `map`,
`priority_queue`, … — backed by Go generics and type safety.

## Design goals

- 🎯 **C++ method parity** — methods are named and behave like their `std::`
  counterparts (`PushBack`, `PopFront`, `LowerBound`, `EqualRange`, `Top`,
  `Insert`, `Erase`, …), so porting C++ code is mechanical.
- ⚡ **Optimal by construction** — every container uses the asymptotically
  optimal data structure for its contract (red-black tree for ordered maps/sets,
  circular ring buffer for `deque`, binary heap for `priority_queue`, in-place
  merge sort for `list`). See the complexity table below.
- 🧪 **Thoroughly tested** — table-driven unit tests, runnable `Example` tests,
  benchmarks, **randomized differential tests** against reference models, and a
  full **red-black invariant validator** exercised over tens of thousands of
  randomized operations.
- 🧩 **Generics + type safety** — compile-time checked element types.
- 📦 **Zero dependencies**, builds on **Go 1.19+**.

## Installation

```bash
go get github.com/wiraphatys/gostd
```

## Containers & complexity

| GoSTD package | C++ equivalent | Backing structure | Key complexities |
|---|---|---|---|
| `vector` | `std::vector` | growable array (×2 growth) | PushBack *O(1)* amortized, index *O(1)*, Insert/Erase *O(n)* |
| `deque` | `std::deque` | circular ring buffer | push/pop both ends *O(1)* amortized, index *O(1)* |
| `list` | `std::list` | doubly linked (sentinel) | insert/erase/splice *O(1)*, sort *O(n log n)* |
| `forwardlist` | `std::forward_list` | singly linked | insert/erase after *O(1)*, sort *O(n log n)* |
| `stack` | `std::stack` | slice | push/pop/top *O(1)* amortized |
| `queue` | `std::queue` | slice | push/pop/front/back *O(1)* amortized |
| `priorityqueue` | `std::priority_queue` | binary heap | push/pop *O(log n)*, top *O(1)* |
| `set` | `std::set` | red-black tree | insert/erase/find/bounds *O(log n)* |
| `multiset` | `std::multiset` | red-black tree | insert/erase/find/bounds *O(log n)* |
| `treemap` | `std::map` | red-black tree | insert/erase/find/bounds *O(log n)* |
| `multimap` | `std::multimap` | red-black tree | insert/erase/find/bounds *O(log n)* |
| `unorderedset` | `std::unordered_set` | hash table (Go map) | insert/erase/find *O(1)* average |
| `unorderedmap` | `std::unordered_map` | hash table (Go map) | insert/erase/find *O(1)* average |
| `bitset` | `std::bitset` | `[]uint64` words | bit ops *O(1)*, bulk ops *O(n/64)* |
| `pair` | `std::pair` | struct | — |
| `constraints` | (`<concepts>` helpers) | type constraints | `Ordered`, `Integer`, `Min`, `Max`, … |

> **Naming note:** `std::map` maps to package **`treemap`** because `map` is a
> reserved word in Go. Multi-word headers drop the underscore
> (`priority_queue` → `priorityqueue`, `forward_list` → `forwardlist`,
> `unordered_map` → `unorderedmap`). Methods keep their exact C++ names.

## Conventions

- **Constructors:** `New[T]()`; plus `NewFromSlice`, `NewFunc` (custom
  comparator), `NewWithCapacity`, etc. where relevant.
- **Ordering:** ordered containers default to natural `<` ordering (like
  `std::less`). Pass a comparator `func(a, b T) bool` (meaning "a before b") to
  the `NewFunc` constructors for any other order — e.g. `std::greater` for a
  descending set or a min-heap.
- **Fallible access** returns Go-idiomatic results instead of throwing:
  `Top()/Front()/At()` return `(T, error)`; `Pop()` returns `error`;
  lookups like `Map.Get` return `(V, bool)`. Out-of-range `bitset` positions
  and unchecked `Get`/`Set` panic, mirroring C++ undefined behaviour /
  exceptions.
- **Iterators:** tree-based containers expose bidirectional iterators with
  `Begin()/End()/RBegin()/REnd()`, `Next()/Prev()`, `Value()` (and `Key()`),
  and `Equal()`/`IsEnd()`. A convenience `ForEach` is also provided.

## Quick examples

### vector

```go
v := vector.New[int]()
v.PushBack(10)
v.PushBack(20)
v.Insert(1, 15)        // [10 15 20]
x, _ := v.At(1)        // 15, bounds-checked
v.Set(0, 99)           // operator[]= (unchecked)
```

### deque

```go
d := deque.New[int]()
d.PushBack(1)
d.PushFront(0)         // [0 1]
front, _ := d.Front()  // 0
```

### priority_queue

```go
pq := priorityqueue.New[int]()        // max-heap (C++ default)
pq.Push(3); pq.Push(9); pq.Push(1)
top, _ := pq.Top()                    // 9
min := priorityqueue.NewMin[int]()    // == priority_queue<T, …, greater<T>>
```

### set / map (ordered, red-black tree)

```go
s := set.New[int]()
s.Insert(3); s.Insert(1); s.Insert(2)
fmt.Println(s.ToSlice())               // [1 2 3] always sorted
lo, hi := s.EqualRange(2)              // C++ equal_range

m := treemap.New[string, int]()        // std::map
m.Set("apple", 5)
*m.Ref("apple")++                      // operator[]: returns a mutable *V
n, _ := m.At("apple")                  // 6
```

### unordered_map (hash)

```go
counts := unorderedmap.New[string, int]()
for _, w := range words {
    counts.Set(w, counts.GetOr(w, 0)+1)
}
```

### bitset

```go
b := bitset.New(8)
b.Set(0); b.Set(7)
fmt.Println(b.ToString())              // 10000001  (MSB first, like std::bitset)
b.Clone().And(other)                   // bitwise &
```

## Testing

```bash
go test ./...              # unit + example tests
go test -bench=. ./...     # benchmarks
go test -cover ./...       # coverage
```

The ordered containers are validated by a red-black invariant checker (root
black, no red-red, equal black-height on every path, BST order, parent-link
consistency) run after randomized insert/erase sequences; `deque` and `bitset`
are checked with randomized differential tests against reference models.

## Roadmap

- ✅ Queue (`std::queue`)
- ✅ Stack (`std::stack`)
- ✅ Vector (`std::vector`)
- ✅ Deque (`std::deque`)
- ✅ List (`std::list`)
- ✅ Forward List (`std::forward_list`)
- ✅ Priority Queue (`std::priority_queue`)
- ✅ Set / MultiSet (`std::set` / `std::multiset`)
- ✅ Map / MultiMap (`std::map` / `std::multimap`)
- ✅ Unordered Set / Map (`std::unordered_set` / `std::unordered_map`)
- ✅ Bitset (`std::bitset`)
- ✅ Pair (`std::pair`)

## License

This project is licensed under the MIT License.
