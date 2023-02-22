package timsort

import (
	"math/rand"
	"sort"
	"testing"
)

type record struct {
	key, order int
}

type records []*record

func LessThanByKey(a, b *record) bool {
	return a.key < b.key
}

type RecordSlice []record

func (s *RecordSlice) Len() int {
	return len(*s)
}

func (s *RecordSlice) Swap(i, j int) {
	(*s)[i], (*s)[j] = (*s)[j], (*s)[i]
}

func (s *RecordSlice) Less(i, j int) bool {
	return (*s)[i].key < (*s)[j].key
}

func makeVector(size int, shape string) records {
	result := make(records, size)
	switch shape {

	case "xor":
		for i := 0; i < size; i++ {
			result[i] = &record{0xff & (i ^ 0xab), i}
		}

	case "sorted":
		for i := 0; i < size; i++ {
			result[i] = &record{i, i}
		}

	case "revsorted":
		for i := 0; i < size; i++ {
			result[i] = &record{size - i, i}
		}

	case "random":
		rand.New(rand.NewSource(1))
		for i := 0; i < size; i++ {
			result[i] = &record{rand.Int(), i}
		}

	default:
		panic(shape)
	}

	return result
}

func makeRecords(size int, shape string) RecordSlice {
	result := make(RecordSlice, size)
	switch shape {

	case "xor":
		for i := 0; i < size; i++ {
			result[i] = record{0xff & (i ^ 0xab), i}
		}

	case "sorted":
		for i := 0; i < size; i++ {
			result[i] = record{i, i}
		}

	case "revsorted":
		for i := 0; i < size; i++ {
			result[i] = record{size - i, i}
		}

	case "random":
		rand.New(rand.NewSource(1))
		for i := 0; i < size; i++ {
			result[i] = record{rand.Int(), i}
		}

	default:
		panic(shape)
	}

	return result
}

func benchmarkTimsort(b *testing.B, size int, shape string) {
	b.StopTimer()

	for j := 0; j < b.N; j++ {
		v := makeVector(size, shape)

		b.StartTimer()
		Sort(v, LessThanByKey)
		b.StopTimer()
	}
}

func benchmarkTimsortInterface(b *testing.B, size int, shape string) {
	b.StopTimer()

	for j := 0; j < b.N; j++ {
		v := makeRecords(size, shape)

		b.StartTimer()
		TimSort(&v)
		b.StopTimer()
	}
}

func benchmarkStandardSort(b *testing.B, size int, shape string) {
	b.StopTimer()

	for j := 0; j < b.N; j++ {
		v := makeRecords(size, shape)

		b.StartTimer()
		sort.Stable(&v)
		b.StopTimer()
	}
}

func BenchmarkTimsortStructsXOR100(b *testing.B) {
	benchmarkTimsort(b, 100, "xor")
}

func BenchmarkTimsortInterfacesXOR100(b *testing.B) {
	benchmarkTimsortInterface(b, 100, "xor")
}

func BenchmarkStandardStructsXOR100(b *testing.B) {
	benchmarkStandardSort(b, 100, "xor")
}

func BenchmarkTimsortStructsSorted100(b *testing.B) {
	benchmarkTimsort(b, 100, "sorted")
}

func BenchmarkTimsortInterfacesSorted100(b *testing.B) {
	benchmarkTimsortInterface(b, 100, "sorted")
}

func BenchmarkStandardStructsSorted100(b *testing.B) {
	benchmarkStandardSort(b, 100, "sorted")
}

func BenchmarkTimsortStructsReverseSorted100(b *testing.B) {
	benchmarkTimsort(b, 100, "revsorted")
}

func BenchmarkTimsortInterfacesReverseSorted100(b *testing.B) {
	benchmarkTimsortInterface(b, 100, "revsorted")
}

func BenchmarkStandardStructsReverseSorted100(b *testing.B) {
	benchmarkStandardSort(b, 100, "revsorted")
}

