package command

import (
	"os"
	"strings"

	"github.com/aymericbeaumet/pimp/normalize"
	"github.com/urfave/cli/v2"
)

var renderCommand = &cli.Command{
	Name:  "--render",
	Usage: "Render the template(s) and exit",
	Action: func(c *cli.Context) error {
		var templates []string
		for i := len(os.Args) - 1; i >= 0; i-- {
			arg := os.Args[i]
			if strings.HasPrefix(arg, "-") {
				break
			}
			templates = append(templates, arg)
		}

		for i := len(templates) - 1; i >= 0; i-- {
			renderFilepath := templates[i]

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
