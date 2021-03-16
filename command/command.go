package command

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"text/template"

	"github.com/aymericbeaumet/pimp/engine"
	perrors "github.com/aymericbeaumet/pimp/errors"
	"github.com/aymericbeaumet/pimp/funcs"
	"github.com/aymericbeaumet/pimp/normalize"
	"github.com/urfave/cli/v2"
)

var Commands = []*cli.Command{
	dumpCommand,
	evalCommand,
	renderCommand,
	runCommand,
	shellCommand,
	zshCommand,
	zshCompletionCommand,
}

func CommandsFlags() []cli.Flag {
	out := make([]cli.Flag, 0, len(Commands))
	for _, c := range Commands {
		out = append(out, &cli.BoolFlag{
			Name:   strings.TrimPrefix(c.Name, "--"),
			Value:  false,
			Hidden: true,
		})
	}
	return out
}

func DefaultCommand(c *cli.Context) error {
	eng, err := initializeEngine(c, true, true)
	if err != nil {
		return err
	}

	env, args, files := eng.Map(os.Environ(), c.Args().Slice())
	if len(args) == 0 {
		return replCommand.Action(c)
	}

	args, err = renderStrings(args, c.String("ldelim"), c.String("rdelim"))
	if err != nil {
		return err
	}
	args = filterEmptyStrings(args)

	for filename, data := range files {
		rendered, err := render(data, c.String("ldelim"), c.String("rdelim"))
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
}

var fm = funcs.FuncMap()

func render(text string, ldelim, rdelim string) (string, error) {
	var sb strings.Builder

	t, err := template.New(text).Funcs(fm).Delims(ldelim, rdelim).Parse(text)
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
		case *perrors.FatalError:
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
func renderStrings(texts []string, ldelim, rdelim string) ([]string, error) {
	const SEP = "\x00pimp\x00"

	joined := strings.Join(texts, SEP)

	rendered, err := render(joined, ldelim, rdelim)
	if err != nil {
		return nil, err
	}

	return strings.Split(rendered, SEP), nil
}

func initializeEngine(c *cli.Context, loadConfig, loadPimpfile bool) (*engine.Engine, error) {
	eng := engine.New()

	type source struct{ flagName, fallback string }
	sources := []*source{}
	if loadPimpfile {
		sources = append(sources, &source{
			flagName: "file",
			fallback: "./Pimpfile",
		})
	}
	if loadConfig {
		sources = append(sources, &source{
			flagName: "config",
			fallback: "~/.pimprc",
		})
	}

	for _, s := range sources {
		file, err := openWithFallback(c.String(s.flagName), s.fallback)
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

func openWithFallback(filename, fallback string) (*os.File, error) {
	allowErrNotExist := false

	if len(filename) == 0 && len(fallback) > 0 {
		allowErrNotExist = true
		filename = fallback
	}

	filename, err := normalize.Path(filename)
	if err != nil {
		return nil, err
	}

	f, err := os.Open(filename)
	if err != nil && errors.Is(err, fs.ErrNotExist) && allowErrNotExist {
		return nil, nil
	}
	return f, err
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

func getFlagUsage(flag cli.Flag) string {
	switch flag := flag.(type) {
	case *cli.BoolFlag:
		return flag.Usage
	case *cli.StringFlag:
		return flag.Usage
	case *cli.StringSliceFlag:
		return flag.Usage
	default:
		return ""
	}
}

func isFlagTakesFile(flag cli.Flag) bool {
	f, ok := flag.(*cli.StringFlag)
	return ok && f.TakesFile
}

func isFlagAllowedMultipleTimes(flag cli.Flag) bool {
	_, ok := flag.(*cli.StringSliceFlag)
	return ok
}
