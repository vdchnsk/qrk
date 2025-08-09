package parser

import (
	"fmt"
	"testing"

	"github.com/vdchnsk/qrk/src/ast"
	"github.com/vdchnsk/qrk/src/lexer"
	"github.com/vdchnsk/qrk/src/utils"
)

func TestLetStatement(t *testing.T) {
	tests := []struct {
		input              string
		expectedIdentifier string
		expectedValue      interface{}
	}{
		{"let x = 5;", "x", 5},
		{"let y = true;", "y", true},
		{"let foo_bar = pi;", "foo_bar", "pi"},
	}

	for _, tt := range tests {
		lexer := lexer.NewLexer(tt.input)
		parser := NewParser(lexer)
		program := parser.ParseProgram()

		checkParserErrors(t, parser)

		if len(program.Statements) != 1 {
			t.Fatalf(
				"program.Statements does not contain 1 statements. got=%d",
				len(program.Statements),
			)
		}

		statement := program.Statements[0]
		if !testLetStatement(t, statement, tt.expectedIdentifier) {
			return
		}

		val := statement.(*ast.LetStatement).Value

		if !testLiteralExpression(t, val, tt.expectedValue) {
			return
		}
	}
}

func testLetStatement(t *testing.T, statement ast.Statement, identifier string) bool {
	if statement.TokenLiteral() != "let" {
		t.Errorf(
			"s.TokenLiteral not 'let'. got=%q",
			statement.TokenLiteral(),
		)
		return false
	}

	letStatement, ok := statement.(*ast.LetStatement)

	if !ok {
		t.Errorf(
			"statement is not *ast.Statement. got=%T",
			statement,
		)
	}

	if letStatement.Identifier.Value != identifier {
		t.Errorf(
			"letStatement.Identifier.Value not '%s'. got=%s",
			identifier, letStatement.Identifier.Value,
		)
		return false
	}

	return true
}

func checkParserErrors(t *testing.T, parser *Parser) {
	errors := parser.Errors()

	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors,", len(errors))

	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}

	t.FailNow()
}

func TestReturnStatement(t *testing.T) {
	tests := []struct {
		input         string
		expectedValue interface{}
	}{
		{"return 5;", 5},
		{"return 10;", 10},
		{"return true;", true},
	}

	for _, tt := range tests {
		lexer := lexer.NewLexer(tt.input)
		parser := NewParser(lexer)
		program := parser.ParseProgram()

		checkParserErrors(t, parser)

		if len(program.Statements) != 1 {
			t.Fatalf(
				"program.Statements does not contain 1 statements. got=%d",
				len(program.Statements),
			)
		}

		statement := program.Statements[0]

		returnStatement, ok := statement.(*ast.ReturnStatement)
		if !ok {
			t.Errorf(
				"statement is not *ast.ReturnStatement. got=%T",
				statement,
			)
		}

		if returnStatement.TokenLiteral() != "return" {
			t.Errorf(
				"returnStmt.TokenLiteral() not 'return', got %q",
				returnStatement.TokenLiteral(),
			)
		}
		val := returnStatement.Value
		if !testLiteralExpression(t, val, tt.expectedValue) {
			return
		}
	}

}

func TestIdentifierExpression(t *testing.T) {
	input := "foobar;"

	lexer := lexer.NewLexer(input)
	parser := NewParser(lexer)
	program := parser.ParseProgram()

	checkParserErrors(t, parser)

	expectedAmountOfStatements := 1
	if len(program.Statements) != expectedAmountOfStatements {
		t.Fatalf(
			"program has wrong amount of statements, expected %d, got=%d",
			expectedAmountOfStatements, (program.Statements),
		)
	}
	statement, ok := program.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.ExpressionStatement")
	}

	ident, ok := statement.Value.(*ast.Identifier)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.Identifier")
	}
	if ident.Value != "foobar" {
		t.Fatalf(
			"ident.Value not %s. got=%s",
			"foobar", ident.Value,
		)
	}
	if ident.TokenLiteral() != "foobar" {
		t.Errorf(
			"ident.TokenLiteral() is not %s, got=%s",
			"foobar", ident.TokenLiteral(),
		)
	}
}

