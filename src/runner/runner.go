package runner

import (
	"fmt"
	"io"
	"os"

	"github.com/vdchnsk/qrk/src/compiler"
	"github.com/vdchnsk/qrk/src/evaluator"
	"github.com/vdchnsk/qrk/src/fs"
	"github.com/vdchnsk/qrk/src/lexer"
	"github.com/vdchnsk/qrk/src/object"
	"github.com/vdchnsk/qrk/src/parser"
	"github.com/vdchnsk/qrk/src/vm"
)

type RunMode int

const (
	Compile RunMode = iota
	Interpret
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
	Run(string(data), env, out, Interpret)
}

func Run(input string, env *object.Environment, out io.Writer, mode RunMode) object.Object {
	line := string(input)
	lexer := lexer.NewLexer(line)
	parser := parser.NewParser(lexer)

	program := parser.ParseProgram()

	if len(parser.Errors()) != 0 {
		parser.PrettyPrintErrors(out)
		return nil
	}

	switch mode {
	case Interpret:
		evalRes := evaluator.Eval(program, env)
		return evalRes

	case Compile:
		compiler := compiler.NewCompiler()
		err := compiler.Compile(program)
		if err != nil {
			fmt.Fprintf(out, "compilation failed: %s\n", err)
		}

		bytecode := compiler.Bytecode()
		vm := vm.NewVm(bytecode)
		err = vm.Run()
		if err != nil {
			fmt.Fprintf(out, "vm error: %s\n", err)
		}

		stackTopElem := vm.StackTop()
		return stackTopElem
	}

	return nil
}
