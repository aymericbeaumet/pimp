package command

import (
	"flag"
	"fmt"
	"reflect"
	"strings"
	"unsafe"

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
		eng, err := initializeEngine(c, true, true)
		if err != nil {
			return err
		}

		args := c.Args().Slice() // pimp [OPTION]... CMD [ARG]...
		lastArgIndex := len(args) - 1

		cmdargs, pendingFlag := skipFlags(c, args[:lastArgIndex]) // CMD [ARG]...
		cmdargs = append(cmdargs, args[lastArgIndex])

		// If a CMD is detected, delegate to the appropriate completion function
		if len(cmdargs) > 1 {
			_, expandedArgs, _ := eng.Map(nil, cmdargs)

			if len(expandedArgs[len(expandedArgs)-1]) != 0 && len(args[len(args)-1]) == 0 {
				expandedArgs = append(expandedArgs, "")
			}
			current := len(expandedArgs) // $CURRENT counts from 1, so len is the index of the last element

			fmt.Fprintln(c.App.Writer, "words=(")
			for _, word := range expandedArgs {
				fmt.Fprintf(c.App.Writer, "%#v\n", word)
			}
			fmt.Fprintln(c.App.Writer, ")")
			fmt.Fprintf(c.App.Writer, "(( CURRENT = %d ))\n", current)
			fmt.Fprintf(c.App.Writer, "_normal -p %#v\n", expandedArgs[0])

			return nil
		}

		// If a flag is currently pending and expecting a parameter
		if pendingFlag != nil {
			if isFlagTakesFile(pendingFlag) {
				fmt.Fprintln(c.App.Writer, "_files")
			}
			return nil
		}

		// If current arg is not an option, then offer to complete pimp commands +
		// path commands
		if !strings.HasPrefix(args[len(args)-1], "-") {
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

		// By default we print completion for the options

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
					} else if _, ok := flag.(*cli.StringSliceFlag); ok {
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
				if isFlagAllowedMultipleTimes(flag) {
					fmt.Fprint(c.App.Writer, `*`)
				}
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

func skipFlags(c *cli.Context, args []string) ([]string, cli.Flag) {
	// get the private flagSet field
	// https://stackoverflow.com/a/43918797/1071486
	flagSetField := reflect.ValueOf(c).Elem().FieldByName("flagSet")
	flagSetValue := reflect.NewAt(flagSetField.Type(), unsafe.Pointer(flagSetField.UnsafeAddr())).Elem()
	flagSet := flagSetValue.Interface().(*flag.FlagSet)

	if len(args) < 1 || args[0] != "pimp" {
		return nil, nil
	}

	if err := flagSet.Parse(args[1:]); err != nil {
		errStr := err.Error()
		if strings.HasPrefix(errStr, "flag needs an argument:") {
			split := strings.Split(errStr, "-")
			flagName := strings.TrimSpace(split[len(split)-1])

			for _, flag := range c.App.Flags {
				for _, name := range flag.Names() {
					if name == flagName {
						return nil, flag
					}
				}
			}
		}

		return nil, nil
	}

	return flagSet.Args(), nil
}

func contains(stack []string, needle string) bool {
	for _, s := range stack {
		if s == needle {
			return true
		}
	}
	return false
}
