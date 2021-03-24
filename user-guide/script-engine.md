# Script Engine \(PimpScript\)

pimp embeds a scripting language \(conveniently named _PimpScript_\). This language is actually just the Go template languages without delimiters \(`{{` and `}}`\), significant whitespaces nor tabs; but with added support for newlines and `;` as end of expressions. You can use `pimp --transpile` to see how a piece of PimpScript would be evaluated:

```bash
$ pimp --transpile <<EOF
range list 1 2 3
  if eq . 1; println "Winner!"
  else; printf "You are in position %d\n" .
  end
end
EOF
```

This produces the following output:

```text
{{- range list 1 2 3 -}}
{{- if eq . 1 -}}
{{- println "Winner!" -}}
{{- else -}}
{{- printf "You are in position %d\n" . -}}
{{- end -}}
{{- end -}}
```

You can execute:

```bash
pimp --exec <<EOF
range list 1 2 3
  if eq . 1
    println "Winner!"
  else
    printf "You are in position %d\n" .
  end
end
EOF
```

Which produces the following output:

```text
Winner!
You are in position 2
You are in position 3
```

### REPL

pimp comes with an REPL with history, autocompletion, etc. It's a good way to experiment short scripts:

```text
$ pimp
pimp> if eq (mul 3 3) 9; println "ok"; end
ok
```

