package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

// Zero-values should represent the default values
type Flags struct {
	DryRun  bool
	Help    bool
	Version bool
	Zsh     bool
}

func ParseFlagsArgs() (Flags, []string, error) {
	var flags Flags
	flag.BoolVar(&flags.DryRun, "dry-run", false, "Print the command and exit with status code 0")
	flag.BoolVar(&flags.Help, "help", false, "Print the help and exit with status code 0")
	flag.BoolVar(&flags.Version, "version", false, "Print the version and exit with status code 0")
	flag.BoolVar(&flags.Zsh, "zsh", false, "Output Zsh config")

	if len(os.Args) < 2 {
		return flags, nil, nil
	}

	firstFlag := -1
	lastFlag := firstFlag
	for i := 1; i < len(os.Args) && strings.HasPrefix(os.Args[i], "-"); i++ {
		if firstFlag < 0 {
			firstFlag = i
		}
		lastFlag = i
	}

	var args = os.Args[1:]
	if firstFlag >= 0 && lastFlag >= 0 {
		if err := flag.CommandLine.Parse(os.Args[firstFlag : lastFlag+1]); err != nil {
			return flags, nil, err
		}
		args = os.Args[lastFlag+1:]
	}

	return flags, args, nil
}

func PrintUsage() {
	fmt.Printf("Usage: pimp [OPTION]... CMD [ARG]...\n\nOptions:\n")
	flag.PrintDefaults()
}
