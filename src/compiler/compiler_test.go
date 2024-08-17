package compiler

import (
	"fmt"
	"testing"

	"github.com/vdchnsk/qrk/src/ast"
	"github.com/vdchnsk/qrk/src/code"
	"github.com/vdchnsk/qrk/src/lexer"
	"github.com/vdchnsk/qrk/src/object"
	"github.com/vdchnsk/qrk/src/parser"
)

func parse(input string) *ast.Program {
	l := lexer.NewLexer(input)
	p := parser.NewParser(l)

	return p.ParseProgram()
}

func flattenInstructions(instructions []code.Instructions) code.Instructions {
	out := code.Instructions{}

	for _, instruction := range instructions {
		out = append(out, instruction...)
	}

	return out
}

func testInstructions(expected []code.Instructions, actual code.Instructions) error {
	flattenedExpected := flattenInstructions(expected)

	if len(flattenedExpected) != len(actual) {
		return fmt.Errorf("wrong instructions length.\nexpected=%q\ngot=%q", flattenedExpected, actual)
	}

	for index, instruction := range flattenedExpected {
		if actual[index] != instruction {
			return fmt.Errorf("wrong instruction at %d.\nexpected=%q\ngot=%q", index, flattenedExpected, actual)
		}
	}

	return nil
}

func testConstants(expected []interface{}, actual []object.Object) error {
	if len(expected) != len(actual) {
		return fmt.Errorf("wrong number of constants. got=%d, expected=%d", len(actual), len(expected))
	}

	for index, constant := range expected {
		switch constant := constant.(type) {
		case int:
			err := testIntegerObject(int64(constant), actual[index])
			if err != nil {
				return fmt.Errorf("constant %d - testIntegerObject failed: %s", index, err)
			}
		}
	}

	return nil
}

func testIntegerObject(expected int64, actual object.Object) error {
	result, ok := actual.(*object.Integer)

	if !ok {
		return fmt.Errorf("object is not Integer. got=%T (%+v)", actual, actual)
	}

	if result.Value != expected {
		return fmt.Errorf("object has wrong value. got=%d, expected=%d", result.Value, expected)
	}

	return nil
}

type compilerTestCase struct {
	input                string
	expectedConstants    []interface{}
	expectedInstructions []code.Instructions
}

func TestIntegerArithmetic(t *testing.T) {
	tests := []compilerTestCase{
		{
			input:             "1 + 2",
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []code.Instructions{
				// 0 and 1 point to indexes in the constants slice
				code.MakeInstruction(code.OpConstant, 0),
				code.MakeInstruction(code.OpConstant, 1),
			},
		},
	}

	runCompilerTests(t, tests)
}

func runCompilerTests(t *testing.T, tests []compilerTestCase) {
	t.Helper()

	for _, tt := range tests {
		program := parse(tt.input)

		compiler := NewCompiler()
		err := compiler.Compile(program)
		if err != nil {
			t.Fatalf("compiler error: %s", err)
		}

		bytecode := compiler.Bytecode()

		err = testInstructions(tt.expectedInstructions, bytecode.Instructions)
		if err != nil {
			t.Fatalf("testInstructions failed: %s", err)
		}

		err = testConstants(tt.expectedConstants, bytecode.Constants)
		if err != nil {
			t.Fatalf("testConstants failed: %s", err)
		}
	}
}
