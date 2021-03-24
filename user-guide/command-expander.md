# Command Expander

When pimp is used as a command expander, it's going to try to match the command and args it is being given with the patterns you have defined in your `~/.Pimpfile`.

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

It is tedious to type `pimp git` instead of `git`, look at the [shell integration](../integrations/bash-zsh-fish.md) to make your life easier.

