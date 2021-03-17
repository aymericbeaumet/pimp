package prelude_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/aymericbeaumet/pimp/pkg/funcs/prelude"
)

type customSortableType struct {
	arr []string
}

func (t customSortableType) Len() int {
	return len(t.arr)
}

func (t customSortableType) Less(i, j int) bool {
	return strings.Compare(t.arr[i], t.arr[j]) > 0
}

func (t customSortableType) Swap(i, j int) {
	t.arr[i], t.arr[j] = t.arr[j], t.arr[i]
}

func TestSort(t *testing.T) {
	tt := []struct {
		in       interface{}
		expected interface{}
	}{
		{in: []string{"c", "b", "a"}, expected: []string{"a", "b", "c"}},
		{in: 1, expected: 1},
		{in: "foobar", expected: "foobar"},
		{in: customSortableType{[]string{"x", "y", "z"}}, expected: customSortableType{[]string{"z", "y", "x"}}},
	}

	for _, test := range tt {
		out := prelude.Sort(test.in)
		if !reflect.DeepEqual(out, test.expected) {
			t.Errorf("expected %#v for input %#v, but got %#v", test.expected, test.in, out)
		}
	}
}
