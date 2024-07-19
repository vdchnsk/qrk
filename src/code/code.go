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
	OpConstant: {"OpConstant", []int{2}},
}

func LookupDefinition(opcode byte) (*Definition, error) {
	def, ok := operandDefinitions[Opcode(opcode)]
	if !ok {
		return nil, fmt.Errorf("opcode %d undefined", opcode)
	}
	return def, nil
}
func MakeInstruction(op Opcode, opernds ...int) []byte {
	operandDefinition, ok := operandDefinitions[op]
	if !ok {
		return []byte{}
	}

	bytesInInstruction := 1
	for _, operandWidth := range operandDefinition.OperandWidths {
		bytesInInstruction += operandWidth
	}

	instruction := make([]byte, bytesInInstruction)

	instruction[0] = byte(op)

	offset := 1

	for index, operand := range opernds {
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
