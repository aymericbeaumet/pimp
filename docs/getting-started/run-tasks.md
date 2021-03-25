# Run tasks

Similar to expanding commands, you can use [Pimpfiles](../user-guide/pimpfile.md) to run tasks anywhere on your system.

For example, you might want to run tests in your Go project. To do so create a `Pimpfile` in your project with the following content:

```yaml
# ./Pimpfile
test: go test -v ./...
```

Now you can execute `pimp test`, it will perform an exact match to the _test_ command and execute `go test -v ./...`. Note that the local Pimpfiles take precedence over your global Pimpfile in case the same command would be defined twice. 

Read more about the [Task Runner](../user-guide/task-runner.md) in the documentation.

Next, let's see how you can actually render template files.

