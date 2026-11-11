package dayxskiplist

import (
	"cmp"
	"math/rand/v2"
	"sync"

	"github.com/Noobmaster-IIIT/gocamp/katas/day02-dataStructures/deque"
)

// will implement a skip list in this file. A skip list is a data structure that allows
// fast search within an ordered sequence of elements. It allows O(log n) average time
// complexity for search, insert, and delete operations.
// SafeMap helps us check association in O(1) time
type Node[T cmp.Ordered, V any] struct {
	key  T
	val  V
	next []*Node[T, V] // slice of pointers to next nodes at different levels

}
type SafeMap[T cmp.Ordered] struct {
	// map to save k,v pairs for checking association in O(1) time

	lock sync.Mutex
	data map[T]struct{}
}

func (sm *SafeMap[T]) Contains(key T) bool {
	sm.lock.Lock()
	defer sm.lock.Unlock()

	_, exists := sm.data[key]
	return exists
}

func (sm *SafeMap[T]) Add(key T) bool {
	sm.lock.Lock()
	defer sm.lock.Unlock()

	if _, exists := sm.data[key]; exists {
		return false
	}

	sm.data[key] = struct{}{}
	return true
}

func (sm *SafeMap[T]) Delete(key T) {
	sm.lock.Lock()
	defer sm.lock.Unlock()

	delete(sm.data, key)
}

type SkipList[T cmp.Ordered, V any] struct {
	head        *Node[T, V] // pointer to the head node
	maxLevel    int         // maximum level of the skip list
	probability float64     // probability for random level generation
	size        int         // number of elements in the skip list
	mp          *SafeMap[T]
}

func NewSkipList[T cmp.Ordered, V any](maxLevel int, probability float64) *SkipList[T, V] {
	return &SkipList[T, V]{
		head:        &Node[T, V]{next: make([]*Node[T, V], maxLevel)},
		maxLevel:    maxLevel,
		probability: probability,
		size:        0,
		mp:          &SafeMap[T]{data: make(map[T]struct{})},
	}
}

func (sl *SkipList[T, V]) Size() int {
	return sl.size
}

// find returns the node holding k, or nil if k is not present.
// It's the shared walk used by Search and Insert.
func (sl *SkipList[T, V]) find(k T) *Node[T, V] {
	if sl.maxLevel == 0 || !sl.mp.Contains(k) {
		return nil
	}

	node := sl.head

	for level := sl.maxLevel - 1; level >= 0; level-- {
		for node.next[level] != nil &&
			node.next[level].key < k {
			node = node.next[level]
		}
	}

	node = node.next[0]

	if node != nil && node.key == k {
		return node
	}

	return nil
}

func (sl *SkipList[T, V]) Search(k T) (value V, found bool) {
	if node := sl.find(k); node != nil {
		return node.val, true
	}
	var zero V
	return zero, false
}

// Insert adds k with value v, or updates the value if k already exists.
// The list is ordered by KEY, so changing an existing node's value never
// changes its position: update in place instead of delete + reinsert.
func (sl *SkipList[T, V]) Insert(k T, v V) bool {
	if node := sl.find(k); node != nil {
		node.val = v // same key, same position: one walk, no new node, tower height kept
		return true
	}
	return sl.InsertKey(k, v)
}

func (sl *SkipList[T, V]) InsertKey(k T, v V) bool {
	if !sl.mp.Add(k) {
		return false
	}

	levels := 1
	for rand.Float64() > 0.5 {
		levels++
	}

	// Ensure head.next has enough slots for the new highest level.
	if levels > len(sl.head.next) {
		sl.head.next = append(
			sl.head.next,
			make([]*Node[T, V], levels-len(sl.head.next))...,
		)
	}
	sl.maxLevel = max(sl.maxLevel, levels)

	d := deque.NewRingBufferDeque[*Node[T, V]](levels)
	node := sl.head

	for lvl := sl.maxLevel - 1; lvl >= 0; lvl-- {
		for node.next[lvl] != nil && node.next[lvl].key < k {
			node = node.next[lvl]
		}

		if lvl < levels {
			d.PushBack(node)
		}
	}

	curr := &Node[T, V]{
		key:  k,
		val:  v,
		next: make([]*Node[T, V], levels),
	}

	for lvl := 0; lvl < levels; lvl++ {
		prev, ok := d.PopBack()
		if !ok {
			return false // unexpected: no predecessor for this level
		}

		curr.next[lvl] = prev.next[lvl]
		prev.next[lvl] = curr
	}

	sl.size++
	return true
}

// DeleteKey removes k from the skip list. It returns false if k is not present.
func (sl *SkipList[T, V]) DeleteKey(k T) bool {
	if sl.maxLevel == 0 {
		return false
	}

	// Find the predecessor of k at every level.
	update := make([]*Node[T, V], sl.maxLevel)
	node := sl.head
	for level := sl.maxLevel - 1; level >= 0; level-- {
		for node.next[level] != nil && node.next[level].key < k {
			node = node.next[level]
		}
		update[level] = node
	}

	target := update[0].next[0]
	if target == nil || target.key != k {
		return false
	}

	// Unlink the target wherever it appears. Its next slice determines its height.
	for level := 0; level < len(target.next); level++ {
		if update[level].next[level] == target {
			update[level].next[level] = target.next[level]
		}
	}

	sl.mp.Delete(k)
	sl.size--

	// Drop empty top levels from the active level count.
	for sl.maxLevel > 0 && sl.head.next[sl.maxLevel-1] == nil {
		sl.maxLevel--
	}

	return true
}
