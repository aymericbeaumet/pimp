package main

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"text/template"

	"github.com/aymericbeaumet/pimp/engine"
	"github.com/aymericbeaumet/pimp/funcmap"
	fmerrors "github.com/aymericbeaumet/pimp/funcmap/errors"
	"github.com/mitchellh/go-homedir"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:      "pimp",
		Usage:     "Command line expander",
		UsageText: "pimp [OPTION]... [--] BIN [ARG]...\n   pimp [OPTION]... COMMAND [ARG]...",
		Version:   "0.0.1", // TODO: use -ldflags to embed the git commit hash
		Description: strings.TrimSpace(`
Command expander. Shipped with a template engine, and more. Providing no
COMMAND is the default and most common behavior, in this case BIN will be
executed and given ARG as parameters.
    `),

		Reader:          os.Stdin,
		Writer:          os.Stdout,
		ErrWriter:       os.Stderr,
		HideHelpCommand: true,

		Before: func(c *cli.Context) error {
			for _, flagName := range []string{"config", "file", "input", "output"} {
				if s := c.String(flagName); len(s) > 0 {
					expanded, err := normalizePath(s)
					if err != nil {
						return fmt.Errorf("error for `%s` flag: %v", flagName, err)
					}
					if err := c.Set(flagName, expanded); err != nil {
						return fmt.Errorf("error for `%s` flag: %v", flagName, err)
					}
				}
			}

			if filename := c.String("input"); len(filename) > 0 {
				f, err := os.Open(filename)
				if err != nil {
					return fmt.Errorf("error for `input` flag: %v", err)
				}
				c.App.Reader = f
			}

			if filename := c.String("output"); len(filename) > 0 {
				f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
				if err != nil {
					return fmt.Errorf("error for `output` flag: %v", err)
				}
				c.App.Writer = f
			}

			if c.Bool("zsh-completion") {
				word := os.Args[len(os.Args)-1]

				fmt.Fprintln(c.App.Writer, "local -a flags")
				for _, flag := range c.App.Flags {
					for _, name := range flag.Names() {
						var prefixedFlag string
						if len(name) == 1 {
							prefixedFlag = "-" + name
						} else {
							prefixedFlag = "--" + name
						}
						if strings.HasPrefix(prefixedFlag, word) && prefixedFlag != word {
							fmt.Fprintf(c.App.Writer, "flags+=(%#v)\n", prefixedFlag)
						}
					}
				}
				fmt.Fprintln(c.App.Writer, "_describe flags flags")

				fmt.Fprintln(c.App.Writer, "_files")

				syscall.Exit(0)
			}

			for _, flagName := range []string{"dump", "render", "shell", "bash", "zsh"} {
				if c.Bool(flagName) {
					commandName := "--" + flagName
					if command := c.App.Command(commandName); command != nil {
						if err := command.Run(c); err != nil {
							return fmt.Errorf("command %s failed: %s", commandName, err)
						}
						syscall.Exit(0)
					}
					panic(fmt.Errorf("implementation error: command %s is missing", commandName))
				}
			}

			return nil
		},

		Action: func(c *cli.Context) error {
			eng, err := initializeEngine(c)
			if err != nil {
				return err
			}

			env, args, files := eng.Map(os.Environ(), c.Args().Slice())
			if len(args) == 0 {
				_ = cli.ShowAppHelp(c)
				return nil
			}

			args, err = renderStrings(args)
			if err != nil {
				return err
			}
			args = filterEmptyStrings(args)

			for filename, data := range files {
				rendered, err := render(data)
				if err != nil {
					return err
				}
				if err := os.WriteFile(filename, []byte(rendered), 0400); err != nil {
					return err
				}
				if !c.Bool("keep") {
					defer os.Remove(filename)
				}
			}

			if c.Bool("expand") {
				for i, arg := range args {
					if i > 0 {
						fmt.Print(" ")
					}
					fmt.Printf("%#v", arg)
				}
				fmt.Print("\n")
				return nil
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
				Usage: "Dump the engine as JSON and exit",
				Action: func(c *cli.Context) error {
					eng, err := initializeEngine(c)
					if err != nil {
						return err
					}
					return eng.JSON(c.App.Writer)
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

					renderFilepath, err := normalizePath(renderFilepath)
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
				Usage: "Print the shell config and exit (aliases only)",
				Action: func(c *cli.Context) error {
					eng, err := initializeEngine(c)
					if err != nil {
						return err
					}
					for _, executable := range eng.Executables() {
						if _, err := fmt.Fprintf(
							c.App.Writer, "alias %#v=%#v\n", executable, "pimp "+executable,
						); err != nil {
							return err
						}
					}
					return nil
				},
			},

			{
				Name:  "--zsh",
				Usage: "Print the Zsh config and exit (aliases and completion)",
				Action: func(c *cli.Context) error {
					if err := c.App.Command("--shell").Run(c); err != nil {
						return err
					}
					_, err := fmt.Fprintln(c.App.Writer, `
            _pimp() {
              eval "$(pimp --zsh-completion -- "${words[@]}")"
            }

            compdef _pimp pimp
          `)
					return err
				},
			},
		},

		Flags: []cli.Flag{
			// Register hidden flags that are used to trigger the corresponding
			// commands as --[command] is not supported by the parser except for
			// flags
			&cli.BoolFlag{Name: "dump", Value: false, Hidden: true},
			&cli.StringFlag{Name: "render", Value: "", Hidden: true},
			&cli.BoolFlag{Name: "shell", Value: false, Hidden: true},
			&cli.BoolFlag{Name: "zsh", Value: false, Hidden: true},
			&cli.BoolFlag{Name: "zsh-completion", Value: false, Hidden: true},

			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Value:   "",
				Usage:   "Provide a different config `FILE`",
				EnvVars: []string{"PIMP_CONFIG"},
			},

			&cli.BoolFlag{
				Name:  "expand",
				Value: false,
				Usage: "Expand and print the command instead of running it",
			},

			&cli.StringFlag{
				Name:    "file",
				Aliases: []string{"f"},
				Value:   "",
				Usage:   "Read `FILE` as a pimpfile",
			},

			&cli.StringFlag{
				Name:  "input",
				Value: "",
				Usage: "Read the input from `FILE` instead of stdin",
			},

			&cli.BoolFlag{
				Name:  "keep",
				Value: false,
				Usage: "Keep the temporary files",
			},

			&cli.StringFlag{
				Name:    "output",
				Aliases: []string{"o"},
				Value:   "",
				Usage:   "Write the output to `FILE` instead of stdout",
			},
		},
	}

	if err := app.RunContext(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
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

// renderStrings renders several strings in a single context. This makes it
// possible to interact between several templates with variable declarations,
// etc. This could generate empty strings in the output that have to be dealt
// with.
func renderStrings(texts []string) ([]string, error) {
	const SEP = "\x00pimp\x00"

	joined := strings.Join(texts, SEP)

	rendered, err := render(joined)
	if err != nil {
		return nil, err
	}

	return strings.Split(rendered, SEP), nil
}

func initializeEngine(c *cli.Context) (*engine.Engine, error) {
	eng := engine.New()

	for flagName, fallback := range map[string]string{
		"file":   "./Pimpfile",
		"config": "~/.pimprc",
	} {
		file, err := openFallback(c.String(flagName), fallback)
		if err != nil {
			return nil, err
		}
		if file != nil {
			defer file.Close()
			if err := eng.Append(file); err != nil {
				return nil, err
			}
		}
	}

	return eng, nil
}

func openFallback(filename, fallback string) (*os.File, error) {
	allowOpenError := false

	if len(filename) == 0 && len(fallback) > 0 {
		allowOpenError = true
		filename = fallback
	}

	filename, err := normalizePath(filename)
	if err != nil {
		return nil, err
	}

	f, err := os.Open(filename)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) && allowOpenError {
			return nil, nil
		}
		return nil, err
	}
	return f, nil
}

func normalizePath(input string) (string, error) {
	input = strings.TrimSpace(input)
	if len(input) == 0 {
		return input, nil
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

func filterEmptyStrings(input []string) []string {
	out := make([]string, 0, len(input))
	for _, i := range input {
		if trimmed := strings.TrimSpace(i); len(trimmed) > 0 {
			out = append(out, trimmed)
		}
	}
	return out
}
