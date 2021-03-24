# Execute scripts

pimp scripts follow the template syntax, but without the `{{` and `}}` delimiters.

```go
{{- /* script.pimp */ -}}

printf "Git branches in %s:\n" pwd
range GitBranches
  printf "* %s\n" .
end
```

Read more about the [Script Engine](../user-guide-1/script-engine.md) in the documentation.

Let's finally see the next steps you can follow to learn more about pimp.

