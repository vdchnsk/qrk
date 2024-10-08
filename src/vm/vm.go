package vm

import (
	"fmt"

	"github.com/vdchnsk/qrk/src/code"
	"github.com/vdchnsk/qrk/src/compiler"
	"github.com/vdchnsk/qrk/src/object"
	"github.com/vdchnsk/qrk/src/utils"
)

const StackSize = 2048

type VM struct {
	constants    []object.Object
	instructions code.Instructions

	stack []object.Object

	// Always points to the first free slot on the stack
	stackPointer int
}

var (
	True  = &object.Boolean{Value: true}
	False = &object.Boolean{Value: false}
)

func nativeToObjectBoolean(nativeValue bool) object.Object {
	if nativeValue {
		return True
	}
	return False
}

func NewVm(bytecode *compiler.Bytecode) *VM {
	stack := make([]object.Object, StackSize)

	return &VM{
		instructions: bytecode.Instructions,
		constants:    bytecode.Constants,

		stack:        stack,
		stackPointer: 0,
	}
}

func (vm *VM) Run() error {
	for instructionPointer := 0; instructionPointer < len(vm.instructions); instructionPointer++ {
		instructionByte := vm.instructions[instructionPointer]

		opcode := code.Opcode(instructionByte)

		switch opcode {
		case code.OpConstant:
			constantIndex := utils.ReadUint16(vm.instructions[instructionPointer+1:])
			readBytes := 2
			instructionPointer += readBytes

			err := vm.stackPush(vm.constants[constantIndex])

			if err != nil {
				return err
			}

		case code.OpAdd, code.OpDiv, code.OpMul, code.OpSub:
			err := vm.executeBinaryOperation(opcode)
			if err != nil {
				return err
			}

		case code.OpTrue:
			err := vm.stackPush(True)
			if err != nil {
				return err
			}

		case code.OpFalse:
			err := vm.stackPush(False)
			if err != nil {
				return err
			}

		case code.OpEqual, code.OpNotEqual, code.OpGreaterThan, code.OpAnd, code.OpOr:
			err := vm.executeComparisonOperation(opcode)
			if err != nil {
				return nil
			}

		case code.OpBang:
			err := vm.executeBangOperation()
			if err != nil {
				return nil
			}

		case code.OpMinus:
			err := vm.executeMinusOperation()
			if err != nil {
				return nil
			}

		case code.OpPop:
			vm.stackPop()
		}
	}

	return nil
}

func (vm *VM) executeBinaryOperation(op code.Opcode) error {
	right := vm.stackPop()
	left := vm.stackPop()

	rightType := right.Type()
	leftType := left.Type()

	if leftType == object.INTEGER_OBJ && rightType == object.INTEGER_OBJ {
		return vm.executeBinaryIntOperation(op, left, right)
	}

	return fmt.Errorf("unsuppoerted type for binary operation: %s %s", leftType, rightType)
}

func (vm *VM) executeComparisonOperation(op code.Opcode) error {
	right := vm.stackPop()
	left := vm.stackPop()

	rightType := right.Type()
	leftType := left.Type()

	if leftType == object.INTEGER_OBJ && rightType == object.INTEGER_OBJ {

		return vm.executeIntegerComparison(op, left, right)
	}

	switch op {
	case code.OpEqual:
		return vm.stackPush(nativeToObjectBoolean(right == left))
	case code.OpNotEqual:
		return vm.stackPush(nativeToObjectBoolean(right != left))
	case code.OpAnd:
		return vm.stackPush(nativeToObjectBoolean(left.(*object.Boolean).Value && right.(*object.Boolean).Value))
	case code.OpOr:
		return vm.stackPush(nativeToObjectBoolean(left.(*object.Boolean).Value || right.(*object.Boolean).Value))
	default:
		return fmt.Errorf("unsuppoerted type for binary operation: %s %s", leftType, rightType)
	}

}

func (vm *VM) executeIntegerComparison(opcode code.Opcode, left, right object.Object) error {
	rightValue := right.(*object.Integer).Value
	leftValue := left.(*object.Integer).Value

	switch opcode {
	case code.OpEqual:
		return vm.stackPush(nativeToObjectBoolean(leftValue == rightValue))
	case code.OpNotEqual:
		return vm.stackPush(nativeToObjectBoolean(leftValue != rightValue))
	case code.OpGreaterThan:
		return vm.stackPush(nativeToObjectBoolean(leftValue > rightValue))
	default:
		return fmt.Errorf("unknown integer operator %d", opcode)
	}

}

func (vm *VM) executeBinaryIntOperation(op code.Opcode, left, right object.Object) error {
	leftValue := left.(*object.Integer).Value
	rightValue := right.(*object.Integer).Value

	var result int64

	switch op {
	case code.OpAdd:
		result = leftValue + rightValue
	case code.OpSub:
		result = leftValue - rightValue
	case code.OpDiv:
		result = leftValue / rightValue
	case code.OpMul:
		result = leftValue * rightValue
	default:
		return fmt.Errorf("unknown integer operator %d", op)
	}

	vm.stackPush(&object.Integer{Value: result})

	return nil
}

func (vm *VM) executeBangOperation() error {
	operand := vm.stackPop()

	switch operand {
	case True:
		return vm.stackPush(False)
	case False:
		return vm.stackPush(True)
	default:
		return vm.stackPush(False)
	}
}

func (vm *VM) executeMinusOperation() error {
	operand := vm.stackPop()

	if operand.Type() != object.INTEGER_OBJ {
		return fmt.Errorf("unsupported type for minus operator: %s", operand.Type())
	}

	currentValue := operand.(*object.Integer).Value
	oppositeInt := &object.Integer{Value: -currentValue}

	return vm.stackPush(oppositeInt)
}

func (vm *VM) StackTop() object.Object {
	if vm.stackPointer == 0 {
		return nil
	}

	topStackElem := vm.stack[vm.stackPointer-1]

	return topStackElem
}

// Use only for tests
func (vm *VM) LastPoppedStackElem() object.Object {
	// since, when popping we only decrement the pointer and not actually overriding top of the stack with nil
	// we can assume that after performing a pop() operation stack pointer will be pointing to the element that was just popped
	return vm.stack[vm.stackPointer]
}

func (vm *VM) stackPush(elem object.Object) error {
	isStackOverflow := vm.stackPointer >= StackSize
	if isStackOverflow {
		return fmt.Errorf("stack overflow")
	}

	firstStackFreeSlot := vm.stackPointer

	vm.stack[firstStackFreeSlot] = elem
	vm.stackPointer++

	return nil
}

func (vm *VM) stackPop() object.Object {
	topStackElem := vm.stack[vm.stackPointer-1]

	vm.stackPointer--

	return topStackElem
}
