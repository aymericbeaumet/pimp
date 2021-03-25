# Render templates

pimp template system is based on [Go template system](https://golang.org/pkg/text/template/). There are no fundamental differences between them.

```go
{{- /* template.tmpl */ -}}

Git branches in {{pwd}}:
{{- range GitBranches}}
  * {{.}}
{{- end}}
```

{% hint style="info" %}
You can use `pimp --eval 'funcs | toPrettyJSON'` to list all the functions available when rendering a template.
{% endhint %}

Read more about the [Template Engine](../user-guide/template-engine/) in the documentation.

Next, let's see how you can run scripts with Pimp.

