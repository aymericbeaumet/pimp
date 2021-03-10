package command

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

var shellCommand = &cli.Command{
	Name:  "--shell",
	Usage: "Print the shell config (aliases only)",
	Action: func(c *cli.Context) error {
		eng, err := initializeEngine(c, true, false)
		if err != nil {
			return err
		}
		for _, executable := range eng.Executables() {
			if _, err := fmt.Fprintf(
				c.App.Writer, "alias %#v=%#v\n", executable, "pimp "+executable,
			); err != nil {
				return err
			}
		}
		return nil
	},
}
