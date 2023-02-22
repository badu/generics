package timsort_test

import (
	"fmt"
	"sort"

	"github.com/badu/generics/timsort"
)

type Record struct {
	ssn  int
	name string
}

func BySsn(a, b Record) bool {
	return a.ssn < b.ssn
}

func ByName(a, b Record) bool {
	return a.name < b.name
}

func Example() { // example 1
	ExampleSort()
	ExampleLessFunc()
	// Output: sorted array: [a b c]
	// sorted by ssn: [{101765430 sue} {123456789 joe} {345623452 mary}]
	// sorted by name: [{123456789 joe} {345623452 mary} {101765430 sue}]
}

func ExampleSort() {
	l := []string{"c", "a", "b"}
	timsort.TimSort(sort.StringSlice(l))
	fmt.Printf("sorted array: %+v\n", l)

}

func ExampleLessFunc() { // example 2

	db := make([]Record, 3)
	db[0] = Record{123456789, "joe"}
	db[1] = Record{101765430, "sue"}
	db[2] = Record{345623452, "mary"}

	// sorts array by ssn (ascending)
	timsort.Sort(db, BySsn)
	fmt.Printf("sorted by ssn: %v\n", db)

	// now re-sort same array by name (ascending)
	timsort.Sort(db, ByName)
	fmt.Printf("sorted by name: %v\n", db)
}
