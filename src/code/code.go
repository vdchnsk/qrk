package code

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/vdchnsk/qrk/src/utils"
)

type Opcode byte

type Instructions []byte

func (instructions Instructions) fmt(opcodeDefinition *Operation, operands []int) string {
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
		op, err := LookupOperation(instructions[instructionByteIndex])
		if err != nil {
			fmt.Fprintf(&out, "ERROR: %s", err)
			continue
		}

		operands, readBytes := ReadOperands(op, instructions[instructionByteIndex+1:])
		formattedInstruction := instructions.fmt(op, operands)

		fmt.Fprintf(&out, "%04d %s\n", instructionByteIndex, formattedInstruction)

		instructionByteIndex += 1 + readBytes
	}

	return out.String()
}

type Operation struct {
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

	OpGotoNotTruthy // goto only if value on top of the stack is not truthy
	OpGoto

	OpNull

	OpGetGlobal
	OpSetGlobal

	OpGetLocal
	OpSetLocal

	OpGetStdlib

	OpArray
	OpHashMap

	OpIndex

	OpCall

	OpReturnValue
	OpReturn
)

var operations = map[Opcode]*Operation{
	OpConstant: {Name: "OpConstant", OperandWidths: []int{2}},

	OpAdd: {Name: "OpAdd"},
	OpSub: {Name: "OpSub"},
	OpMul: {Name: "OpMul"},
	OpDiv: {Name: "OpDiv"},

	OpTrue:  {Name: "OpTrue"},
	OpFalse: {Name: "OpFalse"},

	OpEqual:       {Name: "OpEqual"},
	OpNotEqual:    {Name: "OpNotEqual"},
	OpGreaterThan: {Name: "OpGreaterThan"},
	OpAnd:         {Name: "OpAnd"},
	OpOr:          {Name: "OpOr"},

	OpMinus: {Name: "OpMinus"},
	OpBang:  {Name: "OpBang"},

	OpPop: {Name: "OpPop"},

	OpGotoNotTruthy: {Name: "OpGotoNotTruthy", OperandWidths: []int{2}},
	OpGoto:          {Name: "OpGoto", OperandWidths: []int{2}},

	OpNull: {Name: "OpNull"},

	OpGetGlobal: {Name: "OpGetGlobal", OperandWidths: []int{2}},
	OpSetGlobal: {Name: "OpSetGlobal", OperandWidths: []int{2}},

	OpGetLocal: {Name: "OpGetLocal", OperandWidths: []int{1}},
	OpSetLocal: {Name: "OpSetLocal", OperandWidths: []int{1}},

	OpGetStdlib: {Name: "OpGetStdlibFunc", OperandWidths: []int{1}},

	OpArray:   {Name: "OpArray", OperandWidths: []int{2}},
	OpHashMap: {Name: "OpHashMap", OperandWidths: []int{2}},

	OpIndex: {Name: "OpIndex"},

	OpCall: {Name: "OpCall", OperandWidths: []int{1}},

	OpReturnValue: {Name: "OpReturnValue"},
	OpReturn:      {Name: "OpReturn"},
}

func LookupOperation(opcode byte) (*Operation, error) {
	op, ok := operations[Opcode(opcode)]
	if !ok {
		return nil, fmt.Errorf("opcode %d undefined", opcode)
	}

	return op, nil
}

func MakeInstruction(op Opcode, operands ...int) []byte {
	operation, ok := operations[op]
	if !ok {
		fmt.Printf("ERROR: opcode %d is not defined\n", op)
		return []byte{}
	}

	if len(operands) != len(operation.OperandWidths) {
		fmt.Printf("ERROR: operand count %d does not match definition %d for operation %d\n", len(operands), len(operation.OperandWidths), op)
		return []byte{}
	}

	instructionSize := 1
	for _, operandWidth := range operation.OperandWidths {
		instructionSize += operandWidth
	}

	instruction := make([]byte, instructionSize)

	instruction[0] = byte(op)

	offset := 1

	for index, operand := range operands {
		operandWidth := operation.OperandWidths[index]

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

func ReadOperands(operation *Operation, instructionBytes []byte) ([]int, int) {
	operands := make([]int, len(operation.OperandWidths))

	offset := 0

	for i, width := range operation.OperandWidths {
		switch width {
		case 1:
			operands[i] = int(instructionBytes[offset])
		case 2:
			operands[i] = int(utils.ReadUint16(instructionBytes[offset:]))
		case 4:
			operands[i] = int(utils.ReadUint32(instructionBytes[offset:]))
		}

		offset += width
	}

	return operands, offset
}
