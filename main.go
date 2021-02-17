package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
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

	fmt.Println(name, args)

	cmd := exec.CommandContext(context.Background(), name, args...)

	stdin, err := cmd.StdinPipe()
	if err != nil {
		panic(err)
	}
	go io.Copy(stdin, os.Stdin)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		panic(err)
	}
	go io.Copy(os.Stdout, stdout)

	stderr, err := cmd.StderrPipe()
	if err != nil {
		panic(err)
	}
	go io.Copy(os.Stderr, stderr)

	if err := cmd.Start(); err != nil {
		panic(err)
	}

	state, err := cmd.Process.Wait()
	if err != nil {
		panic(err)
	}

	os.Exit(state.ExitCode())
}
