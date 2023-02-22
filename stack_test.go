package generics

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewStack(t *testing.T) {
	t.Parallel()

	s := NewStack[int]()
	s.Push(1)
	s.Push(2)
	s.Push(3)

	assert.False(t, s.Empty())
	assert.Equal(t, s.Count(), 3)
	assert.Equal(t, *s.Pop(), 3)
	assert.Equal(t, s.Peek(), 2)
	assert.Equal(t, *s.Pop(), 2)
	assert.Equal(t, *s.Pop(), 1)
	assert.True(t, s.Empty())
}

func TestStackChannel(t *testing.T) {
	t.Parallel()

	s := NewStack[chan string]()

	c1 := make(chan string)
	c2 := make(chan string)

	s.Push(c1)
	s.Push(c2)

	assert.False(t, s.Empty())
	assert.Equal(t, s.Count(), 2)
	assert.True(t, s.Contains(c1))
	assert.True(t, s.Contains(c2))
	assert.Equal(t, s.Pop(), &c2)
	assert.Equal(t, s.Pop(), &c1)
	assert.True(t, s.Empty())
}

func TestStackArray(t *testing.T) {
	t.Parallel()

	s := NewStack[[3]bool]()

	a1 := [3]bool{true, false, true}
	a2 := [3]bool{false, false, true}

	s.Push(a1)
	s.Push(a2)

	assert.False(t, s.Empty())
	assert.Equal(t, s.Count(), 2)
	assert.True(t, s.Contains(a1))
	assert.True(t, s.Contains(a2))
	assert.Equal(t, s.Pop(), &a2)
	assert.Equal(t, s.Pop(), &a1)
	assert.True(t, s.Empty())
}

func TestStackStruct(t *testing.T) {
	t.Parallel()

	type Foo struct {
		Foo string
		Bar string
	}

	s := NewStack[Foo]()

	f1 := Foo{Foo: "foo", Bar: "bar"}
	f2 := Foo{Foo: "baz", Bar: "qux"}

	s.Push(f1)
	s.Push(f2)

	assert.False(t, s.Empty())
	assert.Equal(t, s.Count(), 2)
	assert.True(t, s.Contains(f1))
	assert.True(t, s.Contains(f2))
	assert.Equal(t, s.Pop(), &f2)
	assert.Equal(t, s.Pop(), &f1)
	assert.True(t, s.Empty())
}

func TestStackLargeCapacity(t *testing.T) {
	t.Parallel()

	s := NewStack[int]()
	limit := 1_000_000

	for i := 1; i <= limit; i++ {
		s.Push(i)
		assert.Equal(t, s.Peek(), i)
		assert.Equal(t, s.Count(), i)
	}

	assert.Equal(t, s.Peek(), limit)
	assert.Equal(t, s.Count(), limit)
	assert.True(t, s.Contains(1))
	assert.True(t, s.Contains(limit))
	assert.False(t, s.Empty())

	for i := limit; i >= 1; i-- {
		assert.Equal(t, s.Count(), i)
		assert.Equal(t, s.Peek(), i)
		assert.Equal(t, *s.Pop(), i)
	}

	assert.True(t, s.Empty())
	assert.Equal(t, s.Count(), 0)
}

func TestStackFailedPop(t *testing.T) {
	t.Parallel()

	s := NewStack[string]()
	assert.Nil(t, s.Pop())
}

func TestStackContains(t *testing.T) {
	t.Parallel()

	s := NewStack[string]()
	s.Push("foo")
	s.Push("bar")
	s.Push("baz")
	s.Push("qux")

	assert.True(t, s.Contains("bar"))
	assert.False(t, s.Contains("fuz"))
}

func TestStackClearing(t *testing.T) {
	t.Parallel()

	s := NewStack[int]()
	s.Push(1)
	s.Push(2)
	s.Push(3)
	assert.Equal(t, s.Count(), 3)

	s.Clear()
	assert.Equal(t, s.Count(), 0)

	s.Push(1)
	s.Push(2)
	assert.Equal(t, s.Count(), 2)
}

func TestStackForEach(t *testing.T) {
	t.Parallel()

	s := NewStack[int]()

	s.Push(1)
	s.Push(2)
	s.Push(3)
	s.Push(4)
	s.Push(5)

	i := s.Count()
	s.ForEach(func(val int) {
		assert.Equal(t, val, i)
		i--
	})

	s.Clear()
	s.ForEach(func(val int) {
		t.Errorf("stack not cleared")
	})
}

func BenchmarkStackPushAndPop(b *testing.B) {
	s := NewStack[int]()

	b.ReportAllocs()
	b.ResetTimer()

	for x := 0; x < b.N; x++ {
		s.Push(x)
		s.Pop()
	}
}

func BenchmarkStackForEach(b *testing.B) {
	s := NewStack[int]()

	for x := 0; x < 1_000_000; x++ {
		s.Push(x)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for x := 0; x < b.N; x++ {
		s.ForEach(func(val int) {})
	}
}
