# Contributing

The [Go toolchain](https://golang.org/doc/install) is required to work on this project.

## Building

```bash
$ go build -o pimp .
```

## Installing

```text
go install .
```

## Testing

```text
go test -v ./...
```

## Releasing

Pimp tasks are defined in the [./Pimpfile](https://github.com/aymericbeaumet/pimp/blob/master/Pimpfile) to release new versions. The release process is entirely automated and is being taken care of by the CI.

```text
pimp major
pimp minor
pimp patch
```

