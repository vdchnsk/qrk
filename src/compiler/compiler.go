package compiler

import (
	"fmt"

	"github.com/vdchnsk/qrk/src/ast"
	"github.com/vdchnsk/qrk/src/code"
	"github.com/vdchnsk/qrk/src/object"
	"github.com/vdchnsk/qrk/src/token"
)

type Compiler struct {
	instructions        code.Instructions
	constants           []object.Object
	lastInstruction     EmittedInstruction
	previousInstruction EmittedInstruction
}

// Its what we pass to VM
type Bytecode struct {
	Instructions code.Instructions
	Constants    []object.Object
}

type EmittedInstruction struct {
	Opcode   code.Opcode
	Position int
}

func NewCompiler() *Compiler {
	return &Compiler{
		instructions:        code.Instructions{},
		constants:           []object.Object{},
		lastInstruction:     EmittedInstruction{},
		previousInstruction: EmittedInstruction{},
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

	case *ast.PrefixExpression:
		err := c.Compile(node.Right)
		if err != nil {
			return err
		}

		err = c.compilePrefixOperator(node.Operator)
		if err != nil {
			return err
		}

	case *ast.InfixExpression:
		operator, left, right, err := inferInfixExpressionComponents(node)
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

	case *ast.IfExpression:
		err := c.Compile(node.Condition)
		if err != nil {
			return err
		}

		gotoNotTruthyPosition := c.emit(code.OpGotoNotTruthy, 9999)

		err = c.Compile(node.Consequence)
		if err != nil {
			return err
		}

		if c.isLastInstructionPop() {
			c.removeLastInstruction()
		}

		afterConsequencePosition := len(c.instructions)
		c.replaceOperand(gotoNotTruthyPosition, afterConsequencePosition)

	case *ast.BlockStatement:
		for _, statement := range node.Statements {
			err := c.Compile(statement)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func inferInfixExpressionComponents(node *ast.InfixExpression) (string, ast.Expression, ast.Expression, error) {
	if node.Operator == token.LT {
		return token.GT, node.Right, node.Left, nil
	}

	return node.Operator, node.Left, node.Right, nil
}

func (c *Compiler) compilePrefixOperator(operator string) error {
	switch operator {
	case token.BANG:
		c.emit(code.OpBang)
	case token.MINUS:
		c.emit(code.OpMinus)
	default:
		return fmt.Errorf("unknown prefix operator %s", operator)
	}
	return nil
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
		return fmt.Errorf("unknown infix operator %s", operator)
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

	c.setLastInstruction(opcode, position)

	return position
}

func (c *Compiler) removeLastInstruction() {
	c.instructions = c.instructions[:c.lastInstruction.Position]
	c.lastInstruction = c.previousInstruction
}

func (c *Compiler) replaceInstruction(position int, newInstruction []byte) {
	for i := 0; i < len(newInstruction); i++ {
		c.instructions[position+i] = newInstruction[i]
	}
}

// we can only replace operands of the same type, with the same non-variable length
func (c *Compiler) replaceOperand(replaceAt int, operand int) {
	opcode := code.Opcode(c.instructions[replaceAt])

	newInstruction := code.MakeInstruction(opcode, operand) // 15, 1 -> 15, 7

	c.replaceInstruction(replaceAt, newInstruction)
}

func (c *Compiler) setLastInstruction(opcode code.Opcode, position int) {
	previous := c.lastInstruction
	last := EmittedInstruction{Opcode: opcode, Position: position}

	c.previousInstruction = previous
	c.lastInstruction = last
}

func (c *Compiler) isLastInstructionPop() bool {
	return c.lastInstruction.Opcode == code.OpPop
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
