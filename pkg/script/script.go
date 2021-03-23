package script

import (
	"io"
	"strings"
	"text/template"

	ptemplate "github.com/aymericbeaumet/pimp/pkg/template"
)

func Execute(w io.Writer, script, ldelim, rdelim string, fm template.FuncMap) error {
	var sb strings.Builder
	if err := Transpile(&sb, script, ldelim, rdelim); err != nil {
		return err
	}
	return ptemplate.Render(w, sb.String(), ldelim, rdelim, fm)
}

func Transpile(w io.Writer, script, ldelim, rdelim string) error {
	expressions, err := intoExpressions(script)
	if err != nil {
		return err
	}

	for _, expression := range expressions {
		if _, err := w.Write([]byte(ldelim)); err != nil {
			return nil
		}
		if _, err := w.Write([]byte{'-', ' '}); err != nil {
			return err
		}
		if _, err := w.Write([]byte(expression)); err != nil {
			return err
		}
		if _, err := w.Write([]byte{' ', '-'}); err != nil {
			return err
		}
		if _, err := w.Write([]byte(rdelim)); err != nil {
			return err
		}
		if _, err := w.Write([]byte{'\n'}); err != nil {
			return err
		}
	}

	return nil
}

func intoExpressions(input string) ([]string, error) {
	var out []string
	var sb strings.Builder
	var isString, isEscaped bool

	for _, r := range input {
		if !isString && (r == ';' || r == '\n') {
			out = appendNonEmpty(out, sb.String())
			sb.Reset()
			continue
		}
		sb.WriteRune(r)

		if r == '\'' || r == '"' || r == '`' {
			if !isString {
				isString = true
			} else if !isEscaped || r == '`' {
				isString = false
			}
		}

		isEscaped = isString && r == '\\'
	}

	out = appendNonEmpty(out, sb.String())

	return out, nil
}

func appendNonEmpty(slice []string, element string) []string {
	if trimmed := strings.TrimSpace(element); len(trimmed) > 0 {
		return append(slice, trimmed)
	}
	return slice
}
