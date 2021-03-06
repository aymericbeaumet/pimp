package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mitchellh/go-homedir"
)

type Flags struct {
	Config  string
	Dump    bool
	Expand  bool
	Help    bool
	Input   string
	Output  string
	Render  string
	Shell   bool
	Version bool
}

func ParseFlagsArgs() (*Flags, []string, error) {
	var flags Flags
	flag.StringVar(&flags.Config, "config", "~/.pimprc", "Provide a different config file")
	flag.BoolVar(&flags.Dump, "dump", false, "Dump the config on stdout and exit with status code 0")
	flag.BoolVar(&flags.Expand, "expand", false, "Expand the command and exit with status code 0")
	flag.BoolVar(&flags.Help, "help", false, "Print the help and exit with status code 0")
	flag.StringVar(&flags.Input, "input", "", "Read from the input file instead of stdin")
	flag.StringVar(&flags.Output, "output", "", "Write the output to this file instead of stdout")
	flag.StringVar(&flags.Render, "render", "", "Run the template and print to stdout")
	flag.BoolVar(&flags.Shell, "shell", false, "Output shell config (bash, zsh, fish, ...)")
	flag.BoolVar(&flags.Version, "version", false, "Print the version and exit with status code 0")

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

	// Expand config file
	config, err := expand(flags.Config)
	if err != nil {
		return nil, nil, err
	}
	flags.Config = config

	// Expand input path
	input, err := expand(flags.Input)
	if err != nil {
		return nil, nil, err
	}
	flags.Input = input

	// Expand output path
	output, err := expand(flags.Output)
	if err != nil {
		return nil, nil, err
	}
	flags.Output = output

	// Expand render path
	render, err := expand(flags.Render)
	if err != nil {
		return nil, nil, err
	}
	flags.Render = render

	return &flags, argsSlice, nil
}

func PrintUsage() {
	fmt.Printf("Usage: pimp [OPTION]... [--] CMD [ARG]...\n\nOptions:\n")
	flag.PrintDefaults()
}

func expand(input string) (string, error) {
	if len(input) == 0 {
		return "", nil
	}

	expanded, err := homedir.Expand(input)
	if err != nil {
		return "", err
	}
	if strings.HasPrefix(expanded, "/") {
		return expanded, nil
	}

	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return filepath.Join(wd, expanded), nil
}
