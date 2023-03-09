package generics

import (
	"fmt"
	"math"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type data struct {
	child string
}

func (f data) Clone() data {
	return data{child: f.child}
}

type SliceData struct {
	Title    string `json:"title"`
	Body     string `json:"body"`
	ID       int    `json:"id"`
	OwnerID  int    `json:"userId"`
	Quantity int    `json:"quantity"`
}

type List []*SliceData

func TestFirstIndexOf(t *testing.T) {
	is := assert.New(t)

	// found
	t1 := NumericFirstIndexOf([]int{0, 1, 2, 1, 2, 3}, 2)
	is.Equal(t1, 2)

	// not found
	t2 := NumericFirstIndexOf([]int{0, 1, 2, 1, 2, 3}, 6)
	is.Equal(t2, -1)

	// found
	t3 := FirstIndexOf(
		List{
			&SliceData{
				ID:      1,
				OwnerID: 1,
				Title:   "Title",
				Body:    "Body",
			},
			&SliceData{
				ID:      2,
				OwnerID: 2,
				Title:   "Another Title",
				Body:    "Another Body",
			},
		},
		SliceData{
			ID:      2,
			OwnerID: 2,
			Title:   "Another Title",
			Body:    "Another Body",
		},
	)
	is.Equal(t3, 1)

	// not found
	t4 := FirstIndexOf(
		List{
			&SliceData{
				ID:      1,
				OwnerID: 1,
				Title:   "Title",
				Body:    "Body",
			},
			&SliceData{
				ID:      2,
				OwnerID: 2,
				Title:   "Another Title",
				Body:    "Another Body",
			},
		},
		SliceData{
			ID:      1,
			OwnerID: 3,
			Title:   "Title",
			Body:    "Body",
		},
	)
	is.Equal(t4, -1)
}

func TestLastIndexOf(t *testing.T) {
	is := assert.New(t)

	// found
	t1 := NumericLastIndexOf([]int{0, 1, 2, 1, 2, 3}, 2)
	is.Equal(t1, 4)

	// not found
	t2 := NumericLastIndexOf([]int{0, 1, 2, 1, 2, 3}, 6)
	is.Equal(t2, -1)

	// found
	t4 := LastIndexOf(
		List{
			&SliceData{
				ID:      1,
				OwnerID: 1,
				Title:   "Title",
				Body:    "Body",
			},
			&SliceData{
				ID:      2,
				OwnerID: 2,
				Title:   "Another Title",
				Body:    "Another Body",
			},
			&SliceData{
				ID:      2,
				OwnerID: 2,
				Title:   "Another Title",
				Body:    "Another Body",
			},
		},
		SliceData{
			ID:      2,
			OwnerID: 2,
			Title:   "Another Title",
			Body:    "Another Body",
		},
	)
	is.Equal(t4, 2) // not element 1

	// not found
	t5 := LastIndexOf(
		List{
			&SliceData{
				ID:      1,
				OwnerID: 1,
				Title:   "Title",
				Body:    "Body",
			},
			&SliceData{
				ID:      2,
				OwnerID: 2,
				Title:   "Another Title",
				Body:    "Another Body",
			},
		},
		SliceData{
			ID:      1,
			OwnerID: 3,
			Title:   "Title",
			Body:    "Body",
		},
	)
	is.Equal(t5, -1)
}

func TestIndexes(t *testing.T) {
	is := assert.New(t)

	// found
	t1 := NumericIndexes([]int{0, 1, 2, 1, 2, 3}, 2)
	is.Equal(t1, []int{2, 4})

	// not found
	t2 := NumericIndexes([]int{0, 1, 2, 1, 2, 3}, 6)
	is.Equal(t2, []int(nil))

	// found
	t3 := Indexes(
		List{
			&SliceData{
				ID:      1,
				OwnerID: 1,
				Title:   "Title",
				Body:    "Body",
			},
			&SliceData{
				ID:      2,
				OwnerID: 2,
				Title:   "Another Title",
				Body:    "Another Body",
			},
			&SliceData{
				ID:      2,
				OwnerID: 2,
				Title:   "Another Title",
				Body:    "Another Body",
			},
		},
		SliceData{
			ID:      2,
			OwnerID: 2,
			Title:   "Another Title",
			Body:    "Another Body",
		},
	)
	is.Equal(t3, []int{1, 2}) // not element 1

	// not found
	t4 := Indexes(
		List{
			&SliceData{
				ID:      1,
				OwnerID: 1,
				Title:   "Title",
				Body:    "Body",
			},
			&SliceData{
				ID:      2,
				OwnerID: 2,
				Title:   "Another Title",
				Body:    "Another Body",
			},
		},
		SliceData{
			ID:      1,
			OwnerID: 3,
			Title:   "Title",
			Body:    "Body",
		},
	)
	is.Equal(t4, []int(nil))
}

func TestFind(t *testing.T) {
	is := assert.New(t)

	t1, ok1 := FindString([]string{"ntt", "b", "c", "d"}, func(i string) bool {
		return i == "b"
	})
	is.Equal(ok1, true)
	is.Equal(t1, "b")

	t2, ok2 := FindString([]string{"yabadabadoo"}, func(i string) bool {
		return i == "b"
	})
	is.Equal(ok2, false)
	is.Equal(t2, "")

	t3, ok3 := Find(List{
		&SliceData{
			ID:      1,
			OwnerID: 1,
			Title:   "Title",
			Body:    "Body",
		},
		&SliceData{
			ID:      2,
			OwnerID: 2,
			Title:   "Another Title",
			Body:    "Another Body",
		},
	}, func(i *SliceData) bool {
		return i.Title == "Another Title"
	})
	is.Equal(ok3, true)
	is.Equal((*t3).Body, "Another Body")
}

func TestWhereFirst(t *testing.T) {
	is := assert.New(t)

	el1, idx1, ok1 := WhereFirst([]string{"ntt", "b", "c", "d", "b"}, func(i string) bool {
		return i == "b"
	})
	is.Equal(el1, "b")
	is.Equal(ok1, true)
	is.Equal(idx1, 1)

	el2, idx2, ok2 := WhereFirst([]string{"yabadabadoo"}, func(i string) bool {
		return i == "b"
	})
	is.Equal(el2, "")
	is.Equal(ok2, false)
	is.Equal(idx2, -1)

	el3, idx3, ok3 := WhereFirst(List{
		&SliceData{
			ID:      1,
			OwnerID: 1,
			Title:   "Title",
			Body:    "Body",
		},
		&SliceData{
			ID:      2,
			OwnerID: 2,
			Title:   "Another Title",
			Body:    "Another Body",
		},
	}, func(i *SliceData) bool {
		return i.Title == "Another Title"
	})
	is.Equal(ok3, true)
	is.Equal(idx3, 1)
	is.Equal((*el3).Body, "Another Body")
}

func TestWhereLast(t *testing.T) {
	is := assert.New(t)

	el1, index1, ok1 := WhereLast([]string{"ntt", "b", "c", "d", "b"}, func(i string) bool {
		return i == "b"
	})
	is.Equal(*el1, "b")
	is.Equal(ok1, true)
	is.Equal(index1, 4)

	el2, idx2, ok2 := WhereLast([]string{"yabadabadoo"}, func(i string) bool {
		return i == "b"
	})
	is.Equal(el2, (*string)(nil))
	is.Equal(ok2, false)
	is.Equal(idx2, -1)

	el3, idx3, ok3 := WhereLast(List{
		&SliceData{
			ID:      1,
			OwnerID: 1,
			Title:   "Title",
			Body:    "Body",
		},
		&SliceData{
			ID:      2,
			OwnerID: 2,
			Title:   "Another Title",
			Body:    "Another Body",
		},
		&SliceData{
			ID:      3,
			OwnerID: 3,
			Title:   "Another Title",
			Body:    "Another Body",
		},
	}, func(i *SliceData) bool {
		return i.Title == "Another Title"
	})
	is.Equal(ok3, true)
	is.Equal(idx3, 2)
	is.Equal((*el3).Body, "Another Body")
	is.Equal((*el3).ID, 3)
}

func TestWhereElse(t *testing.T) {
	is := assert.New(t)

	t1 := WhereElse([]string{"ntt", "b", "c", "d"}, "x", func(i string) bool {
		return i == "b"
	})
	is.Equal(*t1, "b")

	t2 := WhereElse([]string{"yabadabadoo"}, "x", func(i string) bool {
		return i == "b"
	})
	is.Equal(*t2, "x")

	t3 := WhereElse(List{
		&SliceData{
			ID:      1,
			OwnerID: 1,
			Title:   "Title",
			Body:    "Body",
		},
		&SliceData{
			ID:      2,
			OwnerID: 2,
			Title:   "Another Title",
			Body:    "Another Body",
		},
		&SliceData{
			ID:      3,
			OwnerID: 3,
			Title:   "Another Title",
			Body:    "Another Body",
		},
	}, &SliceData{
		ID:      0,
		OwnerID: 0,
		Title:   "New",
		Body:    "New",
	}, func(i *SliceData) bool {
		return i.Title == "Not to be found"
	})
	is.Equal((*t3).Title, "New")
}

func TestMin(t *testing.T) {
	is := assert.New(t)

	t1 := Min([]int{1, 2, 3})
	is.Equal(t1, 1)

	t2 := Min([]int{3, 2, 1})
	is.Equal(t2, 1)

	t3 := Min([]int{})
	is.Equal(t3, 0) // well, this is not correct, but...

	t32 := Min([]int{-1, -3, 0, 1, 3})
	is.Equal(t32, -3)

	t4 := Min([]string{"unu", "doi", "trei", "patru", "aa"})
	is.Equal(t4, "aa")
}

func TestMinWhere(t *testing.T) {
	is := assert.New(t)

	t1 := MinWhere([]string{"s1", "string2", "s3"}, func(item string, min string) bool {
		return len(item) < len(min)
	})
	is.Equal(t1, "s1")

	t2 := MinWhere([]string{"string1", "string2", "s3"}, func(item string, min string) bool {
		return len(item) < len(min)
	})
	is.Equal(t2, "s3")

	t3 := MinWhere([]string{}, func(item string, min string) bool {
		return len(item) < len(min)
	})
	is.Equal(t3, "")

	t4 := MinWhere(List{
		&SliceData{
			ID:      1,
			OwnerID: 1,
			Title:   "Title",
			Body:    "Body",
		},
		&SliceData{
			ID:      2,
			OwnerID: 2,
			Title:   "Another Title",
			Body:    "Another Body",
		},
		&SliceData{
			ID:      3,
			OwnerID: 3,
			Title:   "Another Title",
			Body:    "Another Body",
		},
	}, func(d1, d2 *SliceData) bool {
		return len(d1.Title) < len(d2.Title)
	})
	is.Equal(t4.Title, "Title")
}

func TestMax(t *testing.T) {
	is := assert.New(t)

	t1 := Max([]int{1, 2, 3})
	is.Equal(t1, 3)

	t2 := Max([]int{3, 2, 1})
	is.Equal(t2, 3)

	t3 := Max([]int{})
	is.Equal(t3, 0)
}

func TestMaxWhere(t *testing.T) {
	is := assert.New(t)

	t1 := MaxWhere([]string{"s1", "string2", "s3"}, func(item string, max string) bool {
		return len(item) > len(max)
	})
	is.Equal(t1, "string2")

	t2 := MaxWhere([]string{"string1", "string2", "s3"}, func(item string, max string) bool {
		return len(item) > len(max)
	})
	is.Equal(t2, "string1")

	t3 := MaxWhere([]string{}, func(item string, max string) bool {
		return len(item) > len(max)
	})
	is.Equal(t3, "")

	t4 := MaxWhere(List{
		&SliceData{
			ID:      1,
			OwnerID: 1,
			Title:   "Title",
			Body:    "Body",
		},
		&SliceData{
			ID:      2,
			OwnerID: 2,
			Title:   "Another Longer Title",
			Body:    "Another Body",
		},
		&SliceData{
			ID:      3,
			OwnerID: 3,
			Title:   "Another Title",
			Body:    "Another Body",
		},
	}, func(d1, d2 *SliceData) bool {
		return len(d1.Title) > len(d2.Title)
	})
	is.Equal(t4.Title, "Another Longer Title")
}

func TestFilter(t *testing.T) {
	is := assert.New(t)

	t1 := Filter([]int{1, 2, 3, 4}, func(x int, _ int) bool {
		return x%2 == 0
	})
	is.Equal(t1, []int{2, 4})

	t2 := Filter([]string{"", "data", "", "dogo", ""}, func(x string, _ int) bool {
		return len(x) > 0
	})
	is.Equal(t2, []string{"data", "dogo"})

	t3 := Filter(
		List{
			&SliceData{
				ID:      1,
				OwnerID: 1,
				Title:   "Title",
				Body:    "Body",
			},
			&SliceData{
				ID:      2,
				OwnerID: 2,
				Title:   "Another Longer Title",
				Body:    "Another Body",
			},
			&SliceData{
				ID:      3,
				OwnerID: 3,
				Title:   "Another Title",
				Body:    "Another Body",
			},
		}, func(d *SliceData, i int) bool {
			return strings.HasPrefix(d.Title, "Another")
		})
	is.Equal(len(t3), 2)
	for i, el := range t3 {
		if i == 0 {
			is.Equal(el.Title, "Another Longer Title")
		} else {
			is.Equal(el.Title, "Another Title")
		}
	}
}

func TestMap(t *testing.T) {
	is := assert.New(t)

	t1 := Map([]int{1, 2, 3, 4}, func(x int, _ int) string {
		return "Hellow"
	})
	is.Equal(len(t1), 4)
	is.Equal(t1, []string{"Hellow", "Hellow", "Hellow", "Hellow"})

	t2 := Map([]int64{10, 11, 12, 13}, func(x int64, _ int) string {
		return strconv.FormatInt(x, 16)
	})
	is.Equal(len(t2), 4)
	is.Equal(t2, []string{"a", "b", "c", "d"})

	t3 := Map(List{
		&SliceData{
			ID:      1,
			OwnerID: 1,
			Title:   "Title",
			Body:    "Body",
		},
		&SliceData{
			ID:      2,
			OwnerID: 2,
			Title:   "Another Longer Title",
			Body:    "Another Body",
		},
		&SliceData{
			ID:      3,
			OwnerID: 3,
			Title:   "Another Title",
			Body:    "Another Body",
		},
	}, func(d *SliceData, i int) int {
		return 0
	})
	is.Equal(len(t3), 3)
	is.Equal(t3, []int{0, 0, 0})
}

func TestMapWhere(t *testing.T) {
	is := assert.New(t)

	t1 := MapWhere([]int64{1, 2, 3, 4, 5, 6}, func(x int64, _ int) (string, bool) {
		if x%2 == 0 {
			return strconv.FormatInt(x, 16), true
		}
		return "", false
	})
	is.Equal(len(t1), 3)
	is.Equal(t1, []string{"2", "4", "6"})

	t2 := MapWhere([]string{"papu", "pupu", "nopoo", "lotsofpoo"}, func(x string, _ int) (string, bool) {
		if strings.HasSuffix(x, "pu") {
			return "haspu", true
		}
		return "", false
	})

	is.Equal(len(t2), 2)
	is.Equal(t2, []string{"haspu", "haspu"})

	t3 := MapWhere(List{
		&SliceData{
			ID:      1,
			OwnerID: 1,
			Title:   "Title",
			Body:    "Body",
		},
		&SliceData{
			ID:      2,
			OwnerID: 2,
			Title:   "Another Longer Title",
			Body:    "Another Body",
		},
		&SliceData{
			ID:      3,
			OwnerID: 3,
			Title:   "Another Title",
			Body:    "Another Body",
		},
	}, func(d *SliceData, i int) (*SliceData, bool) {
		if strings.HasPrefix(d.Title, "Another") {
			return d, true
		}
		return nil, false
	})
	is.Equal(len(t3), 2)

	rt3 := Map(t3, func(d *SliceData, i int) string {
		return d.Title
	})
	is.Equal(rt3, []string{"Another Longer Title", "Another Title"})
}

func TestFlatMap(t *testing.T) {
	is := assert.New(t)

	t1 := FlatMap([]int{0, 1, 2, 3, 4}, func(x int, _ int) []string {
		return []string{"Goodbye"}
	})
	is.Equal(len(t1), 5)
	is.Equal(t1, []string{"Goodbye", "Goodbye", "Goodbye", "Goodbye", "Goodbye"})

	t2 := FlatMap([]int64{0, 1, 2, 3, 4}, func(x int64, _ int) []string {
		result := make([]string, 0, x)
		for i := int64(0); i < x; i++ {
			result = append(result, strconv.FormatInt(x, 16))
		}
		return result
	})
	is.Equal(len(t2), 10)
	is.Equal(t2, []string{"1", "2", "2", "3", "3", "3", "4", "4", "4", "4"})

	t3 := FlatMap(List{
		&SliceData{
			ID:      1,
			OwnerID: 1,
			Title:   "Title",
			Body:    "Body",
		},
		&SliceData{
			ID:      2,
			OwnerID: 2,
			Title:   "Another Longer Title",
			Body:    "Another Body",
		},
		&SliceData{
			ID:      3,
			OwnerID: 3,
			Title:   "Another Title",
			Body:    "Another Body",
		},
	}, func(d *SliceData, _ int) []string {
		return []string{d.Title}
	})

	is.Equal(len(t3), 3)
	is.Equal(t3, []string{"Title", "Another Longer Title", "Another Title"})
}

func TestProduce(t *testing.T) {
	is := assert.New(t)

	t1 := Produce(16, func(i int) string {
		return strconv.FormatInt(int64(i), 16)
	})
	is.Equal(len(t1), 16)
	is.Equal(t1, []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "a", "b", "c", "d", "e", "f"})

	t2 := Produce(3, func(i int) *SliceData {
		return &SliceData{Title: fmt.Sprintf("Title %d", i)}
	})
	is.Equal(len(t2), 3)
	is.Equal(t2, []*SliceData{{Title: "Title 0"}, {Title: "Title 1"}, {Title: "Title 2"}})
}

