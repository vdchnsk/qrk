package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/vdchnsk/qrk/src/object"
	"github.com/vdchnsk/qrk/src/runner"
)

func Start(in io.Reader, out io.Writer) {
	fmt.Println(REPL_WELCOME_MESSAGE)
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()

	for {
		fmt.Print(REPL_PROMPT_MESSAGE)
		scanned := scanner.Scan()

		if !scanned {
			return
		}

		evaluated := runner.Run(scanner.Text(), env, out)

		if evaluated == nil {
			continue
		}

		io.WriteString(out, evaluated.Inspect())
		io.WriteString(out, "\n")
	}
}
