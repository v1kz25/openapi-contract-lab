// Package cli implements local contract comparison, help and version commands.
package cli

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"

	"github.com/v1kz25/openapi-contract-lab/internal/explain"
	"github.com/v1kz25/openapi-contract-lab/internal/oasdiff"
)

const help = `OpenAPI Contract Lab

Usage:
  openapi-contract-lab [help | version]
  openapi-contract-lab compare [--oasdiff PATH] BASE HEAD

Commands:
  help       Show this help (also -h, --help)
  version    Show the development version (also --version)
  compare    Compare local contracts and emit findings with verified explanations
`

const compareHelp = `Usage:
  openapi-contract-lab compare [--oasdiff PATH] BASE HEAD

BASE and HEAD must be local files. Requires oasdiff 1.39.0.
--oasdiff PATH selects the executable (default: oasdiff on PATH).
Documentation-only edits do not produce changelog findings.
Verified rules include deterministic Spanish explanations and local source links.
`

// Run writes command output and returns 0 on success, 1 on an output error,
// or comparison error, or 2 on invalid arguments. Findings do not change the exit code.
func Run(args []string, stdout, stderr io.Writer) int {
	if len(args) > 0 && args[0] == "compare" {
		return compare(args[1:], stdout, stderr)
	}
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

func compare(args []string, stdout, stderr io.Writer) int {
	flags := flag.NewFlagSet("compare", flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	binary := flags.String("oasdiff", "oasdiff", "local oasdiff executable")
	if err := flags.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			if _, err := io.WriteString(stdout, compareHelp); err != nil {
				fmt.Fprintf(stderr, "error: could not write output: %v\n", err)
				return 1
			}
			return 0
		}
		fmt.Fprintf(stderr, "error: %v; use compare --help for usage\n", err)
		return 2
	}
	if flags.NArg() != 2 || *binary == "" {
		fmt.Fprintln(stderr, "error: compare requires BASE and HEAD and a nonempty --oasdiff path; use compare --help for usage")
		return 2
	}
	report, err := oasdiff.New(*binary).Compare(context.Background(), flags.Arg(0), flags.Arg(1))
	if err != nil {
		fmt.Fprintf(stderr, "error: %v\n", err)
		return 1
	}
	if err := writeJSON(stdout, explain.Enrich(report)); err != nil {
		fmt.Fprintf(stderr, "error: could not write output: %v\n", err)
		return 1
	}
	return 0
}