func TestReduce(t *testing.T) {
	is := assert.New(t)

	t0 := Reduce([]string{}, func(p, item string, _ int) string { return p + item }, "")
	is.Equal(t0, "")

	t1 := Reduce([]int{1, 2, 3, 4}, func(agg int, item int, _ int) int {
		return agg + item
	}, 0)
	is.Equal(t1, 10)

	t2 := Reduce([]int{1, 2, 3, 4}, func(agg int, item int, _ int) int {
		return agg + item
	}, 10)
	is.Equal(t2, 20)

	d0 := SliceData{}
	t3 := Reduce(List{
		&SliceData{
			ID:       1,
			OwnerID:  1,
			Title:    "Title",
			Body:     "Body",
			Quantity: 10,
		},
		&SliceData{
			ID:       2,
			OwnerID:  2,
			Title:    "Another Longer Title",
			Body:     "Another Body",
			Quantity: 10,
		},
		&SliceData{
			ID:       3,
			OwnerID:  3,
			Title:    "Another Title",
			Body:     "Another Body",
			Quantity: 10,
		},
	}, func(d1, d2 *SliceData, idx int) *SliceData {
		d1.Quantity += d2.Quantity
		return d1
	}, &d0)

	is.Equal(t3, &d0)
	is.Equal(t3.Quantity, 30)
}

