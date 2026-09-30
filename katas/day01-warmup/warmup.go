package warmup

// 1. Reverse[T any](s []T) reverses the slice in place. That means no new slice and no return value, just swapping elements. Hint: two indices moving toward each other.

// 2. Stack[T any] needs:

// Push(v T)
// Pop() (T, bool): returns false on an empty stack, never panics
// Peek() (T, bool)
// Len() int

type Stack[T any] struct {
	elements []T
	// no need to store the length, slices already have a length
}

// Push adds an element to the top of the stack.
func (s *Stack[T]) Push(v T) {
	s.elements = append(s.elements, v)
}

// Pop removes and returns the top element of the stack. It returns false if the stack is empty.
func (s *Stack[T]) Pop() (T, bool) {
	if len(s.elements) == 0 {
		var zero T
		return zero, false
	}
	result := s.elements[len(s.elements)-1]
	s.elements = s.elements[:len(s.elements)-1]
	return result, true
}

// Peek returns the top element of the stack without removing it. It returns false if the stack is empty.
func (s *Stack[T]) Peek() (T, bool) {
	if len(s.elements) == 0 {
		var zero T
		return zero, false
	}
	return s.elements[len(s.elements)-1], true
}

// Len returns the number of elements in the stack.
func (s *Stack[T]) Len() int {
	return len(s.elements)
}

// Reverse reverses the elements of the slice in place.
func Reverse[T any](s []T) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}
