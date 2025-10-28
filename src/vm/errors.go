package vm

import "fmt"

var (
	ErrWrongNumberOfArguments = func(expected, got int) error {
		return fmt.Errorf("wrong number of arguments, expected=%d, got=%d", expected, got)
	}
)
