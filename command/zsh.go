package command

import (
	"fmt"
	"strings"

	"github.com/urfave/cli/v2"
)

var zshCommand = &cli.Command{
	Name:  "--zsh",
	Usage: "Print the Zsh config (aliases and completion)",
	Action: func(c *cli.Context) error {
		eng, err := initializeEngine(c, true, false)
		if err != nil {
			return err
		}

		for _, executable := range eng.Executables() {
			fmt.Fprintf(c.App.Writer, "alias %#v=%#v\n", executable, "pimp "+executable)
		}

		fmt.Fprintln(c.App.Writer, `
_pimp() {
  eval "$(pimp --zsh-completion -- "${words[@]}")"
}

compdef _pimp pimp
    `)

		return nil
	},
}

// Writing Zsh completion functions it not an easy task, here are some
// resources that were very helpful:
// - http://zsh.sourceforge.net/Doc/Release/Completion-System.html#Completion-System
// - https://blog.kloetzl.info/how-to-write-a-zsh-completion/
// - https://stackoverflow.com/a/13547531/1071486
var zshCompletionCommand = &cli.Command{
	Name:   "--zsh-completion",
	Hidden: true,
	Action: func(c *cli.Context) error {
		args := c.Args().Slice()

		if len(args) >= 2 && args[0] == "pimp" && !strings.HasPrefix(args[1], "-") {
			// TODO: parse args + expannd BIN ARGS to provide accurate completion
			if len(args) == 2 {
				fmt.Fprintln(c.App.Writer, "_path_commands")
			} else {
				delta := 1
				fmt.Fprintf(c.App.Writer, "shift %d words\n", delta)
				fmt.Fprintf(c.App.Writer, "(( CURRENT -= %d ))\n", delta)
				fmt.Fprintf(c.App.Writer, "_normal -p %#v\n", args[1])
			}
			return nil
		}

		sharedExclusionList := []string{"-h", "--help", "-v", "--version"}

		fmt.Fprintln(c.App.Writer, `
local ret=1
local -a options

options+=(
    `)

		fmt.Fprintln(c.App.Writer, `+ '(commands)'`) //
		for _, command := range c.App.VisibleCommands() {
			fmt.Fprintln(c.App.Writer, `'(- *)`+command.Name+"["+command.Usage+`]'`)
		}

		fmt.Fprintln(c.App.Writer, "")

		fmt.Fprintln(c.App.Writer, `+ 'flags'`) //
		for _, flag := range c.App.VisibleFlags() {
			exclusionList := append([]string{}, sharedExclusionList...)
			flagList := []string{}

			for _, name := range flag.Names() {
				var flagName, suffix string
				if len(name) == 1 {
					flagName = "-" + name
				} else {
					flagName = "--" + name
					if _, ok := flag.(*cli.StringFlag); ok {
						suffix = "="
					}
				}

				if contains(sharedExclusionList, flagName) {
					exclusionList = []string{"-", "*"}
				} else {
					exclusionList = append(exclusionList, flagName)
				}

				flagList = append(flagList, flagName+suffix)
			}

			fmt.Fprint(c.App.Writer, `'`)

			if len(exclusionList) > 1 {
				fmt.Fprint(c.App.Writer, `(`+strings.Join(exclusionList, " ")+`)`)
			}

			if len(flagList) > 1 {
				fmt.Fprint(c.App.Writer, `'{`+strings.Join(flagList, ",")+`}'`)
			} else {
				fmt.Fprint(c.App.Writer, flagList[0])
			}

			fmt.Fprint(c.App.Writer, "["+getFlagUsage(flag)+"]")

			if isFlagTakesFile(flag) {
				fmt.Fprint(c.App.Writer, ":file:_files")
			}

			fmt.Fprintln(c.App.Writer, `'`)
		}

		fmt.Fprintln(c.App.Writer, `
)

_arguments -S : "$options[@]" && ret=0
return ret`)

		return nil
	},
}

func contains(stack []string, needle string) bool {
	for _, s := range stack {
		if s == needle {
			return true
		}
	}
	return false
}
