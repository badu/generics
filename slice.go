package generics

import (
	"fmt"
	"math/rand"
	"sync"

	"golang.org/x/exp/constraints"
)

// NumericFirstIndexOf returns the index at which the first occurrence of a value is found in a slice
// or return -1 if the value cannot be found.
func NumericFirstIndexOf[T comparable](source []T, el T) int {
	for i, item := range source {
		if item == el {
			return i
		}
	}
	return -1
}

// FirstIndexOf returns the index at which the first occurrence of a value is found in a slice
// or return -1 if the value cannot be found.
func FirstIndexOf[T comparable](source []*T, el T) int {
	for i, item := range source {
		if *item == el {
			return i
		}
	}
	return -1
}

// NumericLastIndexOf returns the index at which the last occurrence of a value is found in a slice
// or return -1 if the value cannot be found.
func NumericLastIndexOf[T comparable](source []T, el T) int {
	for i := len(source) - 1; i >= 0; i-- {
		if source[i] == el {
			return i
		}
	}
	return -1
}

// LastIndexOf returns the index at which the last occurrence of a value is found in a slice
// or return -1 if the value cannot be found.
func LastIndexOf[T comparable](source []*T, el T) int {
	for i := len(source) - 1; i >= 0; i-- {
		if *source[i] == el {
			return i
		}
	}
	return -1
}

// NumericIndexes returns the indexes at which the last occurrence of a value is found in a slice
func NumericIndexes[T comparable](source []T, el T) []int {
	var result []int
	for i, item := range source {
		if item == el {
			result = append(result, i)
		}
	}
	return result
}

// Indexes returns the indexes at which the last occurrence of a value is found in a slice
func Indexes[T comparable](source []*T, el T) []int {
	var result []int
	for i, item := range source {
		if *item == el {
			result = append(result, i)
		}
	}
	return result
}

// FindString search an element in a slice based on a callback. It returns element and true if element was found.
func FindString[T comparable](source []T, fn func(T) bool) (T, bool) {
	for _, item := range source {
		if fn(item) {
			return item, true
		}
	}
	return *new(T), false
}

// Find search an element in a slice based on a callback. It returns element and true if element was found.
func Find[T comparable](source []*T, fn func(*T) bool) (*T, bool) {
	for _, item := range source {
		if fn(item) {
			return item, true
		}
	}
	return nil, false
}

// WhereFirst searches an element in a slice based on a callback and returns the index and true.
// It returns -1 and false if the element is not found.
func WhereFirst[T any](source []T, fn func(T) bool) (T, int, bool) {
	for i, item := range source {
		if fn(item) {
			return item, i, true
		}
	}
	return *new(T), -1, false
}

// WhereLast searches last element in a slice based on a callback and returns the index and true.
// It returns -1 and false if the element is not found.
func WhereLast[T any](source []T, fn func(T) bool) (*T, int, bool) {
	for i := len(source) - 1; i >= 0; i-- {
		if fn(source[i]) {
			return &source[i], i, true
		}
	}
	return nil, -1, false
}

// WhereElse search an element in a slice based on a callback. It returns the element if found or a given fallback value otherwise.
func WhereElse[T any](source []T, fallback T, fn func(T) bool) *T {
	for _, item := range source {
		if fn(item) {
			return &item
		}
	}

	return &fallback
}

// Min search the minimum value of a slice.
func Min[T constraints.Ordered](source []T) T {
	if len(source) == 0 {
		return *new(T)
	}

	min := source[0]
	for i := 1; i < len(source); i++ {
		if source[i] < min {
			min = source[i]
		}
	}

	return min
}

// MinWhere search the minimum value of a slice using the given comparison function.
// If several values of the slice are equal to the smallest value, returns the first such value.
func MinWhere[T any](source []T, fn func(T, T) bool) T {
	if len(source) == 0 {
		return *new(T)
	}

	min := source[0]
	for i := 1; i < len(source); i++ {
		if fn(source[i], min) {
			min = source[i]
		}
	}

	return min
}

// Max searches the maximum value of a slice.
func Max[T constraints.Ordered](source []T) T {
	if len(source) == 0 {
		return *new(T)
	}

	max := source[0]
	for i := 1; i < len(source); i++ {
		if source[i] > max {
			max = source[i]
		}
	}

	return max
}

