package deque

// Deque is a double-ended queue.
type Deque[T any] interface {
	// primitives like push, pop, peek, etc.
	PushFront(value T)
	PushBack(value T)
	PopFront() (T, bool)
	PopBack() (T, bool)
	PeekFront() (T, bool)
	PeekBack() (T, bool)

	// Size returns the number of elements in the deque
	Size() int

	// IsEmpty returns true if the deque is empty
	IsEmpty() bool

	// Clear removes all elements from the deque
	Clear()
}

// Compile-time checks: the build fails if either type stops satisfying Deque.
// (*X)(nil) is a typed nil pointer; assigning it to the blank identifier costs nothing at runtime.
var (
	_ Deque[int] = (*RingBufferDeque[int])(nil)
	_ Deque[int] = (*LinkedListDeque[int])(nil)
)

// First deque impl: A ring buffer that grows as needed. It is not thread-safe for now.
// After learning concurrency in go, we can add a thread-safe version of this deque. For now, we will focus on the basic implementation.

// RingBufferDeque is a growable circular-buffer deque. All operations are O(1) amortized.
type RingBufferDeque[T any] struct {
	buffer   []T
	head     int // index of the front element
	tail     int // index of the next free slot at the back
	size     int // number of elements; needed because head == tail means BOTH empty and full
	capacity int
}

// NewRingBufferDeque returns an empty deque. A non-positive initialCapacity defaults to 16.
func NewRingBufferDeque[T any](initialCapacity int) *RingBufferDeque[T] {
	if initialCapacity <= 0 {
		initialCapacity = 16 // default capacity
	}
	// head, tail and size start at their zero value 0, so they don't need to be listed.
	// tail is the next available slot. head == tail happens both when empty and when full,
	// so size is what tells those two states apart.
	return &RingBufferDeque[T]{
		buffer:   make([]T, initialCapacity),
		capacity: initialCapacity,
	}
}

// PushFront adds value at the front.
func (dq *RingBufferDeque[T]) PushFront(value T) {
	if dq.size == dq.capacity {
		dq.resize() // allocates a 2x buffer and copies the elements over in order
	}
	// push front adds one element to the front. so (head-1)%capacity is the new head index. we need to wrap around if head is 0, so we use (head-1+capacity)%capacity
	// (Go's % keeps the sign of the dividend: -1 % 8 == -1, not 7. That's why the +capacity is needed.)
	dq.head = (dq.head - 1 + dq.capacity) % dq.capacity
	dq.buffer[dq.head] = value
	dq.size++
}

// PushBack adds value at the back.
func (dq *RingBufferDeque[T]) PushBack(value T) {
	if dq.size == dq.capacity {
		dq.resize()
	}
	dq.buffer[dq.tail] = value
	dq.tail = (dq.tail + 1) % dq.capacity
	dq.size++
}

// PopFront removes and returns the front element, or false if the deque is empty.
func (dq *RingBufferDeque[T]) PopFront() (T, bool) {
	var zero T
	if dq.size == 0 {
		return zero, false
	}
	value := dq.buffer[dq.head]
	dq.buffer[dq.head] = zero // drop the reference so the GC can reclaim it (same lesson as Stack.Pop)
	dq.head = (dq.head + 1) % dq.capacity
	dq.size--
	return value, true
}

// PopBack removes and returns the back element, or false if the deque is empty.
func (dq *RingBufferDeque[T]) PopBack() (T, bool) {
	var zero T
	if dq.size == 0 {
		return zero, false
	}
	dq.tail = (dq.tail - 1 + dq.capacity) % dq.capacity
	value := dq.buffer[dq.tail]
	dq.buffer[dq.tail] = zero
	dq.size--
	return value, true
}

// PeekFront returns the front element without removing it, or false if the deque is empty.
func (dq *RingBufferDeque[T]) PeekFront() (T, bool) {
	if dq.size == 0 {
		var zero T
		return zero, false
	}
	return dq.buffer[dq.head], true
}

// PeekBack returns the back element without removing it, or false if the deque is empty.
func (dq *RingBufferDeque[T]) PeekBack() (T, bool) {
	if dq.size == 0 {
		var zero T
		return zero, false
	}
	backIndex := (dq.tail - 1 + dq.capacity) % dq.capacity
	return dq.buffer[backIndex], true
}

// Size returns the number of elements.
func (dq *RingBufferDeque[T]) Size() int {
	return dq.size
}

// IsEmpty reports whether the deque has no elements.
func (dq *RingBufferDeque[T]) IsEmpty() bool {
	return dq.size == 0
}

// resize doubles the capacity. It is lowercase (unexported) because it's an internal helper.
// Go treats resize and Resize as two different names.
func (dq *RingBufferDeque[T]) resize() {
	newCapacity := dq.capacity << 1 // double the capacity
	newBuffer := make([]T, newCapacity)

	// A plain copy would lose the order of elements (they may wrap around the end), so copy them
	// in logical order, front to back. Loop by COUNT, not until i == tail: resize is only called
	// when the deque is full, and when it's full head == tail, so that loop would copy nothing.
	for k := 0; k < dq.size; k++ {
		newBuffer[k] = dq.buffer[(dq.head+k)%dq.capacity]
	}
	// No need to zero the old buffer. Once dq.buffer points at newBuffer, nothing references the
	// old array, so the GC frees all of it. Zeroing only matters for slots in an array you KEEP.

	// 	dq:= &newBuffer // reassign the dq struct fields to the new buffer and capacity
	// 	THIS IS INCORRECT!!!
	// 	Go passes everything by value, including pointers, so passing *T copies the pointer/address rather than passing the pointer variable itself by reference.
	// Because the copied pointer still points to the same object, modifying *p or p.field inside a function changes the caller's object, but assigning p = ... only changes the local pointer.
	// A slice is also passed by value: its header (pointer, len, cap) is copied, while the copied header and caller's header initially point to the same backing array.
	// append may reuse that backing array or allocate a new one, so you must use its returned slice (s = append(s, x)) if you want the caller's slice header to reflect the new length/capacity.
	// Similarly, an interface is a (dynamic type, dynamic value) pair, so an interface containing a typed nil pointer is not nil because its dynamic type is still present.
	dq.buffer = newBuffer
	dq.capacity = newCapacity
	dq.head = 0
	dq.tail = dq.size
}

