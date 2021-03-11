package command

import (
	"fmt"
	"strings"

	"github.com/urfave/cli/v2"
)

var runCommand = &cli.Command{
	Name:      "--run",
	ArgsUsage: "[ARG]...",
	Usage:     "Render ARGS as a template",
	Action: func(c *cli.Context) error {
		out, err := render(strings.Join(c.Args().Slice(), " "))
		if err != nil {
			return err
		}

		fmt.Fprintln(c.App.Writer, out)

		return nil
	},
}
