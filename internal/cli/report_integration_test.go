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
)

func TestReportIntegration(t *testing.T) {
	binary := os.Getenv("OASDIFF_BINARY")
	if binary == "" {
		t.Skip("set OASDIFF_BINARY to local oasdiff 1.39.0")
	}
	base, err := filepath.Abs("../../spikes/oasdiff/fixtures/base.yaml")
	if err != nil {
		t.Fatal(err)
	}
	head, err := filepath.Abs("../../spikes/oasdiff/fixtures/required.yaml")
	if err != nil {
		t.Fatal(err)
	}
	for _, format := range []string{"json", "markdown", "html"} {
		t.Run(format, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "informe español."+format)
			var stdout, stderr bytes.Buffer
			code := cli.Run([]string{"compare", "--oasdiff", binary, "--format", format, "--output", path, base, head}, &stdout, &stderr)
			if code != 0 || stdout.Len() != 0 || stderr.Len() != 0 {
				t.Fatalf("exit=%d stdout=%q stderr=%q", code, stdout.String(), stderr.String())
			}
			data, err := os.ReadFile(path)
			if err != nil || !utf8.Valid(data) {
				t.Fatalf("invalid file: %v", err)
			}
			if !strings.Contains(string(data), "request-parameter-became-required") {
				t.Fatal("rule ID missing from report")
			}
			if format == "json" {
				var report map[string]any
				if err := json.Unmarshal(data, &report); err != nil {
					t.Fatal(err)
				}
			} else if !strings.Contains(string(data), "El parámetro pasa de opcional a obligatorio.") {
				t.Fatal("Unicode explanation missing")
			}
		})
	}
}
