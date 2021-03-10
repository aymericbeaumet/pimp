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

func init() {
	debug.SetGCPercent(-1)
}

func main() {
	app := &cli.App{
		Name:    "pimp",
		Version: "0.0.1", // TODO: use -ldflags to embed the version and git commit hash
		Description: strings.TrimSpace(`
pimp is a command-line expander and template engine that increases your command
line productivity.
		`),
		Authors: []*cli.Author{
			{
				Name:  "Aymeric Beaumet",
				Email: "hi@aymericbeaumet.com",
			},
		},
		Metadata: map[string]interface{}{
			"Website": "https://github.com/aymericbeaumet/pimp",
		},

		CustomAppHelpTemplate: `{{.Name}} {{.Version}}
{{with (index .Authors 0)}}{{.Name}} <{{.Email}}>{{end}}

{{.Description}}

Project home page: {{index .Metadata "Website"}}

USAGE:
    pimp [OPTION]... COMMAND [ARG]...{{ "\t\t" }}Expand the command and its arguments, execute and exit
{{ range .VisibleCommands}}
{{- if not .HideHelp}}
    pimp [OPTION]... {{ join .Names ", "}}{{ "\t"}}{{.Usage}}
{{- end}}
{{- end}}

OPTIONS:
{{- range .VisibleFlags}}
    {{ . -}}
{{- end}}

EXAMPLES:
    pimp git log{{ "\t\t" }}Expand and execute the 'git log' command
    pimp --render readme.md.tmpl > readme.md{{ "\t\t" }}Render the readme template and write it to readme.md
    pimp --run {{ "'{{GitBranches | JSON}}'" }}{{ "\t\t" }}Print the current git repository branches as JSON
`,

		Reader:          os.Stdin,
		Writer:          os.Stdout,
		ErrWriter:       os.Stderr,
		HideHelpCommand: true,

		Before: func(c *cli.Context) error {
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
		},

		Action: command.MainAction,

		Commands: command.Commands,

		Flags: append(command.CommandsFlags(),
			&cli.BoolFlag{
				Name:  "append",
				Usage: "Append to the output file instead of truncating",
			},

			&cli.StringFlag{
				Name:      "config",
				Aliases:   []string{"c"},
				EnvVars:   []string{"PIMP_CONFIG"},
				Usage:     "Provide a different config `FILE`",
				TakesFile: true,
			},

			&cli.BoolFlag{
				Name:  "expand",
				Usage: "Expand and print the command instead of running it",
			},

			&cli.StringFlag{
				Name:      "file",
				Aliases:   []string{"f"},
				Usage:     "Read `FILE` as a pimpfile",
				TakesFile: true,
			},

			&cli.StringFlag{
				Name:      "input",
				Usage:     "Read the input from `FILE` instead of stdin",
				TakesFile: true,
			},

			&cli.BoolFlag{
				Name:  "keep",
				Usage: "Keep the temporary files",
			},

			&cli.StringFlag{
				Name:      "output",
				Aliases:   []string{"o"},
				Usage:     "Write the output to `FILE` instead of stdout",
				TakesFile: true,
			},
		),
	}

	if err := app.RunContext(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
