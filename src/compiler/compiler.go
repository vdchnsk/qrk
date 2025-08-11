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
	constants []object.Object

	symbolTable *SymbolTable

	scopes     []CompilationScope
	scopeIndex int
}

// Bytecode is the result of the compilation phase.
// It contains executable instructions and a constant pool.
// This is passed directly to the VM for execution.
type Bytecode struct {
	Instructions code.Instructions
	Constants    []object.Object
}

type EmittedInstruction struct {
	Opcode   code.Opcode
	Position int
}

type CompilationScope struct {
	instructions    code.Instructions
	lastInstruction EmittedInstruction
	prevInstruction EmittedInstruction
}

func New() *Compiler {
	mainScope := CompilationScope{
		instructions:    code.Instructions{},
		lastInstruction: EmittedInstruction{},
		prevInstruction: EmittedInstruction{},
	}

	return &Compiler{
		symbolTable: NewSymbolTable(),
		constants:   []object.Object{},
		scopes:      []CompilationScope{mainScope},
		scopeIndex:  0,
	}
}

func (c *Compiler) curScope() *CompilationScope {
	return &c.scopes[c.scopeIndex]
}

func (c *Compiler) curInstructions() code.Instructions {
	return c.curScope().instructions
}

func (c *Compiler) lastInstructionIs(op code.Opcode) bool {
	if len(c.curInstructions()) == 0 {
		return false
	}

	return c.curScope().lastInstruction.Opcode == op
}

func (c *Compiler) replaceLastPopWithReturn() {
	lastPos := c.curScope().lastInstruction.Position
	c.replaceInstruction(lastPos, code.MakeInstruction(code.OpReturnValue))

	c.curScope().lastInstruction.Opcode = code.OpReturnValue
}

func (c *Compiler) enterScope() {
	scope := CompilationScope{
		instructions:    code.Instructions{},
		lastInstruction: EmittedInstruction{},
		prevInstruction: EmittedInstruction{},
	}

	c.scopes = append(c.scopes, scope)
	c.scopeIndex++

	c.symbolTable = NewEnclosedSymbolTable(c.symbolTable)
}

func (c *Compiler) removeLastScope() {
	c.scopes = c.scopes[:len(c.scopes)-1]
	c.scopeIndex--
}

func (c *Compiler) leaveScope() code.Instructions {
	instructions := c.curInstructions()

	c.removeLastScope()

	c.symbolTable = c.symbolTable.Outer

	return instructions
}

func NewWithState(symbolTable *SymbolTable, constants []object.Object) *Compiler {
	compiler := New()
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
		if err := c.Compile(node.Value); err != nil {
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

		if err := c.Compile(left); err != nil {
			return err
		}

		if err := c.Compile(right); err != nil {
			return err
		}

		if err := c.compileInfixOperator(operator); err != nil {
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
		if err := c.Compile(node.Condition); err != nil {
			return err
		}

		gotoElseIns := c.emit(code.OpGotoNotTruthy, -1)

		if err := c.compileBranch(node.Consequence); err != nil {
			return err
		}

		// must be the last instruction in the `if`-block to jump over `else`
		skipElseIns := c.emit(code.OpGoto, -1)

		elseBlockStart := len(c.curInstructions())
		c.replaceOperand(gotoElseIns, elseBlockStart)

		if node.Alternative != nil {
			if err := c.compileBranch(node.Alternative); err != nil {
				return err
			}
		} else {
			c.emit(code.OpNull)
		}

		elseBlockEnd := len(c.curInstructions())
		c.replaceOperand(skipElseIns, elseBlockEnd)

	case *ast.BlockStatement:
		for _, statement := range node.Statements {
			if err := c.Compile(statement); err != nil {
				return err
			}
		}

	case *ast.LetStatement:
		if err := c.Compile(node.Value); err != nil {
			return err
		}

		symbol := c.symbolTable.Define(node.Identifier.Value)

		if symbol.Scope == GlobalScope {
			c.emit(code.OpSetGlobal, symbol.Index)
		} else {
			c.emit(code.OpSetLocal, symbol.Index)
		}

	case *ast.Identifier:
		symbol, ok := c.symbolTable.Resolve(node.Value)
		if !ok {
			return fmt.Errorf("undefined variable %s", node.Value)
		}

		if symbol.Scope == GlobalScope {
			c.emit(code.OpGetGlobal, symbol.Index)
		} else {
			c.emit(code.OpGetLocal, symbol.Index)
		}

	case *ast.ArrayLiteral:
		for _, element := range node.Elements {
			if err := c.Compile(element); err != nil {
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

	case *ast.FuncLiteral:
		c.enterScope()

		if err := c.Compile(node.Body); err != nil {
			return err
		}

		if c.lastInstructionIs(code.OpPop) {
			c.replaceLastPopWithReturn()
		}

		if !c.lastInstructionIs(code.OpReturnValue) {
			c.emit(code.OpReturn)
		}

		instructions := c.leaveScope()

		compiledFunc := &object.CompiledFunction{
			Instructions: instructions,
		}

		c.emit(code.OpConstant, c.addConstant(compiledFunc))

	case *ast.ReturnStatement:
		if err := c.Compile(node.Value); err != nil {
			return err
		}

		c.emit(code.OpReturnValue)

	case *ast.CallExpression:
		if err := c.Compile(node.Function); err != nil {
			return err
		}

		c.emit(code.OpCall)
	}

	return nil
}

func (c *Compiler) compileBranch(branch *ast.BlockStatement) error {
	if err := c.Compile(branch); err != nil {
		return err
	}

	if c.lastInstructionIs(code.OpPop) {
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
	c.curScope().instructions = c.curInstructions()[:c.curScope().lastInstruction.Position]
	c.curScope().lastInstruction = c.curScope().prevInstruction
}

func (c *Compiler) replaceInstruction(position int, newInstruction []byte) {
	for i := range len(newInstruction) {
		c.curInstructions()[position+i] = newInstruction[i]
	}
}

// we can only replace operands of the same type, with the same non-variable length
func (c *Compiler) replaceOperand(replaceAt int, operand int) {
	opcode := code.Opcode(c.curInstructions()[replaceAt])

	newInstruction := code.MakeInstruction(opcode, operand) // 15, 1 -> 15, 7

	c.replaceInstruction(replaceAt, newInstruction)
}

func (c *Compiler) setLastInstruction(opcode code.Opcode, position int) {
	previous := c.curScope().lastInstruction
	last := EmittedInstruction{Opcode: opcode, Position: position}

	c.curScope().prevInstruction = previous
	c.curScope().lastInstruction = last
}

func (c *Compiler) addInstruction(instruction []byte) int {
	positionNewInstruction := len(c.curInstructions())
	c.curScope().instructions = append(c.curInstructions(), instruction...)

	return positionNewInstruction
}

func (c *Compiler) Bytecode() *Bytecode {
	return &Bytecode{
		Instructions: c.curInstructions(),
		Constants:    c.constants,
	}
}
