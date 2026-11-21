package oasdiff

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Opt-in: no test downloads or installs an engine.
func TestIntegration(t *testing.T) {
	binary := os.Getenv("OASDIFF_BINARY")
	if binary == "" {
		t.Skip("set OASDIFF_BINARY to the local oasdiff 1.39.0 executable")
	}
	binary, err := filepath.Abs(binary)
	if err != nil {
		t.Fatal(err)
	}
	fixtures, err := filepath.Abs("../../spikes/oasdiff/fixtures")
	if err != nil {
		t.Fatal(err)
	}
	// Ambient config must not alter results or cause the engine to fail.
	dir := t.TempDir()
	poison := filepath.Join(dir, ".oasdiff.yaml")
	if err := os.WriteFile(poison, []byte("invalid: ["), 0600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("OASDIFF_CONFIG", poison)
	t.Chdir(dir)
	client := New(binary)
	for _, test := range []struct {
		name, base, head, id, severity, compatibility string
		count                                         int
	}{
		{"required", "base", "required", "request-parameter-became-required", "error", "breaking", 1},
		{"description", "base", "description", "", "", "", 0},
		{"optional", "required", "base", "request-parameter-became-optional", "info", "non-breaking", 1},
	} {
		t.Run(test.name, func(t *testing.T) {
			base := filepath.Join(fixtures, test.base+".yaml")
			head := filepath.Join(fixtures, test.head+".yaml")
			report, err := client.Compare(context.Background(), base, head)
			if err != nil {
				t.Fatal(err)
			}
			if len(report.Findings) != test.count || report.Base != base || report.Head != head {
				t.Fatalf("findings=%+v", report.Findings)
			}
			if test.count == 0 {
				return
			}
			finding := report.Findings[0]
			if finding.RuleID != test.id || finding.Severity != test.severity || finding.Compatibility != test.compatibility || finding.Operation != "GET" || finding.Path != "/items" {
				t.Fatalf("finding=%+v", finding)
			}
			if finding.BaseSource == nil || finding.HeadSource == nil || !strings.HasSuffix(filepath.ToSlash(finding.BaseSource.File), "/"+test.base+".yaml") || !strings.HasSuffix(filepath.ToSlash(finding.HeadSource.File), "/"+test.head+".yaml") || finding.BaseSource.Line != 12 || finding.HeadSource.Line != 12 || finding.BaseSource.Column != 11 || finding.HeadSource.Column != 11 {
				t.Fatalf("source locations lost: %+v / %+v", finding.BaseSource, finding.HeadSource)
			}
			// Compare to the recorded engine event after replacing only the local
			// absolute paths. All other upstream fields, including fingerprint, match.
			if test.name == "required" {
				for actual, relative := range map[string]string{finding.BaseSource.File: "spikes/oasdiff/fixtures/base.yaml", finding.HeadSource.File: "spikes/oasdiff/fixtures/required.yaml"} {
					quotedActual, _ := json.Marshal(actual)
					quotedRelative, _ := json.Marshal(relative)
					report.Findings[0].Upstream = bytes.ReplaceAll(report.Findings[0].Upstream, quotedActual, quotedRelative)
				}
				report.Findings[0].BaseSource.File = "spikes/oasdiff/fixtures/base.yaml"
				report.Findings[0].HeadSource.File = "spikes/oasdiff/fixtures/required.yaml"
				report.Base, report.Head = "", ""
				// assertGolden uses a path relative to the package; t.Chdir isolated the engine.
				actual, _ := json.MarshalIndent(report, "", "  ")
				golden, err := os.ReadFile(filepath.Join(fixtures, "../../../internal/oasdiff/testdata/required.golden.json"))
				if err != nil {
					t.Fatal(err)
				}
				if !bytes.Equal(bytes.TrimSpace(actual), bytes.TrimSpace(bytes.ReplaceAll(golden, []byte("\r\n"), []byte("\n")))) {
					t.Fatalf("live report differs from golden:\n%s", actual)
				}
			}
		})
	}
	base := filepath.Join(fixtures, "base.yaml")
	for _, test := range []struct{ name, content, want string }{
		{"malformed", "openapi: [", "failed"},
		{"external-ref", "openapi: 3.0.3\ninfo: {title: Synthetic, version: 1.0.0}\npaths:\n  /items:\n    get:\n      responses:\n        '200':\n          description: Synthetic\n          content:\n            application/json:\n              schema:\n                $ref: './external.yaml'\n", "external"},
	} {
		t.Run(test.name, func(t *testing.T) {
			head := filepath.Join(dir, test.name+".yaml")
			if err := os.WriteFile(head, []byte(test.content), 0600); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(filepath.Join(dir, "external.yaml"), []byte("type: string\n"), 0600); err != nil {
				t.Fatal(err)
			}
			_, err := client.Compare(context.Background(), base, head)
			var adapterErr *Error
			if !errors.As(err, &adapterErr) || adapterErr.Stage != "comparison" || !strings.Contains(strings.ToLower(err.Error()), test.want) {
				t.Fatalf("unexpected engine error: %v", err)
			}
		})
	}
}
