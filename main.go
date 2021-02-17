package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"sync"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: prox cmd [argument]...")
		return
	}

	name := os.Args[1]
	var args []string
	if len(os.Args) >= 3 {
		args = os.Args[2:]
	}

	name, args = remap(name, args)

	var wg sync.WaitGroup
	cmd := exec.CommandContext(context.Background(), name, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		panic(err)
	}

	state, err := cmd.Process.Wait()
	if err != nil {
		panic(err)
	}

	wg.Wait()
	os.Exit(state.ExitCode())
}

func remap(name string, args []string) (string, []string) {
	if len(args) == 0 {
		return "hub", []string{"status", "-sb"}
	}
	return "hub", args
}