// MaxWhere search the maximum value of a slice using the given comparison function.
// If several values of the slice are equal to the greatest value, returns the first such value.
func MaxWhere[T any](source []T, fn func(T, T) bool) T {
	if len(source) == 0 {
		return *new(T)
	}

	max := source[0]
	for i := 1; i < len(source); i++ {
		if fn(source[i], max) {
			max = source[i]
		}
	}

	return max
}

// Filter iterates over elements of slice, returning a slice of all elements callback returns truthy for.
func Filter[V any](source []V, fn func(V, int) bool) []V {
	var result []V

	for i, elem := range source {
		if fn(elem, i) {
			result = append(result, elem)
		}
	}

	return result
}

// Map manipulates a slice and transforms it to a slice of another type.
func Map[T, R any](source []T, fn func(T, int) R) []R {
	result := make([]R, len(source))

	for i, item := range source {
		result[i] = fn(item, i)
	}

	return result
}

// MapWhere returns a slice which obtained after both filtering and mapping using the given callback function.
func MapWhere[T, R any](source []T, fn func(T, int) (R, bool)) []R {
	var result []R

	for i, item := range source {
		if r, ok := fn(item, i); ok {
			result = append(result, r)
		}
	}

	return result
}

// FlatMap manipulates a slice and transforms and flattens it.
func FlatMap[T, R any](source []T, fn func(T, int) []R) []R {
	var result []R

	for i, item := range source {
		result = append(result, fn(item, i)...)
	}

	return result
}

// Produce invokes the callback n times, returning a slice of the results of each invocation.
func Produce[T any](count int, fn func(int) T) []T {
	result := make([]T, count)
	for i := 0; i < count; i++ {
		result[i] = fn(i) // it's like decorating
	}
	return result
}

// Reduce reduces slice to a value which is the accumulated result of running each element in slice
// through a function, where each successive invocation is supplied the return value of the previous.
func Reduce[T comparable, R any](source []T, fn func(R, T, int) R, result R) R {
	if len(source) == 0 {
		return *new(R)
	}

	for i, item := range source {
		result = fn(result, item, i)
	}
	return result
}

// Unique returns a duplicate-free version of a slice, in which only the first occurrence of each element is kept, while preserving the order in which they occur in the slice.
func Unique[T comparable](source []T) []T {
	if len(source) == 0 {
		return nil
	}

	result := make([]T, 0, len(source))
	seen := make(map[T]struct{}, len(source))

	for _, item := range source {
		if _, ok := seen[item]; ok {
			continue
		}

		seen[item] = struct{}{}
		result = append(result, item)
	}

	return result
}

// UniqueWhere returns a duplicate-free version of a slice, in which only the FIRST occurrence of each element is kept, while preserving the order in which they occur in the slice.
func UniqueWhere[T any, U comparable](source []T, fn func(T) U) []T {
	if len(source) == 0 {
		return nil
	}

	result := make([]T, 0, len(source))
	seen := make(map[U]struct{}, len(source))

	for _, item := range source {
		key := fn(item)

		if _, ok := seen[key]; ok {
			continue
		}

		seen[key] = struct{}{}
		result = append(result, item)
	}

	return result
}

// GroupWhere returns an object composed of keys generated from the results of running each element of slice.
func GroupWhere[T any, U comparable](source []T, fn func(T) U) map[U][]T {
	if len(source) == 0 {
		return nil
	}

	result := map[U][]T{}
	for _, item := range source {
		key := fn(item)
		result[key] = append(result[key], item)
	}

	return result
}

// Partition returns a slice of elements split into groups the length of size. If slice can't be split evenly,
// the final chunk will be the remaining elements.
func Partition[T any](source []T, size int) ([][]T, error) {
	if len(source) == 0 {
		return nil, nil
	}

	if size <= 0 {
		return nil, fmt.Errorf("second parameter must be greater than 0")
	}

	result := make([][]T, 0, len(source)/2+1)
	for i := 0; i < len(source); i++ {
		if i%size == 0 {
			result = append(result, make([]T, 0, size))
		}
		result[i/size] = append(result[i/size], source[i])
	}

	return result, nil
}

