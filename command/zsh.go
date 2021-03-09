package command

import (
	"fmt"
	"os"
	"strings"

	"github.com/urfave/cli/v2"
)

var zshCommand = &cli.Command{
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
}

var zshCompletionCommand = &cli.Command{
	Name:   "--zsh-completion",
	Hidden: true,
	Action: func(c *cli.Context) error {
		dashdashIndex := -1
		for i, arg := range os.Args {
			if dashdashIndex == -1 && arg == "--" {
				dashdashIndex = i
				break
			}
		}

		args := os.Args[dashdashIndex+1:]
		lastArg := args[len(args)-1]

		var bin string
		binIndex := dashdashIndex + 2
		if binIndex < len(os.Args) { // ... -- pimp BIN
			bin = os.Args[binIndex]
		}

		if binIndex < len(os.Args)-1 { // if bin is present, but not the latest arg
			fmt.Fprintf(c.App.Writer, "_%s\n", bin)
		} else if len(lastArg) == 0 || strings.HasPrefix(lastArg, "-") {
			fmt.Fprintln(c.App.Writer, "local -a flags")
			for _, flag := range c.App.Flags {
				for _, name := range flag.Names() {
					var prefixedFlag string
					if len(name) == 1 {
						prefixedFlag = "-" + name
					} else {
						prefixedFlag = "--" + name
					}
					if strings.HasPrefix(prefixedFlag, lastArg) && prefixedFlag != lastArg {
						fmt.Fprintf(c.App.Writer, "flags+=(%#v)\n", prefixedFlag)
					}
				}
			}
			fmt.Fprintln(c.App.Writer, "_describe flags flags")
		} else if len(bin) > 0 {
			fmt.Fprintln(c.App.Writer, "_files")
		} else if len(lastArg) >= 1 {
			fmt.Fprintln(c.App.Writer, "_path_commands")
		}

		return nil
	},
}
