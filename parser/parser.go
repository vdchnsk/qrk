package parser

import (
	"fmt"

	"github.com/vdchnsk/i-go/ast"
	"github.com/vdchnsk/i-go/lexer"
	"github.com/vdchnsk/i-go/token"
)

type Parser struct {
	lexer *lexer.Lexer

	currToken token.Token
	peekToken token.Token

	errors []string
}

func NewParser(l *lexer.Lexer) *Parser {
	p := &Parser{
		lexer:  l,
		errors: []string{},
	}

	// call NextToken twice to have currToken and peekToken both set
	p.NextToken()
	p.NextToken()

	return p
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) PeekError(t token.TokenType) {
	msg := fmt.Sprintf(
		"expected next token to be %s, got %s instead",
		t, p.peekToken.Type,
	)
	p.errors = append(p.errors, msg)
}

func (p *Parser) NextToken() {
	p.currToken = p.peekToken
	p.peekToken = p.lexer.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for !p.currTokenIs(token.EOF) {
		statement := p.parseStatement()

		if statement != nil {
			program.Statements = append(program.Statements, statement)
		}
		p.NextToken()
	}

	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.currToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return nil
	}
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	statement := &ast.LetStatement{Token: p.currToken}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	statement.Identifier = &ast.Identifier{Token: p.currToken, Value: p.currToken.Literal}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	for !p.currTokenIs(token.SEMICOLON) {
		p.NextToken()
	}

	return statement
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	statement := &ast.ReturnStatement{Token: p.currToken}

	p.NextToken()

	for !p.currTokenIs(token.SEMICOLON) {
		p.NextToken()
	}

	return statement
}

func (p *Parser) currTokenIs(expectedCurrentToken token.TokenType) bool {
	return p.currToken.Type == expectedCurrentToken
}

func (p *Parser) expectPeek(expectedToken token.TokenType) bool {
	if p.peekToken.Type == expectedToken {
		p.NextToken()
		return true
	}

	p.PeekError(expectedToken)
	return false
}
