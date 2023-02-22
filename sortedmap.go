package generics

import (
	"fmt"
	"strings"
)

func NewSortedMap[K ordered, V any]() *SortedMap[K, V] { return new(SortedMap[K, V]) }

type ordered interface {
	int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | uintptr | float32 | float64 | string
}

type SortedMap[K ordered, V any] struct {
	root *nodeMap[K, V]
}

func (m *SortedMap[K, V]) Get(k K) V {
	return m.root.get(k)
}

func (m *SortedMap[K, V]) Set(k K, v V) {
	store(&m.root, k, v)
}

func (m *SortedMap[K, V]) Keys() []K {
	return m.root.appendKeys(nil)
}

func (m *SortedMap[K, V]) String() string {
	buf := strings.Builder{}
	buf.WriteRune('{')
	for i, key := range m.Keys() {
		if i > 0 {
			buf.WriteRune(',')
		}

		fmt.Fprintf(&buf, "%v:%v", key, m.Get(key))
	}
	buf.WriteRune('}')
	return buf.String()
}

// unbalanced binary tree
type nodeMap[K ordered, V any] struct {
	left, right *nodeMap[K, V]
	key         K
	value       V
}

func store[K ordered, V any](naddr **nodeMap[K, V], key K, value V) {
	if *naddr == nil {
		*naddr = &nodeMap[K, V]{key: key, value: value}
		return
	}

	n := *naddr
	if key < n.key {
		store(&n.left, key, value)
		return
	}

	if key > n.key {
		store(&n.right, key, value)
		return
	}

	n.value = value
}

func (n *nodeMap[K, V]) get(k K) V {
	if n == nil {
		var v V
		return v // empty
	}

	if k < n.key {
		return n.left.get(k)
	}

	if k > n.key {
		return n.right.get(k)
	}

	return n.value
}

func (n *nodeMap[K, V]) appendKeys(out []K) []K {
	if n != nil {
		out = n.left.appendKeys(out)
		out = append(out, n.key)
		out = n.right.appendKeys(out)
	}

	return out
}
