package errors

import "strings"

type FatalError struct {
	exitCode int
	message  string
}

func NewFatalError(exitCode int, message string) *FatalError {
	if !strings.HasSuffix(message, "\n") {
		message += "\n"
	}

	return &FatalError{
		exitCode: exitCode,
		message:  message,
	}
}

func (e FatalError) Error() string {
	return e.message
}

func (e FatalError) ExitCode() int {
	return e.exitCode
}
