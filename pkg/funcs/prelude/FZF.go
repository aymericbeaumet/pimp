package prelude

import (
	"context"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"time"

	perrors "github.com/aymericbeaumet/pimp/pkg/errors"
)

type FZFRet struct {
	Stdout string `json:"stdout"`
}

func (ret FZFRet) String() string {
	return ret.Stdout
}

func FZF(input interface{}) (*FZFRet, error) {
	s := ToString(input)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "fzf")

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	cmd.Stderr = os.Stderr

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

	if _, err := stdin.Write([]byte(s)); err != nil {
		return nil, err
	}
	if _, err := stdin.Write([]byte{'\n'}); err != nil {
		return nil, err
	}

	state, err := cmd.Process.Wait()
	if err != nil {
		return nil, err
	}

	if ec := state.ExitCode(); ec != 0 {
		return nil, perrors.NewFatalError(ec, "")
	}

	outbytes, err := io.ReadAll(stdout)
	if err != nil {
		return nil, err
	}

	return &FZFRet{
		Stdout: string(outbytes),
	}, nil
}
