package command

import (
	"encoding/json"

	"github.com/urfave/cli/v2"
)

var dumpCommand = &cli.Command{
	Name:  "--dump",
	Usage: "Dump the config and engine in JSON format",
	Action: func(c *cli.Context) error {
		conf, eng, err := initializeConfigEngine(c, true)
		if err != nil {
			return err
		}
		return json.NewEncoder(c.App.Writer).Encode(map[string]interface{}{
			"CONFIG": conf,
			"ENGINE": eng,
		})
	},
}
