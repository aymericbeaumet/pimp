# pimp

[![test status](https://img.shields.io/github/workflow/status/aymericbeaumet/pimp/Continuous%20Integration?style=flat-square&logo=github)](https://github.com/aymericbeaumet/pimp/actions)
[![github](https://img.shields.io/github/issues/aymericbeaumet/pimp?style=flat-square&logo=github)](https://github.com/aymericbeaumet/pimp/issues)
[![go.dev](https://img.shields.io/github/v/release/aymericbeaumet/pimp?style=flat-square&logo=go&label=go.dev&logoColor=white)](https://pkg.go.dev/github.com/aymericbeaumet/pimp)

pimp is a shell-agnostic command-line expander and command runner with
pattern matching and templating capabilities that increases your
productivity.

## Table of Contents

1. [Installation](#installation)
   1. [Pre-built binaries](#pre-built-binaries)
   1. [Using the Go toolchain](#using-the-go-toolchain)
1. [Usage](#usage)
   1. [Command expander](#command-expander)
   1. [Command runner](#command-runner)
   1. [Template engine](#template-engine)
   1. [Script engine](#script-engine)
1. [Documentation](#documentation)
1. [Examples](#examples)
1. [Development](#development)
   1. [Building](#building)
   1. [Testing](#testing)

## Installation

### Pre-built binaries

Download the appropriate binary from the [latest
release](https://github.com/aymericbeaumet/pimp/releases/latest) and install it
where you see fit.

_Note: the macOS binaries are not signed yet, you will have to open them through the GUI first._

### Using the Go toolchain

```
go install github.com/aymericbeaumet/pimp@latest
```

## Usage

### Command expander

[Read more](./command-expander.md) in the documentation.

### Command runner

[Read more](./command-runner.md) in the documentation.

### Template engine

[Read more](./template-engine.md) in the documentation.

### Script engine

[Read more](./script-engine.md) in the documentation.

## Documentation

The [documentation](./docs) contains information regarding how you can get the
most out of pimp.

## Examples

The [examples](./examples) are available as references to help you start using
pimp.

## Development

### Building

The [Go toolchain](https://golang.org/doc/install) is required to work on this
project.

```
go build -o pimp .
```

### Testing

```
go test -v ./...
```
