package parser

import (
	"fmt"
	"testing"

	"github.com/vdchnsk/i-go/ast"
	"github.com/vdchnsk/i-go/lexer"
)

func TestLetStatement(t *testing.T) {
	input := `
		let x = 5;
		let y = 5;
		let fooBar = 5;
	`
	lexer := lexer.NewLexer(input)
	parser := NewParser(lexer)

	program := parser.ParseProgram()
	checkParserErrors(t, parser)
	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}
	if len(program.Statements) != 3 {
		t.Fatalf(
			"program.Statements does not contain 3 statements. got=%d",
			len(program.Statements),
		)
	}

	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"fooBar"},
	}

	for i, tt := range tests {
		statement := program.Statements[i]
		if !testLetStatement(t, statement, tt.expectedIdentifier) {
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
	input := `
		return 5;
		return 10;
		return 42;
	`
	lexer := lexer.NewLexer(input)
	parser := NewParser(lexer)

	program := parser.ParseProgram()
	checkParserErrors(t, parser)
	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}
	if len(program.Statements) != 3 {
		t.Fatalf(
			"program.Statements does not contain 3 statements. got=%d",
			len(program.Statements),
		)
	}

	for _, statement := range program.Statements {
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
			intLit,
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
			"!-a",
			"(!(-a))",
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
	}

	for _, tt := range tests {
		lexer := lexer.NewLexer(tt.input)
		parser := NewParser(lexer)

		program := parser.ParseProgram()
		checkParserErrors(t, parser)

		actualProgram := program.ToString()

		if actualProgram != tt.expectedProgram {
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
