package generics

// Queue is the main queue type.
type Queue[T comparable] struct {
	data []T
}

const defaultCapacity = 16

// NewQueue is used to create a new queue.
func NewQueue[T comparable](capacity ...int) Queue[T] {
	if len(capacity) == 1 {
		return Queue[T]{
			data: make([]T, 0, capacity[0]),
		}
	}
	return Queue[T]{
		data: make([]T, 0, defaultCapacity),
	}
}

// Count returns the amount of entries in the queue.
func (q *Queue[T]) Count() int {
	return len(q.data)
}

// Empty returns true if the queue is empty, false if not.
func (q *Queue[T]) Empty() bool {
	return len(q.data) == 0
}

// Enqueue add a value to the queue.
func (q *Queue[T]) Enqueue(val T) {
	q.data = append(q.data, val)
}

// Peek returns the first value.
func (q *Queue[T]) Peek() T {
	return q.data[0]
}

// Dequeue returns the first value and removes it.
func (q *Queue[T]) Dequeue() *T {
	if len(q.data) == 0 {
		return nil
	}
	result := q.data[0]
	q.data = q.data[1:len(q.data)]
	return &result
}

// Contains returns true if the value exists in the queue, false if not.
func (q *Queue[T]) Contains(val T) bool {
	for _, v := range q.data {
		if v == val {
			return true
		}
	}
	return false
}

// Clear empties the entire queue.
func (q *Queue[T]) Clear() {
	q.data = q.data[:0:0]
}

// ForEach iterates over the dataset within the queue, calling the passed function for each value.
func (q *Queue[T]) ForEach(predicate func(val T)) {
	for _, v := range q.data {
		predicate(v)
	}
}
