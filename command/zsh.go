package command

import (
	"fmt"
	"os"
	"strings"

	"github.com/urfave/cli/v2"
)

var zshCommand = &cli.Command{
	Name:            "--zsh",
	Usage:           "Print the Zsh config and exit (aliases and completion)",
	SkipFlagParsing: true,
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
}

// Amazing resource: https://blog.kloetzl.info/how-to-write-a-zsh-completion/
var zshCompletionCommand = &cli.Command{
	Name:            "--zsh-completion",
	Hidden:          true,
	SkipFlagParsing: true,
	Action: func(c *cli.Context) error {
		var args []string
		for i, arg := range os.Args {
			if arg == "--" {
				args = os.Args[i+1:]
				break
			}
		}

		lastArg := args[len(args)-1]

		fmt.Fprintln(c.App.Writer, `
local ret=1
local -a options

options+=(
    `)

		for _, command := range c.App.VisibleCommands() {
			if strings.HasPrefix(command.Name, lastArg) {
				s := command.Name + "[" + command.Usage + "]"
				fmt.Fprintf(c.App.Writer, "'%s'\n", s)
			}
		}

		for _, flag := range c.App.VisibleFlags() {
			for _, name := range flag.Names() {
				var prefixedFlag string
				var suffix string
				if len(name) == 1 {
					prefixedFlag = "-" + name
					suffix = "+"
				} else {
					prefixedFlag = "--" + name
					suffix = "="
				}
				if strings.HasPrefix(prefixedFlag, lastArg) {
					s := prefixedFlag + suffix + "[" + getFlagUsage(flag) + "]"
					if isFlagTakesFile(flag) {
						s += ":file:_files:"
					}
					fmt.Fprintf(c.App.Writer, "'%s'\n", s)
				}
			}
		}

		fmt.Fprintln(c.App.Writer, `
)

_arguments -w -s -S $options[@] && ret=0
return ret
    `)

		return nil
	},
}
