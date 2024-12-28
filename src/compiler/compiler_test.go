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
		return fmt.Errorf("wrong instructions length.\nexpected=%q\ngot=%q", flattenedExpected.String(), actual.String())
	}

	for index, instruction := range flattenedExpected {
		if actual[index] != instruction {
			return fmt.Errorf("wrong instruction at %d.\nexpected=%q\ngot=%q", index, flattenedExpected.String(), actual.String())
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
				code.MakeInstruction(code.OpAdd),
				code.MakeInstruction(code.OpPop),
			},
		},
		{
			input:             "1 * 2",
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []code.Instructions{
				code.MakeInstruction(code.OpConstant, 0),
				code.MakeInstruction(code.OpConstant, 1),
				code.MakeInstruction(code.OpMul),
				code.MakeInstruction(code.OpPop),
			},
		},
		{
			input:             "2 / 2",
			expectedConstants: []interface{}{2, 2},
			expectedInstructions: []code.Instructions{
				code.MakeInstruction(code.OpConstant, 0),
				code.MakeInstruction(code.OpConstant, 1),
				code.MakeInstruction(code.OpDiv),
				code.MakeInstruction(code.OpPop),
			},
		},
		{
			input:             "1 - 2",
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []code.Instructions{
				code.MakeInstruction(code.OpConstant, 0),
				code.MakeInstruction(code.OpConstant, 1),
				code.MakeInstruction(code.OpSub),
				code.MakeInstruction(code.OpPop),
			},
		},
		{
			input:             "1; 2",
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []code.Instructions{
				code.MakeInstruction(code.OpConstant, 0),
				code.MakeInstruction(code.OpPop),
				code.MakeInstruction(code.OpConstant, 1),
				code.MakeInstruction(code.OpPop),
			},
		},
		{
			input:             "-1",
			expectedConstants: []interface{}{1},
			expectedInstructions: []code.Instructions{
				code.MakeInstruction(code.OpConstant, 0),
				code.MakeInstruction(code.OpMinus),
				code.MakeInstruction(code.OpPop),
			},
		},
	}

	runCompilerTests(t, tests)
}

func TestBooleanExpressions(t *testing.T) {
	tests := []compilerTestCase{
		{
			input:             "true;",
			expectedConstants: []interface{}{},
			expectedInstructions: []code.Instructions{
				code.MakeInstruction(code.OpTrue),
				code.MakeInstruction(code.OpPop),
			},
		},
		{
			input:             "false;",
			expectedConstants: []interface{}{},
			expectedInstructions: []code.Instructions{
				code.MakeInstruction(code.OpFalse),
				code.MakeInstruction(code.OpPop),
			},
		},
		{
			input:             "1 > 2",
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []code.Instructions{
				code.MakeInstruction(code.OpConstant, 0),
				code.MakeInstruction(code.OpConstant, 1),
				code.MakeInstruction(code.OpGreaterThan),
				code.MakeInstruction(code.OpPop),
			},
		},
		{
			input:             "1 < 2",
			expectedConstants: []interface{}{2, 1},
			expectedInstructions: []code.Instructions{
				code.MakeInstruction(code.OpConstant, 0),
				code.MakeInstruction(code.OpConstant, 1),
				code.MakeInstruction(code.OpGreaterThan),
				code.MakeInstruction(code.OpPop),
			},
		},
		{
			input:             "1 == 2",
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []code.Instructions{
				code.MakeInstruction(code.OpConstant, 0),
				code.MakeInstruction(code.OpConstant, 1),
				code.MakeInstruction(code.OpEqual),
				code.MakeInstruction(code.OpPop),
			},
		},
		{
			input:             "1 != 2",
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []code.Instructions{
				code.MakeInstruction(code.OpConstant, 0),
				code.MakeInstruction(code.OpConstant, 1),
				code.MakeInstruction(code.OpNotEqual),
				code.MakeInstruction(code.OpPop),
			},
		},
		{
			input:             "true == true",
			expectedConstants: []interface{}{},
			expectedInstructions: []code.Instructions{
				code.MakeInstruction(code.OpTrue),
				code.MakeInstruction(code.OpTrue),
				code.MakeInstruction(code.OpEqual),
				code.MakeInstruction(code.OpPop),
			},
		},
		{
			input:             "true && true",
			expectedConstants: []interface{}{},
			expectedInstructions: []code.Instructions{
				code.MakeInstruction(code.OpTrue),
				code.MakeInstruction(code.OpTrue),
				code.MakeInstruction(code.OpAnd),
				code.MakeInstruction(code.OpPop),
			},
		},
		{
			input:             "true || false",
			expectedConstants: []interface{}{},
			expectedInstructions: []code.Instructions{
				code.MakeInstruction(code.OpTrue),
				code.MakeInstruction(code.OpFalse),
				code.MakeInstruction(code.OpOr),
				code.MakeInstruction(code.OpPop),
			},
		},
		{
			input:             "false || false",
			expectedConstants: []interface{}{},
			expectedInstructions: []code.Instructions{
				code.MakeInstruction(code.OpFalse),
				code.MakeInstruction(code.OpFalse),
				code.MakeInstruction(code.OpOr),
				code.MakeInstruction(code.OpPop),
			},
		},
		{
			input:             "true != false",
			expectedConstants: []interface{}{},
			expectedInstructions: []code.Instructions{
				code.MakeInstruction(code.OpTrue),
				code.MakeInstruction(code.OpFalse),
				code.MakeInstruction(code.OpNotEqual),
				code.MakeInstruction(code.OpPop),
			},
		},
		{
			input:             "!true",
			expectedConstants: []interface{}{},
			expectedInstructions: []code.Instructions{
				code.MakeInstruction(code.OpTrue),
				code.MakeInstruction(code.OpBang),
				code.MakeInstruction(code.OpPop),
			},
		},
	}

	runCompilerTests(t, tests)
}

func TestConditionals(t *testing.T) {
	tests := []compilerTestCase{
		{
			input:             `if (true) { 10 }; 3333;`,
			expectedConstants: []interface{}{10, 3333},
			expectedInstructions: []code.Instructions{
				// 0000
				code.MakeInstruction(code.OpTrue),
				// 0001
				code.MakeInstruction(code.OpGotoNotTruthy, 10),
				// 0004
				code.MakeInstruction(code.OpConstant, 0),
				// 0007
				code.MakeInstruction(code.OpGoto, 11),
				// 0010
				code.MakeInstruction(code.OpNull),
				// 0011
				code.MakeInstruction(code.OpPop),
				// 0012
				code.MakeInstruction(code.OpConstant, 1),
				// 0015
				code.MakeInstruction(code.OpPop),
			},
		},
		{
			input:             `if (true) { 10 } else { 20 }; 3333;`,
			expectedConstants: []interface{}{10, 20, 3333},
			expectedInstructions: []code.Instructions{
				// 0000
				code.MakeInstruction(code.OpTrue),
				// 0001
				code.MakeInstruction(code.OpGotoNotTruthy, 10), // go to else if top-of-the-stack value is not truthy
				// 0004
				code.MakeInstruction(code.OpConstant, 0),
				// 0007
				code.MakeInstruction(code.OpGoto, 13), // skip else
				// 0010
				code.MakeInstruction(code.OpConstant, 1),
				// 0013
				code.MakeInstruction(code.OpPop),
				// 0014
				code.MakeInstruction(code.OpConstant, 2),
				// 0017
				code.MakeInstruction(code.OpPop),
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
