package generics

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewHashMap(t *testing.T) {
	t.Parallel()

	m := NewHashMap[string, int]()
	assert.True(t, m.Empty())

	m.Put("foo", 3)
	val, ok := m.Get("foo")
	assert.Equal(t, val, 3)
	assert.True(t, ok)

	m.Put("foo", 6)
	val, ok = m.Get("foo")
	assert.Equal(t, val, 6)
	assert.True(t, ok)

	assert.Equal(t, m.Count(), 1)
	assert.False(t, m.Empty())

	m.Remove("foo")
	val, ok = m.Get("foo")
	assert.Equal(t, val, 0)
	assert.False(t, ok)

	m.Remove("foo")
	m.Remove("foo")

	assert.Equal(t, m.Count(), 0)
}

func TestHashMapResizing(t *testing.T) {
	t.Parallel()

	m := NewHashMap[string, int]()
	assert.Equal(t, m.capacity, 16)

	m.Put("a", 1)
	m.Put("b", 2)
	m.Put("c", 3)
	m.Put("d", 4)
	m.Put("e", 5)
	m.Put("f", 6)
	m.Put("g", 7)
	m.Put("h", 8)
	m.Put("i", 9)
	m.Put("j", 10)
	m.Put("k", 11)
	assert.Equal(t, m.Count(), 11)
	assert.Equal(t, m.capacity, 16)
	assert.False(t, m.Empty())

	m.Put("l", 12)
	assert.Equal(t, m.Count(), 12)
	assert.Equal(t, m.capacity, 32)

	m.Remove("l")
	assert.Equal(t, m.Count(), 11)
	assert.Equal(t, m.capacity, 16)
}

func TestHashMapLargeCapacity(t *testing.T) {
	t.Parallel()

	m := NewHashMap[int, int]()
	limit := 1_000_000

	for i := 1; i <= limit; i++ {
		m.Put(i, i)
		assert.Equal(t, m.Count(), i)
	}

	assert.Equal(t, m.Count(), limit)
	assert.True(t, m.ContainsKey(1))
	assert.True(t, m.ContainsValue(1))
	assert.True(t, m.ContainsKey(limit))
	assert.True(t, m.ContainsValue(limit))
	assert.False(t, m.Empty())

	for i := 1; i <= limit; i++ {
		m.Remove(i)
		assert.Equal(t, m.Count(), limit-i)
	}

	assert.True(t, m.Empty())
	assert.Equal(t, m.Count(), 0)
}

func TestHashMapContains(t *testing.T) {
	t.Parallel()

	m := NewHashMap[string, int]()
	assert.Equal(t, m.capacity, 16)

	m.Put("a", 1)
	m.Put("b", 2)
	m.Put("c", 3)
	m.Put("d", 4)
	m.Put("e", 5)
	assert.False(t, m.Empty())

	assert.True(t, m.ContainsKey("a"))
	assert.False(t, m.ContainsKey("f"))

	assert.True(t, m.ContainsValue(3))
	assert.False(t, m.ContainsValue(10))
}

func TestHashMapClearing(t *testing.T) {
	t.Parallel()

	m := NewHashMap[string, string]()
	assert.Equal(t, m.capacity, 16)

	m.Put("a", "1")
	m.Put("b", "2")
	m.Put("c", "3")
	m.Put("d", "4")
	m.Put("e", "5")
	m.Put("f", "6")
	m.Put("g", "7")
	m.Put("h", "8")
	m.Put("i", "9")
	m.Put("j", "10")
	m.Put("k", "11")
	m.Put("l", "12")
	assert.Equal(t, m.Count(), 12)
	assert.Equal(t, m.capacity, 32)

	m.Clear()
	assert.Equal(t, m.Count(), 0)
	assert.Equal(t, m.capacity, 16)
}

func TestHashMapForEach(t *testing.T) {
	t.Parallel()

	m := NewHashMap[string, int]()

	m.Put("a", 1)
	m.Put("b", 2)
	m.Put("c", 3)
	m.Put("d", 4)
	m.Put("e", 5)

	assert.Equal(t, m.Count(), 5)

	var i int
	m.ForEach(func(key string, val int) {
		i++
	})

	assert.Equal(t, i, 5)

	m.Clear()
	m.ForEach(func(key string, val int) {
		t.Errorf("hashmap not cleared")
	})
}

func BenchmarkHashMapPut(b *testing.B) {
	m := NewHashMap[string, int]()

	b.ReportAllocs()
	b.ResetTimer()

	for x := 0; x < b.N; x++ {
		m.Put("foo", x)
	}
}

func BenchmarkHashMapGet(b *testing.B) {
	m := NewHashMap[string, int]()

	for x := 0; x < 1_000_000; x++ {
		m.Put(strconv.Itoa(x), x)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for x := 0; x < b.N; x++ {
		m.Get("500000")
	}
}

func BenchmarkHashMapForEach(b *testing.B) {
	m := NewHashMap[int, int]()

	for x := 0; x < 1_000_000; x++ {
		m.Put(x, x)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for x := 0; x < b.N; x++ {
		m.ForEach(func(key int, val int) {})
	}
}
