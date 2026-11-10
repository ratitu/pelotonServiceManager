package main

import (
	"flag"
	"fmt"
	"os"

	"peloton/internal/cli"
)

func main() {
	// Parse global flags without consuming subcommand flags
	var dbPath string
	var showHelp bool
	var showVersion bool

	// Custom flag parsing for global options before subcommand
	args := os.Args[1:]
	var subArgs []string

	for i := 0; i < len(args); i++ {
		arg := args[i]
		if arg == "--db" && i+1 < len(args) {
			dbPath = args[i+1]
			i++
		} else if arg == "--help" || arg == "-h" {
			showHelp = true
		} else if arg == "--version" || arg == "-v" {
			showVersion = true
		} else {
			subArgs = args[i:]
			break
		}
	}

	app, err := cli.NewApp(dbPath, os.Stdout, os.Stderr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal initialization error: %v\n", err)
		os.Exit(1)
	}
	defer app.Close()

	if showVersion {
		app.Run([]string{"--version"})
		return
	}
	if showHelp && len(subArgs) == 0 {
		app.PrintHelp()
		return
	}

	exitCode := app.Run(subArgs)
	if exitCode != 0 {
		os.Exit(exitCode)
	}
}

// Ensure flag package is referenced if needed
var _ = flag.ErrHelp