func TestUnique(t *testing.T) {
	is := assert.New(t)

	t0 := Unique([]int{})
	is.Equal(t0, []int(nil))

	t1 := Unique([]int{1, 2, 2, 1, 1, 1, 2, 2, 2})
	is.Equal(len(t1), 2)
	is.Equal(t1, []int{1, 2})

	t2 := UniqueWhere(List{
		&SliceData{
			ID:       1,
			OwnerID:  1,
			Title:    "Title",
			Body:     "Body",
			Quantity: 10,
		},
		&SliceData{
			ID:       2,
			OwnerID:  2,
			Title:    "Title",
			Body:     "Another Body",
			Quantity: 10,
		},
		&SliceData{
			ID:       3,
			OwnerID:  3,
			Title:    "Title",
			Body:     "Another Body",
			Quantity: 10,
		},
	}, func(d *SliceData) string {
		return d.Title
	})
	is.Equal(len(t2), 1)
	is.Equal(t2, []*SliceData{{ID: 1, OwnerID: 1, Title: "Title", Body: "Body", Quantity: 10}})
}

func TestGroupWhere(t *testing.T) {
	is := assert.New(t)

	t0 := GroupWhere([]int{}, func(i int) int { return 0 })
	is.Equal(t0, map[int][]int(nil))

	t1 := GroupWhere([]int{0, 1, 2, 3, 4, 5}, func(i int) int {
		return i % 3
	})
	is.Equal(len(t1), 3)
	is.Equal(t1, map[int][]int{
		0: {0, 3},
		1: {1, 4},
		2: {2, 5},
	})

	t2 := GroupWhere(List{
		&SliceData{
			ID:       1,
			OwnerID:  1,
			Title:    "Title",
			Body:     "Body",
			Quantity: 10,
		},
		&SliceData{
			ID:       2,
			OwnerID:  2,
			Title:    "Title",
			Body:     "Another Body",
			Quantity: 10,
		},
		&SliceData{
			ID:       3,
			OwnerID:  3,
			Title:    "Title",
			Body:     "Another Body",
			Quantity: 10,
		},
		&SliceData{
			ID:       10,
			OwnerID:  10,
			Title:    "SubTitle",
			Body:     "Body",
			Quantity: 10,
		},
		&SliceData{
			ID:       20,
			OwnerID:  20,
			Title:    "SubTitle",
			Body:     "Another Body",
			Quantity: 10,
		},
		&SliceData{
			ID:       30,
			OwnerID:  30,
			Title:    "SubTitle",
			Body:     "Another Body",
			Quantity: 10,
		},
	}, func(d *SliceData) string {
		return d.Title
	})

	is.Equal(len(t2), 2)
	is.Equal(t2, map[string][]*SliceData{
		"Title": {
			&SliceData{
				ID:       1,
				OwnerID:  1,
				Title:    "Title",
				Body:     "Body",
				Quantity: 10,
			},
			&SliceData{
				ID:       2,
				OwnerID:  2,
				Title:    "Title",
				Body:     "Another Body",
				Quantity: 10,
			},
			&SliceData{
				ID:       3,
				OwnerID:  3,
				Title:    "Title",
				Body:     "Another Body",
				Quantity: 10,
			},
		},
		"SubTitle": {
			&SliceData{
				ID:       10,
				OwnerID:  10,
				Title:    "SubTitle",
				Body:     "Body",
				Quantity: 10,
			},
			&SliceData{
				ID:       20,
				OwnerID:  20,
				Title:    "SubTitle",
				Body:     "Another Body",
				Quantity: 10,
			},
			&SliceData{
				ID:       30,
				OwnerID:  30,
				Title:    "SubTitle",
				Body:     "Another Body",
				Quantity: 10,
			}},
	})
}