// PartitionWhere returns a slice of elements split into groups. The order of grouped values is
// determined by the order they occur in slice. The grouping is generated from the results
// of running each element of slice through callback.
func PartitionWhere[T any, K comparable](source []T, fn func(x T) K) [][]T {
	if len(source) == 0 {
		return nil
	}

	var result [][]T
	seen := map[K]int{}

	for _, item := range source {
		key := fn(item)

		idx, ok := seen[key]
		if !ok {
			idx = len(result)
			seen[key] = idx
			result = append(result, []T{})
		}

		result[idx] = append(result[idx], item)
	}

	return result
}

// Flatten returns a slice a single level deep.
func Flatten[T any](source [][]T) []T {
	if len(source) == 0 {
		return nil
	}

	var result []T

	for _, item := range source {
		result = append(result, item...)
	}

	return result
}

// Shuffle returns a slice of shuffled values. Uses the Fisher-Yates shuffle algorithm.
func Shuffle[T any](source []T) []T {
	if len(source) == 0 {
		return nil
	}

	rand.Shuffle(len(source), func(i, j int) {
		source[i], source[j] = source[j], source[i]
	})

	return source
}

// Reverse reverses slice so that the first element becomes the last, the second element becomes the second to last, and so on.
func Reverse[T any](source []T) []T {
	if len(source) == 0 {
		return nil
	}

	for i := 0; i < len(source)/2; i = i + 1 {
		j := len(source) - 1 - i
		source[i], source[j] = source[j], source[i]
	}

	return source
}

// Clonable defines a constraint of types having Clone() T method.
type Clonable[T any] interface {
	Clone() T
}

// Fill fills elements of slice with `defaults` values.
func Fill[T Clonable[T]](source []T, defaults T) []T {
	if len(source) == 0 {
		return nil
	}

	result := make([]T, 0, len(source))

	for range source {
		result = append(result, defaults.Clone())
	}

	return result
}

// Repeat builds a slice with N copies of initial value.
func Repeat[T Clonable[T]](count int, initial T) []T {
	result := make([]T, 0, count)

	for i := 0; i < count; i++ {
		result = append(result, initial.Clone())
	}

	return result
}

// RepeatWhere builds a slice with values returned by N calls of callback.
func RepeatWhere[T any](count int, fn func(int) T) []T {

	result := make([]T, 0, count)

	for i := 0; i < count; i++ {
		result = append(result, fn(i))
	}

	return result
}

// KeyWhere transforms a slice or a slice of structs to a map based on a pivot callback.
func KeyWhere[K comparable, V any](source []V, fn func(V) K) map[K]V {
	if len(source) == 0 {
		return nil
	}

	result := make(map[K]V, len(source))

	for _, v := range source {
		k := fn(v)
		result[k] = v
	}

	return result
}

// Drop drops n elements from the beginning of a slice or slice.
func Drop[T any](source []T, n int) []T {
	if len(source) == 0 {
		return nil
	}

	if len(source) <= n {
		return make([]T, 0)
	}

	result := make([]T, len(source)-n)
	for i := n; i < len(source); i++ {
		result[i-n] = source[i]
	}

	return result
}

// DropWhile drops elements from the beginning of a slice or slice while the callback returns true.
func DropWhile[T any](source []T, fn func(T) bool) []T {
	if len(source) == 0 {
		return nil
	}

	i := 0
	for ; i < len(source); i++ {
		if !fn(source[i]) {
			break
		}
	}

	result := make([]T, len(source)-i)

	for j := 0; i < len(source); i, j = i+1, j+1 {
		result[j] = source[i]
	}

	return result
}

// DropRight drops n elements from the end of a slice or slice.
func DropRight[T any](source []T, n int) []T {
	if len(source) == 0 {
		return nil
	}

	if len(source) <= n {
		return make([]T, 0)
	}

	result := make([]T, len(source)-n)
	for i := len(source) - 1 - n; i >= 0; i-- {
		result[i] = source[i]
	}

	return result
}

