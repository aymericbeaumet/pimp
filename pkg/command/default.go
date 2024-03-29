package command

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/aymericbeaumet/pimp/pkg/template"
	"github.com/aymericbeaumet/pimp/pkg/util"
	"github.com/mattn/go-isatty"
	"github.com/urfave/cli/v2"
)

var DefaultCommand = &cli.Command{
	Hidden: true,
	Action: func(c *cli.Context) error {
		_, eng, err := initializeConfigEngine(c, true)
		if err != nil {
			return err
		}

		env, args, files, cwd := eng.Map(os.Environ(), c.Args().Slice())
		if len(args) == 0 {
			if len(c.String("input")) == 0 && !isatty.IsTerminal(os.Stdin.Fd()) {
				return execCommand.Action(c)
			}
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
		cmd.Dir = cwd
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
}
