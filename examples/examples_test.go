package examples_test

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
)

func TestExamples(t *testing.T) {
	_, filename, _, _ := runtime.Caller(0)
	dirname := filepath.Dir(filename)
	pattern := filepath.Join(dirname, "*")

	matches, err := filepath.Glob(pattern)
	if err != nil {
		t.Error(err)
	}

	for _, m := range matches {
		if s, err := os.Stat(m); err != nil {
			t.Error(err)
		} else if (s.Mode() & 0100) == 0 { // if non executable by the user
			continue
		}

		cmd := exec.Command(m)
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			t.Error(err)
		}

		if err := cmd.Start(); err != nil {
			t.Error(err)
		}

		out, err := io.ReadAll(stdout)
		if err != nil {
			t.Error(err)
		}

		if err := cmd.Wait(); err != nil {
			t.Error(err)
		}

		expected, err := os.ReadFile(m + ".expected")
		if err != nil {
			t.Error(err)
		}

		if !bytes.Equal(out, expected) {
			t.Errorf("example %s does not match the expected output", m)
		}
	}
}
