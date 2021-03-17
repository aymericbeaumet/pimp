package prelude

import (
	"context"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"time"
)

type ExecRet struct {
	Pid    int    `json:"pid"`
	Status int    `json:"status"`
	Stdout string `json:"stdout"`
	Stderr string `json:"stderr"`
}

func (ret ExecRet) String() string {
	return ret.Stdout
}

func Exec(bin string, args ...string) (*ExecRet, error) {
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

	return &ExecRet{
		Pid:    state.Pid(),
		Status: state.ExitCode(),
		Stdout: string(outbytes),
		Stderr: string(errbytes),
	}, nil
}
