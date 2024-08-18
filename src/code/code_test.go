package code

import (
	"strings"
	"testing"

	"github.com/vdchnsk/qrk/src/utils"
)

func TestMake(t *testing.T) {
	tests := []struct {
		op       Opcode
		operands []int
		expected []byte
	}{
		// 65534 = 11111111(255) + 11111110(254)
		{OpConstant, []int{utils.MaxIntForBytes(2) - 1}, []byte{byte(OpConstant), 255, 254}},

		// max int for 2 bytes is 65535 = 11111111(255) + 11111111(255)
		{OpConstant, []int{utils.MaxIntForBytes(2)}, []byte{byte(OpConstant), 255, 255}},

		{OpAdd, []int{}, []byte{byte(OpAdd)}},
	}

	for _, tt := range tests {
		instructions := MakeInstruction(tt.op, tt.operands...)

		if len(instructions) != len(tt.expected) {
			t.Errorf("instruction has wrong length. expected=%d, got=%d", len(tt.expected), len(instructions))
		}

		for i, b := range tt.expected {
			if instructions[i] != tt.expected[i] {
				t.Errorf("wrong byte at pos %d. expected=%d, got=%d", i, b, instructions[i])
			}
		}
	}
}

func TestInstructionToString(t *testing.T) {
	instructions := []Instructions{
		MakeInstruction(OpConstant, 1),
		MakeInstruction(OpConstant, 2),
		MakeInstruction(OpConstant, utils.MaxIntForBytes(2)),
		MakeInstruction(OpAdd),
	}

	expected := strings.Join([]string{
		"0000 OpConstant 1",
		"0003 OpConstant 2",
		"0006 OpConstant 65535",
		"0009 OpAdd",
	}, "\n") + "\n"

	flattened := Instructions{}
	for _, instructionBytes := range instructions {
		flattened = append(flattened, instructionBytes...)
	}

	if flattened.ToString() != expected {
		t.Errorf("instructions wrongly formatted.\nexpected=%q\ngot=%q", expected, flattened.ToString())
	}
}

func TestReadOperands(t *testing.T) {
	tests := []struct {
		opcode    Opcode
		operands  []int
		bytesRead int
	}{
		{OpConstant, []int{utils.MaxIntForBytes(2)}, 2},
	}

	for _, tt := range tests {
		instruction := MakeInstruction(tt.opcode, tt.operands...)

		opcodeDefinition, err := LookupDefinition(byte(tt.opcode))
		if err != nil {
			t.Fatalf("definition not found: %s", err)
		}

		operandsRead, bytesRead := ReadOperands(opcodeDefinition, instruction[1:])
		if bytesRead != tt.bytesRead {
			t.Fatalf("n wrong. expected=%d, got=%d", tt.bytesRead, bytesRead)
		}

		for i, expectedOperand := range tt.operands {
			if operandsRead[i] != expectedOperand {
				t.Errorf("operand wrong. expected=%d, got=%d", expectedOperand, operandsRead[i])
			}
		}
	}
}
