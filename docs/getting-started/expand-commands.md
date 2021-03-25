# Expand commands

pimp looks for commands to expand in a [Pimpfile](../user-guide/pimpfile.md). Pimpfiles are YAML files where you can define patterns to match, and what they should be matched to.

Let's say you want to show the status of your git repository whenever "git" _alone_ is matched, this is what you would put in your Pimpfile:

```yaml
# ~/.Pimpfile
git : git status -sb
```

You can test this with `pimp git`, this will execute `git status -sb`

Now, what if you also want to add another command? Let's say `git co` _alone_ should offer you a list of branches you can checkout to:

```yaml
# ~/.Pimpfile
git    : git status -sb
git co : git checkout {{GitBranches | fzf}}
```

Great, now what if you want to make `git co` point to `git checkout`? You can use the `...` operator for that, that catches variadic arguments, and pass them along to the end of your expanded command:

```yaml
# ~/.Pimpfile
git        : git status -sb
git co     : git checkout {{GitBranches | fzf}}
git co ... : git checkout
```

`pimp git co readme.md` will actually execute `git checkout readme.md` whereas `git co` will execute `git checkout <branch>`, with &lt;branch&gt; being the result of the `{{GitBranches | fzf}}` template being rendered.

Congratulations, you have expanded your first commands! Read more about the [Command Expander](../user-guide/command-expander.md) to explore the full potential of command expansion with Pimp. The [Shell Integration](https://github.com/aymericbeaumet/pimp/tree/060207933e60cc983a58d90dd5520e56e2c543aa/docs/integrations/bash-zsh-fish.md) section is also helpful to understand how to integrate pimp within your shell.

Next, let's see how we can use Pimp as a task runner.

