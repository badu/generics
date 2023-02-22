package generics

import (
	"fmt"
)

// Node is the node type used within the linked list.
type node[T comparable] struct {
	prev  *node[T]
	next  *node[T]
	value T
}

// LinkedList is the main linked list type.
type LinkedList[T comparable] struct {
	first *node[T]
	last  *node[T]
	count int
}

// NewLinkedList is used to create a new linked list.
func NewLinkedList[T comparable]() LinkedList[T] {
	return LinkedList[T]{}
}

// Count returns the amount of entries in the linked list.
func (l *LinkedList[T]) Count() int {
	return l.count
}

// Empty returns true if the linked list is empty, false if not.
func (l *LinkedList[T]) Empty() bool {
	return l.count == 0
}

// PutOnTop inserts a value at the beginning of the linked list.
func (l *LinkedList[T]) PutOnTop(value T) {
	n := &node[T]{value: value}

	if l.first == nil {
		l.first = n
		l.last = n
		l.count++
		return
	}

	l.first.prev = n
	n.next = l.first
	l.first = n
	l.count++
}

// Top returns the value at the beginning of the linked list.
func (l *LinkedList[T]) Top() *T {
	if l.first == nil {
		return nil
	}

	return &l.first.value
}

// PopTop removes the first value in the linked list.
func (l *LinkedList[T]) PopTop() {
	if l.first != nil {
		if l.first.next == nil {
			l.first = nil
			l.last = nil
		} else {
			l.first = l.first.next
			l.first.prev = nil
		}
	}

	l.count--
}

// PutOnBottom inserts a value at the end of the linked list.
func (l *LinkedList[T]) PutOnBottom(val T) {
	n := &node[T]{value: val}

	if l.last == nil {
		l.first = n
		l.last = n
	} else {
		l.last.next = n
		n.prev = l.last
		l.last = n
	}

	l.count++
}

// Bottom returns the value at the end of the linked list.
func (l *LinkedList[T]) Bottom() *T {
	if l.last == nil {
		return nil
	}

	return &l.last.value
}

// PopBottom removes the last value in the linked list.
func (l *LinkedList[T]) PopBottom() {
	if l.last != nil {
		if l.last.prev == nil {
			l.first = nil
			l.last = nil
		} else {
			l.last = l.last.prev
			l.last.next = nil
		}
	}

	l.count--
}

// AddAt inserts a value at the specified index.
func (l *LinkedList[T]) AddAt(value T, index int) error {
	if index > l.Count() {
		return fmt.Errorf("insertion index invalid %d > %d", index, l.Count())
	}

	if index == 0 {
		l.PutOnTop(value)
		return nil
	}

	if index == l.Count() {
		l.PutOnBottom(value)
		return nil
	}

	n := &node[T]{value: value}

	var listIndex int = 0
	for ln := l.first; ln != nil; ln = ln.next {
		if listIndex == index {
			ln.prev.next = n
			n.prev = ln.prev
			n.next = ln
			ln.prev = n
			break
		}
		listIndex++
	}

	l.count++

	return nil
}

// Get gets a value at the specified index.
func (l *LinkedList[T]) Get(index int) *T {
	if index >= l.Count() {
		return nil //index outside of linked list bounds
	}

	if index == 0 {
		return &l.first.value

	}

	if index == l.Count()-1 {
		return &l.last.value

	}

	var listIndex int = 0
	for ln := l.first; ln != nil; ln = ln.next {
		if listIndex == index {
			return &ln.value
		}
		listIndex++
	}

	return nil
}

// Update updates a value at the specified index.
func (l *LinkedList[T]) Update(index int, val T) {
	if index >= l.Count() {
		panic("index outside of linked list bounds")
	}

	var listIndex int = 0
	for ln := l.first; ln != nil; ln = ln.next {
		if listIndex == index {
			ln.value = val
			break
		}
		listIndex++
	}
}

// Remove removes a value at the specified index.
func (l *LinkedList[T]) Remove(index int) {
	if l.count == 0 {
		return
	}

	if index >= l.Count() {
		return //index outside of linked list bounds
	}

	if index == 0 {
		l.PopTop()
		return
	}

	if index == l.Count()-1 {
		l.PopBottom()
		return
	}

	listIndex := 0
	for ln := l.first; ln != nil; ln = ln.next {
		if listIndex == index {
			ln.prev.next = ln.next
			ln.next.prev = ln.prev
			break
		}
		listIndex++
	}

	l.count--

}

// Has returns true if the value exists in the linked list, false if not.
func (l *LinkedList[T]) Has(value T) bool {
	for ln := l.first; ln != nil; ln = ln.next {
		if ln.value == value {
			return true
		}
	}

	return false
}

// Clear empties the entire linked list.
func (l *LinkedList[T]) Clear() {
	l.first = nil
	l.last = nil
	l.count = 0
}

// ForEach iterates over the dataset within the linked list, calling the passed
// function for each value.
func (l *LinkedList[T]) ForEach(f func(i int, val T)) {
	index := 0

	for ln := l.first; ln != nil; ln = ln.next {
		f(index, ln.value)
		index++
	}
}
