package command

import (
	"fmt"

	"github.com/aymericbeaumet/pimp/pkg/script"
	"github.com/urfave/cli/v2"
)

var evalCommand = &cli.Command{
	Name:      "--eval",
	ArgsUsage: "STRING",
	Usage:     "Evaluate the PimpScript STRING",
	Action: func(c *cli.Context) error {
		args := c.Args().Slice()
		if len(args) != 1 {
			return fmt.Errorf("--eval expects exactly one STRING, got %d", len(args))
		}

		return script.Execute(c.App.Writer, args[0], c.String("ldelim"), c.String("rdelim"), funcmap)
	},
}
