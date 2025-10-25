package vm

import (
	"github.com/vdchnsk/qrk/src/code"
	"github.com/vdchnsk/qrk/src/object"
)

type StackFrame struct {
	fn          *object.CompiledFunction
	ip          int
	basePointer int
}

func NewStackFrame(fn *object.CompiledFunction, basePointer int) *StackFrame {
	return &StackFrame{
		fn:          fn,
		ip:          -1,
		basePointer: basePointer,
	}
}

func (sf *StackFrame) Instructions() code.Instructions {
	return sf.fn.Instructions
}
