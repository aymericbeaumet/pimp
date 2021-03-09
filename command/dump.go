package command

import "github.com/urfave/cli/v2"

var dumpCommand = &cli.Command{
	Name:            "--dump",
	Usage:           "Dump the engine as JSON and exit",
	SkipFlagParsing: true,
	Action: func(c *cli.Context) error {
		eng, err := initializeEngine(c)
		if err != nil {
			return err
		}
		return eng.JSON(c.App.Writer)
	},
}