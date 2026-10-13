package da14heap

import (
	"math/rand"
	"sort"
	"testing"
)

func intLess(a, b int) bool {
	return a < b
}

func intGreater(a, b int) bool {
	return a > b
}
func TestHeapPushPop(t *testing.T) {
	h := NewHeap(intLess)

	values := []int{5, 3, 8, 1, 9, 2, 7}

	for _, v := range values {
		h.Push(v)
	}

	if got, want := h.Size(), len(values); got != want {
		t.Fatalf("Size() = %d, want %d", got, want)
	}

	sort.Ints(values)

	for _, want := range values {
		got, ok := h.Pop()

		if !ok {
			t.Fatalf("Pop() returned ok=false, want %d", want)
		}

		if got != want {
			t.Fatalf("Pop() = %d, want %d", got, want)
		}
	}

	if h.Size() != 0 {
		t.Fatalf("heap size = %d after popping everything, want 0", h.Size())
	}

	_, ok := h.Pop()
	if ok {
		t.Fatal("Pop() on empty heap returned ok=true")
	}
}

func TestHeapEdgeCases(t *testing.T) {
	h1 := NewHeap(intLess)
	h2 := NewHeap(intGreater)
	arr := []int{5, 3, 5, 1, 9, 2, 7}
	//pop from empty heap
	_, ok := h1.Pop()
	if ok {
		t.Fatal("Pop() on empty heap returned ok=true")
	}
	for _, v := range arr {
		h1.Push(v)
		h2.Push(v)
	}
	sort.Ints(arr)
	//test min and max heaps
	if min, ok1 := h1.Pop(); !ok1 || min != 1 {
		t.Fatalf("min heap pop = %d, want 1", min)
	}

	if max, ok2 := h2.Pop(); !ok2 || max != 9 {
		t.Fatalf("max heap pop = %d, want 9", max)
	}

}
func assertHeapInvariant(nodeIndex int, h *Heap[int], t *testing.T) (int, bool) {
	leftChildIndex := 2*nodeIndex + 1
	rightChildIndex := 2*nodeIndex + 2

	if leftChildIndex < h.Size() {
		if h.less(h.data[leftChildIndex], h.data[nodeIndex]) {
			t.Errorf(
				"Heap invariant violated at index %d: parent %d > left child %d",
				nodeIndex,
				h.data[nodeIndex],
				h.data[leftChildIndex],
			)
			return leftChildIndex, false
		}
	}

	if rightChildIndex < h.Size() {
		if h.less(h.data[rightChildIndex], h.data[nodeIndex]) {
			t.Errorf(
				"Heap invariant violated at index %d: parent %d > right child %d",
				nodeIndex,
				h.data[nodeIndex],
				h.data[rightChildIndex],
			)
			return rightChildIndex, false
		}
	}
	//recurse on both children
	if leftChildIndex < h.Size() {
		f1, ok1 := assertHeapInvariant(leftChildIndex, h, t)
		if !ok1 {
			return f1, false
		}
	}
	if rightChildIndex < h.Size() {
		f2, ok2 := assertHeapInvariant(rightChildIndex, h, t)
		if !ok2 {
			return f2, false //we'll miss this if both fail but it's okay for testing purposes
		}
	}
	return -1, true

}

func TestHeapRandomized(t *testing.T) {
	h := NewHeap(intLess)
	rng := rand.New(rand.NewSource(42))
	values := make([]int, 1000)
	for i := range values {
		values[i] = rng.Intn(10000)
		h.Push(values[i])
	}
	//assert heap invariant
	failureIndex, ok := assertHeapInvariant(0, h, t)

	if !ok {
		t.Fatalf("Heap invariant violated at index %d", failureIndex)
	}
	index := 0
	sort.Ints(values)
	for {
		got, ok := h.Pop()
		if !ok {
			break
		}
		want := values[index]
		if got != want {
			t.Fatalf("Pop() = %d, want %d", got, want)
			break
		}
		index++
	}

}
