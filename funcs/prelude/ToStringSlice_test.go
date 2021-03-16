package prelude_test

import (
	"reflect"
	"testing"

	"github.com/aymericbeaumet/pimp/funcs/prelude"
)

func TestToStringSlice(t *testing.T) {
	tt := []struct {
		in       interface{}
		expected []string
	}{
		{in: "1\n2\n3", expected: []string{"1", "2", "3"}},
		{in: 1, expected: []string{"1"}},
		{in: "foobar", expected: []string{"foobar"}},
	}

	for _, test := range tt {
		out := prelude.ToStringSlice(test.in)
		if !reflect.DeepEqual(out, test.expected) {
			t.Errorf("expected %#v for input %#v, but got %#v", test.expected, test.in, out)
		}
	}
}
