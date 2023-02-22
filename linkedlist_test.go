package generics

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewLinkedList(t *testing.T) {
	t.Parallel()

	l := NewLinkedList[int]()

	assert.Equal(t, l.Count(), 0)
	assert.True(t, l.Empty())
}

func TestLinkedListChannel(t *testing.T) {
	t.Parallel()

	s := NewLinkedList[chan string]()

	c1 := make(chan string)
	c2 := make(chan string)

	s.PutOnBottom(c1)
	s.PutOnBottom(c2)

	assert.False(t, s.Empty())
	assert.Equal(t, s.Count(), 2)
	assert.True(t, s.Has(c1))
	assert.True(t, s.Has(c2))
	assert.Equal(t, s.Bottom(), &c2)
	s.PopBottom()
	assert.Equal(t, s.Bottom(), &c1)
	s.PopBottom()
	assert.True(t, s.Empty())
}

func TestLinkedListArray(t *testing.T) {
	t.Parallel()

	s := NewLinkedList[[3]bool]()

	a1 := [3]bool{true, false, true}
	a2 := [3]bool{false, false, true}

	s.PutOnBottom(a1)
	s.PutOnBottom(a2)

	assert.False(t, s.Empty())
	assert.Equal(t, s.Count(), 2)
	assert.True(t, s.Has(a1))
	assert.True(t, s.Has(a2))
	assert.Equal(t, s.Bottom(), &a2)
	s.PopBottom()
	assert.Equal(t, s.Bottom(), &a1)
	s.PopBottom()
	assert.True(t, s.Empty())
}

func TestLinkedListStruct(t *testing.T) {
	t.Parallel()

	type Foo struct {
		Foo string
		Bar string
	}

	s := NewLinkedList[Foo]()

	f1 := Foo{Foo: "foo", Bar: "bar"}
	f2 := Foo{Foo: "baz", Bar: "qux"}

	s.PutOnBottom(f1)
	s.PutOnBottom(f2)

	assert.False(t, s.Empty())
	assert.Equal(t, s.Count(), 2)
	assert.True(t, s.Has(f1))
	assert.True(t, s.Has(f2))
	assert.Equal(t, s.Bottom(), &f2)
	s.PopBottom()
	assert.Equal(t, s.Bottom(), &f1)
	s.PopBottom()
	assert.True(t, s.Empty())
}

func TestLinkedListLargeCapacity(t *testing.T) {
	t.Parallel()

	s := NewLinkedList[uint]()
	limit := uint(1_000_000)

	for i := uint(1); i <= limit; i++ {
		s.PutOnBottom(i)
		assert.Equal(t, s.Bottom(), &i)
		assert.Equal(t, s.Count(), int(i))
	}

	assert.Equal(t, s.Count(), int(limit))
	assert.True(t, s.Has(1))
	assert.True(t, s.Has(limit))
	assert.False(t, s.Empty())

	for i := limit; i >= 1; i-- {
		assert.Equal(t, s.Count(), int(i))
		assert.Equal(t, *s.Bottom(), i)
		s.PopBottom()
	}

	assert.True(t, s.Empty())
	assert.Equal(t, s.Count(), 0)
}

func TestLinkedListInsertFirst(t *testing.T) {
	t.Parallel()

	l := NewLinkedList[int]()

	l.PutOnTop(1)
	assert.Equal(t, *l.Top(), 1)
	assert.Equal(t, *l.Bottom(), 1)

	l.PutOnTop(2)
	assert.Equal(t, *l.Top(), 2)
	assert.Equal(t, *l.Bottom(), 1)

	l.PutOnTop(3)
	assert.Equal(t, *l.Top(), 3)
	assert.Equal(t, *l.Bottom(), 1)

	assert.Equal(t, l.Count(), 3)
	assert.False(t, l.Empty())

	l.PopTop()
	assert.Equal(t, *l.Top(), 2)
	assert.Equal(t, *l.Bottom(), 1)

	l.PopTop()
	assert.Equal(t, *l.Top(), 1)
	assert.Equal(t, *l.Bottom(), 1)

	l.PopTop()
	assert.Equal(t, l.Count(), 0)
	assert.True(t, l.Empty())
}

