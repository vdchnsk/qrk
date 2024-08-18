package compiler

import (
	"fmt"

	"github.com/vdchnsk/qrk/src/ast"
	"github.com/vdchnsk/qrk/src/code"
	"github.com/vdchnsk/qrk/src/object"
	"github.com/vdchnsk/qrk/src/token"
)

type Compiler struct {
	instructions code.Instructions
	constants    []object.Object
}

// Its what we pass to VM
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

func (c *Compiler) Compile(node ast.Node) error {
	switch node := node.(type) {
	case *ast.Program:
		for _, statement := range node.Statements {
			err := c.Compile(statement)
			if err != nil {
				return err
			}
		}

	case *ast.ExpressionStatement:
		err := c.Compile(node.Value)
		if err != nil {
			return err
		}

	case *ast.InfixExpression:
		err := c.Compile(node.Left)
		if err != nil {
			return err
		}

		err = c.Compile(node.Right)
		if err != nil {
			return err
		}

		err = c.compileInfixOperator(node.Operator)
		if err != nil {
			return err
		}

	case *ast.IntegerLiteral:
		integer := &object.Integer{Value: node.Value}

		integerIndex := c.addConstant(integer)
		c.emit(code.OpConstant, integerIndex)
	}

	return nil
}

func (c *Compiler) compileInfixOperator(operator string) error {
	switch operator {
	case token.PLUS:
		c.emit(code.OpAdd)
	default:
		return fmt.Errorf("unknown operator %s", operator)
	}
	return nil
}

func (c *Compiler) addConstant(obj object.Object) int {
	c.constants = append(c.constants, obj)
	index := len(c.constants) - 1

	return index
}

func (c *Compiler) emit(opcode code.Opcode, operands ...int) int {
	instruction := code.MakeInstruction(opcode, operands...)
	position := c.addInstruction(instruction)

	return position
}

func (c *Compiler) addInstruction(instruction []byte) int {
	positionNewInstruction := len(c.instructions)
	c.instructions = append(c.instructions, instruction...)

	return positionNewInstruction
}

func (c *Compiler) Bytecode() *Bytecode {
	return &Bytecode{
		Instructions: c.instructions,
		Constants:    c.constants,
	}
}
