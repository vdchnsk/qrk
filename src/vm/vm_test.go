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
	case *object.Null:
		if actualObj != Null {
			t.Errorf("object is not Null. got=%T (%+v)", actualObj, actualObj)
		}
	}

	if err != nil {
		t.Errorf("testExpectedObject failed: %s", err)
	}
}

func TestIntegerArithmetic(t *testing.T) {
	tests := []vmTestCase{
		{"1", 1},
		{"2", 2},
		{"1 + 2", 3},
		{"1 - 2", -1},
		{"1 * 2", 2},
		{"2 / 1", 2},
		{"50 / 2 * 2 + 10 - 5", 55},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"5 * (2 + 10)", 60},
		{"-50 + 100 + -50", 0},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
		{"-10", -10},
	}

	runVmTests(t, tests)
}

func TestBooleanExpressions(t *testing.T) {
	tests := []vmTestCase{
		{"true;", true},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 == 1", true},
		{"(1 > 2) == true", false},
		{"(1 > 2) == false", true},
		{"true != false", true},
		{"false != true", true},
		{"true == true", true},
		{"true && true", true},
		{"true && false", false},
		{"false && false", false},
		{"false || false", false},
		{"true || false", true},
		{"!!true", true},
		{"!!5", true},
		{"!( if false { 5 } )", true},
	}

	runVmTests(t, tests)
}

func TestConditionals(t *testing.T) {
	tests := []vmTestCase{
		{"if true { 10 }", 10},
		{"if true { 10 } else { 20 }", 10},
		{"if false { 10 } else { 20 }", 20},
		{"if 1 { 10 }", 10},
		{"if 1 < 2 { 10 }", 10},
		{"if 1 < 2 { 10 } else { 20 }", 10},
		{"if 1 > 2 { 10 } else { 20 }", 20},
		{"if 1 > 2 { 10 }", Null},
		{"if false { 10 }", Null},
		{"if (if false { 10 }) { 10 } else { 20 }", 20},
	}

	runVmTests(t, tests)
}
