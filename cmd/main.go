package main

import (
	"os"

	"github.com/vdchnsk/qrk/cmd/repl"
	"github.com/vdchnsk/qrk/src/runner"
)

func main() {
	args := os.Args[1:]
	out := os.Stdout

	shouldRunInRepl := len(args) == 0

	if shouldRunInRepl {
		repl.Start(os.Stdin, out)
		return
	}

	fileToRun := args[0]
	runner.RunFile(fileToRun, out)
}
