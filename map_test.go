package generics

import (
	"sort"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKeys(t *testing.T) {
	is := assert.New(t)

	t1 := Keys(map[string]int{"data": 1, "child": 2})
	sort.Strings(t1)
	is.Equal(t1, []string{"child", "data"})
}

func TestValues(t *testing.T) {
	is := assert.New(t)

	t1 := Values(map[string]int{"data": 1, "child": 2})
	sort.Ints(t1)
	is.Equal(t1, []int{1, 2})
}

func TestFilterMap(t *testing.T) {
	is := assert.New(t)

	t1 := FilterMap(map[string]int{"data": 1, "child": 2, "bazooka": 3}, func(key string, value int) bool {
		return value%2 == 1
	})
	is.Equal(t1, map[string]int{"data": 1, "bazooka": 3})
}

func TestFilterWhereKeys(t *testing.T) {
	is := assert.New(t)

	t1 := FilterWhereKeys(map[string]int{"data": 1, "child": 2, "bazooka": 3}, []string{"data", "bazooka"})
	is.Equal(t1, map[string]int{"data": 1, "bazooka": 3})
}

func TestFilterWhereValues(t *testing.T) {
	is := assert.New(t)

	t1 := FilterWhereValues(map[string]int{"data": 1, "child": 2, "bazooka": 3}, []int{1, 3})
	is.Equal(t1, map[string]int{"data": 1, "bazooka": 3})
}

func TestToEntries(t *testing.T) {
	is := assert.New(t)

	t1 := ToEntries(map[string]int{"data": 1, "child": 2})
	sort.Slice(t1, func(i, j int) bool {
		return t1[i].Value < t1[j].Value
	})
	is.EqualValues(t1, []Entry[string, int]{
		{
			Key:   "data",
			Value: 1,
		},
		{
			Key:   "child",
			Value: 2,
		},
	})
}

func TestFromEntries(t *testing.T) {
	is := assert.New(t)

	t1 := FromEntries([]Entry[string, int]{
		{
			Key:   "data",
			Value: 1,
		},
		{
			Key:   "child",
			Value: 2,
		},
	})
	is.Len(t1, 2)
	is.Equal(t1["data"], 1)
	is.Equal(t1["child"], 2)
}

func TestSwapKeyValue(t *testing.T) {
	is := assert.New(t)

	t1 := SwapKeyValue(map[string]int{"ntt": 1, "b": 2})
	is.Len(t1, 2)
	is.EqualValues(map[int]string{1: "ntt", 2: "b"}, t1)

	t2 := SwapKeyValue(map[string]int{"ntt": 1, "b": 2, "c": 1})
	is.Len(t2, 2)
}

func TestMerge(t *testing.T) {
	is := assert.New(t)

	t1 := Merge(map[string]int{"ntt": 1, "b": 2}, map[string]int{"b": 3, "c": 4})
	is.Len(t1, 3)
	is.Equal(t1, map[string]int{"ntt": 1, "b": 3, "c": 4})
}

func TestMapKeys(t *testing.T) {
	is := assert.New(t)

	t1 := MapKeys(map[int]int{1: 1, 2: 2, 3: 3, 4: 4}, func(x int, _ int) string {
		return "Todo"
	})

	is.Equal(len(t1), 1)
	for k, _ := range t1 {
		if k != "Todo" {
			t.Errorf("key should be Todo, but it's not")
		}
	}

	t2 := MapKeys(map[int]int{1: 1, 2: 2, 3: 3, 4: 4}, func(_ int, v int) string {
		return strconv.FormatInt(int64(v), 16)
	})
	is.Equal(len(t2), 4)
	is.Equal(t2, map[string]int{"1": 1, "2": 2, "3": 3, "4": 4})
}

func TestMapValues(t *testing.T) {
	is := assert.New(t)

	t1 := MapValues(map[int]int{1: 1, 2: 2, 3: 3, 4: 4}, func(x int, _ int) string {
		return "Todo"
	})
	is.Equal(len(t1), 4)
	is.Equal(t1, map[int]string{1: "Todo", 2: "Todo", 3: "Todo", 4: "Todo"})

	t2 := MapValues(map[int]int{1: 1, 2: 2, 3: 3, 4: 4}, func(x int, _ int) string {
		return strconv.FormatInt(int64(x), 10)
	})
	is.Equal(len(t2), 4)
	is.Equal(t2, map[int]string{1: "1", 2: "2", 3: "3", 4: "4"})
}

func TestParallelMap(t *testing.T) {
	is := assert.New(t)

	t1 := ParallelMap([]int{1, 2, 3, 4}, func(x int, _ int) string {
		return "Todo"
	})
	is.Equal(len(t1), 4)
	is.Equal(t1, []string{"Todo", "Todo", "Todo", "Todo"})

	t2 := ParallelMap([]int64{1, 2, 3, 4}, func(x int64, _ int) string {
		return strconv.FormatInt(x, 16)
	})
	is.Equal(len(t2), 4)
	is.Equal(t2, []string{"1", "2", "3", "4"})
}
