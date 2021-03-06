package misc

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"os/signal"
	"reflect"
	"strings"
	"text/template"

	fmerrors "github.com/aymericbeaumet/pimp/funcmap/errors"
)

func FuncMap() template.FuncMap {
	return template.FuncMap{
		"FZF": func(input interface{}) (string, error) {
			var s string
			switch i := input.(type) {
			case string:
				s = i
			case []string:
				s = strings.Join(i, "\n")
			default:
				return "", fmt.Errorf("unsupported type %v", reflect.TypeOf(input))
			}

			cmd := exec.CommandContext(context.Background(), "fzf")

			stdin, err := cmd.StdinPipe()
			if err != nil {
				return "", err
			}

			stdout, err := cmd.StdoutPipe()
			if err != nil {
				return "", err
			}

			cmd.Stderr = os.Stderr

			signalC := make(chan os.Signal, 32)
			signal.Notify(signalC)

			if err := cmd.Start(); err != nil {
				return "", err
			}

			go func() {
				for signal := range signalC {
					_ = cmd.Process.Signal(signal)
				}
			}()

			if _, err := stdin.Write([]byte(s)); err != nil {
				return "", err
			}
			if _, err := stdin.Write([]byte{'\n'}); err != nil {
				return "", err
			}

			state, err := cmd.Process.Wait()
			if err != nil {
				return "", err
			}

			if ec := state.ExitCode(); ec != 0 {
				return "", fmerrors.NewFatalError(ec, "")
			}

			out, err := ioutil.ReadAll(stdout)
			if err != nil {
				return "", err
			}

			return strings.TrimSpace(string(out)), nil
		},

		"Head": func(input []string) (string, error) {
			if len(input) == 0 {
				return "", errors.New("empty input")
			}
			return input[0], nil
		},

		"Tail": func(input []string) (string, error) {
			if len(input) == 0 {
				return "", errors.New("empty input")
			}
			return input[len(input)-1], nil
		},
	}
}
