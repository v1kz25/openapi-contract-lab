package cli_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/v1kz25/openapi-contract-lab/internal/cli"
	"github.com/v1kz25/openapi-contract-lab/internal/explain"
	"github.com/v1kz25/openapi-contract-lab/internal/oasdiff"
)

// Expectations were observed with the pinned engine, not inferred from 3.0
// nullable semantics. Reuse the original 3.1 fixtures for the first three pairs.
func TestOpenAPI31CorpusIntegration(t *testing.T) {
	binary := os.Getenv("OASDIFF_BINARY")
	if binary == "" {
		t.Skip("set OASDIFF_BINARY to local oasdiff 1.39.0")
	}
	var cases []struct {
		Name, Base, Head, RuleID, Severity, Compatibility, Explanation, Message, Fingerprint string
		Exit, Level, BaseLine, BaseColumn, HeadLine, HeadColumn                              int
	}
	data, err := os.ReadFile("testdata/openapi31/cases.json")
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(data, &cases); err != nil {
		t.Fatal(err)
	}
	if len(cases) != 5 {
		t.Fatal("expected five reviewed pairs")
	}
	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			// Exercise drive handling, URI escaping, spaces and Unicode on Windows.
			dir := filepath.Join(t.TempDir(), "contratos español #31")
			if err := os.Mkdir(dir, 0700); err != nil {
				t.Fatal(err)
			}
			original := make(map[string][]byte)
			var paths []string
			for _, name := range []string{c.Base, c.Head} {
				content, err := os.ReadFile(filepath.Join("testdata/openapi31", name))
				if err != nil {
					t.Fatal(err)
				}
				path := filepath.Join(dir, filepath.Base(name))
				if err := os.WriteFile(path, content, 0600); err != nil {
					t.Fatal(err)
				}
				original[path] = content
				paths = append(paths, path)
			}
			base, head := paths[0], paths[1]
			raw, err := oasdiff.New(binary).Compare(context.Background(), base, head)
			if err != nil {
				t.Fatal(err)
			}
			want := explain.Enrich(raw)
			wantCount := 0
			if c.RuleID != "" {
				wantCount = 1
			}
			if len(want.Findings) != wantCount {
				t.Fatalf("observed event count changed: %d", len(want.Findings))
			}
			for _, format := range []string{"json", "html", "markdown"} {
				t.Run(format, func(t *testing.T) {
					output := filepath.Join(dir, "informe."+format)
					var stdout, stderr bytes.Buffer
					args := []string{"compare", "--oasdiff", binary, "--format", format}
					code := cli.Run(append(append([]string{}, args...), "--output", output, base, head), &stdout, &stderr)
					if code != c.Exit || stdout.Len() != 0 || (code == 3 && !strings.Contains(stderr.String(), "policy:")) || (code == 0 && stderr.Len() != 0) {
						t.Fatalf("exit=%d: %s", code, &stderr)
					}
					content, err := os.ReadFile(output)
					if err != nil || !utf8.Valid(content) || bytes.HasPrefix(content, []byte{0xef, 0xbb, 0xbf}) || bytes.Contains(content, []byte("\r\n")) {
						t.Fatalf("invalid UTF-8/LF export: %v", err)
					}
					stderr.Reset()
					if code := cli.Run(append(append([]string{}, args...), "--fail-on", "none", base, head), &stdout, &stderr); code != 0 || stderr.Len() != 0 || !bytes.Equal(content, stdout.Bytes()) {
						t.Fatalf("gate changed report: exit=%d %s", code, &stderr)
					}
					if format != "json" {
						if c.RuleID == "" {
							if !strings.Contains(string(content), "No hay hallazgos que mostrar.") {
								t.Fatal("missing empty state")
							}
							return
						}
						if strings.Count(string(content), c.RuleID) != 1 || !strings.Contains(string(content), "GET /items") || !strings.Contains(string(content), want.Findings[0].Explanation.Summary) || !strings.Contains(string(content), want.Findings[0].Explanation.Impact) || !strings.Contains(string(content), want.Findings[0].Explanation.Action) {
							t.Fatal("human export lost finding or explanation")
						}
						for _, location := range want.Findings[0].Explanation.Locations {
							if !strings.Contains(string(content), location.URL) {
								t.Fatal("source link lost")
							}
						}
						return
					}
					// Compare all JSON fields, including the full upstream event and
					// array order, ignoring only serialization whitespace/Unicode escapes.
					encoded, err := json.Marshal(want)
					if err != nil {
						t.Fatal(err)
					}
					var actualValue, wantValue any
					if err := json.Unmarshal(content, &actualValue); err != nil {
						t.Fatal(err)
					}
					if err := json.Unmarshal(encoded, &wantValue); err != nil {
						t.Fatal(err)
					}
					if !reflect.DeepEqual(actualValue, wantValue) {
						t.Fatal("JSON changed enriched engine fields or order")
					}
					var report explain.Report
					if err := json.Unmarshal(content, &report); err != nil {
						t.Fatal(err)
					}
					if report.Engine != "oasdiff" || report.EngineVersion != oasdiff.Version || report.Base != base || report.Head != head {
						t.Fatal("report metadata changed")
					}
					if c.RuleID == "" {
						if len(report.Findings) != 0 {
							t.Fatal("unexpected changelog event")
						}
						return
					}
					if len(report.Findings) != 1 {
						t.Fatalf("expected one event, got %d", len(report.Findings))
					}
					f := report.Findings[0]
					if f.RuleID != c.RuleID || f.UpstreamLevel != c.Level || f.Severity != c.Severity || f.Compatibility != c.Compatibility || f.Explanation.Kind != c.Explanation || f.Message != c.Message || f.Fingerprint != c.Fingerprint || f.Operation != "GET" || f.Path != "/items" || f.Section != "paths" {
						t.Fatalf("observed engine expectation changed: %+v", f)
					}
					if f.BaseSource == nil || f.HeadSource == nil || f.BaseSource.Line != c.BaseLine || f.BaseSource.Column != c.BaseColumn || f.HeadSource.Line != c.HeadLine || f.HeadSource.Column != c.HeadColumn {
						t.Fatal("engine source positions changed")
					}
					if len(f.Explanation.Locations) != 2 {
						t.Fatal("expected two source locations")
					}
					for i, location := range f.Explanation.Locations {
						source := []*oasdiff.Source{f.BaseSource, f.HeadSource}[i]
						if location.File != source.File || location.Line != source.Line || location.Column != source.Column || location.Side != []string{"base", "head"}[i] {
							t.Fatal("source fields changed")
						}
						u, err := url.Parse(location.URL)
						wantPath := filepath.ToSlash(paths[i])
						if !strings.HasPrefix(wantPath, "/") {
							wantPath = "/" + wantPath
						}
						if err != nil || u.Scheme != "file" || u.Host != "" || u.Path != wantPath || u.Fragment != fmt.Sprintf("L%d", source.Line) {
							t.Fatalf("unsafe or incorrect local link: %s", location.URL)
						}
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