// DropRightWhile drops elements from the end of a slice or slice while the callback returns true.
func DropRightWhile[T any](source []T, fn func(T) bool) []T {
	if len(source) == 0 {
		return nil
	}

	i := len(source) - 1
	for ; i >= 0; i-- {
		if !fn(source[i]) {
			break
		}
	}

	result := make([]T, i+1)

	for ; i >= 0; i-- {
		result[i] = source[i]
	}

	return result
}

// Reject is the opposite of Filter, this method returns the elements of slice that callback does not return truthy for.
func Reject[V any](source []V, fn func(V, int) bool) []V {
	if len(source) == 0 {
		return nil
	}

	var result []V

	for i, item := range source {
		if !fn(item, i) {
			result = append(result, item)
		}
	}

	return result
}

// Count counts the number of elements in the slice that compare equal to value.
func Count[T comparable](source []T, value T) (count int) {
	if len(source) == 0 {
		return 0
	}

	for _, item := range source {
		if item == value {
			count++
		}
	}

	return count
}

// CountWhere counts the number of elements in the slice for which callback is true.
func CountWhere[T any](source []T, fn func(T) bool) (count int) {
	if len(source) == 0 {
		return 0
	}

	for _, item := range source {
		if fn(item) {
			count++
		}
	}

	return count
}

// Subset return part of a slice.
func Subset[T any](source []T, offset int, length uint) []T {
	if len(source) == 0 {
		return nil
	}

	if offset < 0 {
		offset = len(source) + offset
		if offset < 0 {
			offset = 0
		}
	}

	if offset > len(source) {
		return []T{}
	}

	if length > uint(len(source))-uint(offset) {
		length = uint(len(source) - offset)
	}

	return source[offset : offset+int(length)]
}

// Replace returns a copy of the slice with the first n non-overlapping instances of old replaced by new.
func Replace[T comparable](source []T, old T, new T, n int) []T {
	if len(source) == 0 {
		return nil
	}

	size := len(source)
	result := make([]T, 0, size)

	for _, item := range source {
		if item == old && n != 0 {
			result = append(result, new)
			n--
		} else {
			result = append(result, item)
		}
	}

	return result
}

// ReplaceAll returns a copy of the slice with all non-overlapping instances of old replaced by new.
func ReplaceAll[T comparable](source []T, old T, new T) []T {
	return Replace(source, old, new, -1)
}

// Has returns true if an element is present in a slice.
func Has[T comparable](slice []T, element T) bool {
	for _, item := range slice {
		if item == element {
			return true
		}
	}

	return false
}

// HasWhere returns true if callback function return true.
func HasWhere[T any](slice []T, fn func(T) bool) bool {
	for _, item := range slice {
		if fn(item) {
			return true
		}
	}

	return false
}

// Included returns true if all elements of a subset are contained into a slice or if the subset is empty.
func Included[T comparable](slice []T, other []T) bool {
	for _, elem := range other {
		if !Has(slice, elem) {
			return false
		}
	}

	return true
}

// IncludedWhere returns true if the callback returns true for all of the elements in the slice or if the slice is empty.
func IncludedWhere[V any](slice []V, fn func(V) bool) bool {
	for _, v := range slice {
		if !fn(v) {
			return false
		}
	}

	return true
}

// IncludesOne returns true if at least 1 element of a subset is contained into a slice.
func IncludesOne[T comparable](slice []T, other []T) bool {
	for _, elem := range other {
		if Has(slice, elem) {
			return true
		}
	}

	return false
}

// IncludesOneWhere returns true if the callback returns true for any of the elements in the slice.
func IncludesOneWhere[V any](slice []V, fn func(V) bool) bool {
	for _, v := range slice {
		if fn(v) {
			return true
		}
	}

	return false
}

// NotIncludes returns true if no element of a subset are contained into a slice or if the subset is empty.
func NotIncludes[V comparable](slice []V, subset []V) bool {
	for _, elem := range subset {
		if Has(slice, elem) {
			return false
		}
	}

	return true
}

// NotIncludesWhere returns true if the callback returns true for none of the elements in the slice or if the slice is empty.
func NotIncludesWhere[V any](slice []V, fn func(V) bool) bool {
	for _, v := range slice {
		if fn(v) {
			return false
		}
	}

	return true
}

