package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"text/template"

	"github.com/aymericbeaumet/pimp/funcmap"
	fmerrors "github.com/aymericbeaumet/pimp/funcmap/errors"
	"github.com/mitchellh/go-homedir"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:        "pimp",
		Usage:       "Command line expander",
		UsageText:   "pimp [COMMAND] [OPTION]... [--] [BIN [ARG]...]",
		Version:     "0.0.1", // TODO: use -ldflags to embed the git commit hash
		Description: "Command expander. Shipped with a template engine, and more. Providing no COMMAND is the default and most common behavior, in this case BIN will be executed and given ARG as parameters.",

		Reader:          os.Stdin,
		Writer:          os.Stdout,
		ErrWriter:       os.Stderr,
		HideHelpCommand: true,

		Before: func(c *cli.Context) error {
			for _, flagName := range []string{"config", "input", "output"} {
				if s := c.String(flagName); len(s) > 0 {
					expanded, err := expand(s)
					if err != nil {
						return err
					}
					if err := c.Set(flagName, expanded); err != nil {
						return err
					}
				}
			}

			if filename := c.String("input"); len(filename) > 0 {
				f, err := os.Open(filename)
				if err != nil {
					return err
				}
				c.App.Reader = f
			}

			if filename := c.String("output"); len(filename) > 0 {
				f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
				if err != nil {
					return err
				}
				c.App.Writer = f
			}

			// see the corresponding flags (not commands) if you want to know why we
			// need this
			for _, flagName := range []string{"dump", "render", "shell"} {
				if c.IsSet(flagName) {
					commandName := "--" + flagName
					if command := c.App.Command(commandName); command != nil {
						if err := command.Run(c); err != nil {
							_, _ = fmt.Fprintf(c.App.ErrWriter, "Command %s failed: %s\n\n", commandName, err)
							_ = cli.ShowAppHelp(c)
							syscall.Exit(1)
						}
						syscall.Exit(0)
					}
					panic(fmt.Errorf("implementation error: command %s is missing", commandName))
				}
			}

			return nil
		},

		Action: func(c *cli.Context) error {
			engine, err := NewEngineFromFile(c.String("config"))
			if err != nil {
				return err
			}

			env, args, files := engine.Map(os.Environ(), c.Args().Slice())
			if len(args) == 0 {
				_ = cli.ShowAppHelp(c)
				return nil
			}

			for i, arg := range args {
				args[i], err = render(arg)
				if err != nil {
					return err
				}
			}

			if c.IsSet("expand") {
				for i, arg := range args {
					if i > 0 {
						fmt.Print(" ")
					}
					fmt.Printf("%#v", arg)
				}
				fmt.Print("\n")
				return nil
			}

			for filename, data := range files {
				rendered, err := render(data)
				if err != nil {
					return err
				}
				if err := os.WriteFile(filename, []byte(rendered), 0400); err != nil {
					return err
				}
				defer os.Remove(filename)
			}

			cmd := exec.CommandContext(c.Context, args[0], args[1:]...)
			cmd.Env = env
			cmd.Stdin = c.App.Reader
			cmd.Stdout = c.App.Writer
			cmd.Stderr = os.Stderr

			signalC := make(chan os.Signal, 32)
			signal.Notify(signalC)

			if err := cmd.Start(); err != nil {
				return err
			}

			go func() {
				for signal := range signalC {
					_ = cmd.Process.Signal(signal)
				}
			}()

			state, err := cmd.Process.Wait()
			if err != nil {
				return err
			}

			syscall.Exit(state.ExitCode())
			return nil
		},

		Commands: []*cli.Command{
			{
				Name:  "--dump",
				Usage: "Dump the config as JSON and exit",
				Action: func(c *cli.Context) error {
					engine, err := NewEngineFromFile(c.String("config"))
					if err != nil {
						return err
					}
					return engine.JSON(c.App.Writer)
				},
			},

			{
				Name:  "--render",
				Usage: "Render the template and exit",
				Action: func(c *cli.Context) error {
					renderFilepath := c.String("render")
					if len(renderFilepath) == 0 {
						return errors.New("expect one parameter")
					}

					renderFilepath, err := expand(renderFilepath)
					if err != nil {
						return err
					}

					data, err := os.ReadFile(renderFilepath)
					if err != nil {
						return err
					}

					s := string(data)

					// strip shebang if found
					if strings.HasPrefix(s, "#!") {
						if newlineIndex := strings.IndexRune(s, '\n'); newlineIndex > -1 {
							s = s[newlineIndex+1:]
						} else {
							s = ""
						}
					}

					rendered, err := render(s)
					if err != nil {
						return err
					}

					if _, err := c.App.Writer.Write([]byte(rendered)); err != nil {
						return err
					}

					return nil
				},
			},

			{
				Name:  "--shell",
				Usage: "Print the shell config (bash, zsh, fish, ...) and exit",
				Action: func(c *cli.Context) error {
					engine, err := NewEngineFromFile(c.String("config"))
					if err != nil {
						return err
					}
					for _, executable := range engine.Executables() {
						fmt.Fprintf(c.App.Writer, "alias %#v=%#v\n", executable, "pimp "+executable)
					}
					return nil
				},
			},
		},

		Flags: []cli.Flag{
			// Register hidden flags that are used to trigger the --dump, --render and --shell commands
			// This is needed as --[commands] are not supported by the parser
			&cli.BoolFlag{Name: "dump", Value: false, Hidden: true},
			&cli.StringFlag{Name: "render", Value: "", Hidden: true},
			&cli.BoolFlag{Name: "shell", Value: false, Hidden: true},

			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Value:   "~/.pimprc",
				Usage:   "Provide a different config `FILE`",
				EnvVars: []string{"PIMP_CONFIG"},
			},

			&cli.BoolFlag{
				Name:  "expand",
				Value: false,
				Usage: "Expand and print the command instead of running it",
			},

			&cli.StringFlag{
				Name:  "input",
				Value: "",
				Usage: "Read the input from `FILE` instead of stdin",
			},

			&cli.StringFlag{
				Name:    "output",
				Aliases: []string{"o"},
				Value:   "",
				Usage:   "Write the output to `FILE` instead of stdout",
			},
		},
	}

	_ = app.RunContext(context.Background(), os.Args)
}

var fm = funcmap.FuncMap()

func render(text string) (string, error) {
	var sb strings.Builder

	t, err := template.New(text).Funcs(fm).Parse(text)
	if err != nil {
		return "", err
	}

	if err := t.Execute(&sb, nil); err != nil {
		if e, ok := err.(template.ExecError); ok {
			// TODO: wait for this issue to be fixed upstream so that Unwrap()
			// returns the actual error that was returned (probably in Go 1.17).
			// In the meantime we cannot access the underlying error to cleanly
			// write to stderr + exit with the proper status code, so we panic.
			// https://github.com/golang/go/issues/34201
			err = e.Unwrap()
		}
		switch e := err.(type) {
		case *fmerrors.FatalError:
			os.Stderr.WriteString(e.Error())
			syscall.Exit(e.ExitCode())
		default:
			return "", err
		}
	}

	return sb.String(), nil
}

func expand(input string) (string, error) {
	if len(input) == 0 {
		return "", nil
	}

	expanded, err := homedir.Expand(input)
	if err != nil {
		return "", err
	}
	if strings.HasPrefix(expanded, "/") {
		return expanded, nil
	}

	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return filepath.Join(wd, expanded), nil
}
