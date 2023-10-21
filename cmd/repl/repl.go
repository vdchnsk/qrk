package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/vdchnsk/i-go/src/evaluator"
	"github.com/vdchnsk/i-go/src/lexer"
	"github.com/vdchnsk/i-go/src/memory"
	"github.com/vdchnsk/i-go/src/parser"
)

const PRMOPT = "> "
const REPL_WELCOME_MESSAGE = `
   _             
  (_)______  ___ 
 / /___/ _ \/ _ \
/_/    \_  /\___/
      /___/      
`

func Start(in io.Reader, out io.Writer) {
	fmt.Print(REPL_WELCOME_MESSAGE)
	scanner := bufio.NewScanner(in)
	env := memory.NewEnvironment()

	for {
		fmt.Print(PRMOPT)
		scanned := scanner.Scan()

		if !scanned {
			return
		}

		line := scanner.Text()
		lexer := lexer.NewLexer(line)
		parser := parser.NewParser(lexer)

		program := parser.ParseProgram()
		if len(parser.Errors()) != 0 {
			printParserErrors(out, parser.Errors())
			continue
		}

		evaluated := evaluator.Eval(program, env)
		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}
	}
}

func printParserErrors(out io.Writer, errors []string) {
	io.WriteString(out, "Syntax error! \n")
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
