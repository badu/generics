package generics

// Stack is the main stack type.
type Stack[T comparable] struct {
	data []T
}

// NewStack is used to create a new stack.
func NewStack[T comparable]() Stack[T] {
	return Stack[T]{
		data: make([]T, 0, 16),
	}
}

// Count returns the amount of entries in the stack.
func (s *Stack[T]) Count() int {
	return len(s.data)
}

// Empty returns true if the stack is empty, false if not.
func (s *Stack[T]) Empty() bool {
	return len(s.data) == 0
}

// Push adds a value to the stack.
func (s *Stack[T]) Push(val T) {
	s.data = append(s.data, val)
}

// Peek returns the first value.
func (s *Stack[T]) Peek() T {
	return s.data[len(s.data)-1]
}

// Pop returns the first value and removes it.
func (s *Stack[T]) Pop() *T {
	if s.Empty() {
		return nil
	}

	result := s.data[len(s.data)-1]
	s.data = s.data[0 : len(s.data)-1]
	return &result
}

// Contains returns true if the value exists in the stack, false if not.
func (s *Stack[T]) Contains(val T) bool {
	for _, v := range s.data {
		if v == val {
			return true
		}
	}
	return false
}

// Clear empties the entire stack.
func (s *Stack[T]) Clear() {
	s.data = s.data[:0:0]
}

// ForEach iterates over the dataset within the stack, calling the passed
// function for each value.
func (s *Stack[T]) ForEach(f func(val T)) {
	for i := len(s.data) - 1; i >= 0; i-- {
		f(s.data[i])
	}
}
