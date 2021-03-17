package script

import (
	"io"
	"strings"
	"text/template"

	ptemplate "github.com/aymericbeaumet/pimp/pkg/template"
)

func Transpile(w io.Writer, script, ldelim, rdelim string) error {
	for _, line := range strings.Split(script, "\n") {
		line = strings.TrimSpace(line)

		if len(line) == 0 {
			continue
		}

		if _, err := w.Write([]byte(ldelim)); err != nil {
			return nil
		}
		if _, err := w.Write([]byte{'-', ' '}); err != nil {
			return err
		}
		if _, err := w.Write([]byte(line)); err != nil {
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

func Run(w io.Writer, script, ldelim, rdelim string, fm template.FuncMap) error {
	var sb strings.Builder
	if err := Transpile(&sb, script, ldelim, rdelim); err != nil {
		return err
	}
	return ptemplate.Render(w, sb.String(), ldelim, rdelim, fm)
}
