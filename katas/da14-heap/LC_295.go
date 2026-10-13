package da14heap

import "cmp"

type Heap[T any] struct {
	data []T
	less func(a, b T) bool
}

func NewHeap[T any](less func(a, b T) bool) *Heap[T] {
	return &Heap[T]{
		data: []T{},
		less: less,
	}

}

func (h *Heap[T]) Size() int {
	return len(h.data)
}

func (h *Heap[T]) Push(v T) {
	h.data = append(h.data, v)
	h.up(len(h.data) - 1) //sift up the last element
}
func (h *Heap[T]) Pop() (T, bool) {
	var zero T
	if len(h.data) == 0 {
		return zero, false
	}
	v := h.data[0]

	last := len(h.data) - 1
	h.data[0] = h.data[last]
	h.data = h.data[:last]

	if len(h.data) > 0 {
		h.down(0)
	}

	return v, true
}
func (h *Heap[T]) up(i int) {
	//sift up the element at index i
	//parent is at (i-1)/2
	for i > 0 {
		parent := (i - 1) / 2
		if h.less(h.data[i], h.data[parent]) {
			h.swap(i, parent) //swap the child and parent values
			i = parent        // update index to parent and continue sifting up
		} else {
			break
		}
	}
}

func (h *Heap[T]) down(i int) {
	//sift down the element at index i
	n := len(h.data)
	for {
		left := 2*i + 1
		right := 2*i + 2
		smallest := i

		if left < n && h.less(h.data[left], h.data[smallest]) {
			smallest = left
		}
		if right < n && h.less(h.data[right], h.data[smallest]) {
			smallest = right
		}
		if smallest == i {
			break
		}
		h.swap(i, smallest)
		i = smallest
	}
}
func (h *Heap[T]) Top() (T, bool) {
	var zero T
	if len(h.data) == 0 {
		return zero, false
	}
	return h.data[0], true
}
func (h *Heap[T]) Empty() bool {
	return h.Size() == 0
}
func (h *Heap[T]) swap(i, j int) {
	h.data[i], h.data[j] = h.data[j], h.data[i]
}

func less[T cmp.Ordered](a, b T) bool {
	return a < b
}
func greater[T cmp.Ordered](a, b T) bool {
	return a > b
}

type MedianFinder struct {
	minh *Heap[int]
	maxh *Heap[int]
}

func Constructor() MedianFinder {
	return MedianFinder{
		minh: NewHeap[int](less),
		maxh: NewHeap[int](greater),
	}
}

func (mf *MedianFinder) Balance() {
	h1, h2 := mf.maxh, mf.minh
	if h1.Size() == h2.Size() {
		return
	}
	if h1.Size() > h2.Size()+1 {
		//pop one and move to min heap
		a, _ := h1.Pop()
		h2.Push(a)
		return
	}
	if h2.Size() > h1.Size() { // we can always pull fom left if odd (left biased)
		a, _ := h2.Pop()
		h1.Push(a)
		return
	}
}
func (mf *MedianFinder) Size() int {
	return mf.minh.Size() + mf.maxh.Size()
}
func (mf *MedianFinder) AddNum(num int) {
	h1 := mf.maxh
	h2 := mf.minh
	if h1.Empty() && h2.Empty() {
		//pick left
		h1.Push(num) //left biased
		//no balance needed
		return
	}
	tp, _ := h1.Top()
	if num > tp {
		// after last balance size of h1>=h2 hopefully
		h2.Push(num)
		mf.Balance()
		return
	}
	tp, _ = h2.Top()
	if h2.Empty() || num <= tp {
		h1.Push(num)
		mf.Balance()
		return
	}

}

func (mf *MedianFinder) FindMedian() float64 {
	n := mf.Size()
	var res float64
	if n&1 == 1 {
		tp, _ := mf.maxh.Top()
		res = float64(tp)
	} else {
		tp1, _ := mf.maxh.Top()
		tp2, _ := mf.minh.Top()
		res = (float64(tp1) + float64(tp2)) / 2.0
	}

	return res

}

/**
 * Your MedianFinder object will be instantiated and called as such:
 * obj := Constructor();
 * obj.AddNum(num);
 * param_2 := obj.FindMedian();
 */
