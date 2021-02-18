package main

import (
	"errors"
	"flag"
	"os"
	"strings"
)

// Zero-values should represent the default values
type Flags struct {
	DryRun bool
}

func ParseFlagsArgs() (*Flags, []string, error) {
	if len(os.Args) < 2 {
		flag.PrintDefaults()
		return nil, nil, errors.New("usage")
	}

	var flags Flags
	flag.BoolVar(&flags.DryRun, "dry-run", false, "Print the command and exit with status code 0")

	args := os.Args[1:]
	for i, arg := range os.Args {
		if i == 0 {
			continue
		}

		if !strings.HasPrefix(arg, "-") {
			if err := flag.CommandLine.Parse(os.Args[1:i]); err != nil {
				return nil, nil, err
			}
			args = os.Args[i:]
			break
		}
	}

	return &flags, args, nil
}
