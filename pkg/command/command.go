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

	"github.com/aymericbeaumet/pimp/pkg/engine"
	"github.com/aymericbeaumet/pimp/pkg/funcs"
	"github.com/aymericbeaumet/pimp/pkg/template"
	"github.com/aymericbeaumet/pimp/pkg/util"
	"github.com/urfave/cli/v2"
)

var funcmap = funcs.FuncMap()

var Commands = []*cli.Command{
	dumpCommand,
	renderCommand,
	runCommand,
	shellCommand,
	transpileCommand,
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

	args, err = template.RenderStrings(args, c.String("ldelim"), c.String("rdelim"), funcmap)
	if err != nil {
		return err
	}
	args = util.FilterEmptyStrings(args)

	for filename, data := range files {
		f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0400)
		if err != nil {
			return err
		}
		if err := template.Render(f, data, c.String("ldelim"), c.String("rdelim"), funcmap); err != nil {
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

	filename, err := util.NormalizePath(filename)
	if err != nil {
		return nil, err
	}

	f, err := os.Open(filename)
	if err != nil && errors.Is(err, fs.ErrNotExist) && allowErrNotExist {
		return nil, nil
	}
	return f, err
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
