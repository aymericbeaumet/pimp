package command

import (
	"github.com/urfave/cli/v2"
)

var shellCommand = &cli.Command{
	Name:  "--shell",
	Usage: "Print the shell config (aliases only)",
	Action: func(c *cli.Context) error {
		_, eng, err := initializeConfigEngine(c)
		if err != nil {
			return err
		}
		return printShellAliases(c, eng)
	},
}
