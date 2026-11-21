package cli_test

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/v1kz25/openapi-contract-lab/internal/cli"
)

func TestHelp(t *testing.T) {
	for _, args := range [][]string{nil, {"help"}, {"-h"}, {"--help"}} {
		t.Run(strings.Join(args, " "), func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			if code := cli.Run(args, &stdout, &stderr); code != 0 {
				t.Fatalf("exit = %d; stderr = %q", code, stderr.String())
			}
			for _, text := range []string{"OpenAPI Contract Lab", "Usage:", "help", "version", "compare", "BASE HEAD"} {
				if !strings.Contains(stdout.String(), text) {
					t.Errorf("help missing %q: %q", text, stdout.String())
				}
			}
			if stderr.Len() != 0 {
				t.Errorf("unexpected stderr: %q", stderr.String())
			}
		})
	}
}

func TestCompareArguments(t *testing.T) {
	for _, args := range [][]string{
		{"compare"}, {"compare", "base.yaml"}, {"compare", "base", "head", "extra"},
		{"compare", "--unknown"}, {"compare", "--oasdiff"}, {"compare", "--oasdiff=", "base", "head"},
		{"compare", "--format", "xml", "base", "head"},
		{"compare", "--fail-on"},
		{"compare", "--fail-on=", "base", "head"},
		{"compare", "--fail-on", "warn", "base", "head"},
		{"compare", "--fail-on", "ERROR", "base", "head"},
		{"compare", "--output=", "base", "head"},
		{"compare", "--output", "https://example.invalid/report.html", "base", "head"},
		{"compare", "--output", `\\server\share\report.html`, "base", "head"},
	} {
		t.Run(strings.Join(args, " "), func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			if code := cli.Run(args, &stdout, &stderr); code != 2 || stdout.Len() != 0 || !strings.Contains(stderr.String(), "compare --help") {
				t.Fatalf("exit = %d; stdout = %q; stderr = %q", code, stdout.String(), stderr.String())
			}
		})
	}
	var stdout, stderr bytes.Buffer
	if code := cli.Run([]string{"compare", "--help"}, &stdout, &stderr); code != 0 || !strings.Contains(stdout.String(), "oasdiff 1.39.0") || !strings.Contains(stdout.String(), "--fail-on error|warning|info|none") || stderr.Len() != 0 {
		t.Fatalf("exit = %d; stdout = %q; stderr = %q", code, stdout.String(), stderr.String())
	}
}

func TestVersion(t *testing.T) {
	for _, arg := range []string{"version", "--version"} {
		t.Run(arg, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			code := cli.Run([]string{arg}, &stdout, &stderr)
			if code != 0 || stdout.String() != "openapi-contract-lab dev\n" || stderr.Len() != 0 {
				t.Fatalf("exit = %d; stdout = %q; stderr = %q", code, stdout.String(), stderr.String())
			}
		})
	}
}

func TestInvalidArguments(t *testing.T) {
	for _, test := range []struct {
		name string
		args []string
		want string
	}{
		{"unknown command", []string{"diff"}, `unknown command "diff"`},
		{"unknown flag", []string{"--unknown"}, `unknown command "--unknown"`},
		{"empty argument", []string{""}, `unknown command ""`},
		{"extra argument", []string{"version", "extra"}, "expected a single command"},
		{"help with extra argument", []string{"--help", "extra"}, "expected a single command"},
	} {
		t.Run(test.name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			code := cli.Run(test.args, &stdout, &stderr)
			if code != 2 || stdout.Len() != 0 || !strings.Contains(stderr.String(), test.want) || !strings.Contains(stderr.String(), "--help") {
				t.Fatalf("exit = %d; stdout = %q; stderr = %q", code, stdout.String(), stderr.String())
			}
		})
	}
}

type failingWriter struct{}

func (failingWriter) Write([]byte) (int, error) { return 0, errors.New("synthetic write failure") }

func TestOutputFailure(t *testing.T) {
	var stderr bytes.Buffer
	code := cli.Run([]string{"version"}, failingWriter{}, &stderr)
	if code != 1 || !strings.Contains(stderr.String(), "could not write output") || !strings.Contains(stderr.String(), "synthetic write failure") {
		t.Fatalf("exit = %d; stderr = %q", code, stderr.String())
	}
}
