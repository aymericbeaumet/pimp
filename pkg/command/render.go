package command

import (
	"fmt"
	"io"
	"os"

	"github.com/aymericbeaumet/pimp/pkg/template"
	"github.com/aymericbeaumet/pimp/pkg/util"
	"github.com/urfave/cli/v2"
)

var renderCommand = &cli.Command{
	Name:      "--render",
	ArgsUsage: "[FILE]",
	Usage:     "Render the template FILE (use - or omit arg for stdin)",
	Action: func(c *cli.Context) error {
		args := c.Args().Slice()
		if len(args) > 1 {
			return fmt.Errorf("--render expects at most one FILE, got %d", len(args))
		}

		var text string
		if len(args) == 0 || args[0] == "-" {
			data, err := io.ReadAll(c.App.Reader)
			if err != nil {
				return err
			}
			text = string(data)
		} else {
			renderFilepath, err := util.NormalizePath(args[0])
			if err != nil {
				return err
			}
			data, err := os.ReadFile(renderFilepath)
			if err != nil {
				return err
			}
			text = string(data)
		}

		text = util.StripShebang(text)

		// TODO: type check errors.FatalError and exit (when merged https://github.com/golang/go/issues/34201)
		return template.Render(c.App.Writer, text, c.String("ldelim"), c.String("rdelim"), funcmap)
	},
}
