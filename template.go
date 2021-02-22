package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"text/template"
)

var FuncMap = template.FuncMap{
	"FZF": func(values ...interface{}) string {
		input := fmt.Sprintf("%s", values[len(values)-1])

		cmd := exec.CommandContext(context.Background(), "fzf")
		cmd.Stderr = os.Stderr

		stdin, err := cmd.StdinPipe()
		if err != nil {
			panic(err)
		}

		stdout, err := cmd.StdoutPipe()
		if err != nil {
			panic(err)
		}

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

		go func() {
			if _, err := stdin.Write([]byte(input)); err != nil {
				panic(err)
			}
		}()

		state, err := cmd.Process.Wait()
		if err != nil {
			panic(err)
		}

		if state.ExitCode() != 0 {
			panic("fzf failed")
		}

		out, err := ioutil.ReadAll(stdout)
		if err != nil {
			panic(err)
		}

		return strings.TrimSpace(string(out))
	},
}
