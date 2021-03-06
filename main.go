package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"runtime/debug"
	"strings"
	"text/template"

	"github.com/aymericbeaumet/pimp/funcmap"
	fmerrors "github.com/aymericbeaumet/pimp/funcmap/errors"
)

func init() {
	// Disable garbage collection for such a short lived program
	debug.SetGCPercent(-1)
}

func main() {
	flags, args, err := ParseFlagsArgs()
	if err != nil {
		panic(err)
	}

	switch {
	case flags.Help:
		PrintUsage()
		return
	case flags.Version:
		fmt.Println("0.0.1")
		return
	}

	input := os.Stdin
	if len(flags.Input) > 0 {
		f, err := os.Open(flags.Input)
		if err != nil {
			panic(err)
		}
		defer f.Close()
		input = f
	}

	output := os.Stdout
	if len(flags.Output) > 0 {
		f, err := os.OpenFile(flags.Output, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			panic(err)
		}
		defer f.Close()
		output = f
	}

	if len(flags.Render) > 0 {
		data, err := os.ReadFile(flags.Render)
		if err != nil {
			panic(err)
		}
		rendered, err := renderTemplate(string(data))
		if err != nil {
			panic(err)
		}
		if _, err := output.Write([]byte(rendered)); err != nil {
			panic(err)
		}
		return
	}

	engine, err := NewEngineFromFile(flags.Config)
	if err != nil {
		panic(err)
	}

	switch {
	case flags.Dump:
		if err := engine.Dump(os.Stdout); err != nil {
			panic(err)
		}
		return

	case flags.Shell:
		for _, executable := range engine.Executables() {
			fmt.Printf("alias %#v=%#v\n", executable, "pimp "+executable)
		}
		return

	default:
		env, args, files := engine.Map(os.Environ(), args)
		if len(args) == 0 {
			PrintUsage()
			return
		}

		for i, arg := range args {
			args[i], err = renderTemplate(arg)
			if err != nil {
				panic(err)
			}
		}

		if flags.Expand {
			for i, arg := range args {
				if i > 0 {
					fmt.Print(" ")
				}
				fmt.Printf("%#v", arg)
			}
			fmt.Print("\n")
			return
		}

		for name, data := range files {
			rendered, err := renderTemplate(data)
			if err != nil {
				panic(err)
			}
			if err := os.WriteFile(name, []byte(rendered), 0400); err != nil {
				panic(err)
			}
			defer os.Remove(name)
		}

		cmd := exec.CommandContext(context.Background(), args[0], args[1:]...)
		cmd.Env = env
		cmd.Stdin = input
		cmd.Stdout = output
		cmd.Stderr = os.Stderr

		signalC := make(chan os.Signal, 32)
		signal.Notify(signalC)

		if err := cmd.Start(); err != nil {
			panic(err)
		}

		go func() {
			for signal := range signalC {
				_ = cmd.Process.Signal(signal)
			}
		}()

		state, err := cmd.Process.Wait()
		if err != nil {
			panic(err)
		}

		os.Exit(state.ExitCode())
	}
}

var fm = funcmap.FuncMap()

func renderTemplate(text string) (string, error) {
	var sb strings.Builder

	t, err := template.New(text).Funcs(fm).Parse(text)
	if err != nil {
		return "", err
	}

	if err := t.Execute(&sb, nil); err != nil {
		if e, ok := err.(template.ExecError); ok {
			// TODO: wait for this issue to be fixed upstream so that Unwrap()
			// returns the actual error that was returned (probably in Go 1.17).
			// In the meantime we cannot access the underlying error to cleanly
			// write to stderr + exit with the proper status code, so we panic.
			// https://github.com/golang/go/issues/34201
			err = e.Unwrap()
		}
		switch e := err.(type) {
		case *fmerrors.FatalError:
			os.Stderr.WriteString(e.Error())
			os.Exit(e.ExitCode())
		default:
			return "", err
		}
	}

	return sb.String(), nil
}
