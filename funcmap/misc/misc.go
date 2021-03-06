package misc

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"reflect"
	"strings"
	"text/template"
	"time"

	fmerrors "github.com/aymericbeaumet/pimp/funcmap/errors"
)

func FuncMap() template.FuncMap {
	return template.FuncMap{
		"At": func(selector string, input interface{}) (interface{}, error) {
			switch i := input.(type) {
			case map[string]interface{}:
				return i[selector], nil
			default:
				return nil, fmt.Errorf("don't know how to apply selector `%s` on input %#v", selector, input)
			}
		},

		"Exec": func(bin string, args ...string) (map[string]interface{}, error) {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			cmd := exec.CommandContext(ctx, bin, args...)

			stdout, err := cmd.StdoutPipe()
			if err != nil {
				return nil, err
			}

			stderr, err := cmd.StderrPipe()
			if err != nil {
				return nil, err
			}

			signalC := make(chan os.Signal, 32)
			signal.Notify(signalC)

			if err := cmd.Start(); err != nil {
				return nil, err
			}

			go func() {
				for signal := range signalC {
					_ = cmd.Process.Signal(signal)
				}
			}()

			state, err := cmd.Process.Wait()
			if err != nil {
				return nil, err
			}

			outbytes, err := io.ReadAll(stdout)
			if err != nil {
				return nil, err
			}

			errbytes, err := io.ReadAll(stderr)
			if err != nil {
				return nil, err
			}

			return map[string]interface{}{
				"pid":    state.Pid(),
				"status": state.ExitCode(),
				"stdout": string(outbytes),
				"stderr": string(errbytes),
			}, nil
		},

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

			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			cmd := exec.CommandContext(ctx, "fzf")

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

			out, err := io.ReadAll(stdout)
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
