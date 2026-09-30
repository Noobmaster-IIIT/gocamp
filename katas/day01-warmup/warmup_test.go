package warmup

import (
	"testing"

	"slices"
)

// Functions to test the warmup package. The tests are written using the Go testing framework
// and cover the functionality of the Stack type and the Reverse function.

// TestStack tests the Stack type's methods: Push, Pop, Peek, and Len.
func TestStack(t *testing.T) {
	s := &Stack[int]{}
	if s.Len() != 0 {
		t.Errorf("Len() = %d, want 0", s.Len())
	}
	s.Push(1)
	s.Push(2)
	if s.Len() != 2 {
		t.Errorf("Len() = %d, want 2", s.Len())
	}
	v, ok := s.Peek()
	if !ok || v != 2 {
		t.Errorf("Peek() = %v, %v, want 2, true", v, ok)
	}
	v, ok = s.Pop()
	if !ok || v != 2 {
		t.Errorf("Pop() = %v, %v, want 2, true", v, ok)
	}
	v, ok = s.Pop()
	if !ok || v != 1 {
		t.Errorf("Pop() = %v, %v, want 1, true", v, ok)
	}
	v, ok = s.Pop()
	if ok {
		t.Errorf("Pop() = %v, %v, want 0, false", v, ok)
	}
}

// TestReverse tests the Reverse function with various input slices.
func TestReverse(t *testing.T) {
	tests := []struct {
		name string
		in   []int
		want []int
	}{
		{"empty", []int{}, []int{}},
		{"nil", nil, nil},
		{"one", []int{1}, []int{1}},
		{"two", []int{1, 2}, []int{2, 1}},
		{"three", []int{1, 2, 3}, []int{3, 2, 1}},
		{"four", []int{1, 2, 3, 4}, []int{4, 3, 2, 1}},
		{"five", []int{1, 2, 3, 4, 5}, []int{5, 4, 3, 2, 1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := slices.Clone(tt.in) // Clone the input slice to avoid modifying the original test case
			Reverse(tt.in)
			if !slices.Equal(tt.in, tt.want) {
				t.Errorf("Reverse(%v) = %v, want %v", in, tt.in, tt.want)
			}
		})
	}
}