func TestPartition(t *testing.T) {
	is := assert.New(t)
	t0, _ := Partition([]int{}, 30)
	is.Nil(t0)

	t1, _ := Partition([]int{0, 1, 2, 3, 4, 5}, 2)
	is.Equal(t1, [][]int{{0, 1}, {2, 3}, {4, 5}})

	t2, _ := Partition([]int{0, 1, 2, 3, 4, 5, 6}, 2)
	is.Equal(t2, [][]int{{0, 1}, {2, 3}, {4, 5}, {6}})

	t4, _ := Partition([]int{0}, 2)
	is.Equal(t4, [][]int{{0}})
}

func TestPartitionBy(t *testing.T) {
	is := assert.New(t)

	oddEvenFn := func(w int) string {
		if w < 0 {
			return "negative"
		} else if w%2 == 0 {
			return "even"
		}
		return "odd"
	}

	t0 := PartitionWhere([]int{}, oddEvenFn)
	is.Nil(t0)

	t1 := PartitionWhere([]int{-2, -1, 0, 1, 2, 3, 4, 5}, oddEvenFn)
	is.Equal(t1, [][]int{{-2, -1}, {0, 2, 4}, {1, 3, 5}})
}

func TestFlatten(t *testing.T) {
	is := assert.New(t)

	t0 := Flatten([][]int{})
	is.Nil(t0)

	t1 := Flatten([][]int{{0, 1}, {2, 3, 4, 5}})
	is.Equal(t1, []int{0, 1, 2, 3, 4, 5})

	t2 := Flatten([][]*SliceData{
		{
			&SliceData{
				ID:       1,
				OwnerID:  1,
				Title:    "Title",
				Body:     "Body",
				Quantity: 10,
			},
			&SliceData{
				ID:       2,
				OwnerID:  2,
				Title:    "Title",
				Body:     "Another Body",
				Quantity: 10,
			},
			&SliceData{
				ID:       3,
				OwnerID:  3,
				Title:    "Title",
				Body:     "Another Body",
				Quantity: 10,
			},
		},
		{
			&SliceData{
				ID:       10,
				OwnerID:  10,
				Title:    "Title 10",
				Body:     "Body",
				Quantity: 10,
			},
			&SliceData{
				ID:       20,
				OwnerID:  20,
				Title:    "Title 20",
				Body:     "Another Body",
				Quantity: 10,
			},
			&SliceData{
				ID:       30,
				OwnerID:  30,
				Title:    "SubTitle 30",
				Body:     "Another Body",
				Quantity: 10,
			},
		},
	},
	)

	is.Equal(t2, []*SliceData{
		&SliceData{
			ID:       1,
			OwnerID:  1,
			Title:    "Title",
			Body:     "Body",
			Quantity: 10,
		},
		&SliceData{
			ID:       2,
			OwnerID:  2,
			Title:    "Title",
			Body:     "Another Body",
			Quantity: 10,
		},
		&SliceData{
			ID:       3,
			OwnerID:  3,
			Title:    "Title",
			Body:     "Another Body",
			Quantity: 10,
		},
		&SliceData{
			ID:       10,
			OwnerID:  10,
			Title:    "Title 10",
			Body:     "Body",
			Quantity: 10,
		},
		&SliceData{
			ID:       20,
			OwnerID:  20,
			Title:    "Title 20",
			Body:     "Another Body",
			Quantity: 10,
		},
		&SliceData{
			ID:       30,
			OwnerID:  30,
			Title:    "SubTitle 30",
			Body:     "Another Body",
			Quantity: 10,
		},
	})
}