func TestLinkedListFailedFirst(t *testing.T) {
	t.Parallel()

	q := NewLinkedList[string]()
	assert.Nil(t, q.Top())
}

func TestLinkedListInsertLast(t *testing.T) {
	t.Parallel()

	l := NewLinkedList[int]()

	l.PutOnBottom(1)
	assert.Equal(t, *l.Top(), 1)
	assert.Equal(t, *l.Bottom(), 1)

	l.PutOnBottom(2)
	assert.Equal(t, *l.Top(), 1)
	assert.Equal(t, *l.Bottom(), 2)

	l.PutOnBottom(3)
	assert.Equal(t, *l.Top(), 1)
	assert.Equal(t, *l.Bottom(), 3)

	assert.Equal(t, l.Count(), 3)
	assert.False(t, l.Empty())

	l.PopBottom()
	assert.Equal(t, *l.Top(), 1)
	assert.Equal(t, *l.Bottom(), 2)

	l.PopBottom()
	assert.Equal(t, *l.Top(), 1)
	assert.Equal(t, *l.Bottom(), 1)

	l.PopBottom()
	assert.Equal(t, l.Count(), 0)
	assert.True(t, l.Empty())
}

func TestLinkedListFailedLast(t *testing.T) {
	t.Parallel()

	q := NewLinkedList[string]()
	assert.Nil(t, q.Bottom())
}

func TestLinkedListInsert(t *testing.T) {
	t.Parallel()

	l := NewLinkedList[int]()

	l.AddAt(1, 0)
	assert.Equal(t, *l.Top(), 1)
	assert.Equal(t, *l.Bottom(), 1)

	l.AddAt(2, 1)
	assert.Equal(t, *l.Top(), 1)
	assert.Equal(t, *l.Bottom(), 2)

	l.AddAt(3, 0)
	assert.Equal(t, *l.Top(), 3)
	assert.Equal(t, *l.Bottom(), 2)

	l.AddAt(4, 1)
	assert.Equal(t, *l.Top(), 3)
	assert.Equal(t, *l.Bottom(), 2)

	assert.Equal(t, *l.Get(0), 3)
	assert.Equal(t, *l.Get(1), 4)
	assert.Equal(t, *l.Get(2), 1)
	assert.Equal(t, *l.Get(3), 2)

	assert.Equal(t, l.Count(), 4)
	assert.False(t, l.Empty())
}

func TestLinkedListFailedInsert(t *testing.T) {
	t.Parallel()

	q := NewLinkedList[string]()
	err := q.AddAt("foo", 1)
	if err == nil {
		t.Fatal("should fail inserting, but we didn't")
	}
}

func TestLinkedListFailedGet(t *testing.T) {
	t.Parallel()

	q := NewLinkedList[string]()
	assert.Nil(t, q.Get(1))
}

func TestLinkedListUpdate(t *testing.T) {
	t.Parallel()

	l := NewLinkedList[int]()

	l.PutOnBottom(1)
	l.PutOnBottom(2)
	l.PutOnBottom(3)
	l.PutOnBottom(4)

	assert.Equal(t, *l.Get(0), 1)
	assert.Equal(t, *l.Get(1), 2)
	assert.Equal(t, *l.Get(2), 3)
	assert.Equal(t, *l.Get(3), 4)

	l.Update(2, 5)

	assert.Equal(t, *l.Get(0), 1)
	assert.Equal(t, *l.Get(1), 2)
	assert.Equal(t, *l.Get(2), 5)
	assert.Equal(t, *l.Get(3), 4)

	assert.Equal(t, l.Count(), 4)
	assert.False(t, l.Empty())
}

func TestLinkedListFailedUpdate(t *testing.T) {
	t.Parallel()

	defer func() {
		if r := recover(); r == nil {
			t.Fatal("no panic detected")
		}
	}()

	q := NewLinkedList[string]()
	q.Update(1, "foo")
}

