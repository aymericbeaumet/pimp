# pimp

[![test status](https://img.shields.io/github/workflow/status/aymericbeaumet/pimp/Continuous%20Integration?style=flat-square&logo=github)](https://github.com/aymericbeaumet/pimp/actions)
[![github](https://img.shields.io/github/issues/aymericbeaumet/pimp?style=flat-square&logo=github)](https://github.com/aymericbeaumet/pimp/issues)
[![go.dev](https://img.shields.io/github/v/release/aymericbeaumet/pimp?style=flat-square&logo=go&label=go.dev&logoColor=white)](https://pkg.go.dev/github.com/aymericbeaumet/pimp)

pimp is a shell-agnostic command-line expander and command runner with
pattern matching and templating capabilities that increases your
productivity.

## Table of Contents

1. [Install](#install)
   1. [Pre-built binaries](#pre-built-binaries)
   1. [Using the Go toolchain](#using-the-go-toolchain)
1. [Testing](#testing)

## Install

### Pre-built binaries

Download the binary from the [latest
releases](https://github.com/aymericbeaumet/pimp/releases/latest)
and install it where you see fit.

### Using the Go toolchain

```
go install github.com/aymericbeaumet/pimp@latest
```

## Testing

```
go test -v ./...
```
