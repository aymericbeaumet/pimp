# Command Expander

When pimp is used as a command expander, it's going to try to match the command and args it is being given with the patterns you have defined in your Pimpfiles.

### Syntax

For example in this case, when `git` is passed \(with no extra arguments\), it is going to be expanded into `git status -sb`. If some arguments are passed, then it is going to be expanded to `git <args>`.

```yaml
# ~/.Pimpfile
git     : git status -sb
git ... : git
```

```bash
$ pimp git     # equivalent to `git status -sb`
$ pimp git log # equivalent to `git log`
```

It is tedious to type `pimp git` instead of `git`, look at the [shell integration](https://github.com/aymericbeaumet/pimp/tree/060207933e60cc983a58d90dd5520e56e2c543aa/docs/integrations/bash-zsh-fish.md) to make your life easier.

