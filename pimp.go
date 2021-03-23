package pimp

import (
	"io"
	ttemplate "text/template"

	"github.com/aymericbeaumet/pimp/pkg/funcs"
	"github.com/aymericbeaumet/pimp/pkg/script"
	"github.com/aymericbeaumet/pimp/pkg/template"
)

// ExecuteScript executes a PimpScript and writes the output to the given
// writer. Any error from a template function immediately aborts and is
// returned.
func ExecuteScript(w io.Writer, text string) error {
	return script.Execute(w, text, "{{", "}}", FuncMap())
}

// RenderTemplate renders the template and writes the output to the given
// writer.
func RenderTemplate(w io.Writer, text string) error {
	return template.Render(w, text, "{{", "}}", FuncMap())
}

// FuncMap returns a newly allocated funcmap containing all the pimp template
// functions. This can be used with the standard Go module `text/template`.
func FuncMap() ttemplate.FuncMap {
	return funcs.FuncMap()
}
