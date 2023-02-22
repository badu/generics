package generics

import (
	"golang.org/x/exp/slices"
)

// BinaryHeap is the main heap type.
type BinaryHeap[T comparable] struct {
	data      []T
	predicate func(a T, b T) bool
	isSorted  bool
}

// NewBinaryHeap is used to create a new heap.
// The passed function is a predicate that returns true if the first parameter is greater than the second. This predicate defines the sorting order between the heap items and is called during insertion and extraction.
func NewBinaryHeap[T comparable](predicate func(a, b T) bool, capacity ...int) BinaryHeap[T] {
	if len(capacity) == 1 {
		return BinaryHeap[T]{
			data:      make([]T, 0, capacity[0]),
			predicate: predicate,
		}
	}
	return BinaryHeap[T]{
		data:      make([]T, 0, defaultCapacity),
		predicate: predicate,
	}
}

// Count returns the amount of entries in the heap.
func (b *BinaryHeap[T]) Count() int {
	return len(b.data)
}

// Empty returns true if the heap is empty, false if not.
func (b *BinaryHeap[T]) Empty() bool {
	return len(b.data) == 0
}

// siftUp sifts the value at the passed index up through the heap until it finds its correct position.
func (b *BinaryHeap[T]) siftUp(childIndex int) {
	var parent T
	var child T
	var parentIndex int

	if childIndex > 0 {
		if childIndex > 2 {
			if childIndex%2 == 0 {
				parentIndex = (childIndex - 2) / 2
			} else {
				parentIndex = (childIndex - 1) / 2
			}
		} else {
			parentIndex = 0
		}

		parent = b.data[parentIndex]
		child = b.data[childIndex]

		if b.predicate(child, parent) {
			b.data[parentIndex] = child
			b.data[childIndex] = parent

			if parentIndex > 0 {
				b.siftUp(parentIndex)
			}
		}
	}
}

// SiftDown sifts the value at the passed index down through the heap until it finds its correct position.
func (b *BinaryHeap[T]) siftDown(parentIndex int) {
	var (
		parent, child1, child2 T
	)

	child1Index := (2 * parentIndex) + 1
	child2Index := (2 * parentIndex) + 2

	if len(b.data) <= child1Index { // The parent has no children.
		return

	}

	if len(b.data) == child2Index { // The parent has one child.
		parent = b.data[parentIndex]
		child1 = b.data[child1Index]

		if b.predicate(child1, parent) {
			b.data[parentIndex] = child1
			b.data[child1Index] = parent
			b.siftDown(child1Index)
		}
		return
	}

	// parent has two children.
	parent = b.data[parentIndex]
	child1 = b.data[child1Index]
	child2 = b.data[child2Index]

	// compare the parent against the greater child
	if b.predicate(child1, child2) {
		if b.predicate(child1, parent) {
			b.data[parentIndex] = child1
			b.data[child1Index] = parent
			b.siftDown(child1Index)
		}
		return
	}

	if b.predicate(child2, parent) {
		b.data[parentIndex] = child2
		b.data[child2Index] = parent
		b.siftDown(child2Index)
	}
}

// Push inserts a new value into the heap.
func (b *BinaryHeap[T]) Push(val T) {
	b.data = append(b.data, val)
	b.siftUp(len(b.data) - 1)
	b.isSorted = false
}

// Peek returns the first value at the top of the heap.
func (b *BinaryHeap[T]) Peek() T {
	return b.data[0]
}

// Pop returns and removes the first value from the heap.
func (b *BinaryHeap[T]) Pop() *T {
	if b.Empty() {
		return nil
	}
	result := b.data[0]
	b.data[0] = b.data[len(b.data)-1]
	b.data = b.data[0 : len(b.data)-1]
	b.siftDown(0)
	b.isSorted = false
	return &result
}

// Has returns true if the value exists in the heap, false if not.
func (b *BinaryHeap[T]) Has(val T) bool {
	for _, v := range b.data {
		if v == val {
			return true
		}
	}
	return false
}

// Clear empties the entire heap.
func (b *BinaryHeap[T]) Clear() {
	b.data = b.data[:0:0]
}

// Sort the heap ready for iterating.
// The heap needs to be sorted to be iterated correctly using a loop. This is to make sure values are delivered in the correct order. Sorted or unsorted, the
// heap data structure will always be correct because of the algorithms used when sifting.
func (b *BinaryHeap[T]) sort() {
	if !b.isSorted {
		slices.SortFunc(b.data, b.predicate)
		b.isSorted = true
	}
}

// ForEach iterates over the dataset within the heap, calling the passed function for each value.
func (b *BinaryHeap[T]) ForEach(predicate func(val T)) {
	b.sort()
	for _, v := range b.data {
		predicate(v)
	}
}
