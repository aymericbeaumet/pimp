package template

import (
	"io"
	"strings"
	"text/template"
)

func Render(w io.Writer, text, ldelim, rdelim string, fm template.FuncMap) error {
	t, err := template.New("").Funcs(fm).Delims(ldelim, rdelim).Parse(text)
	if err != nil {
		return err
	}
	return t.Execute(w, nil)
}

func RenderString(text, ldelim, rdelim string, fm template.FuncMap) (string, error) {
	var sb strings.Builder
	err := Render(&sb, text, ldelim, rdelim, fm)
	return sb.String(), err
}

// RenderStrings renders several strings in a single context. This makes it
// possible to interact between several templates with variable declarations,
// etc. This could generate empty strings in the output that have to be dealt
// with.
func RenderStrings(texts []string, ldelim, rdelim string, fm template.FuncMap) ([]string, error) {
	const SEP = "\x00pimp\x00"

	text := strings.Join(texts, SEP)
	rendered, err := RenderString(text, ldelim, rdelim, fm)
	if err != nil {
		return nil, err
	}

	return strings.Split(rendered, SEP), nil
}
