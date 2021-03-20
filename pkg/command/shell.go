package command

import (
	"github.com/urfave/cli/v2"
)

var shellCommand = &cli.Command{
	Name:  "--shell",
	Usage: "Print the shell config (aliases only)",
	Action: func(c *cli.Context) error {
		eng, err := initializeEngine(c)
		if err != nil {
			return err
		}

		_, err = printAliases(c, eng)
		return err
	},
}
