package command

import (
	"fmt"
	"strings"

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

		var sb strings.Builder
		sb.WriteString(c.String("ldelim"))
		sb.WriteRune(' ')
		sb.WriteString(args[0])
		sb.WriteRune(' ')
		sb.WriteString(c.String("rdelim"))

		rendered, err := render(sb.String(), c.String("ldelim"), c.String("rdelim"))
		if err != nil {
			return err
		}

		_, err = c.App.Writer.Write([]byte(rendered))
		return err
	},
}
