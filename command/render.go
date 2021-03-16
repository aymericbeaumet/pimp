package command

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/aymericbeaumet/pimp/normalize"
	"github.com/urfave/cli/v2"
)

var renderCommand = &cli.Command{
	Name:      "--render",
	ArgsUsage: "FILE",
	Usage:     "Render the template FILE (- for stdin)",
	Action: func(c *cli.Context) error {
		args := c.Args().Slice()
		if len(args) != 1 {
			return fmt.Errorf("--render expects exactly one FILE, got %d", len(args))
		}

		var text string

		if args[0] == "-" {
			bytes, err := io.ReadAll(c.App.Reader)
			if err != nil {
				return err
			}
			text = string(bytes)
		} else {
			renderFilepath, err := normalize.Path(args[0])
			if err != nil {
				return err
			}

			data, err := os.ReadFile(renderFilepath)
			if err != nil {
				return err
			}

			text = string(data)
		}

		// strip shebang if found
		if strings.HasPrefix(text, "#!") {
			if newlineIndex := strings.IndexRune(text, '\n'); newlineIndex > -1 {
				text = text[newlineIndex+1:]
			} else {
				text = ""
			}
		}

		rendered, err := render(text, c.String("ldelim"), c.String("rdelim"))
		if err != nil {
			return err
		}

		if _, err := c.App.Writer.Write([]byte(rendered)); err != nil {
			return err
		}

		return nil
	},
}
