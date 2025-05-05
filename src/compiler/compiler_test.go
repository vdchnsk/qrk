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
		return fmt.Errorf("wrong instructions length.\nexpected=\n%q\ngot=\n%q", flattenedExpected.String(), actual.String())
	}

	for index, instruction := range flattenedExpected {
		if actual[index] != instruction {
			return fmt.Errorf("wrong instruction at %d.\nexpected=\n%q\ngot=\n%q", index, flattenedExpected.String(), actual.String())
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

		case string:
			err := testStringObject(constant, actual[index])
			if err != nil {
				return fmt.Errorf("constant %d - testIntegerObject failed: %s", index, err)
			}

		case []code.Instructions:
			fn, ok := actual[index].(*object.CompiledFunction)
			if !ok {
				return fmt.Errorf("object is not CompiledFunction. got=%T (%+v)", actual[index], actual[index])
			}

			err := testInstructions(constant, fn.Instructions)
			if err != nil {
				return fmt.Errorf("constant %d - testInstructions failed: %s", index, err)
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

func testStringObject(expected string, actual object.Object) error {
	result, ok := actual.(*object.String)
	if !ok {
		return fmt.Errorf("object is not Integer. got=%T (%+v)", actual, actual)
	}

	if result.Value != expected {
		return fmt.Errorf("object has wrong value. got=%s, expected=%s", result.Value, expected)
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

	for i, tt := range tests {
		program := parse(tt.input)

		compiler := New()
		err := compiler.Compile(program)
		if err != nil {
			t.Fatalf("compiler error: %s", err)
		}

		bytecode := compiler.Bytecode()

		err = testInstructions(tt.expectedInstructions, bytecode.Instructions)
		if err != nil {
			t.Fatalf("[%d] testInstructions failed: %s", i+1, err)
		}

		err = testConstants(tt.expectedConstants, bytecode.Constants)
		if err != nil {
			t.Fatalf("[%d] testConstants failed: %s", i+1, err)
		}
	}
}

func TestGlobalLetStatement(t *testing.T) {
	tests := []compilerTestCase{
		{
			input: `
				let one = 1;
				let two = 2;
			`,
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []code.Instructions{
				code.MakeInstruction(code.OpConstant, 0),
				code.MakeInstruction(code.OpSetGlobal, 0),
				code.MakeInstruction(code.OpConstant, 1),
				code.MakeInstruction(code.OpSetGlobal, 1),
			},
		},
		{
			input: `
				let one = 1;
				one;
			`,
			expectedConstants: []interface{}{1},
			expectedInstructions: []code.Instructions{
				code.MakeInstruction(code.OpConstant, 0),
				code.MakeInstruction(code.OpSetGlobal, 0),
				code.MakeInstruction(code.OpGetGlobal, 0),
				code.MakeInstruction(code.OpPop),
			},
		},
		{
			input: `
				let one = 1;
				let two = one;
				two;
			`,
			expectedConstants: []interface{}{1},
			expectedInstructions: []code.Instructions{
				code.MakeInstruction(code.OpConstant, 0),
				code.MakeInstruction(code.OpSetGlobal, 0),
				code.MakeInstruction(code.OpGetGlobal, 0),
				code.MakeInstruction(code.OpSetGlobal, 1),
				code.MakeInstruction(code.OpGetGlobal, 1),
				code.MakeInstruction(code.OpPop),
			},
		},
	}

	runCompilerTests(t, tests)
}

func TestStringExpression(t *testing.T) {
	tests := []compilerTestCase{
		{
			input:             `"test"`,
			expectedConstants: []interface{}{"test"},
			expectedInstructions: []code.Instructions{
				code.MakeInstruction(code.OpConstant, 0),
				code.MakeInstruction(code.OpPop),
			},
		},
		{
			input:             `"te" + "st"`,
			expectedConstants: []interface{}{"te", "st"},
			expectedInstructions: []code.Instructions{
				code.MakeInstruction(code.OpConstant, 0),
				code.MakeInstruction(code.OpConstant, 1),
				code.MakeInstruction(code.OpAdd),
				code.MakeInstruction(code.OpPop),
			},
		},
	}

	runCompilerTests(t, tests)
}

func TestArrayExpression(t *testing.T) {
	tests := []compilerTestCase{
		{
			input:             `[]`,
			expectedConstants: []interface{}{},
			expectedInstructions: []code.Instructions{
				code.MakeInstruction(code.OpArray, 0),
				code.MakeInstruction(code.OpPop),
			},
		},
		{
			input:             `[1, 2, 3]`,
			expectedConstants: []interface{}{1, 2, 3},
			expectedInstructions: []code.Instructions{
				code.MakeInstruction(code.OpConstant, 0),
				code.MakeInstruction(code.OpConstant, 1),
				code.MakeInstruction(code.OpConstant, 2),
				code.MakeInstruction(code.OpArray, 3),
				code.MakeInstruction(code.OpPop),
			},
		},
		{
			input:             `[1+1, 2*2, 3*3-3]`,
			expectedConstants: []interface{}{1, 1, 2, 2, 3, 3, 3},
			expectedInstructions: []code.Instructions{
				code.MakeInstruction(code.OpConstant, 0),
				code.MakeInstruction(code.OpConstant, 1),
				code.MakeInstruction(code.OpAdd),

				code.MakeInstruction(code.OpConstant, 2),
				code.MakeInstruction(code.OpConstant, 3),
				code.MakeInstruction(code.OpMul),

				code.MakeInstruction(code.OpConstant, 4),
				code.MakeInstruction(code.OpConstant, 5),
				code.MakeInstruction(code.OpMul),
				code.MakeInstruction(code.OpConstant, 6),
				code.MakeInstruction(code.OpSub),

				code.MakeInstruction(code.OpArray, 3),
				code.MakeInstruction(code.OpPop),
			},
		},
	}

	runCompilerTests(t, tests)
}

func TestHashMapExpression(t *testing.T) {
	tests := []compilerTestCase{
		{
			input:             `{}`,
			expectedConstants: []interface{}{},
			expectedInstructions: []code.Instructions{
				code.MakeInstruction(code.OpHashMap, 0),
				code.MakeInstruction(code.OpPop),
			},
		},
		{
			input:             `{1: 2, 2: 3, 3: 4}`,
			expectedConstants: []interface{}{1, 2, 2, 3, 3, 4},
			expectedInstructions: []code.Instructions{
				code.MakeInstruction(code.OpConstant, 0),
				code.MakeInstruction(code.OpConstant, 1),
				code.MakeInstruction(code.OpConstant, 2),
				code.MakeInstruction(code.OpConstant, 3),
				code.MakeInstruction(code.OpConstant, 4),
				code.MakeInstruction(code.OpConstant, 5),
				code.MakeInstruction(code.OpHashMap, 6),
				code.MakeInstruction(code.OpPop),
			},
		},
		{
			input:             `{1: 2, 2: 3+3, 3: 4-1}`,
			expectedConstants: []interface{}{1, 2, 2, 3, 3, 3, 4, 1},
			expectedInstructions: []code.Instructions{
				code.MakeInstruction(code.OpConstant, 0),
				code.MakeInstruction(code.OpConstant, 1),

				code.MakeInstruction(code.OpConstant, 2),
				code.MakeInstruction(code.OpConstant, 3),
				code.MakeInstruction(code.OpConstant, 4),
				code.MakeInstruction(code.OpAdd),

				code.MakeInstruction(code.OpConstant, 5),
				code.MakeInstruction(code.OpConstant, 6),
				code.MakeInstruction(code.OpConstant, 7),
				code.MakeInstruction(code.OpSub),

				code.MakeInstruction(code.OpHashMap, 6),
				code.MakeInstruction(code.OpPop),
			},
		},
	}

	runCompilerTests(t, tests)
}

func TestIndexExpression(t *testing.T) {
	tests := []compilerTestCase{
		{
			input:             `[1][0]`,
			expectedConstants: []interface{}{1, 0},
			expectedInstructions: []code.Instructions{
				code.MakeInstruction(code.OpConstant, 0),
				code.MakeInstruction(code.OpArray, 1),
				code.MakeInstruction(code.OpConstant, 1),
				code.MakeInstruction(code.OpIndex),
				code.MakeInstruction(code.OpPop),
			},
		},
		{
			input:             `[1, 2, 3][1 + 1]`,
			expectedConstants: []interface{}{1, 2, 3, 1, 1},
			expectedInstructions: []code.Instructions{
				code.MakeInstruction(code.OpConstant, 0),
				code.MakeInstruction(code.OpConstant, 1),
				code.MakeInstruction(code.OpConstant, 2),
				code.MakeInstruction(code.OpArray, 3),

				code.MakeInstruction(code.OpConstant, 3),
				code.MakeInstruction(code.OpConstant, 4),
				code.MakeInstruction(code.OpAdd),
				code.MakeInstruction(code.OpIndex),
				code.MakeInstruction(code.OpPop),
			},
		},
		{
			input:             `{1: 2}[1]`,
			expectedConstants: []interface{}{1, 2, 1},
			expectedInstructions: []code.Instructions{
				code.MakeInstruction(code.OpConstant, 0),
				code.MakeInstruction(code.OpConstant, 1),
				code.MakeInstruction(code.OpHashMap, 2),

				code.MakeInstruction(code.OpConstant, 2),
				code.MakeInstruction(code.OpIndex),
				code.MakeInstruction(code.OpPop),
			},
		},
	}

	runCompilerTests(t, tests)
}

func TestFunctions(t *testing.T) {
	tests := []compilerTestCase{
		{
			input: `fn() { return 5 + 10 }`,
			expectedConstants: []interface{}{
				5,
				10,
				[]code.Instructions{
					code.MakeInstruction(code.OpConstant, 0),
					code.MakeInstruction(code.OpConstant, 1),
					code.MakeInstruction(code.OpAdd),
					code.MakeInstruction(code.OpReturnValue),
				},
			},
			expectedInstructions: []code.Instructions{
				code.MakeInstruction(code.OpConstant, 2),
				code.MakeInstruction(code.OpPop),
			},
		},
	}

	runCompilerTests(t, tests)
}

func TestCompilerScopes(t *testing.T) {
	compiler := New()
	if compiler.scopeIndex != 0 {
		t.Fatalf("initial scope index is not 0. got=%d", compiler.scopeIndex)
	}

	compiler.emit(code.OpMul)

	compiler.enterScope()
	if compiler.scopeIndex != 1 {
		t.Fatalf("scope index is not 1 after entering scope. got=%d", compiler.scopeIndex)
	}

	compiler.emit(code.OpSub)

	if len(compiler.scopes[compiler.scopeIndex].instructions) != 1 {
		t.Fatalf("wrong number of instructions in scope. got=%d", len(compiler.scopes[compiler.scopeIndex].instructions))
	}

	last := compiler.scopes[compiler.scopeIndex].lastInstruction
	if last.Opcode != code.OpSub {
		t.Fatalf("wrong last instruction in scope. got=%d", last.Opcode)
	}

	compiler.leaveScope()
	if compiler.scopeIndex != 0 {
		t.Fatalf("scope index is not 0 after leaving scope. got=%d", compiler.scopeIndex)
	}

	compiler.emit(code.OpAdd)

	if len(compiler.scopes[compiler.scopeIndex].instructions) != 2 {
		t.Fatalf("wrong number of instructions in scope. got=%d", len(compiler.scopes[compiler.scopeIndex].instructions))
	}

	last = compiler.scopes[compiler.scopeIndex].lastInstruction
	if last.Opcode != code.OpAdd {
		t.Fatalf("wrong last instruction in scope. got=%d", last.Opcode)
	}

	prev := compiler.scopes[compiler.scopeIndex].prevInstruction

	if prev.Opcode != code.OpMul {
		t.Fatalf("wrong previous instruction in scope. got=%d", prev.Opcode)
	}
}
