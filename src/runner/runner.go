package runner

import (
	"fmt"
	"io"
	"os"

	"github.com/vdchnsk/qrk/src/evaluator"
	"github.com/vdchnsk/qrk/src/fs"
	"github.com/vdchnsk/qrk/src/lexer"
	"github.com/vdchnsk/qrk/src/object"
	"github.com/vdchnsk/qrk/src/parser"
)

func RunFile(path string, out io.Writer) {
	canRunFile := fs.CanRunFile(path)

	if !canRunFile {
		panic(fmt.Sprintf("Can't run file '%s'", path))
	}

	data, err := os.ReadFile(path)

	if err != nil {
		fmt.Println(err)
		return
	}

	env := object.NewEnvironment()
	Run(string(data), env, out)
}

func Run(input string, env *object.Environment, out io.Writer) object.Object {
	line := string(input)
	lexer := lexer.NewLexer(line)
	parser := parser.NewParser(lexer)

	program := parser.ParseProgram()

	if len(parser.Errors()) != 0 {
		parser.PrettyPrintErrors(out)
		return nil
	}

	evalRes := evaluator.Eval(program, env)

	return evalRes
}
