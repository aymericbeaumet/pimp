---
description: >-
  pimp is a powerful command expander in itself, but it can be quite tedious to
  actually type `pimp ...` everything time you want to expand something. This is
  where the shell integrations come in.
---

# Bash, Zsh, Fish

## Bash, Fish, etc

This integration provides aliases for all the commands defined in your global Pimpfile. Add this to your shell configuration file:

```text
eval $(pimp --shell)
```

## Zsh

The Zsh integration is more evolved. It support aliases, but also brings a more advanced completion mecanism system. Add this to your `~/.zshrc`:

```text
eval $(pimp --zsh)
```

