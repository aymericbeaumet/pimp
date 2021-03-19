package template

import (
	"io"
	"strings"
	"text/template"

	"github.com/sebdah/markdown-toc/toc"
)

type afterFunc func(string) (string, error)

// afters contains functions that know how to deal with special placeholders
// that have been placed by template functions. This is necessary as some
// template functions need to perform actions but are missing context at the
// time they are executed.
var afters = map[string]afterFunc{
	// ./pkg/funcs/markdown/MarkdownTOC.go
	"\x00MarkdownTOC\x00": func(rendered string) (string, error) {
		built, err := toc.Build([]byte(rendered), "Table of Contents", 0, 1, true)
		if err != nil {
			return "", err
		}
		for i, b := range built {
			built[i] = strings.TrimPrefix(b, "   ")
		}
		return "## " + strings.Join(built[1:len(built)-1], "\n"), nil
	},
}

func Render(w io.Writer, text, ldelim, rdelim string, fm template.FuncMap) error {
	t, err := template.New("").Funcs(fm).Delims(ldelim, rdelim).Parse(text)
	if err != nil {
		return err
	}

	var sb strings.Builder
	if err := t.Execute(&sb, nil); err != nil {
		return err
	}
	out := sb.String()

	if strings.ContainsRune(out, '\x00') {
		for placeholder, afterFunc := range afters {
			if !strings.Contains(out, placeholder) {
				continue
			}
			replacement, err := afterFunc(out)
			if err != nil {
				return err
			}
			out = strings.ReplaceAll(out, placeholder, replacement)
		}
	}

	_, err = w.Write([]byte(out))
	return err
}

// RenderStrings renders several strings in a single context. This makes it
// possible to interact between several templates with variable declarations,
// etc. This could generate empty strings in the output that have to be dealt
// with.
func RenderStrings(texts []string, ldelim, rdelim string, fm template.FuncMap) ([]string, error) {
	const SEP = "\x00pimp\x00"
	text := strings.Join(texts, SEP)

	var sb strings.Builder
	if err := Render(&sb, text, ldelim, rdelim, fm); err != nil {
		return nil, err
	}

	return strings.Split(sb.String(), SEP), nil
}