// Clear removes all elements.
func (dq *RingBufferDeque[T]) Clear() {
	clear(dq.buffer) // built-in (Go 1.21+): sets every element to the zero value so the GC can reclaim them
	dq.head = 0
	dq.tail = 0
	dq.size = 0
}

// Now for the second deque impl: a doubly linked list of fixed-size blocks (an "unrolled" linked
// list, the same idea as C++ std::deque). It is not thread-safe for now. After learning concurrency
// in go, we can add a thread-safe version of this deque.

const blockSize = 1024

// block holds up to blockSize elements in buffer[head:tail].
type block[T any] struct {
	buffer [blockSize]T // fixed-size array: one allocation holds 1024 elements
	next   *block[T]
	prev   *block[T]
	head   int // first used slot
	tail   int // one past the last used slot
}

// LinkedListDeque is a deque made of linked fixed-size blocks.
//
// Invariant: every block in the list holds at least one element, EXCEPT when the deque is empty.
// Then there is exactly one empty block, with head == tail == blockSize/2 so it has room on both sides.
type LinkedListDeque[T any] struct {
	headBlock *block[T]
	tailBlock *block[T]
	size      int
}

// newBlock returns an empty block whose cursor sits at pos.
// The buffer doesn't need initializing: an array's zero value is all zero elements.
func newBlock[T any](pos int) *block[T] {
	return &block[T]{head: pos, tail: pos}
}

// NewLinkedListDeque returns an empty deque.
func NewLinkedListDeque[T any]() *LinkedListDeque[T] {
	// Start in the MIDDLE so the first PushFront doesn't immediately need a new block.
	b := newBlock[T](blockSize / 2)
	return &LinkedListDeque[T]{headBlock: b, tailBlock: b}
}

// PushFront adds value at the front.
func (dq *LinkedListDeque[T]) PushFront(value T) {
	b := dq.headBlock
	if b.head == 0 { // no space on the left: link a new block in front, filled from its right end
		nb := newBlock[T](blockSize)
		nb.next = b
		b.prev = nb
		dq.headBlock = nb
		b = nb
	}
	b.head--
	b.buffer[b.head] = value
	dq.size++
}

// PushBack adds value at the back.
func (dq *LinkedListDeque[T]) PushBack(value T) {
	b := dq.tailBlock
	if b.tail == blockSize { // no space on the right: link a new block behind, filled from its left end
		nb := newBlock[T](0)
		nb.prev = b
		b.next = nb
		dq.tailBlock = nb
		b = nb
	}
	b.buffer[b.tail] = value
	b.tail++
	dq.size++
}

// PopFront removes and returns the front element, or false if the deque is empty.
func (dq *LinkedListDeque[T]) PopFront() (T, bool) {
	var zero T
	if dq.size == 0 {
		return zero, false
	}
	b := dq.headBlock
	value := b.buffer[b.head]
	b.buffer[b.head] = zero // set to zero value of T so that the GC can reclaim the memory
	b.head++
	dq.size--

	// Pops DO need block logic: when a block runs empty, unlink it (or reset it if it's the last one).
	if b.head == b.tail {
		if b.next == nil { // the only block: re-center it for future pushes
			b.head, b.tail = blockSize/2, blockSize/2
		} else {
			dq.headBlock = b.next
			dq.headBlock.prev = nil
			b.next = nil // fully detach so the GC can free the block
		}
	}
	return value, true
}

// PopBack removes and returns the back element, or false if the deque is empty.
func (dq *LinkedListDeque[T]) PopBack() (T, bool) {
	var zero T
	if dq.size == 0 {
		return zero, false
	}
	b := dq.tailBlock
	b.tail-- // tail is one PAST the last element, so step back first, then read
	value := b.buffer[b.tail]
	b.buffer[b.tail] = zero // set to zero value of T so that the GC can reclaim the memory
	dq.size--

	if b.head == b.tail {
		if b.prev == nil {
			b.head, b.tail = blockSize/2, blockSize/2
		} else {
			dq.tailBlock = b.prev
			dq.tailBlock.next = nil
			b.prev = nil
		}
	}
	return value, true
}

// PeekFront returns the front element without removing it, or false if the deque is empty.
func (dq *LinkedListDeque[T]) PeekFront() (T, bool) {
	if dq.size == 0 {
		var zero T
		return zero, false
	}
	return dq.headBlock.buffer[dq.headBlock.head], true
}

// PeekBack returns the back element without removing it, or false if the deque is empty.
func (dq *LinkedListDeque[T]) PeekBack() (T, bool) {
	if dq.size == 0 {
		var zero T
		return zero, false
	}
	return dq.tailBlock.buffer[dq.tailBlock.tail-1], true
}

// Size returns the number of elements.
func (dq *LinkedListDeque[T]) Size() int {
	return dq.size
}

// IsEmpty reports whether the deque has no elements.
func (dq *LinkedListDeque[T]) IsEmpty() bool {
	return dq.size == 0
}

// Clear removes all elements. The old blocks become unreachable and are freed by the GC.
func (dq *LinkedListDeque[T]) Clear() {
	b := newBlock[T](blockSize / 2)
	dq.headBlock, dq.tailBlock, dq.size = b, b, 0
}
