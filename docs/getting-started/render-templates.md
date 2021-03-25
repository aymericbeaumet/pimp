# Render templates

pimp template system is based on [Go template system](https://golang.org/pkg/text/template/). There are no fundamental differences between them. The following piece of code behaves as you would expect:

```go
{{- /* template.tmpl */ -}}

Git branches in {{pwd}}:
{{- range GitBranches}}
  * {{.}}
{{- end}}
```

Note the use of the `GitBranches` function. Refer to the [template functions](../user-guide/template-engine/functions.md) documentation to learn more about all the functions available.

{% hint style="info" %}
You can use `pimp --eval 'funcs | toPrettyJSON'` to list all the functions available when rendering a template.
{% endhint %}

Read more about the [Template Engine](../user-guide/template-engine/) in the documentation.

Next, let's see how you can run scripts with Pimp.

