package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"os/signal"
	"sort"
	"strings"
	"text/template"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

var FuncMap = template.FuncMap{
	"FZF": func(values ...interface{}) string {
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
			if _, err := stdin.Write([]byte(getString(values...))); err != nil {
				panic(err)
			}
			if _, err := stdin.Write([]byte{'\n'}); err != nil {
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

	"GitBranches": func(values ...interface{}) []string {
		path, err := os.Getwd()
		if err != nil {
			panic(err)
		}

		repo, err := git.PlainOpen(path)
		if err != nil {
			panic(err)
		}

		out := []string{}

		iter, err := repo.Branches()
		if err != nil {
			panic(err)
		}
		if err := iter.ForEach(func(branch *plumbing.Reference) error {
			out = append(out, branch.Name().Short())
			return nil
		}); err != nil {
			panic(err)
		}

		sort.Strings(out)

		return out
	},

	"Head": func(values ...interface{}) string {
		rows := getArray(values...)
		return rows[0]
	},

	"JSON": func(values ...interface{}) string {
		out, err := json.Marshal(get(values...))
		if err != nil {
			panic(err)
		}
		return string(out)
	},

	"Tail": func(values ...interface{}) string {
		rows := getArray(values...)
		return rows[len(rows)-1]
	},
}

func get(values ...interface{}) interface{} {
	if len(values) != 1 {
		panic("expect exactly one arg")
	}
	return values[0]
}

func getArray(values ...interface{}) []string {
	switch value := get(values...).(type) {
	case string:
		return strings.Split(value, "\n")
	case []string:
		return value
	default:
		return strings.Split(fmt.Sprintf("%s", value), "\n")
	}
}

func getString(values ...interface{}) string {
	switch value := get(values...).(type) {
	case string:
		return value
	case []string:
		return strings.Join(value, "\n")
	default:
		return fmt.Sprintf("%s", value)
	}
}
