package cli_test

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/v1kz25/openapi-contract-lab/internal/cli"
	"github.com/v1kz25/openapi-contract-lab/internal/explain"
)

// Fixtures and expected classifications are reviewed together, never regenerated.
func TestRegressionIntegration(t *testing.T) {
	binary := os.Getenv("OASDIFF_BINARY")
	if binary == "" {
		t.Skip("set OASDIFF_BINARY to local oasdiff 1.39.0")
	}
	dir := filepath.Join("testdata", "regression")
	data, err := os.ReadFile(filepath.Join(dir, "cases.json"))
	if err != nil {
		t.Fatal(err)
	}
	var cases []struct {
		Name, Base, Head, RuleID, Severity, Compatibility string
		Exit, Level                                       int
	}
	if err := json.Unmarshal(data, &cases); err != nil {
		t.Fatal(err)
	}
	if len(cases) == 0 {
		t.Fatal("empty regression corpus")
	}
	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			code := cli.Run([]string{"compare", "--oasdiff", binary, filepath.Join(dir, c.Base), filepath.Join(dir, c.Head)}, &stdout, &stderr)
			if code != c.Exit {
				t.Fatalf("exit=%d, want %d: %s", code, c.Exit, &stderr)
			}
			if code == 1 {
				if stdout.Len() != 0 || !strings.Contains(stderr.String(), "oasdiff comparison") {
					t.Fatalf("unexpected failure: %s", &stderr)
				}
				return
			}
			if code == 0 && stderr.Len() != 0 {
				t.Fatalf("unexpected stderr: %s", &stderr)
			}
			if code == 3 && !strings.Contains(stderr.String(), "policy:") {
				t.Fatal("policy diagnostic missing")
			}
			var report explain.Report
			if err := json.Unmarshal(stdout.Bytes(), &report); err != nil {
				t.Fatal(err)
			}
			if c.RuleID == "" {
				if len(report.Findings) != 0 {
					t.Fatalf("unexpected findings: %+v", report.Findings)
				}
				return
			}
			if len(report.Findings) != 1 {
				t.Fatalf("expected one finding, got %d", len(report.Findings))
			}
			f := report.Findings[0]
			if f.RuleID != c.RuleID || f.Severity != c.Severity || f.Compatibility != c.Compatibility || f.UpstreamLevel != c.Level || f.Operation != "GET" || f.Path != "/items" || f.Explanation.Kind != "template" {
				t.Fatalf("classification or explanation regression: %+v", f)
			}
			if len(f.Explanation.Locations) != 2 {
				t.Fatal("expected both source locations")
			}
			for _, location := range f.Explanation.Locations {
				if location.Line != 12 || location.Column != 11 || location.URL == "" {
					t.Fatalf("source location regression: %+v", location)
				}
			}
		})
	}
}
