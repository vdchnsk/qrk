package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/vdchnsk/i-go/lexer"
	"github.com/vdchnsk/i-go/token"
)

const PRMOPT = "> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Print(PRMOPT)
		scanned := scanner.Scan()

		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.NewLexer(line)

		for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
			fmt.Printf("%+v\n", tok)
		}
	}
}
