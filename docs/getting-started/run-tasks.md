# Run tasks

Similar to expanding commands, you can use [Pimpfiles](../user-guide-1/pimpfile.md) to run tasks anywhere on your system.

For example, you might want to run tests in your Go project. To do so create a `Pimfile` with the following content:

```yaml
# ./Pimpfile
test: go test -v ./...
```

Now you can execute `pimp test`, it will perform an exact match to the _test_ command and execute `go test -v ./...`. Note that if you already have a _test_ command in your

Read more about the [Task Runner](../user-guide-1/task-runner.md) in the documentation.

Next, let's see how you can actually render template files.

