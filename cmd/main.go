package main

import (
	"github.com/henriquemarlon/coprocessor-local-solver/cmd/root"
	"os"
)

func main() {
	err := root.Cmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
