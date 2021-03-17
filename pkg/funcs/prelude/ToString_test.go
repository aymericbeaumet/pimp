package prelude_test

import (
	"fmt"
	"net/url"
	"testing"

	"github.com/aymericbeaumet/pimp/pkg/funcs/prelude"
)

func TestToString(t *testing.T) {
	tt := []struct {
		in       interface{}
		expected string
	}{
		{in: []string{"1", "2", "3"}, expected: "1\n2\n3"},
		{in: []fmt.Stringer{
			&url.URL{Scheme: "https", Host: "github.com"},
			&url.URL{Scheme: "https", Host: "gitlab.com"},
		}, expected: "https://github.com\nhttps://gitlab.com"},
		{in: 1, expected: "1"},
		{in: []int{1, 2, 3}, expected: "1\n2\n3"},
		{in: []interface{}{1, 2, 3}, expected: "1\n2\n3"},
		{in: "foobar", expected: "foobar"},
		{in: &url.URL{Scheme: "https", Host: "bitbucket.com"}, expected: "https://bitbucket.com"},
		{in: 1, expected: "1"},
	}

	for _, test := range tt {
		out := prelude.ToString(test.in)
		if out != test.expected {
			t.Errorf("expected %#v for input %#v, but got %#v", test.expected, test.in, out)
		}
	}
}
