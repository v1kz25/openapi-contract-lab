package cli_test

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/v1kz25/openapi-contract-lab/internal/cli"
)

func TestPolicyIntegration(t *testing.T) {
	binary := os.Getenv("OASDIFF_BINARY")
	if binary == "" {
		t.Skip("set OASDIFF_BINARY to local oasdiff 1.39.0")
	}
	fixtures, err := filepath.Abs("../../spikes/oasdiff/fixtures")
	if err != nil {
		t.Fatal(err)
	}
	for _, test := range []struct {
		name, base, head, policy, severity string
		code                               int
	}{
		{"breaking-default", "base", "required", "", "error", 3},
		{"breaking-error", "base", "required", "error", "error", 3},
		{"breaking-warning", "base", "required", "warning", "error", 3},
		{"breaking-info", "base", "required", "info", "error", 3},
		{"breaking-disabled", "base", "required", "none", "error", 0},
		{"compatible-default", "required", "base", "", "info", 0},
		{"compatible-warning", "required", "base", "warning", "info", 0},
		{"compatible-info", "required", "base", "info", "info", 3},
		{"description", "base", "description", "", "", 0},
		{"empty-info", "base", "base", "info", "", 0},
	} {
		t.Run(test.name, func(t *testing.T) {
			args := []string{"compare", "--oasdiff", binary}
			if test.policy != "" {
				args = append(args, "--fail-on", test.policy)
			}
			args = append(args, filepath.Join(fixtures, test.base+".yaml"), filepath.Join(fixtures, test.head+".yaml"))
			var stdout, stderr bytes.Buffer
			code := cli.Run(args, &stdout, &stderr)
			if code != test.code || (code == 3 && !strings.Contains(stderr.String(), "policy:")) || (code == 0 && stderr.Len() != 0) {
				t.Fatalf("exit=%d stderr=%q", code, stderr.String())
			}
			var report struct {
				Findings []struct {
					Severity string `json:"severity"`
				} `json:"findings"`
			}
			if err := json.Unmarshal(stdout.Bytes(), &report); err != nil {
				t.Fatal(err)
			}
			if test.severity == "" {
				if len(report.Findings) != 0 {
					t.Fatal("expected empty report")
				}
			} else if len(report.Findings) != 1 || report.Findings[0].Severity != test.severity {
				t.Fatalf("unexpected findings: %+v", report.Findings)
			}
		})
	}

	base, head := filepath.Join(fixtures, "base.yaml"), filepath.Join(fixtures, "required.yaml")
	for _, format := range []string{"json", "markdown", "html"} {
		t.Run("unchanged-report-"+format, func(t *testing.T) {
			var blocked, allowed, stderr bytes.Buffer
			args := []string{"compare", "--oasdiff", binary, "--format", format}
			if code := cli.Run(append(append([]string{}, args...), base, head), &blocked, &stderr); code != 3 {
				t.Fatalf("exit=%d: %s", code, &stderr)
			}
			stderr.Reset()
			if code := cli.Run(append(args, "--fail-on", "none", base, head), &allowed, &stderr); code != 0 {
				t.Fatalf("exit=%d: %s", code, &stderr)
			}
			if !bytes.Equal(blocked.Bytes(), allowed.Bytes()) {
				t.Fatal("policy changed report bytes")
			}
		})
	}
	t.Run("write-failure-precedes-policy", func(t *testing.T) {
		var stderr bytes.Buffer
		code := cli.Run([]string{"compare", "--oasdiff", binary, base, head}, failingWriter{}, &stderr)
		if code != 1 || !strings.Contains(stderr.String(), "could not write output") || strings.Contains(stderr.String(), "policy:") {
			t.Fatalf("exit=%d: %s", code, &stderr)
		}
	})
	t.Run("none-does-not-hide-engine-errors", func(t *testing.T) {
		invalid := filepath.Join(t.TempDir(), "invalid.yaml")
		if err := os.WriteFile(invalid, []byte("openapi: ["), 0600); err != nil {
			t.Fatal(err)
		}
		var stdout, stderr bytes.Buffer
		code := cli.Run([]string{"compare", "--oasdiff", binary, "--fail-on", "none", base, invalid}, &stdout, &stderr)
		if code != 1 || stdout.Len() != 0 || !strings.Contains(stderr.String(), "error:") {
			t.Fatalf("exit=%d: %s", code, &stderr)
		}
	})
}
