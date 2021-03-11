package command

import (
	"io"
	"os"
	"strings"

	"github.com/aymericbeaumet/pimp/normalize"
	"github.com/urfave/cli/v2"
)

var renderCommand = &cli.Command{
	Name:      "--render",
	ArgsUsage: "[FILE]...",
	Usage:     "Sequentially open and render the template FILES (- for stdin)",
	Action: func(c *cli.Context) error {
		var readerCache []byte

		for _, arg := range c.Args().Slice() {
			var text string

			if arg == "-" {
				if readerCache == nil {
					bytes, err := io.ReadAll(c.App.Reader)
					if err != nil {
						return err
					}
					readerCache = bytes
				}
				text = string(readerCache)
			} else {
				renderFilepath, err := normalize.Path(arg)
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
		}

		return nil
	},
}
