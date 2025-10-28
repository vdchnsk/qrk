package stdlib

import (
	"fmt"

	"github.com/vdchnsk/qrk/src/object"
)

func newError(format string, args ...any) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, args...)}
}