func TestIntegerExpression(t *testing.T) {
	input := "5;"

	lexer := lexer.NewLexer(input)
	parser := NewParser(lexer)
	program := parser.ParseProgram()

	checkParserErrors(t, parser)

	expectedAmountOfStatements := 1
	if len(program.Statements) != expectedAmountOfStatements {
		t.Fatalf(
			"program has wrong amount of statements, expected %d, got=%d",
			expectedAmountOfStatements, (program.Statements),
		)
	}
	statement, ok := program.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.ExpressionStatement")
	}

	ident, ok := statement.Value.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.IntegerLiteral")
	}
	if ident.Value != 5 {
		t.Fatalf(
			"ident.Value not %d. got=%d",
			5, ident.Value,
		)
	}
	if ident.TokenLiteral() != "5" {
		t.Errorf(
			"ident.TokenLiteral() is not %s, got=%s",
			"5", ident.TokenLiteral(),
		)
	}
}

func TestStringExpression(t *testing.T) {
	input := `"Linus Torvalds"`
	expectedOutput := "Linus Torvalds"

	lexer := lexer.NewLexer(input)
	parser := NewParser(lexer)
	program := parser.ParseProgram()

	checkParserErrors(t, parser)

	expectedAmountOfStatements := 1
	if len(program.Statements) != expectedAmountOfStatements {
		t.Fatalf(
			"program has wrong amount of statements, expected %d, got=%d",
			expectedAmountOfStatements, (program.Statements),
		)
	}
	statement, ok := program.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.ExpressionStatement")
	}

	ident, ok := statement.Value.(*ast.StringLiteral)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.StringLiteral")
	}
	if ident.Value != expectedOutput {
		t.Fatalf(
			"ident.Value not %s. got=%s",
			expectedOutput, ident.Value,
		)
	}
	if ident.TokenLiteral() != expectedOutput {
		t.Errorf(
			"ident.TokenLiteral() is not %s, got=%s",
			expectedOutput, ident.TokenLiteral(),
		)
	}
}

func TestParsingPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input    string
		operator string
		value    interface{}
	}{
		{"!5;", "!", 5},
		{"-42;", "-", 42},
		{"!true", "!", true},
		{"!false", "!", false},
	}

	for _, tt := range prefixTests {
		lexer := lexer.NewLexer(tt.input)
		parser := NewParser(lexer)
		program := parser.ParseProgram()

		checkParserErrors(t, parser)

		expectedAmountOfStatements := 1
		if len(program.Statements) != 1 {
			t.Fatalf(
				"program.Statemtnts length=%d, expected=%d",
				len(program.Statements), expectedAmountOfStatements,
			)
		}

		statement, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf(
				"program.Statements[0] is not ast.ExpressionStatement, got=%T",
				program.Statements[0],
			)
		}

		expression, ok := statement.Value.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf(
				"program.Statements[0] is not ast.PrefixExpression, got=%T",
				statement.Value,
			)
		}
		if expression.Operator != tt.operator {
			t.Fatalf(
				"expression.Operator is not '%s'. got=%s",
				tt.operator, expression.Operator,
			)
		}

		if !testLiteralExpression(t, expression.Right, tt.value) {
			return
		}
	}
}

func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {
	intLit, ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf(
			"il is not *ast.IntegerLiteral, got=%T",
			il,
		)
		return false
	}

	if intLit.Value != value {
		t.Errorf(
			"intlit.Value is not %d, got=%d",
			value, intLit.Value,
		)
	}

	if intLit.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf(
			"intLit.TokenLiteral is not %d, got=%s",
			value, intLit.TokenLiteral(),
		)
		return false
	}

	return true
}

func testIdentifier(t *testing.T, exp ast.Expression, value string) bool {
	ident, ok := exp.(*ast.Identifier)
	if !ok {
		t.Errorf(
			"exp not *ast.Identifier, got=%T",
			exp,
		)
		return false
	}

	if ident.Value != value {
		t.Errorf(
			"ident.Value %s != Value %s",
			ident.Value, value,
		)
		return false
	}

	if ident.TokenLiteral() != value {
		t.Errorf(
			"ident.TokenLiteral %s != Value %s",
			ident.TokenLiteral(), value,
		)
		return false
	}

	return true
}

