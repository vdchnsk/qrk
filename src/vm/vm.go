package vm

import (
	"fmt"

	"github.com/vdchnsk/qrk/src/code"
	"github.com/vdchnsk/qrk/src/compiler"
	"github.com/vdchnsk/qrk/src/object"
	"github.com/vdchnsk/qrk/src/utils"
)

const (
	StackSize      = 2048
	GlobalVarsSize = 65536
)

type VM struct {
	constants    []object.Object
	instructions code.Instructions

	stack   []object.Object
	globals []object.Object

	// Always points to the first free slot on the stack
	stackPointer int
}

var (
	True  = &object.Boolean{Value: true}
	False = &object.Boolean{Value: false}

	Null = &object.Null{}
)

func nativeToObjectBoolean(nativeValue bool) object.Object {
	if nativeValue {
		return True
	}
	return False
}

func NewVm(bytecode *compiler.Bytecode) *VM {
	stack := make([]object.Object, StackSize)
	globals := make([]object.Object, GlobalVarsSize)

	return &VM{
		instructions: bytecode.Instructions,
		constants:    bytecode.Constants,
		stack:        stack,
		globals:      globals,
		stackPointer: 0,
	}
}

func NewVmWithGlobalStore(bytecode *compiler.Bytecode, globals []object.Object) *VM {
	vm := NewVm(bytecode)
	vm.globals = globals

	return vm
}

func (vm *VM) Run() error {
	for instructionPointer := 0; instructionPointer < len(vm.instructions); instructionPointer++ {
		instructionByte := vm.instructions[instructionPointer]

		opcode := code.Opcode(instructionByte)

		switch opcode {
		case code.OpConstant:
			constantIndex := utils.ReadUint16(vm.instructions[instructionPointer+1:])

			def, err := code.LookupDefinition(instructionByte)
			if err != nil {
				return err
			}

			instructionPointer += def.OperandWidths[0]

			err = vm.stackPush(vm.constants[constantIndex])
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

		case code.OpGoto:
			instruction := vm.instructions[instructionPointer+1:]
			newPosOperand := int(utils.ReadUint16(instruction))

			instructionPointer = newPosOperand - 1

		case code.OpGotoNotTruthy:
			condition := vm.stackPop()

			if !isTruthy(condition) {
				instruction := vm.instructions[instructionPointer+1:]
				newPosOperand := int(utils.ReadUint16(instruction))

				instructionPointer = newPosOperand - 1
				continue
			}

			def, err := code.LookupDefinition(instructionByte)
			if err != nil {
				return err
			}

			instructionPointer += def.OperandWidths[0]

		case code.OpNull:
			err := vm.stackPush(Null)
			if err != nil {
				return err
			}

		case code.OpSetGlobal:
			instruction := vm.instructions[instructionPointer+1:]
			globalIndex := utils.ReadUint16(instruction)

			instructionPointer += 2
			value := vm.stackPop()
			vm.globals[globalIndex] = value

		case code.OpGetGlobal:
			instruction := vm.instructions[instructionPointer+1:]
			globalIndex := utils.ReadUint16(instruction)

			instructionPointer += 2

			value := vm.globals[globalIndex]
			err := vm.stackPush(value)
			if err != nil {
				return err
			}

		case code.OpArray:
			instruction := vm.instructions[instructionPointer+1:]
			arraySize := int(utils.ReadUint16(instruction))

			def, err := code.LookupDefinition(instructionByte)
			if err != nil {
				return err
			}

			instructionPointer += def.OperandWidths[0]

			start := vm.stackPointer - arraySize
			end := vm.stackPointer

			// override the array elements on stack with the array object itself
			array := vm.buildArray(start, end)
			vm.stackPointer -= arraySize

			err = vm.stackPush(array)
			if err != nil {
				return err
			}

		case code.OpHashMap:
			instruction := vm.instructions[instructionPointer+1:]
			hashmapSize := int(utils.ReadUint16(instruction))

			def, err := code.LookupDefinition(instructionByte)
			if err != nil {
				return err
			}

			instructionPointer += def.OperandWidths[0]

			start := vm.stackPointer - hashmapSize
			end := vm.stackPointer

			// override the hashmap elements on stack with the hashmap object itself
			hashmap, err := vm.buildHashmap(start, end)
			if err != nil {
				return err
			}
			vm.stackPointer -= hashmapSize

			err = vm.stackPush(hashmap)
			if err != nil {
				return err
			}

		case code.OpIndex:
			err := vm.executeIndexExpression()
			if err != nil {
				return nil
			}

		case code.OpCall:
			fn, ok := vm.stackPop().(*object.CompiledFunction)
			if !ok {
				return fmt.Errorf("calling a non-function object: %s", fn.Type())
			}

			return nil

		case code.OpPop:
			vm.stackPop()
		}
	}

	return nil
}

