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
		c.emit(code.OpPop) // clean up expression statement from the stack, since it cannot be reused by anything else in future

	case *ast.InfixExpression:
		operator, left, right, err := inferInfixExpressionComponets(node)
		if err != nil {
			return err
		}

		err = c.Compile(left)
		if err != nil {
			return err
		}

		err = c.Compile(right)
		if err != nil {
			return err
		}

		err = c.compileInfixOperator(operator)
		if err != nil {
			return err
		}

	case *ast.IntegerLiteral:
		integer := &object.Integer{Value: node.Value}

		integerIndex := c.addConstant(integer)
		c.emit(code.OpConstant, integerIndex)

	case *ast.Boolean:
		if node.Value {
			c.emit(code.OpTrue)
		} else {
			c.emit(code.OpFalse)
		}
	}

	return nil
}

func inferInfixExpressionComponets(node *ast.InfixExpression) (string, ast.Expression, ast.Expression, error) {
	if node.Operator == token.LT {
		return token.GT, node.Right, node.Left, nil
	}

	return node.Operator, node.Left, node.Right, nil
}

func (c *Compiler) compileInfixOperator(operator string) error {
	switch operator {
	case token.PLUS:
		c.emit(code.OpAdd)
	case token.MINUS:
		c.emit(code.OpSub)
	case token.ASTERISK:
		c.emit(code.OpMul)
	case token.SLASH:
		c.emit(code.OpDiv)
	case token.EQ:
		c.emit(code.OpEqual)
	case token.NOT_EQ:
		c.emit(code.OpNotEqual)
	case token.GT:
		c.emit(code.OpGreaterThan)
	case token.AND:
		c.emit(code.OpAnd)
	case token.OR:
		c.emit(code.OpOr)
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