func TestShuffle(t *testing.T) {
	is := assert.New(t)

	t1 := Shuffle([]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
	is.NotEqual(t1, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10})

	t2 := Shuffle([]int{})
	is.Nil(t2)
}

func TestReverse(t *testing.T) {
	is := assert.New(t)

	t1 := Reverse([]int{0, 1, 2, 3, 4, 5})
	is.Equal(t1, []int{5, 4, 3, 2, 1, 0})

	t2 := Reverse([]int{0, 1, 2, 3, 4, 5, 6})
	is.Equal(t2, []int{6, 5, 4, 3, 2, 1, 0})

	t3 := Reverse([]int{})
	is.Nil(t3)
}

func TestFill(t *testing.T) {
	is := assert.New(t)

	t1 := Fill([]data{{"ntt"}, {"ntt"}}, data{"b"})
	is.Equal(t1, []data{{"b"}, {"b"}})

	t2 := Fill([]data{}, data{"ntt"})
	is.Nil(t2)
}

func TestRepeat(t *testing.T) {
	is := assert.New(t)

	t1 := Repeat(2, data{"ntt"})
	is.Equal(t1, []data{{"ntt"}, {"ntt"}})

	t2 := Repeat(0, data{"ntt"})
	is.Equal(t2, []data{})
}

func TestRepeatWhere(t *testing.T) {
	is := assert.New(t)

	coolFn := func(i int) int {
		return int(math.Pow(float64(i), 2))
	}

	t1 := RepeatWhere(0, coolFn)
	is.Equal([]int{}, t1)

	t2 := RepeatWhere(2, coolFn)
	is.Equal([]int{0, 1}, t2)

	t3 := RepeatWhere(5, coolFn)
	is.Equal([]int{0, 1, 4, 9, 16}, t3)
}

func TestKeyWhere(t *testing.T) {
	is := assert.New(t)

	t1 := KeyWhere([]string{"ab", "aba", "ababa"}, func(str string) int {
		return len(str)
	})
	is.Equal(t1, map[int]string{2: "ab", 3: "aba", 5: "ababa"})

	t2 := KeyWhere(List{
		&SliceData{
			ID:       1,
			OwnerID:  1,
			Title:    "Title 20",
			Body:     "Body",
			Quantity: 0,
		},
		&SliceData{
			ID:       20,
			OwnerID:  20,
			Title:    "Title 20",
			Body:     "Another Body",
			Quantity: 10,
		},
		&SliceData{
			ID:       30,
			OwnerID:  30,
			Title:    "SubTitle 30",
			Body:     "Another Body 2",
			Quantity: 10,
		},
	}, func(d *SliceData) string {
		return d.Title
	})

	for k, el := range t2 {
		switch k {
		case "Title 20":
			is.Equal(el.ID, 20)
		case "SubTitle 30":
			is.Equal(el.ID, 30)
		}
	}
}

func TestDrop(t *testing.T) {
	is := assert.New(t)

	is.Equal([]int{1, 2, 3, 4}, Drop([]int{0, 1, 2, 3, 4}, 1))
	is.Equal([]int{2, 3, 4}, Drop([]int{0, 1, 2, 3, 4}, 2))
	is.Equal([]int{3, 4}, Drop([]int{0, 1, 2, 3, 4}, 3))
	is.Equal([]int{4}, Drop([]int{0, 1, 2, 3, 4}, 4))
	is.Equal([]int{}, Drop([]int{0, 1, 2, 3, 4}, 5))
	is.Equal([]int{}, Drop([]int{0, 1, 2, 3, 4}, 6))

	is.Equal([]*SliceData{
		&SliceData{
			ID:       3,
			OwnerID:  3,
			Title:    "Title",
			Body:     "Another Body",
			Quantity: 10,
		},
	}, Drop(List{
		&SliceData{
			ID:       1,
			OwnerID:  1,
			Title:    "Title",
			Body:     "Body",
			Quantity: 10,
		},
		&SliceData{
			ID:       2,
			OwnerID:  2,
			Title:    "Title",
			Body:     "Another Body",
			Quantity: 10,
		},
		&SliceData{
			ID:       3,
			OwnerID:  3,
			Title:    "Title",
			Body:     "Another Body",
			Quantity: 10,
		},
	}, 2))
}

