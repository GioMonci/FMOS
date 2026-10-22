// Command fmos starts FMOS (Fix My OS).
package main

import (
	"fmt"
	"os"

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
		if err := tui.Run(version); err != nil {
			fmt.Fprintf(os.Stderr, "FMOS could not open its menu: %v\n", err)
			fmt.Fprintln(os.Stderr, "Try running FMOS in a terminal window, such as Windows Terminal.")
			return 1
		}
		return 0
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

func printUsage() {
	fmt.Print(`FMOS (Fix My OS)

Usage:
  fmos              Open the FMOS menu.
  fmos --help       Show this help message.
  fmos --version    Show the FMOS version.
`)
}