func isTruthy(obj object.Object) bool {
	switch obj := obj.(type) {
	case *object.Boolean:
		return obj.Value

	case *object.Null:
		return false

	default:
		return true
	}
}

func (vm *VM) executeBinaryOperation(op code.Opcode) error {
	right := vm.stackPop()
	left := vm.stackPop()

	rightType := right.Type()
	leftType := left.Type()

	if leftType == object.INTEGER_OBJ && rightType == object.INTEGER_OBJ {
		return vm.executeBinaryIntOperation(op, left, right)
	}

	if leftType == object.STRING_OBJ && rightType == object.STRING_OBJ {
		return vm.executeBinaryStringOperation(op, left, right)
	}

	return fmt.Errorf("unsupported type for binary operation: %s %s", leftType, rightType)
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
		return fmt.Errorf("unsupported type for binary operation: %s %s", leftType, rightType)
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

func (vm *VM) executeBinaryStringOperation(op code.Opcode, left, right object.Object) error {
	leftValue := left.(*object.String).Value
	rightValue := right.(*object.String).Value

	var result string

	switch op {
	case code.OpAdd:
		result = leftValue + rightValue
	default:
		return fmt.Errorf("unknown integer operator %d", op)
	}

	vm.stackPush(&object.String{Value: result})

	return nil
}

func (vm *VM) executeBangOperation() error {
	operand := vm.stackPop()

	switch operand {
	case True:
		return vm.stackPush(False)
	case False:
		return vm.stackPush(True)
	case Null:
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

func (vm *VM) executeIndexExpression() error {
	index := vm.stackPop()
	left := vm.stackPop()

	switch {
	case left.Type() == object.ARRAY_OBJ && index.Type() == object.INTEGER_OBJ:
		array := left.(*object.Array)
		idx := index.(*object.Integer)

		obj, err := vm.executeArrayIndex(array, idx)
		if err != nil {
			return err
		}

		return vm.stackPush(obj)

	case left.Type() == object.HASH_MAP_OBJ:
		hashMap := left.(*object.HashMap)
		key, ok := index.(object.Hashable)
		if !ok {
			return fmt.Errorf("unusable as hashmap key: %s", index.Type())
		}

		obj, err := vm.executeHashmapIndex(hashMap, key)
		if err != nil {
			return err
		}
		return vm.stackPush(obj)

	default:
		return fmt.Errorf("index operator not supported: %s", left.Type())
	}
}

func (vm *VM) executeArrayIndex(array *object.Array, index *object.Integer) (object.Object, error) {
	max := int64(len(array.Elements) - 1)

	if index.Value < 0 || index.Value > max {
		return Null, nil
	}

	return array.Elements[index.Value], nil
}

func (vm *VM) executeHashmapIndex(hashMap *object.HashMap, key object.Hashable) (object.Object, error) {
	pair, ok := hashMap.Pairs[key.HashKey()]
	if !ok {
		return Null, nil
	}

	return pair.Value, nil
}

func (vm *VM) StackTop() object.Object {
	if vm.stackPointer == 0 {
		return nil
	}

	topStackElem := vm.stack[vm.stackPointer-1]

	return topStackElem
}

// ========== !!! ===========
// === Use only for tests ===
// ========== !!! ===========
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

func (vm *VM) buildArray(startStackPointer, endStackPointer int) object.Object {
	size := endStackPointer - startStackPointer
	elements := make([]object.Object, size)

	for i := startStackPointer; i < endStackPointer; i++ {
		elements[i-startStackPointer] = vm.stack[i]
	}

	return &object.Array{Elements: elements}
}

func (vm *VM) buildHashmap(startStackPointer, endStackPointer int) (object.Object, error) {
	hashmap := &object.HashMap{
		Pairs: make(map[object.HashKey]object.HashPair),
	}

	for i := startStackPointer; i < endStackPointer; i += 2 {
		key := vm.stack[i]
		hashableKey, ok := key.(object.Hashable)
		if !ok {
			return nil, fmt.Errorf("unusable as hashmap key: %s", key.Type())
		}

		value := vm.stack[i+1]

		hashmap.Pairs[hashableKey.HashKey()] = object.HashPair{Key: key, Value: value}
	}

	return hashmap, nil
}
