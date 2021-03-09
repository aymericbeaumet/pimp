package command

import (
	"os"
	"strings"

	"github.com/aymericbeaumet/pimp/normalize"
	"github.com/urfave/cli/v2"
)

var renderCommand = &cli.Command{
	Name:            "--render",
	Usage:           "Render the template file(s) and exit",
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

		for _, renderFilepath := range os.Args[idx:] {
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