func testLiteralExpression(t *testing.T, exp ast.Expression, expected interface{}) bool {
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, exp, int64(v))
	case int64:
		return testIntegerLiteral(t, exp, v)
	case string:
		return testIdentifier(t, exp, v)
	case bool:
		return testBooleanLiteral(t, exp, v)
	}
	t.Errorf(
		"type of exp not handled, got=%T",
		exp,
	)
	return false
}

func testBooleanLiteral(t *testing.T, exp ast.Expression, value bool) bool {
	bool, ok := exp.(*ast.Boolean)
	if !ok {
		t.Errorf(
			"exp not *ast.Boolean. got=%T",
			exp,
		)
		return false
	}
	if bool.Value != value {
		t.Errorf(
			"bo.Value not %t. got=%t",
			value, bool.Value,
		)
		return false
	}
	if bool.Value != value {
		t.Errorf(
			"bool.TokenLiteral() is not %t. got=%s",
			value, bool.TokenLiteral(),
		)
		return false
	}

	return true
}

func testInfixExpression(
	t *testing.T,
	exp ast.Expression,
	left interface{},
	operator string,
	right interface{},
) bool {
	opExp, ok := exp.(*ast.InfixExpression)
	if !ok {
		t.Errorf("exp is not ast.OperatorExpression. got=%T(%s)", exp, exp)
		return false
	}
	if !testLiteralExpression(t, opExp.Left, left) {
		return false
	}
	if opExp.Operator != operator {
		t.Errorf("exp.Operator is not '%s', got=%q", operator, opExp.Operator)
		return false
	}
	if !testLiteralExpression(t, opExp.Right, right) {
		return false
	}

	return true
}

func TestParsingInfixExpressions(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  interface{}
		operator   string
		rightValue interface{}
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
		{"true == true", true, "==", true},
		{"true && false", true, "&&", false},
		{"true || true", true, "||", true},
		{"true || false", true, "||", false},
	}

	for _, tt := range infixTests {
		lexer := lexer.NewLexer(tt.input)
		parser := NewParser(lexer)

		program := parser.ParseProgram()
		checkParserErrors(t, parser)

		expectedAmountOfStatements := 1
		if len(program.Statements) != expectedAmountOfStatements {
			t.Fatalf(
				"program.Statements does not contain %d statements. got=%d\n",
				expectedAmountOfStatements, len(program.Statements),
			)
		}

		statement, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf(
				"program.Statements[0] is not ast.ExpressionStatement. got=%T",
				program.Statements[0],
			)
		}
		expression, ok := statement.Value.(*ast.InfixExpression)
		if !ok {
			t.Fatalf(
				"expression is not ast.InfixExpression, got=%T",
				statement.Value,
			)
		}

		if (!testLiteralExpression(t, expression.Left, tt.leftValue)) ||
			(!testLiteralExpression(t, expression.Right, tt.rightValue)) {
			return
		}

		if expression.Operator != tt.operator {
			t.Fatalf(
				"expression.Operator is not '%s'. got=%s",
				tt.operator, expression.Operator,
			)
		}
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input           string
		expectedProgram string
	}{
		{
			"-a * b",
			"((-a) * b)",
		},
		{
			"(a + b) * c",
			"((a + b) * c)",
		},
		{
			"1 + (2 + 3) + 4",
			"((1 + (2 + 3)) + 4)",
		},
		{
			"1 + (2 + 3) * 4",
			"(1 + ((2 + 3) * 4))",
		},
		{
			"!-a",
			"(!(-a))",
		},
		{
			"-(5 + 5)",
			"(-(5 + 5))",
		},
		{
			"a + b + c",
			"((a + b) + c)",
		},
		{
			"a + b - c",
			"((a + b) - c)",
		},
		{
			"a * b * c",
			"((a * b) * c)",
		},
		{
			"a * b / c",
			"((a * b) / c)",
		},
		{
			"a + b / c",
			"(a + (b / c))",
		},
		{
			"a + b * c + d / e - f",
			"(((a + (b * c)) + (d / e)) - f)",
		},
		{
			"3 + 4; -5 * 5",
			"(3 + 4)((-5) * 5)",
		},
		{
			"5 > 4 == 3 < 4",
			"((5 > 4) == (3 < 4))",
		},
		{
			"5 < 4 != 3 > 4",
			"((5 < 4) != (3 > 4))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			"true",
			"true",
		},
		{
			"false",
			"false",
		},
		{
			"3 > 5 == false",
			"((3 > 5) == false)",
		},
		{
			"3 < 5 == true",
			"((3 < 5) == true)",
		},
		{
			"!(true == true)",
			"(!(true == true))",
		},
		{
			"a + add(b * c, e) + d",
			"((a + add((b * c), e)) + d)",
		},
		{
			"getArray(b * c, e)[1 + 1] + d",
			"((getArray((b * c), e)[(1 + 1)] + d)",
		},
	}

	for _, tt := range tests {
		lexer := lexer.NewLexer(tt.input)
		parser := NewParser(lexer)

		program := parser.ParseProgram()
		checkParserErrors(t, parser)

		actualProgram := program.String()

		if utils.RemoveWhitespaces(actualProgram) != utils.RemoveWhitespaces(tt.expectedProgram) {
			t.Fatalf(
				"got program output=%s, expected=%s",
				actualProgram, tt.expectedProgram,
			)
		}
	}
}

