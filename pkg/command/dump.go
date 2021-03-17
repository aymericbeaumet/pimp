package command

import "github.com/urfave/cli/v2"

var dumpCommand = &cli.Command{
	Name:  "--dump",
	Usage: "Dump the engine in JSON format",
	Action: func(c *cli.Context) error {
		eng, err := initializeEngine(c, true, true)
		if err != nil {
			return err
		}
		return eng.JSON(c.App.Writer)
	},
}
