package prelude

import (
	"io"
	"os"

	"github.com/mattn/go-isatty"
)

func Stdin() (string, error) {
	if isatty.IsTerminal(os.Stdin.Fd()) {
		return "", nil
	}
	bytes, err := io.ReadAll(os.Stdin)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
