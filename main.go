package main

import (
	"os"

	"github.com/flug/persona/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
