# Go Notes

Personal reference for the gocamp course. It has two parts: a standard-library cheat sheet, and factoids
collected during lessons.

> Look up any signature in the terminal: `go doc <pkg>.<Name>` (for example `go doc slices.SortFunc`,
> `go doc container/heap`, `go doc io.Reader`).

## Contents

1. [Standard library cheat sheet](#1-standard-library-cheat-sheet)
   - [1.1 Core interfaces](#11-core-interfaces)
   - [1.2 DSA toolkit](#12-dsa-toolkit)
   - [1.3 Errors and formatting](#13-errors-and-formatting)
   - [1.4 I/O helpers](#14-io-helpers)
2. [Language essentials](#2-language-essentials)
   - [2.1 Composite literals and `struct{}{}`](#21-composite-literals-and-struct)
   - [2.2 Integer sizes and conversions](#22-integer-sizes-and-conversions)
   - [2.3 Zero values](#23-zero-values)
   - [2.4 Pointers: Go vs C](#24-pointers-go-vs-c)
   - [2.5 Visibility is case-sensitive](#25-visibility-is-case-sensitive)
   - [2.6 Compile-time interface checks](#26-compile-time-interface-checks)
   - [2.7 `%` keeps the sign of the dividend](#27--keeps-the-sign-of-the-dividend)
3. [Slices and maps](#3-slices-and-maps)
4. [Interfaces](#4-interfaces)
5. [Memory and the GC](#5-memory-and-the-gc)
6. [Go's design and history](#6-gos-design-and-history)

---

## 1. Standard library cheat sheet

### 1.1 Core interfaces

Learn these method sets exactly. Interfaces compose by embedding: `ReadWriter` is `Reader` plus `Writer`.

```go
// Built-in / fmt
type error interface    { Error() string }
type Stringer interface { String() string } // fmt.Stringer: controls how %v prints your type

// io: the most important interfaces in the stdlib
type Reader interface     { Read(p []byte) (n int, err error) }
type Writer interface     { Write(p []byte) (n int, err error) }
type Closer interface     { Close() error }
type ReadWriter interface { Reader; Writer }
type ByteReader interface { ReadByte() (byte, error) }

// sort
type Interface interface { // sort.Interface
	Len() int
	Less(i, j int) bool
	Swap(i, j int)
}

// container/heap
type Interface interface { // heap.Interface
	sort.Interface
	Push(x any) // append x; heap.Push calls this, then restores heap order
	Pop() any   // remove and return the LAST element; heap.Pop swaps the min there first
}

// context
type Context interface {
	Deadline() (deadline time.Time, ok bool)
	Done() <-chan struct{}
	Err() error
	Value(key any) any
}

// Others
type Handler interface   { ServeHTTP(w http.ResponseWriter, r *http.Request) } // net/http
type Marshaler interface { MarshalJSON() ([]byte, error) }                      // encoding/json
type Locker interface    { Lock(); Unlock() }                                   // sync
```

### 1.2 DSA toolkit

**`slices`** (Go 1.21+). Prefer these over the older `sort` package.

| Function | Signature / notes |
|---|---|
| `slices.Sort(s)` | elements must be `cmp.Ordered` |
| `slices.SortFunc(s, cmp)` | `cmp(a, b T) int` returns negative, 0 or positive, **not** a bool |
| `slices.BinarySearch(s, target)` | `(idx int, found bool)` |
| `slices.Contains(s, v)` / `slices.Index(s, v)` | `bool` / `int` (-1 if missing) |
| `slices.Reverse(s)` | in place |
| `slices.Max(s)` / `slices.Min(s)` | panics on an empty slice |
| `slices.Clone(s)` / `slices.Equal(a, b)` | copy / compare element by element |
| `slices.Insert(s, i, vs...)` / `slices.Delete(s, i, j)` | **return the new slice**: `s = slices.Delete(s, i, j)` |

**`container/heap`**: you implement `heap.Interface` and call these:

```go
heap.Init(h)      // build a heap from arbitrary order, O(n)
heap.Push(h, x)   // O(log n)
heap.Pop(h) any   // removes the min (by Less), O(log n); type-assert the result
heap.Fix(h, i)    // after changing element i's priority, O(log n)
```

**`maps`** (1.23+ iterators):

```go
slices.Sorted(maps.Keys(m)) // sorted keys; map iteration order is random
```

**`strconv`**:

```go
strconv.Itoa(i int) string
strconv.Atoi(s string) (int, error)
strconv.ParseInt(s string, base, bitSize int) (int64, error)
```

**`strings`**:

```go
strings.Split(s, sep) []string     strings.Join(elems, sep) string
strings.Fields(s) []string         strings.TrimSpace(s) string
strings.HasPrefix(s, p) bool       strings.Contains(s, sub) bool

var sb strings.Builder // efficient string building
sb.WriteString("x"); sb.WriteByte('y'); s := sb.String()
```

**`math`** and built-ins:

```go
math.MaxInt   math.MinInt   math.Inf(1)
min(a, b, c)  max(a, b) // built-in since Go 1.21, no import needed
```

> **Two comparison styles.** The older `sort.Slice` takes `less(i, j int) bool`, a yes/no on
> *indices*. The newer `slices.SortFunc` takes `cmp(a, b T) int`, a three-way result on
> *values*. Prefer the newer one. `cmp.Compare(a, b)` gives you the three-way result for free,
> and swapping the arguments sorts in descending order.

### 1.3 Errors and formatting

```go
errors.New(text string) error
fmt.Errorf("reading %s: %w", name, err) error  // %w wraps err so errors.Is/As can find it
errors.Is(err, target error) bool              // is this sentinel error anywhere in the chain?
errors.As(err error, target any) bool          // is an error of this TYPE in the chain? (target is a pointer)
fmt.Sprintf(format, args...) string
fmt.Fprintf(w io.Writer, format, args...) (int, error) // print to any Writer
```

### 1.4 I/O helpers

```go
io.ReadAll(r) ([]byte, error)
io.Copy(dst, src) (int64, error)
os.ReadFile(name) ([]byte, error)
os.Open(name) (*os.File, error)
bufio.NewReader(r).ReadString('\n') (string, error)

sc := bufio.NewScanner(r)
for sc.Scan() {
	line := sc.Text()
}
```

---

## 2. Language essentials

### 2.1 Composite literals and `struct{}{}`

A value of a composite type is written **TYPE + `{...}`**:

```go
Point{1, 2}      // type Point,      literal {1, 2}
[]int{1, 2, 3}   // type []int,      literal {1, 2, 3}
map[string]int{} // type map[..]int, literal {} (an empty map)
```

`struct{}` is a **type** (a struct with no fields), so a *value* of it needs a second pair of braces:

```go
seen := map[int]struct{}  {}
//      └──── type ─────┘ └┘ literal

seen[5] = struct{}  {}
//        └ type ┘ └┘ literal
```

- **Idiomatic set:** `map[K]struct{}`. The empty struct takes **0 bytes**. `map[K]bool` also works and
  reads more nicely (`if seen[k]`).
- Naming the type helps readability: `type void struct{}` → `seen[5] = void{}`.
- `seen := map[int]struct{}` doesn't compile, because the right side of `:=` must be a value, not a type.

Ways to create an empty map:

```go
seen := map[int]struct{}{}      // literal
seen := make(map[int]struct{})  // most common
var seen map[int]struct{}       // ⚠️ nil map: reads work, WRITES PANIC
```

### 2.2 Integer sizes and conversions

- `int` / `uint` are the **machine word size**: 64-bit on 64-bit platforms. C++ usually keeps `int` at
  32 bits (the LP64 model). Go chose word size so that `len()`, `cap()` and indices can address all memory.
- Sized types: `int8/16/32/64`, `uint8/16/32/64` (`byte` = `uint8`, `rune` = `int32`).
- **Use `int` by default.** Use sized types when layout matters: packed arrays, binary protocols, file
  formats.
- **There are no implicit conversions**, even between integer types: `var i int32 = n` (where `n` is
  an `int`) is a compile error. Write `int32(n)`.
- **Slice indices can be any integer type**: `s[i]` works with `i int32`.
- **Untyped constants adapt:** `-1` works as an `int32` with no conversion.

### 2.3 Zero values

Every type has a zero value (`0`, `""`, `false`, `nil`, or a struct of zero values). Three ways to get
the zero value of a type parameter `T`:

```go
var zero T         // most common and most readable
*new(T)            // one-liner; new(T) returns a *T to a zeroed T. No heap allocation (escape analysis)
clear(s[i:j])      // Go 1.21+: sets a range of slice elements to their zero value
```

**Make the zero value useful:** design types so that `var s Stack[int]` works without a constructor.
A nil slice is a valid empty stack.

### 2.4 Pointers: Go vs C

| | C | Go |
|---|---|---|
| Pointer arithmetic | `p + 1` | Not allowed (only via `unsafe`) |
| Freeing memory | manual `free`, dangling pointers | Garbage collected; valid as long as you hold it |
| Returning `&local` | bug: points at a dead stack frame | **Safe**: escape analysis moves it to the heap |
| Dereferencing null | undefined behavior | Always a **runtime panic** |
| Field access | `p->x` | `p.x` (automatic dereference) |
| Method calls | n/a | `v.M()` works on a pointer-receiver `M` if `v` is addressable |
| Casting pointer types | `(int*)fp` | Not allowed (only via `unsafe.Pointer`) |
| Address of anything | yes | **Not map elements**: `&m[k]` is a compile error (entries move as the map grows) |
| Interior pointers | just an address | Keep the **whole** object alive for the GC |

- **Everything is passed by value.** Passing a pointer copies the pointer; passing a slice copies its header.
- Slices, maps, channels, functions and interfaces are small headers that point to data elsewhere, so
  `*[]T` or `*map[K]V` is rarely needed.

### 2.5 Visibility is case-sensitive

`resize` and `Resize` are **different identifiers**. Capitalized = exported (visible outside the package);
lowercase = package-private. Internal helpers should be lowercase.

### 2.6 Compile-time interface checks

```go
var _ Deque[int] = (*RingBufferDeque[int])(nil) // build fails if the type stops satisfying Deque
```

`(*X)(nil)` is a typed nil pointer. Assigning it to `_` costs nothing at runtime.

### 2.7 `%` keeps the sign of the dividend

`-1 % 8 == -1` in Go (as in C), not `7`. For ring-buffer wraparound use `(i - 1 + n) % n`.

---

## 3. Slices and maps

### A slice is a header: `(ptr, len, cap)`

Passing a slice copies the header, but the **backing array is shared**:

```go
func add(s []int) {
	s[0] = 99        // writes into the SHARED array → the caller sees it
	s = append(s, 4) // changes only the local header
}

a := make([]int, 3, 10) // spare capacity
add(a)
fmt.Println(a)     // [99 0 0]: the caller's len is still 3
fmt.Println(a[:4]) // [99 0 0 4]: the 4 IS in the shared array

b := []int{1, 2, 3} // full (len == cap)
add(b)              // append reallocates → later changes never reach the caller
```

- Idiom: `s = append(s, x)`. Functions that grow a slice **return** it.
- Methods that change a slice field (`Push`, `Pop`) need a **pointer receiver**.

### Missing map key → zero value

`m[k]` on a missing key returns the zero value, with no error. Use comma-ok when zero could be a
real value:

```go
v, ok := m[k]
if _, ok := seen[5]; ok { ... }
```

This caused a real bug in the Day 1 LRU: a reverse-map lookup silently returned `0`, so the wrong key
was deleted.

### Dense integer keys → use an array, not a map

When keys are a small dense range `0..n-1`, an array beats a map. It's a single indexed load, with no
hashing, no per-entry overhead (a map costs about 20–50+ bytes per entry), and no silent zero for a
missing key. That's how the flat-array LRU replaced its reverse map with `keys []int`.

---

## 4. Interfaces

### An interface value is a pair: (type, value)

It's `nil` only when **both** halves are nil. A nil pointer inside an interface is **not** nil:

```go
type MyErr struct{}
func (*MyErr) Error() string { return "boom" }

func find() error {
	var e *MyErr = nil
	return e // (type=*MyErr, value=nil)
}

fmt.Println(find() == nil) // false!
```

**Fix:** return a literal `nil` from functions whose return type is an interface:
`if failed { return &MyErr{} }; return nil`.

### Implicit satisfaction: the consumer defines the interface

A type satisfies an interface just by having the methods; there's no `implements`. So the package
that *uses* a behavior defines a small interface for what it needs:

```go
type byteSink interface {
	Write(p []byte) (n int, err error)
}

func dump(w byteSink, data []byte) { w.Write(data) }

dump(os.Stdout, data)       // *os.File fits
dump(&bytes.Buffer{}, data) // *bytes.Buffer fits
dump(conn, data)            // net.Conn fits
```

- Go interfaces tend to be **tiny**, often one method (`io.Reader`, `io.Writer`, `fmt.Stringer`,
  `error`).
- Proverb: **"The bigger the interface, the weaker the abstraction."** (Rob Pike)
- Guideline: **"Accept interfaces, return structs."**
- Capstone use: the Redis server reads from an `io.Reader`, so tests can feed it a `strings.Reader`
  instead of a TCP connection.

---

## 5. Memory and the GC

### Resliced-away elements still pin memory

`s = s[:len(s)-1]` only shrinks `len`. The backing array still holds the old element, and the GC
doesn't know about `len`, so a popped `*BigThing` can never be freed. Zero the slot before reslicing:

```go
n := len(s.elements) - 1
v := s.elements[n]
var zero T
s.elements[n] = zero // or: clear(s.elements[n:])
s.elements = s.elements[:n]
```

`slices.Delete` has done this automatically since Go 1.22. This applies to stacks, LRU eviction, and
heap pops.

### Replaced arrays don't need zeroing

After `dq.buffer = newBuffer`, nothing references the old array, so the GC frees **all of it**.
Zero individual slots only in an array you **keep** (pops, evictions).

### Empty slice ≠ slice of nil pointers

`[]*T{}` has **len 0**: indexing `[0]` panics with *index out of range*. `make([]*T, n)` has n **nil**
pointers: indexing works, but dereferencing panics with *nil pointer dereference*.

### Detach unlinked nodes

When removing a node from a linked structure, also clear its own `next`/`prev`. Otherwise a stale
pointer from a dead node can keep a live chain reachable longer than needed.

### Flat arrays vs pointer-based nodes

Index arrays (`prev []int32`, `next []int32`) instead of `*node` pointers mean no allocation per node,
**no pointers for the GC to scan**, and better cache locality. `int32` indices halve the link memory,
with a limit of about 2.1 billion slots.

---

## 6. Go's design and history

- Designed at Google by **Ken Thompson** (co-creator of Unix and of B, C's predecessor),
  **Rob Pike** and **Robert Griesemer**: roughly "C, fixed for modern work".
  - **Kept from C:** structs, pointers, value semantics, a small language, predictable cost.
  - **Removed:** pointer arithmetic, header files, implicit conversions, macros, manual `free`.
  - **Added:** GC, slices, maps, interfaces, goroutines and channels, a standard toolchain
    (`go fmt`, `go test`, modules).
- **25 keywords** (C has 32, C++ has 90+). "Less is more": usually one obvious way to write something,
  so strangers' Go code looks like your own.
- Built for infrastructure: Docker, Kubernetes, etcd, Prometheus, CockroachDB.
