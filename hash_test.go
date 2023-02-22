package generics

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashBooleans(t *testing.T) {
	t.Parallel()

	assert.True(t, Hash(true) > 0)
	assert.True(t, Hash(false) > 0)
}

func TestHashNumbers(t *testing.T) {
	t.Parallel()

	assert.True(t, Hash(123) > 0)
	assert.True(t, Hash(int(123)) > 0)
	assert.True(t, Hash(int8(123)) > 0)
	assert.True(t, Hash(int16(123)) > 0)
	assert.True(t, Hash(int32(123)) > 0)
	assert.True(t, Hash(int64(123)) > 0)
	assert.True(t, Hash(uint(123)) > 0)
	assert.True(t, Hash(uint8(123)) > 0)
	assert.True(t, Hash(uint16(123)) > 0)
	assert.True(t, Hash(uint32(123)) > 0)
	assert.True(t, Hash(uint64(123)) > 0)
	assert.True(t, Hash(3.1415927) > 0)
	assert.True(t, Hash(float32(3.1415927)) > 0)
	assert.True(t, Hash(float64(3.1415927)) > 0)
	assert.True(t, Hash(complex(3, -5)) > 0)
	assert.True(t, Hash(complex64(complex(3, -5))) > 0)
	assert.True(t, Hash(complex128(complex(3, -5))) > 0)
}

func TestHashStrings(t *testing.T) {
	t.Parallel()

	assert.True(t, Hash("foo") > 0)
	assert.True(t, Hash("bar") > 0)
}

func TestHashPointers(t *testing.T) {
	t.Parallel()

	p1 := true
	p2 := false
	p3 := 246
	p4 := int(123)
	p5 := int8(123)
	p6 := int16(123)
	p7 := int32(123)
	p8 := int64(123)
	p9 := uint(123)
	p10 := uint8(123)
	p11 := uint16(123)
	p12 := uint32(123)
	p13 := uint64(123)
	p14 := "foo bar baz qux"
	p15 := uintptr(0xDEADBEEF)

	assert.True(t, Hash(&p1) > 0)
	assert.True(t, Hash(&p2) > 0)
	assert.True(t, Hash(&p3) > 0)
	assert.True(t, Hash(&p4) > 0)
	assert.True(t, Hash(&p5) > 0)
	assert.True(t, Hash(&p6) > 0)
	assert.True(t, Hash(&p7) > 0)
	assert.True(t, Hash(&p8) > 0)
	assert.True(t, Hash(&p9) > 0)
	assert.True(t, Hash(&p10) > 0)
	assert.True(t, Hash(&p11) > 0)
	assert.True(t, Hash(&p12) > 0)
	assert.True(t, Hash(&p13) > 0)
	assert.True(t, Hash(&p14) > 0)
	assert.True(t, Hash(p15) > 0)
}

type foo struct {
	foo string
	bar string
}

type baz struct {
	baz string
	qux string
}

func (h baz) Hash() uint32 {
	buf := new(bytes.Buffer)
	buf.Write([]byte(h.baz))
	buf.Write([]byte(h.qux))
	return HashBytes(buf.Bytes())
}

func TestHashHashers(t *testing.T) {
	t.Parallel()

	assert.True(t, Hash(baz{}) > 0)
	assert.True(t, Hash(baz{baz: "foo", qux: "bar"}) > 0)
	assert.True(t, Hash(baz{baz: "baz", qux: "qux"}) > 0)
}

func TestHashBytes(t *testing.T) {
	t.Parallel()

	assert.True(t, HashBytes([]byte{1, 2, 3, 4, 5}) > 0)
	assert.True(t, HashBytes([]byte{6, 7, 8, 9, 10}) > 0)
}

func BenchmarkHashInt(b *testing.B) {
	b.ReportAllocs()

	for x := 0; x < b.N; x++ {
		Hash(x)
	}
}

func BenchmarkHashIntPtr(b *testing.B) {
	b.ReportAllocs()

	for x := 0; x < b.N; x++ {
		Hash(&x)
	}
}

func BenchmarkHashUint(b *testing.B) {
	b.ReportAllocs()

	var i uint
	for x := 0; x < b.N; x++ {
		Hash(i)
		i++
	}
}

func BenchmarkHashUintPtr(b *testing.B) {
	b.ReportAllocs()

	var i uint
	for x := 0; x < b.N; x++ {
		Hash(&i)
		i++
	}
}

func BenchmarkHashStrings(b *testing.B) {
	b.ReportAllocs()

	for x := 0; x < b.N; x++ {
		Hash("foo bar baz qux")
	}
}

func BenchmarkHashStringPtr(b *testing.B) {
	b.ReportAllocs()

	str := "foo bar baz qux"
	for x := 0; x < b.N; x++ {
		Hash(&str)
	}
}

func BenchmarkHashStructs(b *testing.B) {
	b.ReportAllocs()

	for x := 0; x < b.N; x++ {
		Hash(foo{foo: "foo", bar: "bar"})
	}
}

func BenchmarkHashHashers(b *testing.B) {
	b.ReportAllocs()

	for x := 0; x < b.N; x++ {
		Hash(baz{baz: "baz", qux: "qux"})
	}
}
