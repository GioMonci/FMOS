// Command fmos starts FMOS (Fix My OS).
package main

import (
	"fmt"
	"os"

	"github.com/GioMonci/FMOS/internal/logging"
	"github.com/GioMonci/FMOS/internal/tui"
)

// version is the FMOS version. Release builds can set it with
// -ldflags "-X main.version=<version>".
var version = "0.1.0-dev"

func main() {
	os.Exit(run(os.Args[1:]))
}

// run handles the command-line arguments and returns the exit code.
func run(args []string) int {
	if len(args) == 0 {
		return openMenu()
	}

	switch args[0] {
	case "help", "--help", "-h":
		printUsage()
		return 0
	case "version", "--version":
		fmt.Printf("FMOS %s\n", version)
		return 0
	default:
		fmt.Fprintf(os.Stderr, "FMOS does not recognize %q.\n\n", args[0])
		printUsage()
		return 1
	}
}

// openMenu starts this session's log, opens the menu, and finishes the log
// when the user quits.
func openMenu() int {
	session, logErr := logging.Start(version)
	menuErr := tui.Run(version, session, logErr)

	if session != nil {
		if err := session.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "FMOS could not finish the log %s: %v\n", session.Path(), err)
		} else {
			fmt.Printf("Log saved to %s\n", session.Path())
		}
	}
	if menuErr != nil {
		fmt.Fprintf(os.Stderr, "FMOS could not open its menu: %v\n", menuErr)
		fmt.Fprintln(os.Stderr, "Try running FMOS in a terminal window, such as Windows Terminal.")
		return 1
	}
	return 0
}

func printUsage() {
	fmt.Print(`FMOS (Fix My OS)

Usage:
  fmos              Open the FMOS menu.
  fmos --help       Show this help message.
  fmos --version    Show the FMOS version.
`)
}