func TestDropRight(t *testing.T) {
	is := assert.New(t)

	is.Equal([]int{0, 1, 2, 3}, DropRight([]int{0, 1, 2, 3, 4}, 1))
	is.Equal([]int{0, 1, 2}, DropRight([]int{0, 1, 2, 3, 4}, 2))
	is.Equal([]int{0, 1}, DropRight([]int{0, 1, 2, 3, 4}, 3))
	is.Equal([]int{0}, DropRight([]int{0, 1, 2, 3, 4}, 4))
	is.Equal([]int{}, DropRight([]int{0, 1, 2, 3, 4}, 5))
	is.Equal([]int{}, DropRight([]int{0, 1, 2, 3, 4}, 6))
}

func TestDropWhile(t *testing.T) {
	is := assert.New(t)

	is.Equal([]int{4, 5, 6}, DropWhile([]int{0, 1, 2, 3, 4, 5, 6}, func(t int) bool {
		return t != 4
	}))

	is.Equal([]int{}, DropWhile([]int{0, 1, 2, 3, 4, 5, 6}, func(t int) bool {
		return true
	}))

	is.Equal([]int{0, 1, 2, 3, 4, 5, 6}, DropWhile([]int{0, 1, 2, 3, 4, 5, 6}, func(t int) bool {
		return t == 10
	}))
}

func TestDropRightWhile(t *testing.T) {
	is := assert.New(t)

	is.Equal([]int{0, 1, 2, 3}, DropRightWhile([]int{0, 1, 2, 3, 4, 5, 6}, func(t int) bool {
		return t != 3
	}))

	is.Equal([]int{0, 1}, DropRightWhile([]int{0, 1, 2, 3, 4, 5, 6}, func(t int) bool {
		return t != 1
	}))

	is.Equal([]int{0, 1, 2, 3, 4, 5, 6}, DropRightWhile([]int{0, 1, 2, 3, 4, 5, 6}, func(t int) bool {
		return t == 10
	}))

	is.Equal([]int{}, DropRightWhile([]int{0, 1, 2, 3, 4, 5, 6}, func(t int) bool {
		return t != 10
	}))
}

func TestReject(t *testing.T) {
	is := assert.New(t)

	r1 := Reject([]int{1, 2, 3, 4}, func(x int, _ int) bool {
		return x%2 == 0
	})

	is.Equal(r1, []int{1, 3})

	r2 := Reject([]string{"Abracadabra", "data", "Badu", "dogo", "NoVision", "#thePlaceNotToBe"}, func(x string, _ int) bool {
		return len(x) > 4
	})

	is.Equal(r2, []string{"data", "Badu", "dogo"})
}

func TestCount(t *testing.T) {
	is := assert.New(t)

	t1 := Count([]int{1, 2, 1}, 1)
	is.Equal(t1, 2)

	t2 := Count([]int{1, 2, 1}, 3)
	is.Equal(t2, 0)

	t3 := Count([]int{}, 1)
	is.Equal(t3, 0)
}

func TestCountBy(t *testing.T) {
	is := assert.New(t)

	t1 := CountWhere([]int{1, 2, 1}, func(i int) bool {
		return i < 2
	})
	is.Equal(t1, 2)

	t2 := CountWhere([]int{1, 2, 1}, func(i int) bool {
		return i > 2
	})
	is.Equal(t2, 0)

	t3 := CountWhere([]int{}, func(i int) bool {
		return i <= 2
	})
	is.Equal(t3, 0)
}

func TestSubset(t *testing.T) {
	is := assert.New(t)

	in := []int{0, 1, 2, 3, 4}

	t1 := Subset(in, 0, 0)
	is.Equal([]int{}, t1)

	t2 := Subset(in, 10, 2)
	is.Equal([]int{}, t2)

	t3 := Subset(in, -10, 2)
	is.Equal([]int{0, 1}, t3)

	t4 := Subset(in, 0, 10)
	is.Equal([]int{0, 1, 2, 3, 4}, t4)

	t5 := Subset(in, 0, 2)
	is.Equal([]int{0, 1}, t5)

	t6 := Subset(in, 2, 2)
	is.Equal([]int{2, 3}, t6)

	t7 := Subset(in, 2, 5)
	is.Equal([]int{2, 3, 4}, t7)

	t8 := Subset(in, 2, 3)
	is.Equal([]int{2, 3, 4}, t8)

	t9 := Subset(in, 2, 4)
	is.Equal([]int{2, 3, 4}, t9)

	t10 := Subset(in, -2, 4)
	is.Equal([]int{3, 4}, t10)

	t11 := Subset(in, -4, 1)
	is.Equal([]int{1}, t11)

	t12 := Subset(in, -4, math.MaxUint)
	is.Equal([]int{1, 2, 3, 4}, t12)
}

