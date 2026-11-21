// Package cli implements local contract comparison, help and version commands.
package cli

import (
	"bytes"
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/v1kz25/openapi-contract-lab/internal/explain"
	"github.com/v1kz25/openapi-contract-lab/internal/oasdiff"
	reports "github.com/v1kz25/openapi-contract-lab/internal/report"
)

const help = `OpenAPI Contract Lab

Usage:
  openapi-contract-lab [help | version]
  openapi-contract-lab compare [options] BASE HEAD

Commands:
  help       Show this help (also -h, --help)
  version    Show the development version (also --version)
  compare    Compare local contracts and emit findings with verified explanations

Use compare --help for export options and severity policy.
`

const compareHelp = `Usage:
  openapi-contract-lab compare [options] BASE HEAD

BASE and HEAD must be local files. Requires oasdiff 1.39.0.
--oasdiff PATH selects the executable (default: oasdiff on PATH).
--format selects json (default), markdown or html.
--output writes a new local UTF-8 file without overwriting an existing file.
--fail-on error|warning|info|none sets the minimum blocking severity (default: error).
Unknown severities block unless --fail-on none. Reports retain all findings.
Exit codes: 0 passed, 1 execution/output failure, 2 invalid arguments, 3 policy failed.
Without --output, the report is written to stdout.
Documentation-only edits do not produce changelog findings.
Verified rules include deterministic Spanish explanations and local source links.
`

// Run returns 0 on success, 1 on execution/output failure, 2 on invalid arguments,
// or 3 when a completed report fails the selected severity policy.
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
	format := flags.String("format", "json", "json, markdown or html")
	output := flags.String("output", "", "new local output file")
	failOn := flags.String("fail-on", "error", "minimum blocking severity or none")
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
	if *format != "json" && *format != "markdown" && *format != "html" {
		fmt.Fprintln(stderr, "error: --format must be json, markdown or html; use compare --help for usage")
		return 2
	}
	policy, err := parsePolicy(*failOn)
	if err != nil {
		fmt.Fprintf(stderr, "error: %v; use compare --help for usage\n", err)
		return 2
	}
	outputSet := false
	flags.Visit(func(f *flag.Flag) {
		if f.Name == "output" {
			outputSet = true
		}
	})
	if outputSet && (*output == "" || *output == "-" || strings.Contains(*output, "://") || strings.HasPrefix(*output, `\\`) || strings.HasPrefix(*output, "//")) {
		fmt.Fprintln(stderr, "error: --output requires a new local file path; use compare --help for usage")
		return 2
	}
	report, err := oasdiff.New(*binary).Compare(context.Background(), flags.Arg(0), flags.Arg(1))
	if err != nil {
		fmt.Fprintf(stderr, "error: %v\n", err)
		return 1
	}
	enriched := explain.Enrich(report)
	var data []byte
	if *format == "json" {
		var buffer bytes.Buffer
		err = writeJSON(&buffer, enriched)
		data = buffer.Bytes()
	} else {
		data, err = reports.Render(*format, enriched)
	}
	if err == nil {
		if outputSet {
			err = writeNewFile(*output, data)
		} else {
			_, err = io.Copy(stdout, bytes.NewReader(data))
		}
	}
	if err != nil {
		fmt.Fprintf(stderr, "error: could not write output: %v\n", err)
		return 1
	}
	if policy.blocks(report.Findings) {
		fmt.Fprintf(stderr, "policy: findings meet --fail-on %s or have unknown severity; report generated\n", *failOn)
		return 3
	}
	return 0
}

func writeNewFile(path string, data []byte) error {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
	if err != nil {
		return fmt.Errorf("create %q (existing files are not overwritten): %w", path, err)
	}
	_, writeErr := io.Copy(file, bytes.NewReader(data))
	closeErr := file.Close()
	if writeErr != nil {
		return fmt.Errorf("write %q: %w", path, writeErr)
	}
	if closeErr != nil {
		return fmt.Errorf("close %q: %w", path, closeErr)
	}
	return nil
}
