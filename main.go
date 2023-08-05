package main

import (
	"os"

	"github.com/vdchnsk/i-go/repl"
)

func main() {
	repl.Start(os.Stdin, os.Stdout)
}
