package main

import (
	"os"

	"github.com/Mugen-Builders/cartesi-coprocessor-nonodox/cmd/root"
)

func main() {
	err := root.Cmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
