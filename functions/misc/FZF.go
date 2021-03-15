package misc

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"reflect"
	"strings"
	"time"

	perrors "github.com/aymericbeaumet/pimp/errors"
)

func FZF(input interface{}) (string, error) {
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
		return "", perrors.NewFatalError(ec, "")
	}

	out, err := io.ReadAll(stdout)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(out)), nil
}
