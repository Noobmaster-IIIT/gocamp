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
   - [2.4.1 Loops, receivers, shadowing](#241-loops-receivers-shadowing)
   - [2.4.2 Recursion in Go](#242-recursion-in-go)
   - [2.4.3 Function literals and closures](#243-function-literals-and-closures)
   - [2.4.4 Packages and file names](#244-packages-and-file-names)
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

**`slices`: full tour** (`go doc slices` lists everything)

| Group | Functions | Notes |
|---|---|---|
| Search | `Contains`, `Index`, `ContainsFunc`, `IndexFunc`, `BinarySearch`, `BinarySearchFunc` | `BinarySearch` needs sorted input and returns `(pos, found)`; `pos` is the insertion point |
| Sort | `Sort`, `SortFunc`, `SortStableFunc`, `IsSorted`, `IsSortedFunc` | `SortFunc` is **not stable**; use `SortStableFunc` when equal elements must keep their order |
| Compare | `Equal`, `EqualFunc`, `Compare`, `CompareFunc` | `Compare` is lexicographic: -1/0/+1 |
| Min/max | `Min`, `Max`, `MinFunc`, `MaxFunc` | **panic on an empty slice** |
| Modify (**return the new slice**) | `Insert`, `Delete`, `DeleteFunc`, `Replace`, `Compact`, `CompactFunc`, `Concat`, `Repeat`, `Grow`, `Clip` | Always reassign: `s = slices.Delete(s, i, j)`. Since 1.22 `Delete` zeroes the freed tail |
| In place | `Reverse` | |
| Copy | `Clone` | shallow copy (pointers inside are shared) |
| Iterators (1.23+) | `All`, `Values`, `Backward`, `Chunk`, `Collect`, `AppendSeq`, `Sorted`, `SortedFunc` | bridge between slices and `iter.Seq` |

Patterns worth remembering:

```go
slices.Sort(s); s = slices.Compact(s)               // dedupe: Compact only removes ADJACENT duplicates
s = slices.DeleteFunc(s, func(x int) bool { return x < 0 }) // filter in place
i, found := slices.BinarySearch(s, x)
if !found { s = slices.Insert(s, i, x) }            // keep a slice sorted (ordered-set style)
for i, v := range slices.Backward(s) { ... }        // iterate in reverse
for chunk := range slices.Chunk(s, 100) { ... }     // batches of 100 (handy for pipelines)
safe := slices.Clip(s)                              // cap = len: a caller's append can't overwrite your backing array
```

**`maps`**

| Function | Notes |
|---|---|
| `Keys(m)`, `Values(m)`, `All(m)` | return **iterators** (1.23+), not slices: `slices.Collect(maps.Keys(m))` or `slices.Sorted(maps.Keys(m))` |
| `Clone(m)`, `Copy(dst, src)` | shallow |
| `Equal(a, b)`, `EqualFunc` | |
| `DeleteFunc(m, pred)` | filter in place |
| `Collect(seq)`, `Insert(m, seq)` | build or extend a map from an iterator |

Built-ins that pair with these: `clear(m)` empties a map; `clear(s)` zeroes a slice's elements (len is unchanged).

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

**`container/heap` gotchas**

- Always call **`heap.Push(h, x)` / `heap.Pop(h)`** (the package functions), never `h.Push` / `h.Pop`. Your
  methods only append and remove the last element; the package functions do the sifting.
- Your `Pop` removes the **last** element: `heap.Pop` has already swapped the min to the end.
- `Push`/`Pop` need **pointer receivers** (they change the slice header); `Len`/`Less`/`Swap` can use value receivers.
- `heap.Pop(h)` returns `any`, so type-assert it: `x := heap.Pop(h).(int)` (panics if the type is wrong).
- Min-heap by default (`Less` = `<`). For a max-heap, flip `Less` to `>`.
- Array layout: children of `i` are `2i+1` and `2i+2`; the parent is `(i-1)/2`.

**Testing a heap: the invariant in one loop.** Every node must not be "less" than its parent:

```go
for i := 1; i < len(h.data); i++ {
	if h.less(h.data[i], h.data[(i-1)/2]) { t.Fatalf("index %d beats its parent", i) }
}
```

**Ordered set / multiset: no `std::set` or `TreeMap` in Go's standard library**

| Need | Option | Cost |
|---|---|---|
| Small n, or mostly reads | Sorted slice + `slices.BinarySearch` / `slices.Insert` / `slices.Delete` | O(log n) search, **O(n)** insert/delete (memmove is fast; fine for about 10⁴ items) |
| Only min/max needed | Heap (`container/heap` or your `Heap[T]`), plus lazy deletion via a `map[T]int` of pending removals | O(log n) |
| Full ordered set | Third-party: `github.com/google/btree` (generic `BTreeG[T]`), `github.com/tidwall/btree`, `github.com/emirpasic/gods/v2` (red-black tree set/map) | O(log n) all ops |
| Write your own | **Treap** or **skip list** (much simpler than red-black, same expected O(log n)); AVL / red-black if you want the exercise | O(log n) |
| Counting queries on known values | Fenwick tree (BIT) + coordinate compression | O(log n) |

- **Multiset** = ordered structure over distinct keys + `map[T]int` for counts. Or store `(value, uniqueID)` pairs so
  duplicates become distinct keys.
- **LeetCode:** no `go get`. Use a sorted slice, a heap with lazy deletion, a BIT, or a hand-written treap. (LeetCode's Go
  environment has historically included `emirpasic/gods`; check its language-info page before relying on it.)
- **Real-world note:** Redis sorted sets (`ZSET`) use a **skip list + hash map**, not a balanced tree. Skip lists are
  easier to implement, and support range queries and rank lookups well.

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

### 2.4.1 Loops, receivers, shadowing

- The `for init; cond; post` **init runs once**. `for v, ok := f(); ok; {}` never calls `f` again, so ok stays
  true forever. Use `for { v, ok := f(); if !ok { break }; ... }`.
- **Be consistent with receivers:** if any method needs a pointer receiver, give *all* methods pointer
  receivers (the Go style guide says so, and linters flag mixed receivers).
- **Don't shadow builtins:** naming a variable `min`, `max`, `len`, `cap`, `new` or `clear` hides the builtin
  for the rest of the scope. It's legal, but it causes confusing errors.
- In tests, **`t.Fatalf` stops the test immediately** (it calls `runtime.Goexit`), so code after it never
  runs. `t.Errorf` records the failure and continues.

### 2.4.2 Recursion in Go

- **No tail-call optimization.** Even tail-recursive functions grow the stack.
- **Goroutine stacks grow dynamically** (they start at a few KB and are copied to bigger ones when needed), so deep
  recursion is fine. The default limit is **1 GB** on 64-bit (`debug.SetMaxStack`), and going past it is a fatal
  `goroutine stack exceeds` error, not a recoverable panic.
- **A recursive closure must be declared before it's assigned**, because it refers to itself:

```go
var walk func(i int) bool
walk = func(i int) bool { /* ... */ return walk(2*i+1) && walk(2*i+2) }
```

### 2.4.3 Function literals and closures

- Functions are first-class values. Pass one inline as a **function literal**: `NewHeap(func(a, b int) bool { return a < b })`.
- **There's no short lambda syntax** (no `=>`). Always write the full `func(params) results { body }`.
- Name function types for readability: `type LessFunc[T any] func(a, b T) bool`.
- **Closures capture variables by reference**, so they can read and change outer variables (e.g. a `count++` inside a
  comparator counts comparisons).
- You can call one immediately: `func() { ... }()`. This is the common form for `go func(){...}()` and `defer func(){...}()`.
- **Go 1.22+: each loop iteration gets a fresh variable**, so closures in loops capture the value you expect. The old
  workaround `v := v` is no longer needed.
- **"Inlining" is a compiler optimization, not a keyword:** small functions are automatically copied into their
  callers. See the decisions with `go build -gcflags=-m ./pkg 2>&1 | grep inline`.

### 2.4.4 Packages and file names

- **One directory = one package.** Every `.go` file in a folder must declare the same `package` name, or the
  folder won't build ("found packages X and Y in dir"). To split code, make a subfolder.
- The exception: test files may use `package foo_test` (an external test package that only sees the exported API).
- Package names are short and lowercase with no underscores, and usually **match the folder name**.
- File names are lowercase; underscores are fine, **but some suffixes are special**: `_test.go` marks a test file;
  `_linux.go`, `_darwin.go`, `_windows.go`, `_amd64.go`, `_arm64.go` are build tags that compile only on that OS or CPU.

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

### What an interface body can contain

- **Method signatures**, and **embedded interfaces** (their methods are merged in; duplicate identical
  methods are fine since Go 1.14).
- **No fields and no method bodies.**
- Since 1.18, **type elements** too: `~int | ~float64`. An interface with these is a **constraint**,
  usable only in type parameters (`func Sum[T Number]`), never as a variable type. `cmp.Ordered` is one.

### Constraints in practice

```go
type Number interface { // user-defined, not built in
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~float32 | ~float64
}
func Sum[T Number](xs ...T) T { var t T; for _, x := range xs { t += x }; return t }
```

- `|` is a union; `~int` means **any type whose underlying type is int** (so `type Celsius float64` qualifies with `~float64`).
- Inside a generic function you can only use operations that **every** type in the constraint supports (`+` is fine
  for `Number`, not for `any`).
- `xs ...T` is **variadic**: call it as `Sum(1, 2, 3)` or `Sum(slice...)`; inside the function it's a `[]T`.
- Ready-made constraints: `any`, `comparable` (built in), `cmp.Ordered` (stdlib), and `golang.org/x/exp/constraints`
  (external: `Integer`, `Float`, `Signed`, `Unsigned`).
- **Explicit instantiation** (`NewHeap[int](...)`) is valid. It's required when there's nothing to infer from
  (`NewStack[int]()`), when untyped constants would default wrong (`Sum[float64](1, 2)`), or when you use a
  generic func as a value (`f := NewHeap[int]`).
- **Go 1.21+ infers generic functions passed as arguments:** with `func less[T cmp.Ordered](a, b T) bool`,
  `NewHeap[int](less)` works, because Go instantiates `less[int]` from the expected parameter type.

### Method sets: who implements what

| Type | Method set |
|---|---|
| `T` | value-receiver methods only |
| `*T` | value **and** pointer-receiver methods |

If any required method has a pointer receiver, only `*T` implements the interface. Example: `IntHeap` has
pointer-receiver `Push`/`Pop`, so `heap.Init(IntHeap{...})` fails with *"method Pop has pointer receiver"*.
Use `&IntHeap{...}`. (Reason: a value inside an interface is a copy with no stable address, so a
pointer-receiver method can't safely be called on it.)

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
