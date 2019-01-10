package text

import (
	"fmt"
)

type ErrMissingArg string

func (arg ErrMissingArg) Error() string {
	return fmt.Sprintf("missing required arg: %s", string(arg))
}

type ErrMissingGroup []string

func (group ErrMissingGroup) Error() string {
	return fmt.Sprintf(
		"missing one of required arguments group: %s",
		[]string(group),
	)
}

type ErrIncorrectType struct {
	Arg   string
	Value interface{}
	Type  string
}

func (err ErrIncorrectType) Error() string {
	return fmt.Sprintf(
		"incorrect argument %q type: given %T, expected %s",
		err.Arg,
		err.Value,
		err.Type,
	)
}

type ErrUnknownArg string

func (arg ErrUnknownArg) Error() string {
	return fmt.Sprintf(
		"unknown arg: %s",
		string(arg),
	)
}
