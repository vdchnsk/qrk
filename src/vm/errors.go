package vm

import (
	"fmt"

	"github.com/vdchnsk/qrk/src/object"
)

var (
	ErrCallingNonFunction = func(got object.ObjectType) error {
		return fmt.Errorf("calling a non-function object: %s", got)
	}

	ErrWrongNumberOfArguments = func(expected, got int) error {
		return fmt.Errorf("wrong number of arguments, expected=%d, got=%d", expected, got)
	}
)
