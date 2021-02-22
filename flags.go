package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/mitchellh/go-homedir"
)

type Flags struct {
	Config  string
	Dump    bool
	DryRun  bool
	Help    bool
	Version bool
	Zsh     bool
}

func ParseFlagsArgs() (*Flags, []string, error) {
	var flags Flags
	flag.StringVar(&flags.Config, "config", "~/.pimprc", "Provide a different config file")
	flag.BoolVar(&flags.Dump, "dump", false, "Dump the config on stdout and exit with status code 0")
	flag.BoolVar(&flags.DryRun, "dry-run", false, "Print the command to be executed and exit with status code 0")
	flag.BoolVar(&flags.Help, "help", false, "Print the help and exit with status code 0")
	flag.BoolVar(&flags.Version, "version", false, "Print the version and exit with status code 0")
	flag.BoolVar(&flags.Zsh, "zsh", false, "Output Zsh config")

	var flagsSlice []string
	argsSlice := os.Args[1:]
	for i := 1; i < len(os.Args) && strings.HasPrefix(os.Args[i], "-"); i++ {
		flagsSlice = os.Args[1 : i+1]
		argsSlice = os.Args[i+1:]
		if os.Args[i] == "--" {
			break
		}
	}

	// Parse flags
	if err := flag.CommandLine.Parse(flagsSlice); err != nil {
		return nil, nil, err
	}

	// Expand paths
	config, err := homedir.Expand(flags.Config)
	if err != nil {
		return nil, nil, err
	}
	flags.Config = config

	return &flags, argsSlice, nil
}

func PrintUsage() {
	fmt.Printf("Usage: pimp [OPTION]... [--] CMD [ARG]...\n\nOptions:\n")
	flag.PrintDefaults()
}
