package command

import (
	"os"
	"strings"

	"github.com/aymericbeaumet/pimp/normalize"
	"github.com/urfave/cli/v2"
)

var renderCommand = &cli.Command{
	Name:            "--render",
	Usage:           "Open and render ARGs files as templates and exit",
	SkipFlagParsing: true,
	Action: func(c *cli.Context) error {
		for _, renderFilepath := range c.Args().Slice() {
			renderFilepath, err := normalize.Path(renderFilepath)
			if err != nil {
				return err
			}

			data, err := os.ReadFile(renderFilepath)
			if err != nil {
				return err
			}

			s := string(data)

			// strip shebang if found
			if strings.HasPrefix(s, "#!") {
				if newlineIndex := strings.IndexRune(s, '\n'); newlineIndex > -1 {
					s = s[newlineIndex+1:]
				} else {
					s = ""
				}
			}

			rendered, err := render(s)
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
