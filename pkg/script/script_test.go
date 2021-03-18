package script_test

import (
	"strings"
	"testing"

	"github.com/aymericbeaumet/pimp/pkg/funcs"
	"github.com/aymericbeaumet/pimp/pkg/script"
)

func TestTranspile(t *testing.T) {
	var out strings.Builder
	if err := script.Transpile(&out, `$a := 1
$b := 2

println $a $b
`, "{{", "}}"); err != nil {
		t.Error(err)
	}

	if out.String() != `{{- $a := 1 -}}
{{- $b := 2 -}}
{{- println $a $b -}}
` {
		t.Errorf("prepared script differs from expected output, got %#v", out.String())
	}
}

func TestRunFunctionCall(t *testing.T) {
	var out strings.Builder
	if err := script.Run(&out, `$a := 1
$b := 2

println $a $b
`, "{{", "}}", funcs.FuncMap()); err != nil {
		t.Error(err)
	}

	if out.String() != "1 2\n" {
		t.Errorf("script output differs from what's expected, got %#v", out.String())
	}
}
