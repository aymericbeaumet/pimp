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

		fmt.Fprintln(c.App.Writer, "local -a commands")
		for _, command := range c.App.VisibleCommands() {
			if strings.HasPrefix(command.Name, lastArg) {
				fmt.Fprintf(c.App.Writer, "commands+=('%s')\n", command.Name+":"+command.Usage)
			}
		}
		fmt.Fprintln(c.App.Writer, "_describe -t commands 'commands' commands")

		fmt.Fprintln(c.App.Writer, "local -a flags")
		for _, flag := range c.App.VisibleFlags() {
			for _, name := range flag.Names() {
				var prefixedFlag string
				if len(name) == 1 {
					prefixedFlag = "-" + name
				} else {
					prefixedFlag = "--" + name
				}
				if strings.HasPrefix(prefixedFlag, lastArg) {
					fmt.Fprintf(c.App.Writer, "flags+=('%s')\n", prefixedFlag+":"+getFlagUsage(flag))
				}
			}
		}
		fmt.Fprintln(c.App.Writer, "_describe -t flags 'flags' flags")

		return nil
	},
}
