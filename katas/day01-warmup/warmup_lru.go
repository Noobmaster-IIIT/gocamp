package warmup

// LRUCache is a fixed-capacity least-recently-used cache backed by flat arrays
// instead of a pointer-based doubly linked list. Slot i holds keys[i], values[i],
// and its list neighbours prev[i] / next[i] (-1 means none).
//
// Same O(1) complexity as the linked-list version, but no per-node allocations
// and no pointers for the GC to scan, plus better cache locality.
type LRUCache struct {
	keys     []int         // keys[i] is the key stored in slot i (replaces a reverse map)
	values   []int         // values[i] is the value stored in slot i
	prev     []int32       // prev[i] is the slot before i in recency order, or -1
	next     []int32       // next[i] is the slot after i in recency order, or -1
	head     int32         // most recently used slot, or -1 if empty
	tail     int32         // least recently used slot, or -1 if empty
	size     int           // number of slots in use
	capacity int           // maximum number of slots
	index    map[int]int32 // key -> slot
}

// Constructor returns an empty cache that holds at most capacity entries (LeetCode 146 signature).
func Constructor(capacity int) LRUCache {
	return LRUCache{
		keys:     make([]int, capacity),
		values:   make([]int, capacity),
		prev:     make([]int32, capacity),
		next:     make([]int32, capacity),
		head:     -1,
		tail:     -1,
		capacity: capacity,
		index:    make(map[int]int32, capacity),
	}
}

// Get returns the value for key and marks it most recently used,
// or -1 if key is not in the cache.
func (c *LRUCache) Get(key int) int {
	i, ok := c.index[key]
	if !ok {
		return -1
	}
	c.moveToFront(i)
	return c.values[i]
}

// Put inserts or updates key and marks it most recently used.
// If the cache is full, the least recently used entry is evicted.
func (c *LRUCache) Put(key, val int) {
	if i, ok := c.index[key]; ok {
		c.values[i] = val
		c.moveToFront(i)
		return
	}

	var i int32
	if c.size < c.capacity {
		// Free slot available: slots are handed out in order 0, 1, 2, ...
		i = int32(c.size) // explicit conversion: Go never converts between integer types implicitly
		c.size++
	} else {
		// Full: reuse the tail's slot for the new entry.
		i = c.tail
		c.unlink(i)
		delete(c.index, c.keys[i])
	}

	c.keys[i] = key
	c.values[i] = val
	c.index[key] = i
	c.pushFront(i)
}

// moveToFront marks slot i as most recently used.
func (c *LRUCache) moveToFront(i int32) {
	if i == c.head {
		return
	}
	c.unlink(i)
	c.pushFront(i)
}

// unlink removes slot i from the list, patching its neighbours and head/tail.
func (c *LRUCache) unlink(i int32) {
	p, n := c.prev[i], c.next[i]
	if p != -1 {
		c.next[p] = n
	} else {
		c.head = n
	}
	if n != -1 {
		c.prev[n] = p
	} else {
		c.tail = p
	}
}

// pushFront makes slot i the new head of the list.
func (c *LRUCache) pushFront(i int32) {
	c.prev[i], c.next[i] = -1, c.head
	if c.head != -1 {
		c.prev[c.head] = i
	}
	c.head = i
	if c.tail == -1 {
		c.tail = i
	}
}
