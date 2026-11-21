// Package cli implements the initial help and version commands.
package cli

import (
	"fmt"
	"io"
)

const help = `OpenAPI Contract Lab

Usage:
  openapi-contract-lab [help | version]

Commands:
  help       Show this help (also -h, --help)
  version    Show the development version (also --version)

Contract comparison is not available yet.
`

// Run writes command output and returns 0 on success, 1 on an output error,
// or 2 on invalid arguments. It does not start services or invoke oasdiff.
func Run(args []string, stdout, stderr io.Writer) int {
	if len(args) > 1 {
		fmt.Fprintln(stderr, "error: expected a single command; use --help for usage")
		return 2
	}

	command := "help"
	if len(args) == 1 {
		command = args[0]
	}

	var output string
	switch command {
	case "help", "-h", "--help":
		output = help
	case "version", "--version":
		output = "openapi-contract-lab dev\n"
	default:
		fmt.Fprintf(stderr, "error: unknown command %q; use --help for usage\n", command)
		return 2
	}

	if _, err := io.WriteString(stdout, output); err != nil {
		fmt.Fprintf(stderr, "error: could not write output: %v\n", err)
		return 1
	}
	return 0
}