func TestReplace(t *testing.T) {
	is := assert.New(t)

	in := []int{0, 1, 0, 1, 2, 3, 0}

	t1 := Replace(in, 0, 42, 2)
	is.Equal([]int{42, 1, 42, 1, 2, 3, 0}, t1)

	t2 := Replace(in, 0, 42, 1)
	is.Equal([]int{42, 1, 0, 1, 2, 3, 0}, t2)

	t3 := Replace(in, 0, 42, 0)
	is.Equal([]int{0, 1, 0, 1, 2, 3, 0}, t3)

	t4 := Replace(in, 0, 42, -1)
	is.Equal([]int{42, 1, 42, 1, 2, 3, 42}, t4)

	t5 := Replace(in, 0, 42, -1)
	is.Equal([]int{42, 1, 42, 1, 2, 3, 42}, t5)

	t6 := Replace(in, -1, 42, 2)
	is.Equal([]int{0, 1, 0, 1, 2, 3, 0}, t6)

	t7 := Replace(in, -1, 42, 1)
	is.Equal([]int{0, 1, 0, 1, 2, 3, 0}, t7)

	t8 := Replace(in, -1, 42, 0)
	is.Equal([]int{0, 1, 0, 1, 2, 3, 0}, t8)

	t9 := Replace(in, -1, 42, -1)
	is.Equal([]int{0, 1, 0, 1, 2, 3, 0}, t9)

	t10 := Replace(in, -1, 42, -1)
	is.Equal([]int{0, 1, 0, 1, 2, 3, 0}, t10)
}

func TestReplaceAll(t *testing.T) {
	is := assert.New(t)

	in := []int{0, 1, 0, 1, 2, 3, 0}

	out1 := ReplaceAll(in, 0, 42)
	out2 := ReplaceAll(in, -1, 42)

	is.Equal([]int{42, 1, 42, 1, 2, 3, 42}, out1)
	is.Equal([]int{0, 1, 0, 1, 2, 3, 0}, out2)
}

func TestHas(t *testing.T) {
	is := assert.New(t)

	t1 := Has([]int{0, 1, 2, 3, 4, 5}, 5)
	t2 := Has([]int{0, 1, 2, 3, 4, 5}, 6)

	is.Equal(t1, true)
	is.Equal(t2, false)
}

func TestHasWhere(t *testing.T) {
	is := assert.New(t)

	type ntt struct {
		B string
		A int
	}

	a1 := []ntt{{A: 1, B: "1"}, {A: 2, B: "2"}, {A: 3, B: "3"}}

	t1 := HasWhere(a1, func(t ntt) bool { return t.A == 1 && t.B == "2" })
	is.Equal(t1, false)

	t2 := HasWhere(a1, func(t ntt) bool { return t.A == 2 && t.B == "2" })
	is.Equal(t2, true)

	a2 := []string{"aaa", "bbb", "ccc"}

	t3 := HasWhere(a2, func(t string) bool { return t == "ccc" })
	is.Equal(t3, true)

	t4 := HasWhere(a2, func(t string) bool { return t == "ddd" })
	is.Equal(t4, false)
}

func TestIncluded(t *testing.T) {
	is := assert.New(t)

	t1 := Included([]int{0, 1, 2, 3, 4, 5}, []int{0, 2})
	is.True(t1)

	t2 := Included([]int{0, 1, 2, 3, 4, 5}, []int{0, 6})
	is.False(t2)

	t3 := Included([]int{0, 1, 2, 3, 4, 5}, []int{-1, 6})
	is.False(t3)

	t4 := Included([]int{0, 1, 2, 3, 4, 5}, []int{})
	is.True(t4)
}

func TestIncludedWhere(t *testing.T) {
	is := assert.New(t)

	t1 := IncludedWhere([]int{1, 2, 3, 4}, func(x int) bool {
		return x < 5
	})
	is.True(t1)

	t2 := IncludedWhere([]int{1, 2, 3, 4}, func(x int) bool {
		return x < 3
	})
	is.False(t2)

	t3 := IncludedWhere([]int{1, 2, 3, 4}, func(x int) bool {
		return x < 0
	})
	is.False(t3)

	t4 := IncludedWhere([]int{}, func(x int) bool {
		return x < 5
	})
	is.True(t4)
}

func TestIncludesOne(t *testing.T) {
	is := assert.New(t)

	t1 := IncludesOne([]int{0, 1, 2, 3, 4, 5}, []int{0, 2})
	is.True(t1)

	t2 := IncludesOne([]int{0, 1, 2, 3, 4, 5}, []int{0, 6})
	is.True(t2)

	t3 := IncludesOne([]int{0, 1, 2, 3, 4, 5}, []int{-1, 6})
	is.False(t3)

	t4 := IncludesOne([]int{0, 1, 2, 3, 4, 5}, []int{})
	is.False(t4)
}

func TestIncludesOneWhere(t *testing.T) {
	is := assert.New(t)

	t1 := IncludesOneWhere([]int{1, 2, 3, 4}, func(x int) bool {
		return x < 5
	})
	is.True(t1)

	t2 := IncludesOneWhere([]int{1, 2, 3, 4}, func(x int) bool {
		return x < 3
	})
	is.True(t2)

	t3 := IncludesOneWhere([]int{1, 2, 3, 4}, func(x int) bool {
		return x < 0
	})
	is.False(t3)

	t4 := IncludesOneWhere([]int{}, func(x int) bool {
		return x < 5
	})
	is.False(t4)
}

func TestNotIncludes(t *testing.T) {
	is := assert.New(t)

	t1 := NotIncludes([]int{0, 1, 2, 3, 4, 5}, []int{0, 2})
	is.False(t1)

	t2 := NotIncludes([]int{0, 1, 2, 3, 4, 5}, []int{0, 6})
	is.False(t2)

	t3 := NotIncludes([]int{0, 1, 2, 3, 4, 5}, []int{-1, 6})
	is.True(t3)

	t4 := NotIncludes([]int{0, 1, 2, 3, 4, 5}, []int{})
	is.True(t4)
}

func TestNotIncludesWhere(t *testing.T) {
	is := assert.New(t)

	t1 := NotIncludesWhere([]int{1, 2, 3, 4}, func(x int) bool {
		return x < 5
	})
	is.False(t1)

	t2 := NotIncludesWhere([]int{1, 2, 3, 4}, func(x int) bool {
		return x < 3
	})
	is.False(t2)

	t3 := NotIncludesWhere([]int{1, 2, 3, 4}, func(x int) bool {
		return x < 0
	})
	is.True(t3)

	t4 := NotIncludesWhere([]int{}, func(x int) bool {
		return x < 5
	})
	is.True(t4)
}

