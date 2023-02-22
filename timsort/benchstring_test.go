package timsort

import (
	"math/rand"
	"sort"
	"strconv"
	"testing"
)

func makeStrings(size int, shape string) sort.StringSlice {
	result := make(sort.StringSlice, 0, size)
	switch shape {

	case "xor":
		for i := 0; i < size; i++ {
			result = append(result, strconv.Itoa(0xff&(i^0xab)))
		}

	case "sorted":
		for i := 0; i < size; i++ {
			result = append(result, strconv.Itoa(i))
		}

	case "revsorted":
		for i := 0; i < size; i++ {
			result = append(result, strconv.Itoa(size-i))
		}

	case "random":
		rand.New(rand.NewSource(1))
		for i := 0; i < size; i++ {
			result = append(result, strconv.Itoa(rand.Int()))
		}

	default:
		panic(shape)
	}

	return result
}

func benchmarkTimsortStr(b *testing.B, size int, shape string) {
	b.StopTimer()

	for j := 0; j < b.N; j++ {
		v := makeStrings(size, shape)

		b.StartTimer()
		TimSort(v)
		b.StopTimer()
	}
}

func benchmarkStandardSortStr(b *testing.B, size int, shape string) {
	b.StopTimer()

	for j := 0; j < b.N; j++ {
		v := makeStrings(size, shape)

		b.StartTimer()
		sort.Sort(&v)
		b.StopTimer()
	}
}

func BenchmarkTimsortStringsXOR100(b *testing.B) {
	benchmarkTimsortStr(b, 100, "xor")
}

func BenchmarkStandardStringsXOR100(b *testing.B) {
	benchmarkStandardSortStr(b, 100, "xor")
}

func BenchmarkTimsortStringsSorted100(b *testing.B) {
	benchmarkTimsortStr(b, 100, "sorted")
}

func BenchmarkStandardStringsSorted100(b *testing.B) {
	benchmarkStandardSortStr(b, 100, "sorted")
}

func BenchmarkTimsortStringsReverseSorted100(b *testing.B) {
	benchmarkTimsortStr(b, 100, "revsorted")
}

func BenchmarkStandardStringsReversedSorted100(b *testing.B) {
	benchmarkStandardSortStr(b, 100, "revsorted")
}

func BenchmarkTimsortStringsRandom100(b *testing.B) {
	benchmarkTimsortStr(b, 100, "random")
}

func BenchmarkStandardStringsRandom100(b *testing.B) {
	benchmarkStandardSortStr(b, 100, "random")
}

func BenchmarkTimsortStringsXOR1K(b *testing.B) {
	benchmarkTimsortStr(b, 1024, "xor")
}

func BenchmarkStandardStringsXOR1K(b *testing.B) {
	benchmarkStandardSortStr(b, 1024, "xor")
}

func BenchmarkTimsortStringsSorted1K(b *testing.B) {
	benchmarkTimsortStr(b, 1024, "sorted")
}

func BenchmarkStandardStringsSorted1K(b *testing.B) {
	benchmarkStandardSortStr(b, 1024, "sorted")
}

func BenchmarkTimsortStringsReverseSorted1K(b *testing.B) {
	benchmarkTimsortStr(b, 1024, "revsorted")
}

func BenchmarkStandardStringsReverseSorted1K(b *testing.B) {
	benchmarkStandardSortStr(b, 1024, "revsorted")
}

func BenchmarkTimsortStringsRandom1K(b *testing.B) {
	benchmarkTimsortStr(b, 1024, "random")
}

func BenchmarkStandardStringsRandom1K(b *testing.B) {
	benchmarkStandardSortStr(b, 1024, "random")
}

func BenchmarkTimsortStringsXOR1M(b *testing.B) {
	benchmarkTimsortStr(b, 1024*1024, "xor")
}

func BenchmarkStandardStringsXOR1M(b *testing.B) {
	benchmarkStandardSortStr(b, 1024*1024, "xor")
}

func BenchmarkTimsortStringsSorted1M(b *testing.B) {
	benchmarkTimsortStr(b, 1024*1024, "sorted")
}

func BenchmarkStandardStringsSorted1M(b *testing.B) {
	benchmarkStandardSortStr(b, 1024*1024, "sorted")
}

func BenchmarkTimsortStringsReverseSorted1M(b *testing.B) {
	benchmarkTimsortStr(b, 1024*1024, "revsorted")
}

func BenchmarkStandardStringsReverseSorted1M(b *testing.B) {
	benchmarkStandardSortStr(b, 1024*1024, "revsorted")
}

func BenchmarkTimsortStringsRandom1M(b *testing.B) {
	benchmarkTimsortStr(b, 1024*1024, "random")
}

func BenchmarkStandardStringsRandom1M(b *testing.B) {
	benchmarkStandardSortStr(b, 1024*1024, "random")
}
