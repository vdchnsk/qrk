package ast

import (
	"bytes"
	"strings"

	"github.com/vdchnsk/quasark/src/token"
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

type AssignStatement struct {
	Token      token.Token // "=" token
	Identifier *Identifier
	Value      Expression
}

func (as *AssignStatement) TokenLiteral() string { return as.Token.Literal }
func (as *AssignStatement) statementNode()       {}
func (as *AssignStatement) ToString() string {
	var out bytes.Buffer

	out.WriteString(as.Identifier.ToString())
	out.WriteString(" = ")

	if as.Value != nil {
		out.WriteString(as.Value.ToString())
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

type BlockStatement struct {
	Token      token.Token // the { token
	Statements []Statement
}

func (bs *BlockStatement) statementNode()       {}
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *BlockStatement) ToString() string {
	var out bytes.Buffer

	for _, s := range bs.Statements {
		out.WriteString(s.ToString())
	}

	return out.String()
}

type IntegerLiteral struct {
	Token token.Token // INT token
	Value int64
}

func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) ToString() string     { return il.Token.Literal }

type StringLiteral struct {
	Token token.Token // STRING token
	Value string
}

func (sl *StringLiteral) TokenLiteral() string { return sl.Token.Literal }
func (sl *StringLiteral) expressionNode()      {}
func (sl *StringLiteral) ToString() string     { return sl.Token.Literal }

type Boolean struct {
	Token token.Token
	Value bool
}

func (b *Boolean) TokenLiteral() string { return b.Token.Literal }
func (b *Boolean) expressionNode()      {}
func (b *Boolean) ToString() string     { return b.Token.Literal }

type PrefixExpression struct {
	Token    token.Token // The prfix token e.g "-"
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Literal }
func (pe *PrefixExpression) expressionNode()      {}
func (pe *PrefixExpression) ToString() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.ToString())
	out.WriteString(")")

	return out.String()
}

type InfixExpression struct {
	Token    token.Token // The prfix token e.g "-"
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *InfixExpression) expressionNode()      {}
func (ie *InfixExpression) ToString() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ie.Left.ToString())
	out.WriteString(" " + ie.Operator + " ")
	out.WriteString(ie.Right.ToString())
	out.WriteString(")")

	return out.String()
}

type IfExpression struct {
	Token       token.Token // The token "if"
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ife *IfExpression) TokenLiteral() string { return ife.Token.Literal }
func (ife *IfExpression) expressionNode()      {}
func (ife *IfExpression) ToString() string {
	var out bytes.Buffer

	out.WriteString("if ")
	out.WriteString(ife.Condition.ToString())
	out.WriteString(" ")
	out.WriteString(ife.Consequence.ToString())

	if ife.Alternative != nil {
		out.WriteString("else ")
		out.WriteString(ife.Alternative.ToString())
	}

	return out.String()
}

type FuncLiteral struct {
	Token      token.Token // The token "fn"
	Parameters []*Identifier
	Body       *BlockStatement
}

func (fl *FuncLiteral) TokenLiteral() string { return fl.Token.Literal }
func (fl *FuncLiteral) expressionNode()      {}
func (fl *FuncLiteral) ToString() string {
	var out bytes.Buffer

	out.WriteString("fn ")
	out.WriteString(fl.TokenLiteral())

	params := []string{}
	for _, parameter := range fl.Parameters {
		params = append(params, parameter.ToString())
	}
	out.WriteString("(")
	out.WriteString(strings.Join(params, ","))
	out.WriteString(")")

	out.WriteString(fl.Body.ToString())

	return out.String()
}

type CallExpression struct {
	Token    token.Token // "("
	Function Expression  // either Identifier or Function declaration
	Argments []Expression
}

func (ce *CallExpression) TokenLiteral() string { return ce.Token.Literal }
func (ce *CallExpression) expressionNode()      {}
func (ce *CallExpression) ToString() string {
	var out bytes.Buffer

	args := []string{}

	for _, arg := range ce.Argments {
		args = append(args, arg.ToString())
	}
	out.WriteString(ce.Function.ToString())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")

	return out.String()
}

type ArrayLiteral struct {
	Token    token.Token // "["
	Elements []Expression
}

func (al *ArrayLiteral) TokenLiteral() string { return al.Token.Literal }
func (al *ArrayLiteral) expressionNode()      {}
func (al *ArrayLiteral) ToString() string {
	var out bytes.Buffer

	elements := []string{}
	for _, elem := range al.Elements {
		elements = append(elements, elem.ToString())
	}
	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}

type HashMapLiteral struct {
	Token token.Token // "{"
	Pairs map[Expression]Expression
}

func (hml *HashMapLiteral) TokenLiteral() string { return hml.Token.Literal }
func (hml *HashMapLiteral) expressionNode()      {}
func (hml *HashMapLiteral) ToString() string {
	var out bytes.Buffer

	pairs := []string{}
	for key, value := range hml.Pairs {
		pairs = append(pairs, key.ToString()+":"+value.ToString())
	}
	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")

	return out.String()
}

type IndexExpression struct {
	Token token.Token // "["
	Left  Expression
	Index Expression
}

func (ie *IndexExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IndexExpression) expressionNode()      {}
func (ie *IndexExpression) ToString() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ie.Left.ToString())
	out.WriteString("[")
	out.WriteString(ie.Index.ToString())
	out.WriteString("]")

	return out.String()
}
