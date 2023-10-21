package evaluator

import (
	"testing"

	"github.com/vdchnsk/i-go/src/lexer"
	"github.com/vdchnsk/i-go/src/object"
	"github.com/vdchnsk/i-go/src/parser"
)

func testEval(input string) object.Object {
	lexer := lexer.NewLexer(input)
	p := parser.NewParser(lexer)
	program := p.ParseProgram()
	env := object.NewEnvironment()

	return Eval(program, env)
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

func testNullObject(t *testing.T, obj object.Object) bool {
	_, ok := obj.(*object.Null)
	if !ok {
		t.Errorf(
			"object is not null. got=%d",
			obj,
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
		{"(1 == 1) == true", true},
		{"(10 > 100) == true", false},
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

func TestEvalIfExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"if (true) { 10 }", 10},
		{"if (false) { 10 } else { 42 }", 42},
		{"if (false) { 10 }", nil},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func TestEvalReturnStatement(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"return 10;", 10},
		{"42 * 5; return 10; 13;", 10},
		{"42 * 5; return 10; return 42;", 10},
		{`
		if (true) {
			return 10;
		}
		`, 10},
		{`
		if (true) {
			return true;
		}
		`, true},
		{`
		if (true) {
			if (true) {
				return 10;
			}
		}
		`, 10},
		{`
		if (false) {
			return 5;
		} else {
			if (true) {
				if (true) {
					return 42;
				}
				if (true) {
					return 10;
				}
			}
		}
		`, 42},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, isInt := tt.expected.(int)
		if isInt {
			testIntegerObject(t, evaluated, int64(integer))
		}
		boolean, isBool := tt.expected.(bool)
		if isBool {
			testBoooleanObject(t, evaluated, boolean)
		}
	}
}

func TestErrorEval(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"1 + true;", "type mismatch: INTEGER + BOOLEAN"},
		{"-true;", "unknown operator: -BOOLEAN"},
		{"true + true;", "unknown operator: BOOLEAN + BOOLEAN"},
		{`
		if (true) {
			true + true;
		}
		return 10;
		`, "unknown operator: BOOLEAN + BOOLEAN"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		err, ok := evaluated.(*object.Error)
		if !ok {
			t.Errorf(
				"no error is returned, got=%T(%+v)",
				evaluated, evaluated,
			)
		}
		if err.Message != tt.expected {
			t.Errorf(
				"wrong error message, got=%s, expected=%s",
				err.Message,
				tt.expected,
			)
		}
	}
}

func TestLetStatement(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let a = 4; a", 4},
		{`
		let a = 4;
		a;
		`, 4},
	}

	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input           string
		expectedMessage string
	}{
		{"foobar", "identifier not found: foobar"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		err, ok := evaluated.(*object.Error)
		if !ok {
			t.Errorf(
				"no error is returned, got=%T(%+v)",
				evaluated, evaluated,
			)
		}
		if err.Message != tt.expectedMessage {
			t.Errorf(
				"wrong error message, got=%s, expected=%s",
				err.Message,
				tt.expectedMessage,
			)
		}
	}
}

func TestFunctionObject(t *testing.T) {
	input := "fn(a, b) { a + b; };"
	expectedBody := "(a + b)"

	evaluated := testEval(input)
	fn, ok := evaluated.(*object.Function)
	if !ok {
		t.Errorf(
			"object is not a function, got=%T(%+v)",
			evaluated, evaluated,
		)
	}

	params := fn.Parameters

	if len(params) != 2 {
		t.Errorf(
			"wrong amount of parameters, got=%T(%+v)",
			evaluated, evaluated,
		)
	}
	if params[0].ToString() != "a" || params[1].ToString() != "b" {
		t.Errorf(
			"wrong parameters, got=%T, expected=%s",
			params, "(a, b)",
		)
	}
	if fn.Body.ToString() != expectedBody {
		t.Errorf(
			"wrong body, got=%s, expected=%s",
			fn.Body.ToString(), expectedBody,
		)
	}
}
