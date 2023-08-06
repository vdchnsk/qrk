package ast

import "github.com/vdchnsk/i-go/token"

type Node interface {
	TokenLiteral() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

type Identifier struct {
	Token token.Token
	Value string
}

func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) expressionNode()      {}

type LetStatement struct {
	Token      token.Token // "let" token
	Identifier *Identifier
	Value      Expression
}

func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }
func (ls *LetStatement) statementNode()       {}

type ReturnStatement struct {
	Token token.Token // "return" token
	Value Expression
}

func (ls *ReturnStatement) TokenLiteral() string { return ls.Token.Literal }
func (ls *ReturnStatement) statementNode()       {}
