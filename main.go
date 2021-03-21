package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"runtime/debug"
	"strings"
	"syscall"

	"github.com/aymericbeaumet/pimp/pkg/command"
	"github.com/aymericbeaumet/pimp/pkg/util"
	"github.com/fatih/color"
	"github.com/mattn/go-shellwords"
	"github.com/urfave/cli/v2"
)

// Populated by Goreleaser via ldflags (https://goreleaser.com/customization/build/)
var version = "$version"
var commit = "$commit"
var date = "$date"
var builtBy = "$builtBy"

var colorExample = color.New(color.FgYellow).SprintFunc()

func init() {
	debug.SetGCPercent(-1)
}

func main() {
	app := cli.NewApp()

	app.Name = "pimp"
	app.Version = version
	app.Description = strings.TrimSpace(`
pimp is a shell-agnostic command-line expander and command runner with pattern
matching and templating capabilities that increases your productivity.
		`)
	app.Authors = []*cli.Author{
		{
			Name:  "Aymeric Beaumet",
			Email: "hi@aymericbeaumet.com",
		},
	}
	app.Metadata = map[string]interface{}{
		"builtBy": builtBy,
		"commit":  commit,
		"date":    date,
		"website": "https://github.com/aymericbeaumet/pimp",
	}

	app.CustomAppHelpTemplate = `{{.Name}} {{.Version}}
{{- range .Authors}}
{{.Name}} <{{.Email}}>{{end}}

{{.Description}}

Project homepage: {{index .Metadata "website"}}

USAGE:
    pimp [OPTION]... COMMAND [ARG]...{{"\t"}}Match COMMAND and ARGS, expand, then execute
                                     {{"\t"}}(resolution order: Pimpfiles, $PATH commands)
{{- range .VisibleCommands}}
{{- if not .HideHelp}}
    pimp [OPTION]... {{.Name}} {{.ArgsUsage}}{{"\t"}}{{.Usage}}
{{- end}}
{{- end}}

    Execute pimp without arguments to start a REPL.
    Pipe a script into pimp without arguments to evaluate it.

OPTIONS:
{{- range .VisibleFlags}}
    {{.}}
{{- end}}

EXAMPLES:

    Let's start with the classic "Hello, World!". This illustrates how pimp
    acts as a fancy command proxy. No expansion is performed here.

        $ ` + colorExample(`pimp echo 'Hello, World!'`) + `
        Hello, World!

    Let's make it a little bit more interesting by adding some mappings to a
    Pimpfile. Pimp stops as soon as a match is found. Note how "..." enables us
    to catch variadic arguments which are automatically appended during the
    expansion process.

        $ ` + colorExample(`cat ./Pimpfile`) + `
        git co     : git checkout {{"{{GitLocalBranches | FZF}}"}}
        git co ... : git checkout
        $ ` + colorExample(`pimp git co`) + `{{"\t"}}# executes "git checkout <branch>" with the branch name chosen in fzf
        $ ` + colorExample(`pimp git co master`) + `{{"\t"}}# executes "git checkout master" ("master" is from the "...")

    You can also leverage the pimp templating system to render arbitrary files.

        $ ` + colorExample(`pimp -o readme.md --render readme.md.tmpl`) + `{{"\t"}}# Overwrite the readme with the rendered template

    See the project homepage and documentation for more advanced examples.
`

	app.Before = func(c *cli.Context) error {
		for _, flag := range c.App.Flags {
			switch flag := flag.(type) {
			case *cli.StringFlag:
				if flag.TakesFile {
					if value := c.String(flag.Name); len(value) > 0 {
						normalized, err := util.NormalizePath(value)
						if err != nil {
							return fmt.Errorf("error when normalizing `%s` flag: %w", flag.Name, err)
						}
						if err := c.Set(flag.Name, normalized); err != nil {
							return fmt.Errorf("error when setting `%s` flag: %w", flag.Name, err)
						}
					}
				}
			case *cli.StringSliceFlag:
				if flag.TakesFile {
					values := c.StringSlice(flag.Name)
					for i := range values {
						normalized, err := util.NormalizePath(values[i])
						if err != nil {
							return fmt.Errorf("error when normalizing `%s` flag: %w", flag.Name, err)
						}
						values[i] = normalized
					}
				}
			}
		}

		for _, s := range c.StringSlice("env") {
			split := strings.SplitN(s, "=", 2)
			if len(split) != 2 {
				return fmt.Errorf("error for `env` flag: %#v should be of length 2", split)
			}
			os.Setenv(split[0], split[1])
		}

		if filename := c.String("input"); len(filename) > 0 {
			f, err := os.Open(filename)
			if err != nil {
				return fmt.Errorf("error for `input` flag: %w", err)
			}
			c.App.Reader = f
		}

		if filename := c.String("output"); len(filename) > 0 {
			var perm os.FileMode
			var flag int

			if c.Bool("frozen") {
				perm = 0
				flag |= os.O_RDONLY
			} else {
				perm = 0644
				flag |= os.O_WRONLY | os.O_CREATE
				if c.Bool("append") {
					flag |= os.O_APPEND
				} else {
					flag |= os.O_TRUNC
				}
			}

			f, err := os.OpenFile(filename, flag, perm)
			if err != nil {
				return fmt.Errorf("error for `output` flag: %w", err)
			}

			if c.Bool("frozen") {
				var out bytes.Buffer
				c.App.Writer = &out
				c.App.After = func(c *cli.Context) error {
					truth, err := io.ReadAll(f)
					if err != nil {
						return err
					}
					if !bytes.Equal(truth, out.Bytes()) {
						return fmt.Errorf("output differs for output %s", filename)
					}
					return nil
				}
			} else {
				c.App.Writer = f
			}
		}

		for _, command := range command.Commands {
			if c.Bool(strings.TrimPrefix(command.Name, "--")) {
				command := c.App.Command(command.Name)
				if command == nil {
					panic(fmt.Errorf("implementation error: command %s is missing", command.Name))
				}
				if err := command.Action(c); err != nil {
					return fmt.Errorf("%s failed: %w", command.Name, err)
				}
				if after := c.App.After; after != nil {
					if err := after(c); err != nil {
						return fmt.Errorf("%s failed: %w", command.Name, err)
					}
				}
				syscall.Exit(0)
			}
		}

		return nil
	}

	app.Action = command.DefaultCommand.Action
	app.Commands = command.Commands
	app.HideHelpCommand = true

	app.Flags = append(command.CommandsFlags(),
		&cli.BoolFlag{
			Name:  "append",
			Usage: "Append to the --output file instead of truncating",
		},

		&cli.StringFlag{
			Name:      "config",
			Aliases:   []string{"c"},
			EnvVars:   []string{"PIMP_CONFIG"},
			Usage:     "Load this configuration `FILE` (default: ~/.pimprc)",
			TakesFile: true,
		},

		&cli.StringSliceFlag{
			Name:    "env",
			Aliases: []string{"e"},
			Usage:   "Define one or more environment variables in the form `KEY=VALUE`",
		},

		&cli.BoolFlag{
			Name:  "expand",
			Usage: "Expand and print the command without executing",
		},

		&cli.StringSliceFlag{
			Name:      "file",
			Aliases:   []string{"f"},
			Usage:     "Load one or more `PIMPFILE`S, bypassing the default resolution mecanism",
			TakesFile: true,
		},

		&cli.BoolFlag{
			Name:  "frozen",
			Usage: "Fail if the output differs from the --output file",
		},

		&cli.StringFlag{
			Name:      "input",
			Usage:     "Read the input from `FILE` instead of stdin",
			TakesFile: true,
		},

		&cli.BoolFlag{
			Name:  "keep",
			Usage: "Keep the temporary shebang files",
		},

		&cli.StringFlag{
			Name:  "ldelim",
			Usage: "Define the left template delimiter",
			Value: "{{",
		},

		&cli.StringFlag{
			Name:      "output",
			Aliases:   []string{"o"},
			Usage:     "Write the output to `FILE` instead of stdout",
			TakesFile: true,
		},

		&cli.StringFlag{
			Name:  "rdelim",
			Usage: "Define the right template delimiter",
			Value: "}}",
		},
	)

	defaultOpts, err := shellwords.Parse(os.Getenv("PIMP_DEFAULT_OPTS"))
	if err != nil {
		log.Fatal(err)
	}
	args := make([]string, 0, len(os.Args)+len(defaultOpts))
	args = append(args, os.Args[0])
	args = append(args, defaultOpts...)
	args = append(args, os.Args[1:]...)

	if err := app.RunContext(context.Background(), args); err != nil {
		log.Fatal(err)
	}
}
