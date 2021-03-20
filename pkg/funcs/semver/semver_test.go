package semver_test

import (
	"testing"

	"github.com/aymericbeaumet/pimp/pkg/funcs/semver"
)

func TestEmpty(t *testing.T) {
	tt := []func(input interface{}) (*semver.Version, error){
		semver.Major,
		semver.Minor,
		semver.Patch,
	}

	for _, fn := range tt {
		_, err := fn("")
		if err == nil {
			t.Error("expecting error but got nil")
		}
	}
}

func TestBump(t *testing.T) {
	const input = "0.0.0"
	tt := map[string]func(input interface{}) (*semver.Version, error){
		"1.0.0": semver.Major,
		"0.1.0": semver.Minor,
		"0.0.1": semver.Patch,
	}

	for expected, fn := range tt {
		out, err := fn(input)
		if err != nil {
			t.Error(err)
		}
		if s := out.String(); s != expected {
			t.Errorf("expected %#v for input %#v, but got %#v", expected, input, s)
		}
	}
}

func TestBumpKeepVPrefix(t *testing.T) {
	const input = "v0.0.0"
	tt := map[string]func(input interface{}) (*semver.Version, error){
		"v1.0.0": semver.Major,
		"v0.1.0": semver.Minor,
		"v0.0.1": semver.Patch,
	}

	for expected, fn := range tt {
		out, err := fn(input)
		if err != nil {
			t.Error(err)
		}
		if s := out.String(); s != expected {
			t.Errorf("expected %#v for input %#v, but got %#v", expected, input, s)
		}
	}
}

func TestBumpKeepArbitraryPrefix(t *testing.T) {
	const input = "foobar 0.0.0"
	tt := map[string]func(input interface{}) (*semver.Version, error){
		"foobar 1.0.0": semver.Major,
		"foobar 0.1.0": semver.Minor,
		"foobar 0.0.1": semver.Patch,
	}

	for expected, fn := range tt {
		out, err := fn(input)
		if err != nil {
			t.Error(err)
		}
		if s := out.String(); s != expected {
			t.Errorf("expected %#v for input %#v, but got %#v", expected, input, s)
		}
	}
}

func TestUnsortedStringSlice(t *testing.T) {
	input := []string{
		"1.0.0",
		"3.0.0",
		"0.0.0",
		"2.0.0",
	}
	tt := map[string]func(input interface{}) (*semver.Version, error){
		"4.0.0": semver.Major,
		"3.1.0": semver.Minor,
		"3.0.1": semver.Patch,
	}

	for expected, fn := range tt {
		out, err := fn(input)
		if err != nil {
			t.Error(err)
		}
		if s := out.String(); s != expected {
			t.Errorf("expected %#v for input %#v, but got %#v", expected, input, s)
		}
	}
}
