# pimp

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
