package lexer

import (
	"github.com/vdchnsk/i-go/token"
	"github.com/vdchnsk/i-go/utils"
)

type Lexer struct {
	input            string
	position         int
	currReadPosition int
	currChar         byte
}

func NewLexer(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()

	return l
}

func (l *Lexer) readChar() {
	if l.currReadPosition >= len(l.input) {
		l.currChar = 0
	} else {
		l.currChar = l.input[l.currReadPosition]
	}

	l.position = l.currReadPosition
	l.currReadPosition += 1
}

func (l *Lexer) peekChar() byte {
	if l.currReadPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.currReadPosition]
	}
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	switch l.currChar {
	case '=':
		if l.peekChar() == '=' {
			currChar := l.currChar
			l.readChar()
			tok = token.Token{Type: token.EQ, Literal: string(currChar) + string(currChar)}
		} else {
			tok = newToken(token.ASSIGN, l.currChar)
		}
	case '+':
		tok = newToken(token.PLUS, l.currChar)
	case '-':
		tok = newToken(token.MINUS, l.currChar)
	case '/':
		tok = newToken(token.SLASH, l.currChar)
	case '*':
		tok = newToken(token.ASTERISK, l.currChar)
	case '<':
		tok = newToken(token.LT, l.currChar)
	case '>':
		tok = newToken(token.GT, l.currChar)
	case '!':
		currChar := l.currChar
		if l.peekChar() == '=' {
			l.readChar()
			tok = token.Token{Type: token.NOT_EQ, Literal: string(currChar) + string(l.currChar)}
		} else {
			tok = newToken(token.BANG, l.currChar)
		}
	case '(':
		tok = newToken(token.LPAREN, l.currChar)
	case ')':
		tok = newToken(token.RPAREN, l.currChar)
	case '{':
		tok = newToken(token.LBRACE, l.currChar)
	case '}':
		tok = newToken(token.RBRACE, l.currChar)
	case ';':
		tok = newToken(token.SEMICOLON, l.currChar)
	case ',':
		tok = newToken(token.COMMA, l.currChar)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isLetter(l.currChar) {
			identifier := l.readIdentifier()

			tok.Type = token.LookupIdentifier(identifier)
			tok.Literal = identifier
			return tok
		} else if isDigit(l.currChar) {
			tok.Type = token.INT
			tok.Literal = l.readNumber()
			return tok
		} else {
			tok = newToken(token.ILLEGAL, l.currChar)
		}
	}

	l.readChar()

	return tok
}

func newToken(tokenType token.TokenType, char byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(char)}
}

func (l *Lexer) skipWhitespace() {
	whitespaces := []byte{' ', '\t', '\r', '\n'}

	for utils.Contains[byte](whitespaces, l.currChar) {
		l.readChar()
	}
}

func (l *Lexer) readNumber() string {
	initialPosition := l.position

	for isDigit(l.currChar) {
		l.readChar()
	}

	return l.input[initialPosition:l.position]
}

func (l *Lexer) readIdentifier() string {
	initialPosition := l.position

	for isLetter(l.currChar) {
		l.readChar()
	}

	return l.input[initialPosition:l.position]
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}