// Common returns the intersection between two slices.
func Common[T comparable](source1 []T, source2 []T) []T {
	seen := map[T]struct{}{}
	for _, elem := range source1 {
		seen[elem] = struct{}{}
	}

	var result []T
	for _, elem := range source2 {
		if _, ok := seen[elem]; ok {
			result = append(result, elem)
		}
	}

	return result
}

// Diff returns the difference between two slices.
// The first return is the slice of elements absent of source1.
// The second return is the slice of elements absent of source2.
func Diff[T comparable](source1 []T, source2 []T) ([]T, []T) {
	seenFirst := map[T]struct{}{}
	for _, elem := range source1 {
		seenFirst[elem] = struct{}{}
	}

	seenSecond := map[T]struct{}{}
	for _, elem := range source2 {
		seenSecond[elem] = struct{}{}
	}

	var notSeenFirst, notSeenSecond []T
	for _, elem := range source1 {
		if _, ok := seenSecond[elem]; !ok {
			notSeenSecond = append(notSeenSecond, elem)
		}
	}

	for _, elem := range source2 {
		if _, ok := seenFirst[elem]; !ok {
			notSeenFirst = append(notSeenFirst, elem)
		}
	}

	return notSeenFirst, notSeenSecond
}

// Union returns all distinct elements from both slices. Resolve DOES NOT change the order of elements relatively.
func Union[T comparable](source1, source2 []T) []T {
	seen := map[T]struct{}{}
	for _, e := range source1 {
		seen[e] = struct{}{}
	}
	for _, e := range source2 {
		seen[e] = struct{}{}
	}

	hasAdd := map[T]struct{}{}
	var result []T
	for _, e := range source1 {
		if _, ok := seen[e]; ok {
			result = append(result, e)
			hasAdd[e] = struct{}{}
		}
	}

	for _, e := range source2 {
		if _, ok := hasAdd[e]; ok {
			continue
		}
		if _, ok := seen[e]; ok {
			result = append(result, e)
		}
	}

	return result
}

// ParallelForEach iterates over elements of slice and invokes callback for each element.
func ParallelForEach[T any](slice []T, fn func(T, int)) {
	var wg sync.WaitGroup
	wg.Add(len(slice))

	for i, item := range slice {
		go func(el T, index int) {
			fn(el, index)
			wg.Done()
		}(item, i)
	}

	wg.Wait()
}

// ParallelDo invokes the callback n times, returning a slice of the results of each invocation.
func ParallelDo[T any](count int, fn func(int) T) []T {
	result := make([]T, count)

	var wg sync.WaitGroup
	wg.Add(count)
	for i := 0; i < count; i++ {
		go func(index int) {
			defer wg.Done()
			item := fn(index)
			result[index] = item
		}(i)
	}
	wg.Wait()

	return result
}

// ParallelGroupWhere returns an object composed of keys generated from the results of running each element of slice through callback.
func ParallelGroupWhere[T any, U comparable](slice []T, fn func(T) U) map[U][]T {
	result := map[U][]T{}

	var (
		mu sync.Mutex
		wg sync.WaitGroup
	)
	wg.Add(len(slice))
	for _, item := range slice {
		go func(el T) {
			defer func() {
				mu.Unlock()
				wg.Done()
			}()

			key := fn(el)
			mu.Lock()
			result[key] = append(result[key], el)
		}(item)
	}
	wg.Wait()

	return result
}

// ParallelPartitionWhere returns a slice of elements split into groups. The order of grouped values is
// determined by the order they occur in slice. The grouping is generated from the results
// of running each element of slice through callback.
// `callback` is call in parallel.
func ParallelPartitionWhere[T any, K comparable](slice []T, fn func(x T) K) [][]T {
	var (
		result [][]T
		mu     sync.Mutex
		wg     sync.WaitGroup
	)

	wg.Add(len(slice))
	seen := map[K]int{}
	for _, item := range slice {
		go func(el T) {
			defer func() {
				mu.Unlock()
				wg.Done()
			}()

			key := fn(el)
			mu.Lock()
			resultIndex, ok := seen[key]
			if !ok {
				resultIndex = len(result)
				seen[key] = resultIndex
				result = append(result, []T{})
			}
			result[resultIndex] = append(result[resultIndex], el)
		}(item)
	}
	wg.Wait()

	return result
}
