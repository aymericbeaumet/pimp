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

		for _, command := range eng.Commands() {
			fmt.Fprintf(c.App.Writer, "alias %#v=%#v\n", command, "pimp "+command)
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
		if args := c.Args().Slice(); len(args) >= 2 && !strings.HasPrefix(args[1], "-") { // pimp CMD [ARG]...
			eng, err := initializeEngine(c, true, true)
			if err != nil {
				return err
			}

			// todo: extract args of this pimp command via the flags parser

			if len(args) == 2 {
				if cmds := eng.Commands(); len(cmds) > 0 {
					fmt.Fprintln(c.App.Writer, "local -a pcmds; pcmds=(")
					for _, cmd := range cmds {
						fmt.Fprintf(c.App.Writer, "%#v\n", cmd)
					}
					fmt.Fprintln(c.App.Writer, ")")
					fmt.Fprintln(c.App.Writer, "_describe 'pimp command' pcmds")
				}
				fmt.Fprintln(c.App.Writer, "_path_commands")
				return nil
			}

			_, expandedArgs, _ := eng.Map(nil, args[1:])
			if len(expandedArgs[len(expandedArgs)-1]) != 0 && len(args[len(args)-1]) == 0 {
				expandedArgs = append(expandedArgs, "")
			}
			current := len(expandedArgs) // index starts at 1
			fmt.Fprintln(c.App.Writer, "words=(")
			for _, word := range expandedArgs {
				fmt.Fprintf(c.App.Writer, "%#v\n", word)
			}
			fmt.Fprintln(c.App.Writer, ")")
			fmt.Fprintf(c.App.Writer, "(( CURRENT = %d ))\n", current)
			fmt.Fprintf(c.App.Writer, "_normal -p %#v\n", expandedArgs[0])

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
