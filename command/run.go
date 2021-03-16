package command

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/aymericbeaumet/pimp/normalize"
	"github.com/urfave/cli/v2"
)

var runCommand = &cli.Command{
	Name:      "--run",
	ArgsUsage: "FILE",
	Usage:     "Run the PimpScript FILE (- for stdin)",
	Action: func(c *cli.Context) error {
		args := c.Args().Slice()
		if len(args) != 1 {
			return fmt.Errorf("--run expects exactly one FILE, got %d", len(args))
		}

		var sb strings.Builder
		sb.WriteString(c.String("ldelim"))
		sb.WriteRune(' ')

		if args[0] == "-" {
			bytes, err := io.ReadAll(c.App.Reader)
			if err != nil {
				return err
			}
			sb.Write(bytes)
		} else {
			renderFilepath, err := normalize.Path(args[0])
			if err != nil {
				return err
			}
			data, err := os.ReadFile(renderFilepath)
			if err != nil {
				return err
			}
			sb.Write(data)
		}

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
