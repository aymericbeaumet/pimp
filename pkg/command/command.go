package command

import (
	"errors"
	"io/fs"
	"os"
	"strings"

	"github.com/aymericbeaumet/pimp/pkg/engine"
	"github.com/aymericbeaumet/pimp/pkg/funcs"
	"github.com/aymericbeaumet/pimp/pkg/util"
	"github.com/urfave/cli/v2"
)

var funcmap = funcs.FuncMap()

var Commands = []*cli.Command{
	dumpCommand,
	evalCommand,
	renderCommand,
	runCommand,
	shellCommand,
	transpileCommand,
	zshCommand,
	zshCompletionCommand,
}

func CommandsFlags() []cli.Flag {
	out := make([]cli.Flag, 0, len(Commands))
	for _, c := range Commands {
		out = append(out, &cli.BoolFlag{
			Name:   strings.TrimPrefix(c.Name, "--"),
			Value:  false,
			Hidden: true,
		})
	}
	return out
}

func initializeEngine(c *cli.Context) (*engine.Engine, error) {
	eng := engine.New()

	allowErrNotExist := false
	pimpfiles := c.StringSlice("file")
	if len(pimpfiles) == 0 {
		allowErrNotExist = true
		pimpfiles = []string{"./Pimpfile.go", "./Pimpfile"}
	}

	for _, pimpfile := range pimpfiles {
		normalized, err := util.NormalizePath(pimpfile)
		if err != nil {
			return nil, err
		}

		file, err := os.Open(normalized)
		if err != nil && !(errors.Is(err, fs.ErrNotExist) && allowErrNotExist) {
			return nil, err
		}

		if file != nil {
			defer file.Close()
			if err := eng.LoadPimpfile(file); err != nil {
				return nil, err
			}
		}
	}

	return eng, nil
}

func getFlagUsage(flag cli.Flag) string {
	switch flag := flag.(type) {
	case *cli.BoolFlag:
		return flag.Usage
	case *cli.StringFlag:
		return flag.Usage
	case *cli.StringSliceFlag:
		return flag.Usage
	default:
		return ""
	}
}

func isFlagTakesFile(flag cli.Flag) bool {
	switch f := flag.(type) {
	case *cli.StringFlag:
		return f.TakesFile
	case *cli.StringSliceFlag:
		return f.TakesFile
	default:
		return false
	}
}

func isFlagAllowedMultipleTimes(flag cli.Flag) bool {
	switch flag.(type) {
	case (*cli.StringSliceFlag):
		return true
	default:
		return false
	}
}
