package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"runtime/debug"
	"strings"
	"syscall"

	"github.com/aymericbeaumet/pimp/command"
	"github.com/aymericbeaumet/pimp/normalize"
	"github.com/urfave/cli/v2"
)

// Populated by Goreleaser via ldflags (https://goreleaser.com/customization/build/)
var version = "$version"
var commit = "$commit"
var date = "$date"
var builtBy = "$builtBy"

func init() {
	debug.SetGCPercent(-1)
}

func main() {
	app := cli.NewApp()

	app.Name = "pimp"
	app.Version = version
	app.Description = strings.TrimSpace(`
pimp is a command-line expander with pattern matching and templating
capabilities that increases your productivity.
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
                                     {{"\t"}}(priority: pimpfile, configuration file, $PATH commands)
{{- range .VisibleCommands}}
{{- if not .HideHelp}}
    pimp [OPTION]... {{.Name}} {{.ArgsUsage}}{{"\t"}}{{.Usage}}
{{- end}}
{{- end}}

OPTIONS:
{{- range .VisibleFlags}}
    {{.}}
{{- end}}

EXAMPLES:

    Let's start with the classic "Hello, World!". This illustrates how pimp
    acts as a fancy command proxy. No expansion is performed here.

        $ pimp echo 'Hello, World!'
        Hello, World!

    Let's make it a little bit more interesting by adding some mappings to the
    ~/.pimprc configuration file. Pimp stops after the first match is found.
    Note how "..." enables us to catch variadic arguments which are
    automatically appended during the expansion process.

        $ cat ~/.pimprc
        git co     : git checkout {{"{{GitLocalBranches | FZF}}"}}
        git co ... : git checkout
        $ pimp git co{{"\t"}}# executes "git checkout <branch>" with the branch name chosen in fzf
        $ pimp git co master{{"\t"}}# executes "git checkout master" ("master" is from the "...")

    To make this more convenient, you can send all the "git" calls to pimp with
    a shell alias (this is automatically done if you place $(pimp --shell) in
    your shell configuration).

        $ alias git='pimp git'
        $ git co{{"\t"}}# same as in the previous example
        $ git co master{{"\t"}}# same as in the previous example

    You can also leverage pimp templating system to render arbitrary files.

        $ pimp -o readme.md --render readme.md.tmpl{{"\t"}}# Overwrite the readme with the rendered template

    It is also possible to render pimp templates from the command-line arguments.

        $ pimp --run {{"'{{GitRemotes | JSON}}'"}}
        ["origin"]

    See the project homepage for more advanced examples.
`

	app.Before = func(c *cli.Context) error {
		for _, flag := range c.App.Flags {
			if flag, ok := flag.(*cli.StringFlag); ok && flag.TakesFile {
				if value := c.String(flag.Name); len(value) > 0 {
					normalized, err := normalize.Path(value)
					if err != nil {
						return fmt.Errorf("error when normalizing `%s` flag: %w", flag.Name, err)
					}
					if err := c.Set(flag.Name, normalized); err != nil {
						return fmt.Errorf("error when setting `%s` flag: %w", flag.Name, err)
					}
				}
			}
		}

		if filename := c.String("input"); len(filename) > 0 {
			f, err := os.Open(filename)
			if err != nil {
				return fmt.Errorf("error for `input` flag: %w", err)
			}
			c.App.Reader = f
		}

		if filename := c.String("output"); len(filename) > 0 {
			flags := os.O_WRONLY | os.O_CREATE
			if c.Bool("append") {
				flags = flags | os.O_APPEND
			} else {
				flags = flags | os.O_TRUNC
			}
			f, err := os.OpenFile(filename, flags, 0644)
			if err != nil {
				return fmt.Errorf("error for `output` flag: %w", err)
			}
			c.App.Writer = f
		}

		for _, command := range command.Commands {
			if c.Bool(strings.TrimPrefix(command.Name, "--")) {
				command := c.App.Command(command.Name)
				if command == nil {
					panic(fmt.Errorf("implementation error: command %s is missing", command.Name))
				}
				if err := command.Action(c); err != nil {
					return fmt.Errorf("command %s failed: %w", command.Name, err)
				}
				syscall.Exit(0)
			}
		}

		return nil
	}

	app.Action = command.MainAction

	app.Commands = command.Commands
	app.HideHelpCommand = true

	app.Flags = append(command.CommandsFlags(),
		&cli.BoolFlag{
			Name:  "append",
			Usage: "Append to the output file instead of truncating",
		},

		&cli.StringFlag{
			Name:      "config",
			Aliases:   []string{"c"},
			EnvVars:   []string{"PIMP_CONFIG"},
			Usage:     "Load this configuration `FILE` (default: ~/.pimprc)",
			TakesFile: true,
		},

		&cli.BoolFlag{
			Name:  "expand",
			Usage: "Expand and print the command without executing",
		},

		&cli.StringFlag{
			Name:      "file",
			Aliases:   []string{"f"},
			EnvVars:   []string{"PIMP_FILE"},
			Usage:     "Load this pimpfile `FILE` (default: ./Pimpfile)",
			TakesFile: true,
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
			Usage: "Left template delimiter",
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
			Usage: "Right template delimiter",
			Value: "}}",
		},
	)

	if err := app.RunContext(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
