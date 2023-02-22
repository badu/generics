package generics

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewQueue(t *testing.T) {
	t.Parallel()

	q := NewQueue[int]()
	q.Enqueue(1)
	q.Enqueue(2)
	q.Enqueue(3)

	assert.False(t, q.Empty())
	assert.Equal(t, q.Count(), 3)
	assert.Equal(t, *q.Dequeue(), 1)
	assert.Equal(t, q.Peek(), 2)
	assert.Equal(t, *q.Dequeue(), 2)
	assert.Equal(t, *q.Dequeue(), 3)
	assert.True(t, q.Empty())
}

func TestQueueChannel(t *testing.T) {
	t.Parallel()

	s := NewQueue[chan string]()

	c1 := make(chan string)
	c2 := make(chan string)

	s.Enqueue(c1)
	s.Enqueue(c2)

	assert.False(t, s.Empty())
	assert.Equal(t, s.Count(), 2)
	assert.True(t, s.Contains(c1))
	assert.True(t, s.Contains(c2))
	assert.Equal(t, s.Dequeue(), &c1)
	assert.Equal(t, s.Dequeue(), &c2)
	assert.True(t, s.Empty())
}

func TestQueueArray(t *testing.T) {
	t.Parallel()

	s := NewQueue[[3]bool]()

	a1 := [3]bool{true, false, true}
	a2 := [3]bool{false, false, true}

	s.Enqueue(a1)
	s.Enqueue(a2)

	assert.False(t, s.Empty())
	assert.Equal(t, s.Count(), 2)
	assert.True(t, s.Contains(a1))
	assert.True(t, s.Contains(a2))
	assert.Equal(t, s.Dequeue(), &a1)
	assert.Equal(t, s.Dequeue(), &a2)
	assert.True(t, s.Empty())
}

func TestQueueStruct(t *testing.T) {
	t.Parallel()

	type Foo struct {
		Foo string
		Bar string
	}

	s := NewQueue[Foo]()

	f1 := Foo{Foo: "foo", Bar: "bar"}
	f2 := Foo{Foo: "baz", Bar: "qux"}

	s.Enqueue(f1)
	s.Enqueue(f2)

	assert.False(t, s.Empty())
	assert.Equal(t, s.Count(), 2)
	assert.True(t, s.Contains(f1))
	assert.True(t, s.Contains(f2))
	assert.Equal(t, s.Dequeue(), &f1)
	assert.Equal(t, s.Dequeue(), &f2)
	assert.True(t, s.Empty())
}

func TestQueueLargeCapacity(t *testing.T) {
	t.Parallel()

	q := NewQueue[int]()
	limit := 1_000_000

	for i := 1; i <= limit; i++ {
		q.Enqueue(i)
		assert.Equal(t, q.Peek(), 1)
		assert.Equal(t, q.Count(), i)
	}

	assert.Equal(t, q.Peek(), 1)
	assert.Equal(t, q.Count(), limit)
	assert.True(t, q.Contains(1))
	assert.True(t, q.Contains(limit))
	assert.False(t, q.Empty())

	for i := 1; i <= limit; i++ {
		assert.Equal(t, q.Peek(), i)
		assert.Equal(t, *q.Dequeue(), i)
		assert.Equal(t, q.Count(), limit-i)
	}

	assert.True(t, q.Empty())
	assert.Equal(t, q.Count(), 0)
}

func TestQueueFailedDequeue(t *testing.T) {
	t.Parallel()

	q := NewQueue[string]()
	assert.Nil(t, q.Dequeue())
}

func TestQueueContains(t *testing.T) {
	t.Parallel()

	q := NewQueue[string]()
	q.Enqueue("foo")
	q.Enqueue("bar")
	q.Enqueue("baz")
	q.Enqueue("qux")

	assert.True(t, q.Contains("bar"))
	assert.False(t, q.Contains("fuz"))
}

func TestQueueClearing(t *testing.T) {
	t.Parallel()

	q := NewQueue[int]()
	q.Enqueue(1)
	q.Enqueue(2)
	q.Enqueue(3)
	assert.Equal(t, q.Count(), 3)

	q.Clear()
	assert.Equal(t, q.Count(), 0)

	q.Enqueue(1)
	q.Enqueue(2)
	assert.Equal(t, q.Count(), 2)
}

func TestQueueForEach(t *testing.T) {
	t.Parallel()

	q := NewQueue[int]()

	q.Enqueue(1)
	q.Enqueue(2)
	q.Enqueue(3)
	q.Enqueue(4)
	q.Enqueue(5)

	i := 1
	q.ForEach(func(val int) {
		assert.Equal(t, val, i)
		i++
	})

	q.Clear()
	q.ForEach(func(val int) {
		t.Errorf("queue not cleared")
	})
}

func BenchmarkQueueEnqueueAndDequeue(b *testing.B) {
	q := NewQueue[int]()

	b.ReportAllocs()
	b.ResetTimer()

	for x := 0; x < b.N; x++ {
		q.Enqueue(x)
		q.Dequeue()
	}
}

func BenchmarkQueueForEach(b *testing.B) {
	q := NewQueue[int]()

	for x := 0; x < 1_000_000; x++ {
		q.Enqueue(x)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for x := 0; x < b.N; x++ {
		q.ForEach(func(val int) {})
	}
}
