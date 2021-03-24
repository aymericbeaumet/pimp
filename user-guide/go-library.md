# Go Library

It is possible to import pimp as a standalone Go library. The 3 examples below respectively show you how to print the git branches of a repository by:

1. rendering a template
2. executing a PimpScript
3. using the native [`text/template`](https://golang.org/pkg/text/template/) package along with the pimp template functions

### 1. Render templates

```go
package main

import (
	"os"

	"github.com/aymericbeaumet/pimp"
)

func main() {
	_ = pimp.RenderTemplate(os.Stdout, `Git branches in {{pwd}}:
{{- range GitBranches}}
  - {{.}}
{{- end}}
`)
}
```

### 2. Execute PimpScript

```go
package main

import (
	"os"

	"github.com/aymericbeaumet/pimp"
)

func main() {
	_ = pimp.ExecuteScript(os.Stdout, `
    printf "Git branches in %s:\n" pwd

    range GitBranches
      printf "- %s\n" .
    end
  `)
}
```

### 3. Use the pimp template functions with \`text/template\`

```go
package main

import (
	"os"
	"text/template"

	"github.com/aymericbeaumet/pimp"
)

func main() {
	t, _ := template.New("git_branches").
		Funcs(pimp.FuncMap()).
		Parse(`Git branches in {{pwd}}:
{{- range GitBranches}}
  - {{.}}
{{- end}}
`)

	_ = t.Execute(os.Stdout, nil)
}

```

