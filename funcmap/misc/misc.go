package misc

import (
	"context"
	"errors"
	"io/ioutil"
	"os"
	"os/exec"
	"os/signal"
	"strings"

	fmerrors "github.com/aymericbeaumet/pimp/funcmap/errors"
)

func FZF(input string) (string, error) {
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

	if _, err := stdin.Write([]byte(input)); err != nil {
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
}

func Head(input []string) (string, error) {
	if len(input) == 0 {
		return "", errors.New("empty input")
	}
	return input[0], nil
}

func Tail(input []string) (string, error) {
	if len(input) == 0 {
		return "", errors.New("empty input")
	}
	return input[len(input)-1], nil
}
