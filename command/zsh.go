package command

import (
	"fmt"
	"strings"

	"github.com/urfave/cli/v2"
)

var zshCommand = &cli.Command{
	Name:            "--zsh",
	Usage:           "Print the Zsh config and exit (aliases and completion)",
	SkipFlagParsing: true,
	Action: func(c *cli.Context) error {
		eng, err := initializeEngine(c)
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
var zshCompletionCommand = &cli.Command{
	Name:            "--zsh-completion",
	Hidden:          true,
	SkipFlagParsing: true,
	Action: func(c *cli.Context) error {
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

_arguments -S "$options[@]" && ret=0
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
