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
		text := strings.Join(c.Args().Slice(), " ")

		out, err := render(text, c.String("ldelim"), c.String("rdelim"))
		if err != nil {
			return err
		}

		fmt.Fprintln(c.App.Writer, out)

		return nil
	},
}
