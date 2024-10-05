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

func testObject[T bool | int64](expected T, actual object.Object) error {
	var result T

	switch v := actual.(type) {
	case *object.Integer:
		if _, ok := any(expected).(int64); !ok {
			return fmt.Errorf("expected int64, got %T", expected)
		}
		result = any(v.Value).(T)
	case *object.Boolean:
		if _, ok := any(expected).(bool); !ok {
			return fmt.Errorf("expected bool, got %T", expected)
		}
		result = any(v.Value).(T)
	default:
		return fmt.Errorf("unsupported object type: %T", actual)
	}

	if result != expected {
		return fmt.Errorf("object has wrong value. got=%v, expected=%v", result, expected)
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

	var err error

	switch expectedObj := expectedObj.(type) {
	case int:
		err = testObject(int64(expectedObj), actualObj)
	case bool:
		err = testObject(expectedObj, actualObj)
	}

	if err != nil {
		t.Errorf("testExpectedObject failed: %s", err)
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

func TestBooleanExpressions(t *testing.T) {
	tests := []vmTestCase{
		{input: "true;", expected: true},
		{input: "1 < 2", expected: true},
		{input: "1 > 2", expected: false},
		{input: "1 < 1", expected: false},
		{input: "1 > 1", expected: false},
		{input: "1 == 1", expected: true},
		{input: "(1 > 2) == true", expected: false},
		{input: "(1 > 2) == false", expected: true},
		{input: "true != false", expected: true},
		{input: "false != true", expected: true},
		{input: "true == true", expected: true},
		{input: "true && true", expected: true},
		{input: "true && false", expected: false},
		{input: "false && false", expected: false},
		{input: "false || false", expected: false},
		{input: "true || false", expected: true},
	}

	runVmTests(t, tests)
}
