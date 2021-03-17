package prelude

import "github.com/aymericbeaumet/pimp/pkg/errors"

func Exit(code int, args ...interface{}) (interface{}, error) {
	var message string
	if len(args) > 0 {
		message = ToString(args[0])
	}

	return nil, errors.NewFatalError(code, message)
}
