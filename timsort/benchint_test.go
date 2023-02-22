package timsort

import (
	"math/rand"
	"sort"
	"testing"
)

type ints []int

func (p *ints) Len() int           { return len(*p) }
func (p *ints) Less(i, j int) bool { return (*p)[i] < (*p)[j] }
func (p *ints) Swap(i, j int)      { (*p)[i], (*p)[j] = (*p)[j], (*p)[i] }

func LessThanInt(a, b int) bool {
	return a < b
}

func makeInts(size int, kind string) ints {
	result := make(ints, 0, size)
	switch kind {

	case "xor":
		for i := 0; i < size; i++ {
			result = append(result, 0xff&(i^0xab))
		}

	case "sorted":
		for i := 0; i < size; i++ {
			result = append(result, i)
		}

	case "revsorted":
		for i := 0; i < size; i++ {
			result = append(result, size-i)
		}

	case "random":
		rand.New(rand.NewSource(1))
		for i := 0; i < size; i++ {
			result = append(result, rand.Int())
		}

	default:
		panic(kind)
	}

	return result
}

func benchmarkTimsortI(b *testing.B, size int, shape string) {
	b.StopTimer()

	for j := 0; j < b.N; j++ {
		v := makeInts(size, shape)

		b.StartTimer()
		Ints(v, LessThanInt)
		b.StopTimer()
	}
}

func benchmarkStandardSortI(b *testing.B, size int, shape string) {
	b.StopTimer()

	for j := 0; j < b.N; j++ {
		v := makeInts(size, shape)

		b.StartTimer()
		sort.Sort(&v)
		b.StopTimer()
	}
}

func BenchmarkTimsortIntsXOR100(b *testing.B) {
	benchmarkTimsortI(b, 100, "xor")
}

func BenchmarkStandardIntsXOR100(b *testing.B) {
	benchmarkStandardSortI(b, 100, "xor")
}

func BenchmarkTimsortIntsSorted100(b *testing.B) {
	benchmarkTimsortI(b, 100, "sorted")
}

func BenchmarkStandardIntsSorted100(b *testing.B) {
	benchmarkStandardSortI(b, 100, "sorted")
}

func BenchmarkTimsortIntsReverseSorted100(b *testing.B) {
	benchmarkTimsortI(b, 100, "revsorted")
}

func BenchmarkStandardIntsReverseSorted100(b *testing.B) {
	benchmarkStandardSortI(b, 100, "revsorted")
}

func BenchmarkTimsortIntsRandom100(b *testing.B) {
	benchmarkTimsortI(b, 100, "random")
}

func BenchmarkStandardIntsRandom100(b *testing.B) {
	benchmarkStandardSortI(b, 100, "random")
}

func BenchmarkTimsortIntsXOR1K(b *testing.B) {
	benchmarkTimsortI(b, 1024, "xor")
}

func BenchmarkStandardIntsXOR1K(b *testing.B) {
	benchmarkStandardSortI(b, 1024, "xor")
}

func BenchmarkTimsortIntsSorted1K(b *testing.B) {
	benchmarkTimsortI(b, 1024, "sorted")
}

func BenchmarkStandardIntsSorted1K(b *testing.B) {
	benchmarkStandardSortI(b, 1024, "sorted")
}

func BenchmarkTimsortIntsReverseSorted1K(b *testing.B) {
	benchmarkTimsortI(b, 1024, "revsorted")
}

func BenchmarkStandardIntsReverseSorted1K(b *testing.B) {
	benchmarkStandardSortI(b, 1024, "revsorted")
}

func BenchmarkTimsortIntsRandom1K(b *testing.B) {
	benchmarkTimsortI(b, 1024, "random")
}

func BenchmarkStandardIntsRandom1K(b *testing.B) {
	benchmarkStandardSortI(b, 1024, "random")
}

func BenchmarkTimsortIntsXOR1M(b *testing.B) {
	benchmarkTimsortI(b, 1024*1024, "xor")
}

func BenchmarkStandardIntsXOR1M(b *testing.B) {
	benchmarkStandardSortI(b, 1024*1024, "xor")
}

func BenchmarkTimsortIntsSorted1M(b *testing.B) {
	benchmarkTimsortI(b, 1024*1024, "sorted")
}

func BenchmarkStandardIntsSorted1M(b *testing.B) {
	benchmarkStandardSortI(b, 1024*1024, "sorted")
}

func BenchmarkTimsortIntsReverseSorted1M(b *testing.B) {
	benchmarkTimsortI(b, 1024*1024, "revsorted")
}

func BenchmarkStandardIntsReverseSorted1M(b *testing.B) {
	benchmarkStandardSortI(b, 1024*1024, "revsorted")
}

func BenchmarkTimsortIntsRandom1M(b *testing.B) {
	benchmarkTimsortI(b, 1024*1024, "random")
}

func BenchmarkStandardIntsRandom1M(b *testing.B) {
	benchmarkStandardSortI(b, 1024*1024, "random")
}
