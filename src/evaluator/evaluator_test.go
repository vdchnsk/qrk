package evaluator

import (
	"fmt"
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

func testStringObject(t *testing.T, obj object.Object, expectedValue string) bool {
	result, ok := obj.(*object.String)
	if !ok {
		t.Errorf(
			"object is not String. got=%T (%+v)",
			obj, obj,
		)
		return false
	}
	if result.Value != expectedValue {
		t.Errorf(
			"integer != expected value. got=%s, expected=%s",
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
		{"if true { 10 }", 10},
		{"if false { 10 }", nil},
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
		{"1 + true;", fmt.Sprintf("%s: INTEGER + BOOLEAN", TYPE_MISMATCH)},
		{"-true;", fmt.Sprintf("%s: -BOOLEAN", UNKNOWN_OPERATOR)},
		{"true + true;", fmt.Sprintf("%s: BOOLEAN + BOOLEAN", UNKNOWN_OPERATOR)},
		{`
		if (true) {
			true + true;
		}
		return 10;
		`, fmt.Sprintf("%s: BOOLEAN + BOOLEAN", UNKNOWN_OPERATOR)},
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
		{"foobar", fmt.Sprintf("%s: foobar", IDENTIFIER_NOT_FOUND)},
		{`"str"-"str"`, fmt.Sprintf("%s: STRING - STRING", UNKNOWN_OPERATOR)},
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

func TestFunctionEval(t *testing.T) {
	tests := []struct {
		input          string
		expectedOutput int64
	}{
		{"let double = fn(x) { x * 2; }; double(5);", 10},
		{"let add = fn(x, y) { x + y; }; add(5, 5);", 10},
		{"let identity = fn(x) { x; }; identity(5);", 5},
		{"let identity = fn(x) { return x; }; identity(5);", 5},
		{"let add = fn(x, y) { x + y; }; add(5 + 5, add(5, 5));", 20},
		{"fn(x) { x; }(5)", 5},
	}

	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expectedOutput)
	}
}

func TestStringLiteral(t *testing.T) {
	input := `"Hello world!"`

	evaluated := testEval(input)
	expectedRes := "Hello world!"
	if !testStringObject(t, evaluated, expectedRes) {
		t.Errorf(
			"String is evaluated incorrectly, expected=%s",
			expectedRes,
		)
	}
}

func TestStringConcatenation(t *testing.T) {
	input := `"Hello " + "world!"`

	evaluated := testEval(input)
	expectedRes := "Hello world!"
	if !testStringObject(t, evaluated, expectedRes) {
		t.Errorf(
			"String is evaluated incorrectly, expected=%s",
			expectedRes,
		)
	}
}

func TestBuiltInFunctions(t *testing.T) {
	tests := []struct {
		input          string
		expectedOutput interface{}
	}{
		{`len("")`, 0},
		{`len("hello world")`, 11},
		{`len(1)`, "argument to `len` is not supported, got INTEGER"},
		{`len(true)`, "argument to `len` is not supported, got BOOLEAN"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		switch expected := tt.expectedOutput.(type) {
		case int64:
			testIntegerObject(t, evaluated, int64(expected))
		case string:
			error, ok := evaluated.(*object.Error)
			if ok && error.Message != tt.expectedOutput {
				t.Errorf(
					"got an error with unexpected message, got=%s, expected=%s",
					error.Message, tt.expectedOutput,
				)
			}
		}
	}
}

func TestArrays(t *testing.T) {
	tests := []struct {
		input          string
		expectedOutput []int64
	}{
		{`[1,2];`, []int64{1, 2}},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		result, ok := evaluated.(*object.Array)
		if !ok {
			t.Fatalf(
				"object is not array, got=%T",
				evaluated,
			)
		}
		if len(result.Elements) != len(tt.expectedOutput) {
			t.Fatalf(
				"array has wrong amount of elements, expected=%d got=%d",
				len(tt.expectedOutput), len(result.Elements),
			)
		}
		for i, element := range result.Elements {
			testIntegerObject(t, element, tt.expectedOutput[i])
		}
	}
}

func TestArraysIndecies(t *testing.T) {
	tests := []struct {
		input          string
		expectedOutput interface{}
	}{
		{`[1,2][0];`, 1},
		{`[1, 2, 3][-1];`, 3},
		{`[1, 2, 3][4];`, nil},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expectedOutput.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func TestHashLiteral(t *testing.T) {
	input := `{"a": 42};`
	expectedOutput := map[object.HashKey]int64{
		(&object.String{Value: "a"}).HashKey(): 42,
	}

	evaluated := testEval(input)
	result, ok := evaluated.(*object.HashMap)
	if !ok {
		t.Errorf("evaluated object is not hashMap")
	}

	if len(result.Pairs) != len(expectedOutput) {
		t.Errorf(
			"wrong amount of pairs, expected=%d, received=%d",
			len(expectedOutput), len(result.Pairs),
		)
	}

	for expectedKey, expectedValue := range expectedOutput {
		pair, ok := result.Pairs[expectedKey]
		if !ok {
			t.Errorf("no pair for given key in Pairs")
		}
		testIntegerObject(t, pair.Value, expectedValue)
	}
}

func TestHashMapIndecies(t *testing.T) {
	tests := []struct {
		input          string
		expectedOutput interface{}
	}{
		{`{"1": 2}["1"];`, 2},
		{`{"1": 3 + 3}["1"];`, 6},
		{`{"1": 2}["check"];`, nil},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expectedOutput.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}
