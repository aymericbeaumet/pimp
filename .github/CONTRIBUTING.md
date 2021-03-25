# Contributing

The [Go toolchain](https://golang.org/doc/install) is required to work on this
project.

[golangci-lint](https://golangci-lint.run/) is used to lint the project.

The [GoReleaser](https://goreleaser.com/) GitHub
[action](https://github.com/goreleaser/goreleaser-action) is used to
automatically create a new release every time a new tag is created.

## Setup

It is easier to contribute to the project if you have pimp installed on your
machine. See the [installation instructions](https://www.pimp.dev/installation).

## Building

To build the latest dev version:

```
pimp build
```

## Installing

To install the latest dev version:

```
pimp install
```

## Testing

To run the test suite:

```
pimp test
```

## Releasing

To create and push a new tag (hence performing a release):

```
pimp major
pimp minor
pimp patch
```
