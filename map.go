package generics

import (
	"sync"
)

// Entry defines a key/value pairs.
type Entry[K comparable, V any] struct {
	Key   K
	Value V
}

// Keys creates an array of the map keys.
func Keys[K comparable, V any](source map[K]V) []K {
	result := make([]K, 0, len(source))

	for key := range source {
		result = append(result, key)
	}

	return result
}

// Values creates an array of the map values.
func Values[K comparable, V any](source map[K]V) []V {
	result := make([]V, 0, len(source))

	for _, value := range source {
		result = append(result, value)
	}

	return result
}

// FilterMap returns same map type filtered by given callback.
func FilterMap[K comparable, V any](source map[K]V, fn func(K, V) bool) map[K]V {
	result := map[K]V{}
	for key, value := range source {
		if fn(key, value) {
			result[key] = value
		}
	}
	return result
}

// FilterWhereKeys returns same map type filtered by given keys.
func FilterWhereKeys[K comparable, V any](source map[K]V, keys []K) map[K]V {
	result := map[K]V{}
	for key, value := range source {
		if Has(keys, key) {
			result[key] = value
		}
	}
	return result
}

// FilterWhereValues returns same map type filtered by given values.
func FilterWhereValues[K comparable, V comparable](source map[K]V, values []V) map[K]V {
	result := map[K]V{}
	for key, value := range source {
		if Has(values, value) {
			result[key] = value
		}
	}
	return result
}

// ToEntries transforms a map into array of key/value pairs.
func ToEntries[K comparable, V any](source map[K]V) []Entry[K, V] {
	result := make([]Entry[K, V], 0, len(source))

	for key, value := range source {
		result = append(result, Entry[K, V]{
			Key:   key,
			Value: value,
		})
	}

	return result
}

// FromEntries transforms an array of key/value pairs into a map.
func FromEntries[K comparable, V any](source []Entry[K, V]) map[K]V {
	result := map[K]V{}

	for _, value := range source {
		result[value.Key] = value.Value
	}

	return result
}

// SwapKeyValue creates a map composed of the inverted keys and values. If map contains duplicate values, of course, subsequent values will overwrite property assignments of previous values.
func SwapKeyValue[K comparable, V comparable](source map[K]V) map[V]K {
	result := map[V]K{}

	for key, value := range source {
		result[value] = key
	}

	return result
}

// Merge merges multiple maps from left to right.
func Merge[K comparable, V any](sources ...map[K]V) map[K]V {
	result := map[K]V{}

	for _, source := range sources {
		for key, value := range source {
			result[key] = value
		}
	}

	return result
}

// MapKeys manipulates a map keys and transforms it to a map of another type.
func MapKeys[K comparable, V any, R comparable](source map[K]V, fn func(V, K) R) map[R]V {
	result := map[R]V{}

	for key, value := range source {
		result[fn(value, key)] = value
	}

	return result
}

// MapValues manipulates a map values and transforms it to a map of another type.
func MapValues[K comparable, V any, R any](source map[K]V, fn func(V, K) R) map[K]R {
	result := map[K]R{}

	for key, value := range source {
		result[key] = fn(value, key)
	}

	return result
}

// ParallelMap manipulates a slice and transforms it to a slice of another type.
// `iteratee` is call in parallel. Resolve keep the same order.
func ParallelMap[T any, R any](collection []T, iteratee func(T, int) R) []R {
	result := make([]R, len(collection))

	var wg sync.WaitGroup
	wg.Add(len(collection))

	for i, item := range collection {
		go func(el T, index int) {
			defer wg.Done()
			res := iteratee(el, index)
			result[index] = res
		}(item, i)
	}

	wg.Wait()

	return result
}
