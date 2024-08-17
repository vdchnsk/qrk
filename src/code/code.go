package code

import (
	"encoding/binary"
	"fmt"
)

type Opcode byte

type Instructions []byte

type Definition struct {
	Name          string
	OperandWidths []int
}

const (
	OpConstant Opcode = iota
)

var operandDefinitions = map[Opcode]*Definition{
	OpConstant: {Name: "OpConstant", OperandWidths: []int{2}},
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
