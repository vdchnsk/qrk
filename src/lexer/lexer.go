package lexer

import (
	"errors"

	"github.com/vdchnsk/qrk/src/token"
	"github.com/vdchnsk/qrk/src/utils"
)

type Lexer struct {
	input            string
	position         int
	currReadPosition int
	currChar         byte
}

func NewLexer(input string) *Lexer {
	l := &Lexer{input: input, currReadPosition: 0, position: 0}
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

func (l *Lexer) NextToken() (token.Token, error) {
	var tok token.Token

	l.skipWhitespace()

	switch l.currChar {
	case '=':
		currChar := l.currChar
		switch l.peekChar() {
		case '=':
			l.readChar()
			tok = token.Token{Type: token.EQ, Literal: string(currChar) + string(currChar)}
		case '>':
			l.readChar()
			tok = token.Token{Type: token.ARROW, Literal: "=>"}
		default:
			tok = newToken(token.ASSIGN, currChar)
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
		isComparison := l.peekChar() == '='
		if isComparison {
			l.readChar()
			tok = token.Token{Type: token.NOT_EQ, Literal: string(currChar) + string(l.currChar)}
		} else {
			tok = newToken(token.BANG, currChar)
		}
	case '&':
		currChar := l.currChar
		isAndOperator := l.peekChar() == currChar
		if isAndOperator {
			l.readChar()
			tok = token.Token{Type: token.AND, Literal: string(currChar) + string(l.currChar)}
		} else {
			tok = newToken(token.ILLEGAL, l.currChar)
		}
	case '|':
		currChar := l.currChar
		isOrOperator := l.peekChar() == currChar
		if isOrOperator {
			l.readChar()
			tok = token.Token{Type: token.OR, Literal: string(currChar) + string(l.currChar)}
		} else {
			tok = newToken(token.ILLEGAL, l.currChar)
		}
	case '(':
		tok = newToken(token.LPAREN, l.currChar)
	case ')':
		tok = newToken(token.RPAREN, l.currChar)
	case '{':
		tok = newToken(token.LBRACE, l.currChar)
	case '}':
		tok = newToken(token.RBRACE, l.currChar)
	case '[':
		tok = newToken(token.LBRACKET, l.currChar)
	case ']':
		tok = newToken(token.RBRACKET, l.currChar)
	case ';':
		tok = newToken(token.SEMICOLON, l.currChar)
	case ':':
		tok = newToken(token.COLON, l.currChar)
	case ',':
		tok = newToken(token.COMMA, l.currChar)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	case '"':
		strValue, err := l.readString()
		if err != nil {
			return token.Token{}, err
		}
		tok.Type = token.STRING
		tok.Literal = strValue
	default:
		if isLetter(l.currChar) {
			identifier := l.readIdentifier()

			tok.Type = token.LookupIdentifier(identifier)
			tok.Literal = identifier
			return tok, nil
		} else if isDigit(l.currChar) {
			tok.Type = token.INT
			tok.Literal = l.readNumber()
			return tok, nil
		} else {
			tok = newToken(token.ILLEGAL, l.currChar)
		}
	}

	l.readChar()

	return tok, nil
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

func (l *Lexer) readString() (string, error) {
	initialPosition := l.position // "

	for {
		l.readChar()
		if l.currChar == 0 {
			return "", errors.New("no closing stirng symbol was found")
		}

		if l.currChar == '"' {
			break
		}
	}

	startValuePosition := initialPosition + 1
	endValuePosition := l.position

	return l.input[startValuePosition:endValuePosition], nil
}

func (l *Lexer) readIdentifier() string {
	initialPosition := l.position

	for isLetter(l.currChar) {
		l.readChar()
	}

	identEndPosition := l.position
	return l.input[initialPosition:identEndPosition]
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}
