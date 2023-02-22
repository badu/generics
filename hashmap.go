package generics

import (
	"math"
)

const (
	minBuckets = 16
	loadFactor = 0.75
)

// Payload is the main payload of the hash map.
type payload[K, V comparable] struct {
	key K
	val V
}

// HashMap is the main hash map type.
type HashMap[K, V comparable] struct {
	capacity int
	data     []LinkedList[payload[K, V]]
	count    int
}

// NewHashMap is used to create a new map.
func NewHashMap[K, V comparable]() HashMap[K, V] {
	return HashMap[K, V]{
		capacity: minBuckets,
		data:     make([]LinkedList[payload[K, V]], minBuckets),
		count:    0,
	}
}

// Count returns the amount of entries in the map.
func (m *HashMap[K, V]) Count() int {
	return m.count
}

// Empty returns true if the map is empty, false if not.
func (m *HashMap[K, V]) Empty() bool {
	return m.count == 0
}

// Resize reallocates and resizes the underlying array to the passed number of buckets.
func (m *HashMap[K, V]) resize(cap int) {
	if cap < minBuckets {
		return // can only resize greater or equal to the minimum bucket size
	}

	data := m.data

	m.capacity = cap
	m.data = make([]LinkedList[payload[K, V]], m.capacity)
	m.count = 0

	for _, ln := range data {
		ln.ForEach(func(i int, p payload[K, V]) {
			m.Put(p.key, p.val)
		})
	}
}

// Put adds a value to the hash map relating to the passed key.
func (m *HashMap[K, V]) Put(key K, val V) {
	if m.count+1 >= int(float64(m.capacity)*loadFactor) {
		m.resize(m.capacity * 2)
	}

	hash := Hash(key)
	bucket := int(math.Mod(float64(hash), float64(m.capacity)))

	var keyExists bool
	m.data[bucket].ForEach(func(i int, p payload[K, V]) {
		if p.key == key {
			keyExists = true
			m.data[bucket].Update(i, payload[K, V]{key: key, val: val})
			return
		}
	})

	if !keyExists {
		m.data[bucket].PutOnBottom(payload[K, V]{key: key, val: val})
		m.count++
	}
}

// Get gets a value from the hash map relating to the passed key.
func (m *HashMap[K, V]) Get(key K) (val V, ok bool) {
	hash := Hash(key)
	bucket := int(math.Mod(float64(hash), float64(m.capacity)))

	m.data[bucket].ForEach(func(i int, p payload[K, V]) {
		if p.key == key {
			val = p.val
			ok = true
			return
		}
	})

	return val, ok
}

// Remove deletes a value from the hash map relating to the passed key.
func (m *HashMap[K, V]) Remove(key K) {
	hash := Hash(key)
	bucket := int(math.Mod(float64(hash), float64(m.capacity)))

	m.data[bucket].ForEach(func(i int, p payload[K, V]) {
		if p.key == key {
			m.data[bucket].Remove(i)
			m.count--
			return
		}
	})

	if (m.capacity/2) >= minBuckets && m.count < int((float64(m.capacity)/2)*loadFactor) {
		m.resize(m.capacity / 2)
	}
}

// ContainsValue returns true if the passed value is present, false if not.
func (m *HashMap[K, V]) ContainsValue(val V) bool {
	var result bool
	for _, ln := range m.data {
		ln.ForEach(func(i int, p payload[K, V]) {
			if p.val == val {
				result = true
				return
			}
		})
	}
	return result
}

// ContainsKey returns true if the passed key is present, false if not.
func (m *HashMap[K, V]) ContainsKey(key K) bool {
	var result bool
	for _, ln := range m.data {
		ln.ForEach(func(i int, p payload[K, V]) {
			if p.key == key {
				result = true
				return
			}
		})
	}
	return result
}

// Clear empties the entire hash map.
func (m *HashMap[K, V]) Clear() {
	m.capacity = minBuckets
	m.data = make([]LinkedList[payload[K, V]], minBuckets)
	m.count = 0
}

// ForEach iterates over the dataset within the hash map, calling the passed
// function for each value.
func (m *HashMap[K, V]) ForEach(f func(key K, val V)) {
	for _, ln := range m.data {
		ln.ForEach(func(i int, p payload[K, V]) {
			f(p.key, p.val)
		})
	}
}
