package generics_test

import (
	"testing"

	"github.com/badu/generics"
)

func TestSortedMapExample(t *testing.T) {
	var m generics.SortedMap[string, int]
	m.Set("doi", -12)
	m.Set("trei", 3)
	m.Set("doi", 2)
	m.Set("unu", 1)
	t.Log(&m) // {doi:2,trei:3,unu:1}
}
