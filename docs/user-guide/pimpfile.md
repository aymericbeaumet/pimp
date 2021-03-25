# Pimpfile

### YAML format

Pimpfile are an important part of pimp. These YAML files allow you to configure command expansion and tasks with a simple and expressive syntax.

The whole file will be rendered with the [Template Engine](template-engine/).

### Go format

_todo_

### Resolution order

The following resolution algorithm is being used to prioritize the order in which the Pimpfiles commands will be attempted to be matched:

1. Are we in a git repository?
   1. Yes -&gt; Sequentially open and append the Pimpfiles commands until the root of the repository
   2. No -&gt; Open and append commands from the Pimpfile in the current directory
2. Open and append the commands from the global Pimpfiles as defined in `~/.pimprc`

