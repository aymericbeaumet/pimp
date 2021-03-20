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

func initializeEngine(c *cli.Context, loadConfig, loadPimpfile bool) (*engine.Engine, error) {
	eng := engine.New()

	type source struct{ flagName, fallback string }
	sources := []*source{}
	if loadPimpfile {
		sources = append(sources, &source{
			flagName: "file",
			fallback: "./Pimpfile",
		})
	}
	if loadConfig {
		sources = append(sources, &source{
			flagName: "config",
			fallback: "~/.pimprc",
		})
	}

	for _, s := range sources {
		file, err := openWithFallback(c.String(s.flagName), s.fallback)
		if err != nil {
			return nil, err
		}
		if file != nil {
			defer file.Close()
			if err := eng.Append(file); err != nil {
				return nil, err
			}
		}
	}

	return eng, nil
}

func openWithFallback(filename, fallback string) (*os.File, error) {
	allowErrNotExist := false

	if len(filename) == 0 && len(fallback) > 0 {
		allowErrNotExist = true
		filename = fallback
	}

	filename, err := util.NormalizePath(filename)
	if err != nil {
		return nil, err
	}

	f, err := os.Open(filename)
	if err != nil && errors.Is(err, fs.ErrNotExist) && allowErrNotExist {
		return nil, nil
	}
	return f, err
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
	f, ok := flag.(*cli.StringFlag)
	return ok && f.TakesFile
}

func isFlagAllowedMultipleTimes(flag cli.Flag) bool {
	_, ok := flag.(*cli.StringSliceFlag)
	return ok
}
