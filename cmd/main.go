package main

import (
	"os"

	"github.com/vdchnsk/i-go/cmd/repl"
)

func main() {
	repl.Start(os.Stdin, os.Stdout)
}
