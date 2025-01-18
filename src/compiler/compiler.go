package compiler

import (
	"fmt"

	"github.com/vdchnsk/qrk/src/ast"
	"github.com/vdchnsk/qrk/src/code"
	"github.com/vdchnsk/qrk/src/object"
	"github.com/vdchnsk/qrk/src/token"
	"github.com/vdchnsk/qrk/src/utils"
)

type Compiler struct {
	instructions        code.Instructions
	constants           []object.Object
	lastInstruction     EmittedInstruction
	previousInstruction EmittedInstruction
	symbolTable         *SymbolTable
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
		symbolTable:         NewSymbolTable(),
	}
}

func NewCompilerWithState(symbolTable *SymbolTable, constants []object.Object) *Compiler {
	compiler := NewCompiler()
	compiler.symbolTable = symbolTable
	compiler.constants = constants

	return compiler
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

	case *ast.StringLiteral:
		str := &object.String{Value: node.Value}

		stringIndex := c.addConstant(str)
		c.emit(code.OpConstant, stringIndex)

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

		gotoElseIns := c.emit(code.OpGotoNotTruthy, -1)

		err = c.compileBranch(node.Consequence)
		if err != nil {
			return err
		}

		// must be the last instruction in the `if`-block to jump over `else`
		skipElseIns := c.emit(code.OpGoto, -1)

		elseBlockStart := len(c.instructions)
		c.replaceOperand(gotoElseIns, elseBlockStart)

		if node.Alternative != nil {
			err := c.compileBranch(node.Alternative)
			if err != nil {
				return err
			}
		} else {
			c.emit(code.OpNull)
		}

		elseBlockEnd := len(c.instructions)
		c.replaceOperand(skipElseIns, elseBlockEnd)

	case *ast.BlockStatement:
		for _, statement := range node.Statements {
			err := c.Compile(statement)
			if err != nil {
				return err
			}
		}

	case *ast.LetStatement:
		err := c.Compile(node.Value)
		if err != nil {
			return err
		}

		symbol := c.symbolTable.Define(node.Identifier.Value)
		c.emit(code.OpSetGlobal, symbol.Index)

	case *ast.Identifier:
		symbol, ok := c.symbolTable.Resolve(node.Value)
		if !ok {
			return fmt.Errorf("undefined variable %s", node.Value)
		}

		c.emit(code.OpGetGlobal, symbol.Index)

	case *ast.ArrayLiteral:
		for _, element := range node.Elements {
			err := c.Compile(element)
			if err != nil {
				return err
			}
		}

		c.emit(code.OpArray, len(node.Elements))

	case *ast.HashMapLiteral:
		keys := make([]ast.Expression, 0, len(node.Pairs))
		for key := range node.Pairs {
			keys = append(keys, key)
		}
		utils.SortByString(keys)

		for _, key := range keys {
			if err := c.Compile(key); err != nil {
				return err
			}

			value := node.Pairs[key]
			if err := c.Compile(value); err != nil {
				return err
			}
		}

		c.emit(code.OpHashMap, len(node.Pairs)*2)

	case *ast.IndexExpression:
		if err := c.Compile(node.Left); err != nil {
			return err
		}

		if err := c.Compile(node.Index); err != nil {
			return err
		}

		c.emit(code.OpIndex)
	}

	return nil
}

func (c *Compiler) compileBranch(branch *ast.BlockStatement) error {
	err := c.Compile(branch)
	if err != nil {
		return nil
	}

	if c.isLastInstructionPop() {
		c.removeLastInstruction()
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
