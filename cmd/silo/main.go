package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/shravan20/silo/internal/backup"
	"github.com/shravan20/silo/internal/config"
)

func main() {
	configPath := flag.String("config", "backup.yaml", "path to config file")
	indexPath := flag.String("index", ".silo/index.db", "path to index database")
	restorePath := flag.String("path", "", "restore destination directory")
	flag.Parse()
	args := flag.Args()

	if len(args) < 1 {
		printUsage()
		os.Exit(1)
	}
	ctx := context.Background()

	switch args[0] {
	case "run":
		cfg, err := config.Load(*configPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "config: %v\n", err)
			os.Exit(1)
		}
		if err := backup.Run(ctx, cfg, *indexPath); err != nil {
			fmt.Fprintf(os.Stderr, "backup: %v\n", err)
			os.Exit(1)
		}
		os.Exit(0)
	case "restore":
		if len(args) < 2 {
			fmt.Fprintf(os.Stderr, "silo restore <backup-id> [flags]\n")
			os.Exit(1)
		}
		backupID := args[1]
		if *restorePath == "" {
			fmt.Fprintf(os.Stderr, "restore: -path is required\n")
			os.Exit(1)
		}
		cfg, err := config.Load(*configPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "config: %v\n", err)
			os.Exit(1)
		}
		if err := backup.Restore(ctx, cfg, *indexPath, backupID, *restorePath); err != nil {
			fmt.Fprintf(os.Stderr, "restore: %v\n", err)
			os.Exit(1)
		}
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
	fmt.Fprintf(os.Stderr, "Usage: silo <command> [args]\n")
	fmt.Fprintf(os.Stderr, "Commands: run, restore <backup-id>, status\n")
	fmt.Fprintf(os.Stderr, "  restore requires -path <dest>\n")
}
