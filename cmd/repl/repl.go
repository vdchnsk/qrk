package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/vdchnsk/qrk/src/evaluator"
	"github.com/vdchnsk/qrk/src/lexer"
	"github.com/vdchnsk/qrk/src/object"
	"github.com/vdchnsk/qrk/src/parser"
)

const PRMOPT = "> "
const REPL_WELCOME_MESSAGE = `
                                        _    
                                       | |   
  ____  _   _   ____  ___   ____  _ __ | | __
 / _  || | | | / _  |/ __| / _  || '__|| |/ /
| (_| || |_| || (_| |\__ \| (_| || |   |   < 
 \__  | \____| \____||___/ \__,_||_|   |_|\_\
    | |                                      
    |_|                                      
`

func Start(in io.Reader, out io.Writer) {
	fmt.Print(REPL_WELCOME_MESSAGE)
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()

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
