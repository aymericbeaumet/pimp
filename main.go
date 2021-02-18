package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"runtime/debug"
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

	engine, err := NewEngineFromHomeConfig()
	if err != nil {
		panic(err)
	}

	env, args := engine.Map(os.Environ(), args)

	if flags.DryRun {
		for i, arg := range args {
			if i > 0 {
				os.Stdout.Write([]byte{' '})
			}
			if _, err := fmt.Fprintf(os.Stdout, "%#v", arg); err != nil {
				panic(err)
			}
		}
		os.Stdout.Write([]byte{'\n'})
		return
	}

	cmd := exec.CommandContext(context.Background(), args[0], args[1:]...)
	cmd.Env = env
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		panic(err)
	}

	signalC := make(chan os.Signal, 32)
	signal.Notify(signalC)
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
