package parser

import (
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
		t.Errorf("s.TokenLiteral not 'let'. got=%q", statement.TokenLiteral())
		return false
	}

	letStatement, ok := statement.(*ast.LetStatement)

	if !ok {
		t.Errorf("statement is not *ast.Statement. got=%T", statement)
	}

	if letStatement.Identifier.Value != identifier {
		t.Errorf("letStatement.Identifier.Value not '%s'. got=%s", identifier, letStatement.Identifier.Value)
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
