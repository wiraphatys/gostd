# GoSTD - Go Standard Data Structures

A Go library providing C++ STL-compatible data structures with 100% API compatibility. This library aims to provide familiar interfaces for developers coming from C++ while leveraging Go's type safety and generics.

## Features

- 🚀 **Type-safe**: Built with Go generics for compile-time type checking
- 🔄 **100% C++ STL API compatible**: All methods match C++ std:: containers exactly
- ⚡ **High performance**: Optimized implementations with minimal overhead
- 🧪 **Well tested**: Comprehensive test suite with examples and benchmarks
- 📖 **Well documented**: Complete documentation with usage examples

## Installation

```bash
go get github.com/wiraphatys/gostd
```

## Available Data Structures

### Queue

A FIFO (First-In-First-Out) container adapter that matches C++ `std::queue` exactly.

#### Usage

```go
package main

import (
    "fmt"
    "github.com/wiraphatys/gostd/queue"
)

func main() {
    // Create a new queue
    q := queue.New[int]()

    // Push elements
    q.Push(10)
    q.Push(20)
    q.Push(30)

    // Access elements
    front, _ := q.Front() // 10
    back, _ := q.Back()   // 30
    size := q.Size()      // 3

    // Pop elements (FIFO order)
    for !q.Empty() {
        front, _ := q.Front()
        fmt.Println(front) // prints: 10, 20, 30
        q.Pop()
    }
}
```

#### API Reference

All methods match C++ `std::queue` exactly:

| Method | Description | C++ Equivalent |
|--------|-------------|----------------|
| `New[T]()` | Create new queue | Constructor |
| `Push(value T)` | Add element to back | `push()` |
| `Pop()` | Remove front element | `pop()` |
| `Front()` | Access front element | `front()` |
| `Back()` | Access back element | `back()` |
| `Empty()` | Check if empty | `empty()` |
| `Size()` | Get number of elements | `size()` |
| `Emplace(value T)` | Construct in place | `emplace()` (C++11) |
| `Swap(other)` | Swap contents | `swap()` (C++11) |

#### Additional Go-friendly methods:

- `NewFromSlice([]T)` - Create queue from slice
- `Clear()` - Remove all elements
- `ToSlice()` - Convert to slice

### Stack

A LIFO (Last-In-First-Out) container adapter that matches C++ `std::stack` exactly.

#### Usage

```go
package main

import (
    "fmt"
    "github.com/wiraphatys/gostd/stack"
)

func main() {
    // Create a new stack
    s := stack.New[int]()

    // Push elements
    s.Push(10)
    s.Push(20)
    s.Push(30)

    // Access top element
    top, _ := s.Top()     // 30
    size := s.Size()      // 3

    // Pop elements (LIFO order)
    for !s.Empty() {
        top, _ := s.Top()
        fmt.Println(top) // prints: 30, 20, 10
        s.Pop()
    }
}
```

#### API Reference

All methods match C++ `std::stack` exactly:

| Method | Description | C++ Equivalent |
|--------|-------------|----------------|
| `New[T]()` | Create new stack | Constructor |
| `Push(value T)` | Add element to top | `push()` |
| `Pop()` | Remove top element | `pop()` |
| `Top()` | Access top element | `top()` |
| `Empty()` | Check if empty | `empty()` |
| `Size()` | Get number of elements | `size()` |
| `Emplace(value T)` | Construct in place | `emplace()` (C++11) |
| `Swap(other)` | Swap contents | `swap()` (C++11) |

#### Additional Go-friendly methods:

- `NewFromSlice([]T)` - Create stack from slice
- `Clear()` - Remove all elements
- `ToSlice()` - Convert to slice (bottom to top order)

## Testing

Run all tests:

```bash
go test ./...
```

Run with benchmarks:

```bash
go test -bench=. ./...
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License.

## Roadmap

- ✅ Queue (`std::queue`)
- ✅ Stack (`std::stack`) 
- ⏳ Priority Queue (`std::priority_queue`)
- ⏳ Vector (`std::vector`)
- ⏳ Deque (`std::deque`)
- ⏳ List (`std::list`)
- ⏳ Set (`std::set`)
- ⏳ Map (`std::map`)