package ast

import (
	"bytes"

	"github.com/vdchnsk/i-go/token"
)

type Node interface {
	TokenLiteral() string
	ToString() string
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
func (p *Program) ToString() string {
	var out bytes.Buffer

	for _, s := range p.Statements {
		out.WriteString(s.ToString())
	}

	return out.String()
}

type Identifier struct {
	Token token.Token
	Value string
}

func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) expressionNode()      {}
func (i *Identifier) ToString() string     { return i.Value }

type LetStatement struct {
	Token      token.Token // "let" token
	Identifier *Identifier
	Value      Expression
}

func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }
func (ls *LetStatement) statementNode()       {}
func (ls *LetStatement) ToString() string {
	var out bytes.Buffer

	out.WriteString(ls.TokenLiteral() + " ")
	out.WriteString(ls.Identifier.ToString())
	out.WriteString(" = ")

	if ls.Value != nil {
		out.WriteString(ls.Value.ToString())
	}

	out.WriteString(";")

	return out.String()
}

type ReturnStatement struct {
	Token token.Token // "return" token
	Value Expression
}

func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }
func (rs *ReturnStatement) statementNode()       {}
func (rs *ReturnStatement) ToString() string {
	var out bytes.Buffer

	out.WriteString(rs.TokenLiteral() + " ")

	if rs.Value != nil {
		out.WriteString(rs.Value.ToString())
	}

	out.WriteString(";")

	return out.String()
}

type ExpressionStatement struct {
	Token token.Token // first token of the expression
	Value Expression
}

func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }
func (es *ExpressionStatement) statementNode()       {}
func (es *ExpressionStatement) ToString() string {
	if es.Value == nil {
		return ""
	}

	return es.Value.ToString()
}

type IntegerLiteral struct {
	Token token.Token // INT token
	Value int64
}

func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) ToString() string     { return il.Token.Literal }