func BenchmarkTimsortStructsRandom100(b *testing.B) {
	benchmarkTimsort(b, 100, "random")
}

func BenchmarkTimsortInterfacesRandom100(b *testing.B) {
	benchmarkTimsortInterface(b, 100, "random")
}

func BenchmarkStandardStructsRandom100(b *testing.B) {
	benchmarkStandardSort(b, 100, "random")
}

func BenchmarkTimsortStructsXOR1K(b *testing.B) {
	benchmarkTimsort(b, 1024, "xor")
}

func BenchmarkTimsortInterfacesXOR1K(b *testing.B) {
	benchmarkTimsortInterface(b, 1024, "xor")
}

func BenchmarkStandardStructsXOR1K(b *testing.B) {
	benchmarkStandardSort(b, 1024, "xor")
}

func BenchmarkTimsortStructsSorted1K(b *testing.B) {
	benchmarkTimsort(b, 1024, "sorted")
}

func BenchmarkTimsortInterfacesSorted1K(b *testing.B) {
	benchmarkTimsortInterface(b, 1024, "sorted")
}

func BenchmarkStandardStructsSorted1K(b *testing.B) {
	benchmarkStandardSort(b, 1024, "sorted")
}

func BenchmarkTimsortStructsReverseSorted1K(b *testing.B) {
	benchmarkTimsort(b, 1024, "revsorted")
}

func BenchmarkTimsortInterfacesReverseSorted1K(b *testing.B) {
	benchmarkTimsortInterface(b, 1024, "revsorted")
}

func BenchmarkStandardStructsReverseSorted1K(b *testing.B) {
	benchmarkStandardSort(b, 1024, "revsorted")
}

func BenchmarkTimsortStructsRandom1K(b *testing.B) {
	benchmarkTimsort(b, 1024, "random")
}

func BenchmarkTimsortInterfacesRandom1K(b *testing.B) {
	benchmarkTimsortInterface(b, 1024, "random")
}

func BenchmarkStandardStructsRandom1K(b *testing.B) {
	benchmarkStandardSort(b, 1024, "random")
}

func BenchmarkTimsortStructsXOR1M(b *testing.B) {
	benchmarkTimsort(b, 1024*1024, "xor")
}

func BenchmarkTimsortInterfacesXOR1M(b *testing.B) {
	benchmarkTimsortInterface(b, 1024*1024, "xor")
}

func BenchmarkStandardStructsXOR1M(b *testing.B) {
	benchmarkStandardSort(b, 1024*1024, "xor")
}

func BenchmarkTimsortStructsSorted1M(b *testing.B) {
	benchmarkTimsort(b, 1024*1024, "sorted")
}

func BenchmarkTimsortInterfacesSorted1M(b *testing.B) {
	benchmarkTimsortInterface(b, 1024*1024, "sorted")
}

func BenchmarkStandardStructsSorted1M(b *testing.B) {
	benchmarkStandardSort(b, 1024*1024, "sorted")
}

func BenchmarkTimsortStructsReverseSorted1M(b *testing.B) {
	benchmarkTimsort(b, 1024*1024, "revsorted")
}

func BenchmarkTimsortInterfacesRevSorted1M(b *testing.B) {
	benchmarkTimsortInterface(b, 1024*1024, "revsorted")
}

func BenchmarkStandardStructsReverseSorted1M(b *testing.B) {
	benchmarkStandardSort(b, 1024*1024, "revsorted")
}

func BenchmarkTimsortStructsRandom1M(b *testing.B) {
	benchmarkTimsort(b, 1024*1024, "random")
}

func BenchmarkTimsortInterfacesRandom1M(b *testing.B) {
	benchmarkTimsortInterface(b, 1024*1024, "random")
}

func BenchmarkStandardStructsRandom1M(b *testing.B) {
	benchmarkStandardSort(b, 1024*1024, "random")
}