func TestBooleanExpression(t *testing.T) {
	input := "true;"

	lexer := lexer.NewLexer(input)
	parser := NewParser(lexer)
	program := parser.ParseProgram()

	checkParserErrors(t, parser)

	expectedAmountOfStatements := 1
	if len(program.Statements) != expectedAmountOfStatements {
		t.Fatalf(
			"program has wrong amount of statements, expected %d, got=%d",
			expectedAmountOfStatements, (program.Statements),
		)
	}
	statement, ok := program.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.ExpressionStatement")
	}

	ident, ok := statement.Value.(*ast.Boolean)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.Identifier")
	}
	if ident.Value != true {
		t.Fatalf(
			"ident.Value not %t. got=%t",
			true, ident.Value,
		)
	}
	if ident.TokenLiteral() != "true" {
		t.Errorf(
			"ident.TokenLiteral() is not %s, got=%s",
			"foobar", ident.TokenLiteral(),
		)
	}
}

func TestIfElseExpression(t *testing.T) {
	input := `if x < y { x } else { y }`

	lexer := lexer.NewLexer(input)
	parser := NewParser(lexer)

	program := parser.ParseProgram()
	checkParserErrors(t, parser)

	expectedAmountOfStatements := 1
	if len(program.Statements) != expectedAmountOfStatements {
		t.Fatalf(
			"program.Statemtnts length=%d, expected=%d",
			len(program.Statements), expectedAmountOfStatements,
		)
	}

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf(
			"program.statements is not ast.expressionstatement, got=%T",
			program.Statements[0],
		)
	}
	expression, ok := statement.Value.(*ast.IfExpression)
	if !ok {
		t.Fatalf(
			"statement.Value is not ast.IfExpression, got=%T",
			statement.Value,
		)
	}

	if len(expression.Consequence.Statements) != 1 {
		t.Fatalf(
			"wrong amount of consequence statemetns got=%d",
			len(expression.Consequence.Statements),
		)
	}
	conseqeneceExpression, ok := expression.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf(
			"conseqence is not expression got=%T",
			expression.Consequence.Statements[0],
		)
	}
	if !testIdentifier(t, conseqeneceExpression.Value, "x") {
		t.Fatalf(
			"identifer of consequence expression in is not, got=%T", // TODO
			conseqeneceExpression.Value,
		)
	}
	if len(expression.Alternative.Statements) != 1 {
		t.Fatalf(
			"expected 1 alternative, got=%T",
			expression.Alternative.Statements,
		)
	}

	alternativeExpression, ok := expression.Alternative.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf(
			"conseqence is not expression got=%T",
			expression.Alternative.Statements[0],
		)
	}
	if !testIdentifier(t, alternativeExpression.Value, "y") {
		t.Fatalf(
			"identifer of consequence expression in is not, got=%T", // TODO
			conseqeneceExpression.Value,
		)
	}
}

