package command

import (
	"fmt"
	"os"
	"strings"

	"github.com/urfave/cli/v2"
)

var runCommand = &cli.Command{
	Name:            "--run",
	Usage:           "Run the ARGs as templates and exit",
	SkipFlagParsing: true,
	Action: func(c *cli.Context) error {
		idx := -1
		for i, arg := range os.Args {
			if arg == "--" {
				idx = i + 1
				break
			}
		}
		if idx == -1 {
			idx = len(os.Args)
			for i := len(os.Args) - 1; i >= 0 && !strings.HasPrefix(os.Args[i], "-"); i-- {
				idx = i
			}
		}

		out, err := render(strings.Join(os.Args[idx:], " "))
		if err != nil {
			return err
		}

		fmt.Fprintln(c.App.Writer, out)

		return nil
	},
}
