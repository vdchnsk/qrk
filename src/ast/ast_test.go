package ast

import (
	"testing"

	"github.com/vdchnsk/qrk/src/token"
)

func TestString(t *testing.T) {
	program := &Program{
		Statements: []Statement{
			&LetStatement{
				Token: token.Token{Type: token.LET, Literal: "let"},
				Identifier: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "myVar"},
					Value: "myVar",
				},
				Value: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "anotherVar"},
					Value: "anotherVar",
				},
			},
		},
	}
	if program.ToString() != "let myVar = anotherVar;" {
		t.Errorf("program.String() wrong. got=%q", program.ToString())
	}
}