func TestFuncLiteral(t *testing.T) {
	input := `fn(x, y) { x + y; }`

	lexer := lexer.NewLexer(input)
	parser := NewParser(lexer)

	program := parser.ParseProgram()
	checkParserErrors(t, parser)

	if len(program.Statements) != 1 {
		t.Fatalf(
			"program.Body does not contain %d statements, got %d\n",
			1, len(program.Statements),
		)
	}

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf(
			"program.Statements[0] is not ast.ExpressionStatement, got=%T",
			program.Statements[0],
		)
	}

	function, ok := statement.Value.(*ast.FuncLiteral)
	if !ok {
		t.Fatalf(
			"function is not ast.FuncLiteral, got=%T",
			statement.Value,
		)
	}
	if len(function.Parameters) != 2 {
		t.Fatalf(
			"function received unexpected amount of parameters, waited 2, got=%d\n",
			len(function.Parameters),
		)
	}

	testLiteralExpression(t, function.Parameters[0], "x")
	testLiteralExpression(t, function.Parameters[1], "y")

	if len(function.Body.Statements) != 1 {
		t.Fatalf(
			"function body contains unexpected amount of statements, waited 1, got=%d\n",
			len(function.Body.Statements),
		)
	}

	functionBodyStatement, ok := function.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf(
			"function.Body.Statements[0] is not ast.BlockStatement, got=%T",
			len(function.Body.Statements),
		)
	}

	testInfixExpression(t, functionBodyStatement.Value, "x", "+", "y")
}

func TestFuncParams(t *testing.T) {
	tests := []struct {
		input          string
		expectedOutput []string
	}{
		{"fn() { }", []string{}},
		{"fn(x) { x; }", []string{"x"}},
		{"fn(x, y) { x + y; }", []string{"x", "y"}},
	}

	for _, tt := range tests {
		lexer := lexer.NewLexer(tt.input)
		parser := NewParser(lexer)
		program := parser.ParseProgram()
		checkParserErrors(t, parser)

		if len(program.Statements) != 1 {
			t.Fatalf(
				"program.Body does not contain %d statements, got %d\n",
				1, len(program.Statements),
			)
		}

		statement, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf(
				"program.Statements[0] is not ast.ExpressionStatement, got=%T",
				program.Statements[0],
			)
		}

		function, ok := statement.Value.(*ast.FuncLiteral)
		if !ok {
			t.Fatalf(
				"function is not ast.FuncLiteral, got=%T",
				statement.Value,
			)
		}
		funcParams := function.Parameters

		if len(funcParams) != len(tt.expectedOutput) {
			t.Fatalf(
				"wrong number of parameters, expectedd=%d, got=%d,",
				len(tt.expectedOutput),
				len(funcParams),
			)
		}

		for index := range funcParams {
			testLiteralExpression(t, funcParams[index], tt.expectedOutput[index])
		}
	}
}

func TestCallExpression(t *testing.T) {
	input := "add(3, 14);"

	lexer := lexer.NewLexer(input)
	parser := NewParser(lexer)
	program := parser.ParseProgram()

	checkParserErrors(t, parser)

	if len(program.Statements) != 1 {
		t.Fatalf(
			"program.Body does not contain %d statements, got %d\n",
			1, len(program.Statements),
		)
	}

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf(
			"program.Statements[0] is not ast.ExpressionStatement, got=%T",
			program.Statements[0],
		)
	}

	callExpression, ok := statement.Value.(*ast.CallExpression)
	if !ok {
		t.Fatalf(
			"function is not ast.CallExpession, got=%T",
			statement.Value,
		)
	}
	if !testIdentifier(t, callExpression.Function, "add") {
		t.Fatalf(
			"expression got wrong identifier, expected=%s",
			"add",
		)
	}

	if len(callExpression.Arguments) != 2 {
		t.Fatalf(
			"call expression got unexpected amount of arguments, expected=%d got=%d",
			2,
			len(callExpression.Arguments),
		)
	}
	testLiteralExpression(t, callExpression.Arguments[0], 3)
	testLiteralExpression(t, callExpression.Arguments[1], 14)
}

