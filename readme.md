# pimp

[![test status](https://img.shields.io/github/workflow/status/aymericbeaumet/pimp/Continuous%20Integration?style=flat-square&logo=github)](https://github.com/aymericbeaumet/pimp/actions)
[![github](https://img.shields.io/github/issues/aymericbeaumet/pimp?style=flat-square&logo=github)](https://github.com/aymericbeaumet/pimp/issues)
[![Go.Dev reference](https://img.shields.io/badge/go.dev-reference-blue?style=flat-square&logo=go&logoColor=white)](https://pkg.go.dev/github.com/aymericbeaumet/pimp)

## Development (local)

```
go run . git co
go run . --dry-run git co
```

## Development (system-wide)

```
go install .

# add to ~/.zshrc
eval "$(pimp --zsh)"

git co
pimp --dry-run git co
```
