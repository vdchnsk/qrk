package code

import "testing"

func TestMake(t *testing.T) {
	tests := []struct {
		op       Opcode
		operands []int
		expected []byte
	}{
		// 65534 = 11111111(255) + 11111110(254)
		{OpConstant, []int{65534}, []byte{byte(OpConstant), 255, 254}},

		// 65535 = 11111111(255) + 11111111(255)
		{OpConstant, []int{65535}, []byte{byte(OpConstant), 255, 255}},
	}

	for _, tt := range tests {
		instructions := MakeInstruction(tt.op, tt.operands...)

		if len(instructions) != len(tt.expected) {
			t.Errorf("instruction has wrong length. want=%d, got=%d", len(tt.expected), len(instructions))
		}

		for i, b := range tt.expected {
			if instructions[i] != tt.expected[i] {
				t.Errorf("wrong byte at pos %d. want=%d, got=%d", i, b, instructions[i])
			}
		}
	}
}
