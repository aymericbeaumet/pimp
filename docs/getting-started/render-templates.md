# Render templates

pimp template system is based on [Go template system](https://golang.org/pkg/text/template/). There are no fundamental differences between.

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

Read more about the [Template Engine](../user-guide-1/template-engine.md) in the documentation.

Next, let's see how you can run scripts with Pimp.

