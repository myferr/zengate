package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) == 1 {
		printHelp()
		return
	}

	err := RunCommand(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
