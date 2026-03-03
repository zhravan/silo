package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/shravan20/silo/internal/config"
)

func main() {
	configPath := flag.String("config", "backup.yaml", "path to config file")
	flag.Parse()
	args := flag.Args()

	if len(args) < 1 {
		printUsage()
		os.Exit(1)
	}
	switch args[0] {
	case "run":
		if _, err := config.Load(*configPath); err != nil {
			fmt.Fprintf(os.Stderr, "config: %v\n", err)
			os.Exit(1)
		}
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
