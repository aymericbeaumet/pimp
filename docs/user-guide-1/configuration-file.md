# Configuration File

The default location for the configuration file is `~/.pimprc` . You can configure this via the `--config` CLI flag or the `PIMP_CONFIG` environment variable. If the configuration file or a specific configuration key is missing, then the corresponding default value as shown below will be used.

```yaml
# ~/.pimprc (with default values)

# pimpfiles ([]string) contains all the global Pimpfiles that should be resolved
# and will act as a fallback when a command is expanded.
pimpfiles:
  - ~/.Pimpfile.go
  - ~/.Pimpfile
```

