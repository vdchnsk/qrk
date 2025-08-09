package vm

import (
	"github.com/vdchnsk/qrk/src/code"
	"github.com/vdchnsk/qrk/src/object"
)

type StackFrame struct {
	fn *object.CompiledFunction
	ip int
}

func NewStackFrame(fn *object.CompiledFunction) *StackFrame {
	return &StackFrame{
		fn: fn,
		ip: -1,
	}
}

func (sf *StackFrame) Instructions() code.Instructions {
	return sf.fn.Instructions
}
