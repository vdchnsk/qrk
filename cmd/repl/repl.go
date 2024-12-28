package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/vdchnsk/qrk/src/compiler"
	"github.com/vdchnsk/qrk/src/object"
	"github.com/vdchnsk/qrk/src/runner"
	"github.com/vdchnsk/qrk/src/vm"
)

func Start(in io.Reader, out io.Writer) {
	fmt.Println(REPL_WELCOME_MESSAGE)
	scanner := bufio.NewScanner(in)

	symbolTable := compiler.NewSymbolTable()
	constants := []object.Object{}
	globals := make([]object.Object, vm.GlobalVarsSize)

	for {
		fmt.Print(REPL_PROMPT_MESSAGE)
		scanned := scanner.Scan()

		if !scanned {
			return
		}

		// TODO: add ability to specify run mode via CLI
		output := runner.Compile(scanner.Text(), out, symbolTable, constants, globals)
		if output == nil {
			continue
		}

		io.WriteString(out, output.Inspect())
		io.WriteString(out, "\n")
	}
}
