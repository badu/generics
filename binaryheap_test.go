package generics

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewBinaryHeap(t *testing.T) {
	t.Parallel()

	b := NewBinaryHeap(func(a, b int) bool { return a < b })
	b.Push(6)
	b.Push(5)
	b.Push(2)
	b.Push(9)
	b.Push(4)
	b.Push(8)
	b.Push(7)
	b.Push(1)
	b.Push(3)
	b.Push(10)

	assert.False(t, b.Empty())
	assert.Equal(t, b.Count(), 10)

	for i := 1; i <= 10; i++ {
		assert.Equal(t, b.Peek(), i)
		assert.Equal(t, b.Pop(), &i)
	}

	assert.True(t, b.Empty())
	assert.Equal(t, b.Count(), 0)
}

func TestBinaryHeapStruct(t *testing.T) {
	t.Parallel()

	type Foo struct {
		Foo int
		Bar string
	}

	b := NewBinaryHeap(func(a, b Foo) bool { return a.Foo < b.Foo })

	f1 := Foo{Foo: 2, Bar: "bar"}
	f2 := Foo{Foo: 4, Bar: "qux"}
	f3 := Foo{Foo: 3, Bar: "baz"}
	f4 := Foo{Foo: 1, Bar: "foo"}

	b.Push(f1)
	b.Push(f2)
	b.Push(f3)
	b.Push(f4)

	assert.False(t, b.Empty())
	assert.Equal(t, b.Count(), 4)

	assert.True(t, b.Has(f1))
	assert.True(t, b.Has(f2))
	assert.True(t, b.Has(f3))
	assert.True(t, b.Has(f4))

	assert.Equal(t, b.Pop(), &f4)
	assert.Equal(t, b.Pop(), &f1)
	assert.Equal(t, b.Pop(), &f3)
	assert.Equal(t, b.Pop(), &f2)

	assert.True(t, b.Empty())
}

func TestBinaryHeapLargeCapacity(t *testing.T) {
	t.Parallel()

	b := NewBinaryHeap(func(a, b int) bool { return a < b })
	limit := 1_000_000

	for i := 1; i <= limit; i++ {
		b.Push(i)
		assert.Equal(t, b.Peek(), 1)
		assert.Equal(t, b.Count(), i)
	}

	assert.Equal(t, b.Peek(), 1)
	assert.Equal(t, b.Count(), limit)
	assert.True(t, b.Has(1))
	assert.True(t, b.Has(limit))
	assert.False(t, b.Empty())

	c := b.Count()
	for i := 1; i <= limit; i++ {
		assert.Equal(t, b.Count(), c)
		assert.Equal(t, b.Peek(), i)
		assert.Equal(t, b.Pop(), &i)
		c--
	}

	assert.True(t, b.Empty())
	assert.Equal(t, b.Count(), 0)
}

func TestBinaryHeapFailedExtract(t *testing.T) {
	t.Parallel()

	b := NewBinaryHeap(func(a, b int) bool { return a < b })
	assert.Nil(t, b.Pop())
}

func TestBinaryHeapContains(t *testing.T) {
	t.Parallel()

	b := NewBinaryHeap(func(a, b string) bool { return a < b })
	b.Push("foo")
	b.Push("bar")
	b.Push("baz")
	b.Push("qux")

	assert.True(t, b.Has("bar"))
	assert.False(t, b.Has("fuz"))
}

func TestBinaryHeapClearing(t *testing.T) {
	t.Parallel()

	b := NewBinaryHeap(func(a, b string) bool { return a < b })
	b.Push("foo")
	b.Push("bar")
	b.Push("baz")
	b.Push("qux")
	assert.Equal(t, b.Count(), 4)

	b.Clear()
	assert.Equal(t, b.Count(), 0)

	b.Push("foo")
	b.Push("bar")
	assert.Equal(t, b.Count(), 2)
}

func TestBinaryHeapForEach(t *testing.T) {
	t.Parallel()

	b := NewBinaryHeap(func(a, b int) bool { return a < b })
	b.Push(6)
	b.Push(5)
	b.Push(2)
	b.Push(9)
	b.Push(4)
	b.Push(8)
	b.Push(7)
	b.Push(1)
	b.Push(3)
	b.Push(10)

	assert.False(t, b.Empty())
	assert.Equal(t, b.Count(), 10)

	i := 1
	b.ForEach(func(val int) {
		assert.Equal(t, val, i)
		i++
	})

	assert.False(t, b.Empty())
	assert.Equal(t, b.Count(), 10)

	b.Clear()
	b.ForEach(func(val int) {
		t.Errorf("stack not cleared")
	})
}

func BenchmarkBinaryHeapInsert(b *testing.B) {
	h := NewBinaryHeap(func(a, b int) bool { return a < b })

	b.ReportAllocs()
	b.ResetTimer()

	for x := 0; x < b.N; x++ {
		h.Push(x)
	}
}

func BenchmarkBinaryHeapForEach(b *testing.B) {
	h := NewBinaryHeap(func(a, b int) bool { return a < b })

	for x := 0; x < 1_000_000; x++ {
		h.Push(x)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for x := 0; x < b.N; x++ {
		h.ForEach(func(val int) {})
	}
}