func TestCallExpressionParams(t *testing.T) {
	tests := []struct {
		input          string
		expectedOutput []int
	}{
		{"call(3, 14);", []int{3, 14}},
	}

	for _, tt := range tests {
		lexer := lexer.NewLexer(tt.input)
		parser := NewParser(lexer)
		program := parser.ParseProgram()
		checkParserErrors(t, parser)

		if len(program.Statements) != 1 {
			t.Fatalf(
				"program.Body does not contain %d statements, got %d\n",
				1, len(program.Statements),
			)
		}

		statement, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf(
				"program.Statements[0] is not ast.ExpressionStatement, got=%T",
				program.Statements[0],
			)
		}

		callExpression, ok := statement.Value.(*ast.CallExpression)
		if !ok {
			t.Fatalf(
				"expression is not ast.CallExpression, got=%T",
				statement.Value,
			)
		}
		callExprArguments := callExpression.Arguments

		if len(callExprArguments) != len(tt.expectedOutput) {
			t.Fatalf(
				"wrong number of arguments, expected=%d, got=%d,",
				len(tt.expectedOutput),
				len(callExprArguments),
			)
		}

		for index := range callExprArguments {
			testLiteralExpression(t, callExprArguments[index], tt.expectedOutput[index])
		}
	}
}

func TestArrayLiteral(t *testing.T) {
	tests := []struct {
		input          string
		expectedOutput []int
	}{
		{"[1, 2];", []int{1, 2}},
	}

	for _, tt := range tests {
		lexer := lexer.NewLexer(tt.input)
		parser := NewParser(lexer)
		program := parser.ParseProgram()
		checkParserErrors(t, parser)

		if len(program.Statements) != 1 {
			t.Fatalf(
				"program.Body does not contain %d statements, got %d\n",
				1, len(program.Statements),
			)
		}

		statement, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf(
				"program.Statements[0] is not ast.ExpressionStatement, got=%T",
				program.Statements[0],
			)
		}

		array, ok := statement.Value.(*ast.ArrayLiteral)
		if !ok {
			t.Fatalf(
				"statement is not array literal, got=%T",
				statement.Value,
			)
		}

		if len(array.Elements) != len(tt.expectedOutput) {
			t.Fatalf(
				"unexpected amount of elements in the array, expected=%d, got=%d",
				len(tt.expectedOutput), len(array.Elements),
			)
		}
	}
}

func TestIndexExpression(t *testing.T) {
	input := "myArray[1 + 1]"

	lexer := lexer.NewLexer(input)
	parser := NewParser(lexer)
	program := parser.ParseProgram()
	checkParserErrors(t, parser)

	statement, _ := program.Statements[0].(*ast.ExpressionStatement)
	indexExp, ok := statement.Value.(*ast.IndexExpression)
	if !ok {
		t.Fatalf(
			"expression is not *ast.IndexExpression, received=%T",
			statement.Value,
		)
	}

	if !testIdentifier(t, indexExp.Left, "myArray") {
		return
	}

	if !testInfixExpression(t, indexExp.Index, 1, "+", 1) {
		return
	}
}

func TestHashMapLiteral(t *testing.T) {
	input := `{ "check": 1, "69": 42 }`
	expectedOutput := map[string]int64{
		"check": 1,
		"69":    42,
	}

	lexer := lexer.NewLexer(input)
	parser := NewParser(lexer)
	program := parser.ParseProgram()
	checkParserErrors(t, parser)

	statement := program.Statements[0].(*ast.ExpressionStatement)
	hashMap, ok := statement.Value.(*ast.HashMapLiteral)
	if !ok {
		t.Errorf(
			"provided statement is not HashMapLiteral, got=%T",
			statement,
		)
	}
	amountOfPairs := len(hashMap.Pairs)

	if amountOfPairs != 2 {
		t.Errorf(
			"wrong amount of has map pairs, got=%d",
			amountOfPairs,
		)
	}

	for key, value := range hashMap.Pairs {
		strLiteral, ok := key.(*ast.StringLiteral)
		if !ok {
			t.Errorf("key is not string")
		}

		expectedValue := expectedOutput[strLiteral.String()]

		testIntegerLiteral(t, value, expectedValue)
	}
}

func TestEmptyHashMap(t *testing.T) {
	input := `{}`

	lexer := lexer.NewLexer(input)
	parser := NewParser(lexer)
	program := parser.ParseProgram()

	statement := program.Statements[0].(*ast.ExpressionStatement)
	hashMap, ok := statement.Value.(*ast.HashMapLiteral)
	if !ok {
		t.Errorf(
			"provided statement is not HashMapLiteral, got=%T",
			statement,
		)
	}
	amountOfPairs := len(hashMap.Pairs)

	if amountOfPairs != 0 {
		t.Errorf(
			"wrong amount of has map pairs, got=%d",
			amountOfPairs,
		)
	}
}
