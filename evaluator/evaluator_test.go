package evaluator

import (
	"testing"

	"github.com/vdchnsk/i-go/lexer"
	"github.com/vdchnsk/i-go/object"
	"github.com/vdchnsk/i-go/parser"
)

func testEval(input string) object.Object {
	lexer := lexer.NewLexer(input)
	p := parser.NewParser(lexer)
	program := p.ParseProgram()

	return Eval(program)
}

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"42", 42},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func testIntegerObject(t *testing.T, obj object.Object, expectedValue int64) bool {
	result, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf(
			"object is not Integer. got=%T (%+v)",
			obj, obj,
		)
		return false
	}
	if result.Value != expectedValue {
		t.Errorf(
			"integer != expected value. got=%d, expected=%d",
			result, expectedValue,
		)
		return false
	}
	return true
}

func TestEvalBooelanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
		{"1 < 2", true},
		{"1 < 1", false},
		{"1 > 2", false},
		{"1 != 2", true},
		{"1 == 2", false},
		{"1 == 1", true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBoooleanObject(t, evaluated, tt.expected)
	}
}

func testBoooleanObject(t *testing.T, obj object.Object, expectedValue bool) bool {
	result, ok := obj.(*object.Boolean)
	if !ok {
		t.Errorf(
			"object is not Booolean. got=%T (%+v)",
			obj, obj,
		)
		return false
	}
	if result.Value != expectedValue {
		t.Errorf(
			"boolean != expected value. got=%T, expected=%T",
			result, expectedValue,
		)
		return false
	}
	return true
}

func TestBangOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"!true", false},
		{"!false", true},
		{"!!true", true},
		{"!!false", false},
		{"!42", false},
		{"!!42", true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBoooleanObject(t, evaluated, tt.expected)
	}
}

func TestMinusperator(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"-5", -5},
		{"-42", -42},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestEvalInfixExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"-5-5", -10},
		{"1000 + 100", 1100},
		{"5 * 5", 25},
		{"25 / 5", 5},
		{"(25 + 5) / 6", 5},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}
