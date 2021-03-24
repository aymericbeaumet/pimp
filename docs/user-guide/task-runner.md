# Task Runner

Following the same concept as command expansion \(see above\), you can also leverage pimp to behave as a task runner for your project. The `./Pimpfile` in your local directory always has the highest priority.

For example, this is how you would do if you wanted to define a `test` task:

```yaml
# ./Pimpfile
test: go test ./...
```

```bash
$ pimp test # runs `go test ./...`
```