func TestLinkedListRemove(t *testing.T) {
	t.Parallel()

	l := NewLinkedList[int]()

	l.PutOnBottom(1)
	l.PutOnBottom(2)
	l.PutOnBottom(3)
	l.PutOnBottom(4)
	l.PutOnBottom(5)

	assert.Equal(t, *l.Get(0), 1)
	assert.Equal(t, *l.Get(1), 2)
	assert.Equal(t, *l.Get(2), 3)
	assert.Equal(t, *l.Get(3), 4)
	assert.Equal(t, *l.Get(4), 5)

	assert.Equal(t, l.Count(), 5)
	assert.False(t, l.Empty())

	l.Remove(4)
	assert.Equal(t, *l.Get(0), 1)
	assert.Equal(t, *l.Get(1), 2)
	assert.Equal(t, *l.Get(2), 3)
	assert.Equal(t, *l.Get(3), 4)

	assert.Equal(t, l.Count(), 4)
	assert.False(t, l.Empty())

	l.Remove(0)
	assert.Equal(t, *l.Get(0), 2)
	assert.Equal(t, *l.Get(1), 3)
	assert.Equal(t, *l.Get(2), 4)

	assert.Equal(t, l.Count(), 3)
	assert.False(t, l.Empty())

	l.Remove(1)
	assert.Equal(t, *l.Get(0), 2)
	assert.Equal(t, *l.Get(1), 4)

	assert.Equal(t, l.Count(), 2)
	assert.False(t, l.Empty())
}

func TestLinkedListFailedRemoveOnEmptyList(t *testing.T) {
	t.Parallel()

	l := NewLinkedList[string]()
	l.Remove(1)
}

func TestLinkedListFailedRemoveOutsideOfBounds(t *testing.T) {
	t.Parallel()

	l := NewLinkedList[string]()
	l.PutOnBottom("foo")
	l.Remove(5)
}

func TestLinkedListContains(t *testing.T) {
	t.Parallel()

	l := NewLinkedList[int]()

	l.PutOnBottom(1)
	l.PutOnBottom(2)
	l.PutOnBottom(3)

	assert.False(t, l.Has(0))
	assert.True(t, l.Has(1))
	assert.True(t, l.Has(2))
	assert.True(t, l.Has(3))
	assert.False(t, l.Has(4))
}

func TestLinkedListClear(t *testing.T) {
	t.Parallel()

	l := NewLinkedList[int]()

	l.PutOnBottom(1)
	l.PutOnBottom(2)
	l.PutOnBottom(3)

	l.Clear()
	assert.Equal(t, l.Count(), 0)
	assert.True(t, l.Empty())

	l.PutOnBottom(4)
	l.PutOnBottom(5)
	l.PutOnBottom(6)

	assert.Equal(t, l.Count(), 3)
	assert.False(t, l.Empty())
}

func TestLinkedListForEach(t *testing.T) {
	t.Parallel()

	l := NewLinkedList[int]()

	l.PutOnBottom(0)
	l.PutOnBottom(1)
	l.PutOnBottom(2)
	l.PutOnBottom(3)
	l.PutOnBottom(4)

	l.ForEach(func(i int, val int) {
		assert.Equal(t, val, i)
	})

	l.Clear()
	l.ForEach(func(i int, val int) {
		t.Errorf("linked list not cleared")
	})
}

func BenchmarkLinkedListInsertAndRemove(b *testing.B) {
	l := NewLinkedList[int]()

	b.ReportAllocs()
	b.ResetTimer()

	for x := 0; x < b.N; x++ {
		l.PutOnBottom(x)
		l.PopBottom()
	}
}

func BenchmarkLinkedListForEach(b *testing.B) {
	l := NewLinkedList[int]()

	for x := 0; x < 1_000_000; x++ {
		l.PutOnBottom(x)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for x := 0; x < b.N; x++ {
		l.ForEach(func(index int, val int) {})
	}
}
