package main

import (
	"os"

	"github.com/vdchnsk/quasark/cmd/repl"
)

func main() {
	repl.Start(os.Stdin, os.Stdout)
}
