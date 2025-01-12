package code

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/vdchnsk/qrk/src/utils"
)

type Opcode byte

type Instructions []byte

func (instructions Instructions) fmt(opcodeDefinition *Definition, operands []int) string {
	operandsCount := len(opcodeDefinition.OperandWidths)

	if len(operands) != operandsCount {
		return fmt.Sprintf("ERROR: operand len %d does not match defined %d\n", len(operands), operandsCount)
	}

	switch operandsCount {
	case 0:
		return opcodeDefinition.Name
	case 1:
		return fmt.Sprintf("%s %d", opcodeDefinition.Name, operands[0])
	}

	return fmt.Sprintf("ERROR: unhandled operandCount for %s\n", opcodeDefinition.Name)
}

func (instructions Instructions) String() string {
	var out bytes.Buffer

	instructionByteIndex := 0

	for instructionByteIndex < len(instructions) {
		opcodeDefinition, err := LookupDefinition(instructions[instructionByteIndex])
		if err != nil {
			fmt.Fprintf(&out, "ERROR: %s", err)
			continue
		}

		operands, readBytes := ReadOperands(opcodeDefinition, instructions[instructionByteIndex+1:])
		formattedInstruction := instructions.fmt(opcodeDefinition, operands)

		fmt.Fprintf(&out, "%04d %s\n", instructionByteIndex, formattedInstruction)

		instructionByteIndex += 1 + readBytes
	}

	return out.String()
}

type Definition struct {
	Name          string
	OperandWidths []int
}

const (
	OpConstant Opcode = iota

	OpAdd
	OpSub
	OpMul
	OpDiv

	OpTrue
	OpFalse

	OpEqual
	OpNotEqual
	OpGreaterThan
	OpAnd
	OpOr

	OpMinus
	OpBang

	OpPop

	OpGotoNotTruthy // goto only if value of top of the stack is not truthy
	OpGoto

	OpNull

	OpGetGlobal
	OpSetGlobal

	OpArray
	OpHashMap
)

var operandDefinitions = map[Opcode]*Definition{
	OpConstant: {Name: "OpConstant", OperandWidths: []int{2}},

	OpAdd: {Name: "OpAdd", OperandWidths: []int{}},
	OpSub: {Name: "OpSub", OperandWidths: []int{}},
	OpMul: {Name: "OpMul", OperandWidths: []int{}},
	OpDiv: {Name: "OpDiv", OperandWidths: []int{}},

	OpTrue:  {Name: "OpTrue", OperandWidths: []int{}},
	OpFalse: {Name: "OpFalse", OperandWidths: []int{}},

	OpEqual:       {Name: "OpEqual", OperandWidths: []int{}},
	OpNotEqual:    {Name: "OpNotEqual", OperandWidths: []int{}},
	OpGreaterThan: {Name: "OpGreaterThan", OperandWidths: []int{}},
	OpAnd:         {Name: "OpAnd", OperandWidths: []int{}},
	OpOr:          {Name: "OpOr", OperandWidths: []int{}},

	OpMinus: {Name: "OpMinus", OperandWidths: []int{}},
	OpBang:  {Name: "OpBang", OperandWidths: []int{}},

	OpPop: {Name: "OpPop", OperandWidths: []int{}},

	OpGotoNotTruthy: {Name: "OpGotoNotTruthy", OperandWidths: []int{2}},
	OpGoto:          {Name: "OpGoto", OperandWidths: []int{2}},

	OpNull: {Name: "OpNull", OperandWidths: []int{}},

	OpGetGlobal: {Name: "OpGetGlobal", OperandWidths: []int{2}},
	OpSetGlobal: {Name: "OpSetGlobal", OperandWidths: []int{2}},

	OpArray:   {Name: "OpArray", OperandWidths: []int{2}},
	OpHashMap: {Name: "OpHashMap", OperandWidths: []int{2}},
}

func LookupDefinition(opcode byte) (*Definition, error) {
	def, ok := operandDefinitions[Opcode(opcode)]
	if !ok {
		return nil, fmt.Errorf("opcode %d undefined", opcode)
	}
	return def, nil
}

func MakeInstruction(op Opcode, operands ...int) []byte {
	operandDefinition, ok := operandDefinitions[op]
	if !ok {
		return []byte{}
	}

	if len(operands) != len(operandDefinition.OperandWidths) {
		return []byte{}
	}

	instructionSize := 1
	for _, operandWidth := range operandDefinition.OperandWidths {
		instructionSize += operandWidth
	}

	instruction := make([]byte, instructionSize)

	instruction[0] = byte(op)

	offset := 1

	for index, operand := range operands {
		operandWidth := operandDefinition.OperandWidths[index]

		switch operandWidth {
		case 1:
			instruction[offset] = byte(operand)
		case 2:
			binary.BigEndian.PutUint16(instruction[offset:], uint16(operand))
		case 4:
			binary.BigEndian.PutUint32(instruction[offset:], uint32(operand))
		case 8:
			binary.BigEndian.PutUint64(instruction[offset:], uint64(operand))
		}

		offset += operandWidth
	}

	return instruction
}

func ReadOperands(opcodeDefinition *Definition, instructionBytes []byte) ([]int, int) {
	operands := make([]int, len(opcodeDefinition.OperandWidths))

	offset := 0

	for i, width := range opcodeDefinition.OperandWidths {
		switch width {
		case 2:
			operands[i] = int(utils.ReadUint16(instructionBytes[offset:]))
		case 4:
			operands[i] = int(utils.ReadUint32(instructionBytes[offset:]))
		}
		offset += width
	}

	return operands, offset
}