func TestCommon(t *testing.T) {
	is := assert.New(t)

	t1 := Common([]int{0, 1, 2, 3, 4, 5}, []int{0, 2})
	is.Equal(t1, []int{0, 2})

	t2 := Common([]int{0, 1, 2, 3, 4, 5}, []int{0, 6})
	is.Equal(t2, []int{0})

	t3 := Common([]int{0, 1, 2, 3, 4, 5}, []int{-1, 6})
	is.Equal(t3, []int(nil))

	t4 := Common([]int{0, 6}, []int{0, 1, 2, 3, 4, 5})
	is.Equal(t4, []int{0})

	t5 := Common([]int{0, 6, 0}, []int{0, 1, 2, 3, 4, 5})
	is.Equal(t5, []int{0})
}

func TestDiff(t *testing.T) {
	is := assert.New(t)

	first1 := []int{0, 1, 2, 3, 4, 5}
	second1 := []int{0, 2, 6}
	notInFirst1, notInSecond1 := Diff(first1, second1)
	is.Equal(notInSecond1, []int{1, 3, 4, 5})
	is.Equal(notInFirst1, []int{6})

	first2 := []int{1, 2, 3, 4, 5}
	second2 := []int{0, 6}
	notInFirst2, notInSecond2 := Diff(first2, second2)
	is.Equal(notInSecond2, []int{1, 2, 3, 4, 5})
	is.Equal(notInFirst2, []int{0, 6})

	first3 := []int{0, 1, 2, 3, 4, 5}
	second3 := []int{0, 1, 2, 3, 4, 5}
	right3, left3 := Diff(first3, second3)
	is.Nil(left3)
	is.Nil(right3)
}

func TestUnion(t *testing.T) {
	is := assert.New(t)

	t1 := Union([]int{0, 1, 2, 3, 4, 5}, []int{0, 2, 10})
	is.Equal(t1, []int{0, 1, 2, 3, 4, 5, 10})

	t2 := Union([]int{0, 1, 2, 3, 4, 5}, []int{6, 7})
	is.Equal(t2, []int{0, 1, 2, 3, 4, 5, 6, 7})

	t3 := Union([]int{0, 1, 2, 3, 4, 5}, []int{})
	is.Equal(t3, []int{0, 1, 2, 3, 4, 5})

	t4 := Union([]int{0, 1, 2}, []int{0, 1, 2})
	is.Equal(t4, []int{0, 1, 2})

	t5 := Union([]int{}, []int{})
	is.Equal(t5, []int(nil))
}

func TestParallelDo(t *testing.T) {
	is := assert.New(t)

	t1 := ParallelDo(3, func(i int) string {
		return strconv.FormatInt(int64(i), 10)
	})
	is.Equal(len(t1), 3)
	is.Equal(t1, []string{"0", "1", "2"})
}

func TestParallelGroupWhere(t *testing.T) {
	is := assert.New(t)

	t1 := ParallelGroupWhere([]int{0, 1, 2, 3, 4, 5}, func(i int) int {
		return i % 3
	})

	// sort
	for x := range t1 {
		sort.Slice(t1[x], func(i, j int) bool {
			return t1[x][i] < t1[x][j]
		})
	}

	is.EqualValues(len(t1), 3)
	is.EqualValues(t1, map[int][]int{
		0: {0, 3},
		1: {1, 4},
		2: {2, 5},
	})
}

func TestParallelPartitionWhere(t *testing.T) {
	is := assert.New(t)

	oddEven := func(x int) string {
		if x < 0 {
			return "negative"
		} else if x%2 == 0 {
			return "even"
		}
		return "odd"
	}

	result1 := ParallelPartitionWhere([]int{-2, -1, 0, 1, 2, 3, 4, 5}, oddEven)
	result2 := ParallelPartitionWhere([]int{}, oddEven)

	// order
	sort.Slice(result1, func(i, j int) bool {
		return result1[i][0] < result1[j][0]
	})
	for x := range result1 {
		sort.Slice(result1[x], func(i, j int) bool {
			return result1[x][i] < result1[x][j]
		})
	}

	is.ElementsMatch(result1, [][]int{{-2, -1}, {0, 2, 4}, {1, 3, 5}})
	is.Equal(result2, [][]int(nil))
}

func TestPermutation(t *testing.T) {
	type testCase struct {
		input []int
	}
	tests := []struct {
		name      string
		test      testCase
		want      [][]int
		expectErr bool
	}{
		{
			name:      "empty_input",
			test:      testCase{input: []int(nil)},
			want:      nil,
			expectErr: true,
		},
		{
			name:      "empty_slice",
			test:      testCase{input: []int{}},
			want:      nil,
			expectErr: true,
		},
		{
			name: "one_int",
			test: testCase{input: []int{1}},
			want: [][]int{{1}},
		},
		{
			name: "two_ints",
			test: testCase{input: []int{1, 2}},
			want: [][]int{{1, 2}, {2, 1}},
		},
		{
			name: "three_ints",
			test: testCase{input: []int{1, 2, 3}},
			want: [][]int{{1, 2, 3}, {1, 3, 2}, {2, 1, 3}, {2, 3, 1}, {3, 1, 2}, {3, 2, 1}},
		},
		{
			name: "four_ints",
			test: testCase{input: []int{1, 2, 3, 4}},
			want: [][]int{
				{1, 2, 3, 4}, {1, 2, 4, 3}, {1, 3, 2, 4}, {1, 3, 4, 2}, {1, 4, 2, 3}, {1, 4, 3, 2},
				{2, 1, 3, 4}, {2, 1, 4, 3}, {2, 3, 1, 4}, {2, 3, 4, 1}, {2, 4, 1, 3}, {2, 4, 3, 1},
				{3, 1, 2, 4}, {3, 1, 4, 2}, {3, 2, 1, 4}, {3, 2, 4, 1}, {3, 4, 1, 2}, {3, 4, 2, 1},
				{4, 1, 2, 3}, {4, 1, 3, 2}, {4, 2, 1, 3}, {4, 2, 3, 1}, {4, 3, 1, 2}, {4, 3, 2, 1},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Permutations(tt.test.input)
			if (err != nil) != tt.expectErr {
				t.Errorf("error. error is = %v, but expects error %v", err, tt.expectErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("error. got = %v, but want %v", got, tt.want)
			}
		})
	}
}
