package command

import (
	"fmt"
	"os"
	"path/filepath"
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

	var pimpfiles []string

	// First load the local Pimpfiles
	pimpfiles = append(pimpfiles, c.StringSlice("file")...)
	if len(pimpfiles) == 0 { // if no Pimpfiles are defined, apply the default resolution mecanism
		p, err := resolvePimpfiles()
		if err != nil {
			return nil, err
		}
		pimpfiles = p
	}

	// Add the global Pimpfiles (with lower priority)
	pimpfiles = append(pimpfiles, c.StringSlice("global")...)

	for _, pimpfile := range pimpfiles {
		normalized, err := util.NormalizePath(pimpfile)
		if err != nil {
			return nil, err
		}

		file, err := os.Open(normalized)
		if err != nil {
			return nil, err
		}
		defer file.Close()
		if err := eng.LoadPimpfile(file); err != nil {
			return nil, err
		}
	}

	return eng, nil
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

func printAliases(c *cli.Context, eng *engine.Engine) (string, error) {
	var flags strings.Builder

	for _, global := range c.StringSlice("global") {
		flags.Write([]byte(fmt.Sprintf(" --global %#v", global)))
	}

	for _, file := range c.StringSlice("file") {
		flags.Write([]byte(fmt.Sprintf(" -f %#v", file)))
	}

	for _, command := range eng.Commands() {
		fmt.Fprintf(c.App.Writer, "alias %#v=%#v\n", command, "pimp"+flags.String()+" "+command)
	}

	return flags.String(), nil
}
