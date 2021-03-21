package command

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/aymericbeaumet/pimp/pkg/config"
	"github.com/aymericbeaumet/pimp/pkg/engine"
	"github.com/aymericbeaumet/pimp/pkg/funcs"
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

func initializeConfigEngine(c *cli.Context) (*config.Config, *engine.Engine, error) {
	eng := engine.New()

	// Load configuration
	conf, err := config.Load(c.String("config"))
	if err != nil {
		return nil, nil, err
	}
	defer conf.Close()

	// Load the local Pimpfiles
	pimpfiles := append([]string{}, c.StringSlice("file")...)
	if len(pimpfiles) == 0 { // if no Pimpfiles are defined, apply the default resolution mecanism
		p, err := resolvePimpfiles()
		if err != nil {
			return nil, nil, err
		}
		pimpfiles = p
	}
	for _, pimpfile := range pimpfiles {
		f, err := os.Open(pimpfile)
		if err != nil {
			return nil, nil, err
		}
		defer f.Close()
		if err := eng.LoadPimpfile(f, true); err != nil {
			return nil, nil, err
		}
	}

	// Load the global Pimpfiles (with lower priority)
	for _, pimpfile := range conf.Pimpfiles {
		if err := eng.LoadPimpfile(pimpfile, false); err != nil {
			return nil, nil, err
		}
	}

	return conf, eng, nil
}

var pimpfileCandidates = []string{"Pimpfile.go", "Pimpfile"}
var rootMarkers = []string{".git"}

func resolvePimpfiles() ([]string, error) {
	var out []string

	currentDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	for {
		// try to find as many pimpfiles as possible
		for _, candidate := range pimpfileCandidates {
			name := filepath.Join(currentDir, candidate)
			if s, err := os.Stat(name); err == nil && !s.IsDir() {
				out = append(out, name)
			}
		}

		// stop here if a root marker is found
		for _, rootMarker := range rootMarkers {
			dirname := filepath.Join(currentDir, rootMarker)
			if s, err := os.Stat(dirname); err == nil && s.IsDir() {
				return out, nil
			}
		}

		// we should not reach the root, discard everything found if that's the case
		if len(strings.TrimRight(currentDir, "/")) == 0 {
			break
		}

		// move up to the parent directory
		currentDir = filepath.Dir(currentDir)
	}

	// if the root directory as been reached (/), then just resolve candidates in
	// the cwd
	out = []string{}
	for _, candidate := range pimpfileCandidates {
		name := filepath.Join(currentDir, candidate)
		if s, err := os.Stat(name); err == nil && !s.IsDir() {
			out = append(out, name)
		}
	}
	return out, nil
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

func printShellAliases(c *cli.Context, eng *engine.Engine) error {
	if len(c.StringSlice("file")) > 0 {
		return errors.New("flag --file is not compatible with this command")
	}

	for _, command := range eng.Commands() {
		var alias string
		if config := c.String("config"); len(config) > 0 {
			alias = fmt.Sprintf("pimp --config=%#v %s", config, command)
		} else {
			alias = fmt.Sprintf("pimp %s", command)
		}
		fmt.Fprintf(c.App.Writer, "alias %#v=%#v\n", command, alias)
	}

	return nil
}
