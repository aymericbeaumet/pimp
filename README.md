# pimp

[![test status](https://img.shields.io/github/workflow/status/aymericbeaumet/pimp/Continuous%20Integration?style=flat-square&logo=github)](https://github.com/aymericbeaumet/pimp/actions) [![github](https://img.shields.io/github/issues/aymericbeaumet/pimp?style=flat-square&logo=github)](https://github.com/aymericbeaumet/pimp/issues) [![go.dev](https://img.shields.io/github/v/release/aymericbeaumet/pimp?style=flat-square&logo=go&label=go.dev&logoColor=white)](https://pkg.go.dev/github.com/aymericbeaumet/pimp)

pimp is a shell-agnostic command expander and task runner with pattern matching, templating and scripting capabilities that increases your productivity.

## Installation

### Pre-built binaries

Download the appropriate binary from the [latest release](https://github.com/aymericbeaumet/pimp/releases/latest) and install it where you see fit.

_Note: the macOS binaries are not signed yet, you will have to open them through the GUI first._

### Via the Go toolchain

If you have the [Go toolchain](https://golang.org/doc/install) on your machine, you can install the latest pimp release with the following command:

```text
go install github.com/aymericbeaumet/pimp@latest
```

## Usage

### Pimpfile

Pimpfile are an important part of pimp. These YAML files allow you to configure command expansion and tasks with a simple and expressive syntax.

Read more about [Pimpfiles](docs/pimpfile.md) in the documentation.

### Command Expander

When pimp is used as a command expander, it's going to try to match the command and args it is being given with the patterns you have defined in your `~/.Pimpfile`.

For example in this case, when `git` is passed \(with no extra arguments\), it is going to be expanded into `git status -sb`. If some arguments are passed, then it is going to be expanded to `git <args>`.

```text
$ cat ~/.Pimpfile
git: git status -sb
git ...: git
```

```text
$ pimp git # equivalent to `git status -sb`
$ pimp git log # equivalent to `git log`
```

Read more about [command expansion](docs/command-expander.md) in the documentation.

### Task Runner

Following the same concept as command expansion \(see above\), you can also leverage pimp to behave as a task runner for your project. The Pimpfile in your local directory always has the highest priority.

For example this is how you would do if you wanted to define a

```text
$ cat ./Pimpfile
test: go test ./...
$ pimp test
```

Read more about [running tasks](docs/task-runner.md) in the documentation.

### Template Engine

```text
$ pimp --render template.tmpl
```

Read more about how to use pimp as a stand-alone [template engine](docs/template-engine.md) in the documentation.

### Script Engine \(PimpScript\)

```text
$ pimp --run script.pimp
```

Read more about how to use pimp as a [script engine](docs/script-engine.md) in the documentation.

### Go Library

Read more about how to import pimp as a [Go library](docs/go-library.md) in the documentation.

## Documentation

The [documentation](docs/) contains information regarding how you can get the most out of pimp.

## Examples

The [examples](https://github.com/aymericbeaumet/pimp/tree/89f6c6236b5942f787c9035ec6b814c3c2a14c28/examples/README.md) are available as references to help you start using pimp.

## Development

The [Go toolchain](https://golang.org/doc/install) is required to work on this project.

### Building

```text
go build -o pimp .
```

### Installing

```text
go install .
```

### Testing

```text
go test -v ./...
```

### Releasing

Pimp tasks are defined in the [./Pimpfile](https://github.com/aymericbeaumet/pimp/tree/89f6c6236b5942f787c9035ec6b814c3c2a14c28/Pimpfile/README.md) to release new versions. The release process is entirely automated and is being taken care of by the CI.

```text
pimp major
pimp minor
pimp patch
```

