package vm

import (
	"fmt"
	"testing"

	"github.com/vdchnsk/qrk/src/ast"
	"github.com/vdchnsk/qrk/src/compiler"
	"github.com/vdchnsk/qrk/src/lexer"
	"github.com/vdchnsk/qrk/src/object"
	"github.com/vdchnsk/qrk/src/parser"
)

func parse(input string) *ast.Program {
	l := lexer.NewLexer(input)
	p := parser.NewParser(l)

	return p.ParseProgram()
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

type vmTestCase struct {
	input    string
	expected interface{}
}

func runVmTests(t *testing.T, tests []vmTestCase) {
	t.Helper()

	for _, tt := range tests {
		program := parse(tt.input)

		compiler := compiler.NewCompiler()
		err := compiler.Compile(program)
		if err != nil {
			t.Fatalf("compiler error %s", err)
		}

		bytecode := compiler.Bytecode()
		vm := NewVm(bytecode)

		err = vm.Run()
		if err != nil {
			t.Fatalf("vm error: %s", err)
		}

		stackElement := vm.LastPoppedStackElem()

		testExpectedObject(t, tt.expected, stackElement)
	}
}

func testExpectedObject(t *testing.T, expectedObj interface{}, actualObj object.Object) {
	t.Helper()

	switch expectedObj := expectedObj.(type) {
	case int:
		err := testIntegerObject(int64(expectedObj), actualObj)
		if err != nil {
			t.Errorf("testIntegerObject failed: %s", err)
		}
	}
}

func TestIntegerArithmetic(t *testing.T) {
	tests := []vmTestCase{
		{input: "1", expected: 1},
		{input: "2", expected: 2},
		{input: "1 + 2", expected: 3},
		{input: "1 - 2", expected: -1},
		{input: "1 * 2", expected: 2},
		{input: "2 / 1", expected: 2},
		{input: "50 / 2 * 2 + 10 - 5", expected: 55},
		{input: "5 + 5 + 5 + 5 - 10", expected: 10},
		{input: "2 * 2 * 2 * 2 * 2", expected: 32},
		{input: "5 * 2 + 10", expected: 20},
		{input: "5 + 2 * 10", expected: 25},
		{input: "5 * (2 + 10)", expected: 60},
	}

	runVmTests(t, tests)
}
