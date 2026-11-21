package cli_test

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/v1kz25/openapi-contract-lab/internal/cli"
	"github.com/v1kz25/openapi-contract-lab/internal/explain"
)

// Exercise the new explanations through all public export paths and the gate.
func TestBreakingExplanationsIntegration(t *testing.T) {
	binary := os.Getenv("OASDIFF_BINARY")
	if binary == "" {
		t.Skip("set OASDIFF_BINARY to local oasdiff 1.39.0")
	}
	for _, test := range []struct{ name, id, summary string }{
		{"parameter", "request-parameter-type-changed", "Cambia de forma incompatible el tipo o formato del parámetro de petición."},
		{"body", "request-body-became-required", "El cuerpo de la petición pasa de opcional a obligatorio."},
		{"response", "response-required-property-removed", "Se retira una propiedad obligatoria del esquema de respuesta."},
	} {
		t.Run(test.name, func(t *testing.T) {
			fixtures := filepath.Join("..", "explain", "testdata", "breaking")
			base, head := filepath.Join(fixtures, test.name+"-base.json"), filepath.Join(fixtures, test.name+"-head.json")
			original := make(map[string][]byte)
			for _, path := range []string{base, head} {
				data, err := os.ReadFile(path)
				if err != nil {
					t.Fatal(err)
				}
				original[path] = data
			}
			for _, format := range []string{"json", "markdown", "html"} {
				t.Run(format, func(t *testing.T) {
					path := filepath.Join(t.TempDir(), "informe español."+format)
					args := []string{"compare", "--oasdiff", binary, "--format", format, "--output", path, base, head}
					var stdout, stderr bytes.Buffer
					if code := cli.Run(args, &stdout, &stderr); code != 3 || stdout.Len() != 0 || !strings.Contains(stderr.String(), "policy:") {
						t.Fatalf("exit=%d stdout=%q stderr=%q", code, &stdout, &stderr)
					}
					data, err := os.ReadFile(path)
					if err != nil || !utf8.Valid(data) {
						t.Fatalf("invalid UTF-8 export: %v", err)
					}
					if format == "json" {
						var report explain.Report
						if err := json.Unmarshal(data, &report); err != nil {
							t.Fatal(err)
						}
						if len(report.Findings) != 1 {
							t.Fatalf("unexpected findings: %+v", report.Findings)
						}
						f := report.Findings[0]
						if f.RuleID != test.id || f.UpstreamLevel != 3 || f.Severity != "error" || f.Compatibility != "breaking" || f.Explanation.Kind != "template" || f.Explanation.Summary != test.summary || f.Explanation.Impact == "" || f.Explanation.Action == "" || len(f.Upstream) == 0 {
							t.Fatalf("incomplete explanation or changed classification: %+v", f)
						}
					} else if !strings.Contains(string(data), test.id) || !strings.Contains(string(data), test.summary) {
						t.Fatal("new explanation missing from human report")
					}
					stdout.Reset()
					stderr.Reset()
					allowed := []string{"compare", "--oasdiff", binary, "--format", format, "--fail-on", "none", base, head}
					if code := cli.Run(allowed, &stdout, &stderr); code != 0 || stderr.Len() != 0 || !bytes.Equal(stdout.Bytes(), data) {
						t.Fatalf("disabled gate changed report or exit: %d %s", code, &stderr)
					}
					stdout.Reset()
					stderr.Reset()
					if code := cli.Run(args, &stdout, &stderr); code != 1 || stdout.Len() != 0 || strings.Contains(stderr.String(), "policy:") {
						t.Fatalf("overwrite not rejected before gate: %d %s", code, &stderr)
					}
					retained, err := os.ReadFile(path)
					if err != nil || !bytes.Equal(retained, data) {
						t.Fatal("existing export changed")
					}
				})
			}
			for path, before := range original {
				after, err := os.ReadFile(path)
				if err != nil || !bytes.Equal(before, after) {
					t.Fatal("input contract changed")
				}
			}
		})
	}
}
