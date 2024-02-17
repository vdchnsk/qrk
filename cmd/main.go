package main

import (
	"os"

	"github.com/vdchnsk/qrk/cmd/repl"
)

func main() {
	repl.Start(os.Stdin, os.Stdout)
}
