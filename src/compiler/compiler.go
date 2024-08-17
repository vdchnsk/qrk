package compiler

import (
	"github.com/vdchnsk/qrk/src/ast"
	"github.com/vdchnsk/qrk/src/code"
	"github.com/vdchnsk/qrk/src/object"
)

type Compiler struct {
	instructions code.Instructions
	constants    []object.Object
}

// Is what we pass to VM
type Bytecode struct {
	Instructions code.Instructions
	Constants    []object.Object
}

func NewCompiler() *Compiler {
	return &Compiler{
		instructions: code.Instructions{},
		constants:    []object.Object{},
	}
}

func (c *Compiler) Compile(program *ast.Program) error {
	return nil
}

func (c *Compiler) Bytecode() *Bytecode {
	return &Bytecode{
		Instructions: c.instructions,
		Constants:    c.constants,
	}
}
