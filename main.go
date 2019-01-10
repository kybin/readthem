package main

import (
	"fmt"
	"os"
)

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "readthem needs a sub command")
		os.Exit(1)
	}
	subcmd := args[0]
	switch subcmd {
	case "server":
		serverMain()
	case "node":
		nodeMain()
	default:
		fmt.Fprintf(os.Stderr, "%s is not a valid sub command\n", subcmd)
	}
}
