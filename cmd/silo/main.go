package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}
	switch os.Args[1] {
	case "run":
		// TODO: run backup
		os.Exit(0)
	case "restore":
		// TODO: restore from backup
		os.Exit(0)
	case "status":
		// TODO: show status
		os.Exit(0)
	default:
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Fprintf(os.Stderr, "Usage: silo <command>\n")
	fmt.Fprintf(os.Stderr, "Commands: run, restore, status\n")
}
