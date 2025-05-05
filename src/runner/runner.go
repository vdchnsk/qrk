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
	CompileMode RunMode = iota
	InterpretMode
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
	Interpret(string(data), env, out)
}

func Interpret(input string, env *object.Environment, out io.Writer) object.Object {
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

func Compile(input string, out io.Writer, symbolTable *compiler.SymbolTable, constants []object.Object, globals []object.Object) object.Object {
	line := string(input)
	lexer := lexer.NewLexer(line)
	parser := parser.NewParser(lexer)

	program := parser.ParseProgram()

	if len(parser.Errors()) != 0 {
		parser.PrettyPrintErrors(out)
		return nil
	}

	compiler := compiler.NewWithState(symbolTable, constants)
	err := compiler.Compile(program)
	if err != nil {
		fmt.Fprintf(out, "compilation failed: %s\n", err)
	}

	bytecode := compiler.Bytecode()
	vm := vm.NewVmWithGlobalStore(bytecode, globals)
	err = vm.Run()
	if err != nil {
		fmt.Fprintf(out, "vm error: %s\n", err)
	}

	stackTopElem := vm.LastPoppedStackElem()
	return stackTopElem
}
