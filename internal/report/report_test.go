package report

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/v1kz25/openapi-contract-lab/internal/explain"
	"github.com/v1kz25/openapi-contract-lab/internal/oasdiff"
)

func snapshotReport() explain.Report {
	result := explain.Enrich(oasdiff.Report{
		Engine: "oasdiff", EngineVersion: oasdiff.Version, Base: "base.yaml", Head: "head.yaml",
		Findings: []oasdiff.Finding{
			{RuleID: "request-parameter-became-required", Severity: "error", Compatibility: "breaking", UpstreamLevel: 3, Operation: "GET", Path: "/items", Section: "paths", Message: "the `query` request parameter `filter` became required"},
			{RuleID: "unverified-warning", Severity: "warning", Compatibility: "potentially-breaking", UpstreamLevel: 2, Operation: "GET", Path: "/items", Message: "synthetic warning"},
			{RuleID: "request-parameter-became-optional", Severity: "info", Compatibility: "non-breaking", UpstreamLevel: 1, Operation: "POST", Path: "/items", Message: "synthetic optional parameter"},
			{RuleID: "future-rule", Severity: "unknown", Compatibility: "unknown", UpstreamLevel: 99, Message: "synthetic unknown event"},
		},
	})
	result.Findings[0].Explanation.Locations = []explain.Location{
		{Side: "base", File: "base.yaml", Line: 12, Column: 11, URL: "file:///synthetic/base.yaml#L12"},
		{Side: "head", File: "head.yaml", Line: 12, Column: 11, URL: "file:///synthetic/head.yaml#L12"},
	}
	return result
}

func TestSnapshots(t *testing.T) {
	for _, test := range []struct {
		name   string
		report explain.Report
	}{
		{"mixed", snapshotReport()},
		{"global", globalReport()},
		{"empty", explain.Enrich(oasdiff.Report{Engine: "oasdiff", EngineVersion: oasdiff.Version, Base: "base.yaml", Head: "head.yaml"})},
	} {
		for _, format := range []string{"markdown", "html"} {
			t.Run(test.name+"/"+format, func(t *testing.T) {
				actual, err := Render(format, test.report)
				if err != nil {
					t.Fatal(err)
				}
				if !utf8.Valid(actual) {
					t.Fatal("report is not valid UTF-8")
				}
				repeated, err := Render(format, test.report)
				if err != nil || !bytes.Equal(actual, repeated) {
					t.Fatal("report is not deterministic")
				}
				want, err := os.ReadFile(filepath.Join("testdata", test.name+"."+format+".golden"))
				if err != nil {
					t.Fatal(err)
				}
				want = bytes.ReplaceAll(want, []byte("\r\n"), []byte("\n"))
				if !bytes.Equal(actual, want) {
					t.Fatalf("%s snapshot mismatch:\n%s", format, actual)
				}
			})
		}
	}
}

func globalReport() explain.Report {
	r := groupedReport()
	r.Findings = []explain.Finding{r.Findings[3], r.Findings[7], r.Findings[8]}
	return r
}

func TestEscapingAndOfflineLinks(t *testing.T) {
	payload := `<script>alert(1)</script><img src=x onerror="alert(2)">` + "\n# injected [link](javascript:alert(3))"
	for _, format := range []string{"html", "markdown"} {
		t.Run(format, func(t *testing.T) {
			report := snapshotReport()
			report.Engine, report.EngineVersion, report.Base, report.Head = payload, payload, payload, payload
			entry := &report.Findings[0]
			entry.RuleID, entry.Message, entry.Operation, entry.Path, entry.Section, entry.Severity, entry.Compatibility = payload, payload, payload, payload, payload, payload, payload
			entry.Explanation.Kind, entry.Explanation.Summary, entry.Explanation.Impact, entry.Explanation.Action = payload, payload, payload, payload
			entry.Explanation.Locations = []explain.Location{{Side: payload, File: payload, URL: "javascript:alert(4)", Line: 12}}
			data, err := Render(format, report)
			if err != nil {
				t.Fatal(err)
			}
			text := string(data)
			for _, unsafe := range []string{"<script>", "<img ", `href="javascript:`} {
				if strings.Contains(text, unsafe) {
					t.Fatalf("active markup leaked: %q", unsafe)
				}
			}
			if format == "markdown" && strings.Contains(text, "\n# injected") {
				t.Fatal("Markdown block injected")
			}
			if strings.Count(text, "&lt;script&gt;") < 10 {
				t.Fatal("untrusted fields were dropped instead of escaped")
			}
		})
	}
	for _, raw := range []string{"javascript:alert(1)", "data:text/html,test", "https://example.invalid/", "//example.invalid/x", "file://server/share/spec.yaml", "file:relative.yaml", "file:////server/share/spec.yaml", "file:///%5C%5Cserver/share/spec.yaml", "file:///local.yaml?query", "file:///local.yaml?", "file:///local.yaml#bad", "file:///local%0aevil"} {
		if localLink(raw) != "" {
			t.Fatalf("unsafe link accepted: %q", raw)
		}
		for _, format := range []string{"html", "markdown"} {
			report := snapshotReport()
			report.Findings[0].Explanation.Locations = []explain.Location{{Side: "head", File: "source label", URL: raw}}
			data, err := Render(format, report)
			if err != nil || strings.Contains(string(data), raw) {
				t.Fatalf("unsafe URL leaked into %s: %q (%v)", format, raw, err)
			}
		}
	}
	for _, raw := range []string{"file:///synthetic/base.yaml#L12", "file:///C:/contracts/spec%20one.yaml#L12"} {
		if string(localLink(raw)) != raw {
			t.Fatalf("safe local link rejected: %q", raw)
		}
	}
	if _, err := Render("xml", snapshotReport()); err == nil {
		t.Fatal("unsupported format accepted")
	}
}

func TestMarkdownInlineEscapes(t *testing.T) {
	got := markdownText("<em>acción</em> [x](javascript:go) `code` | row\n# heading & 'quote'")
	for _, want := range []string{"&lt;em&gt;acción&lt;/em&gt;", `\[x\]\(javascript:go\)`, `\` + "`code" + `\` + "`", `\|`, "&amp;", "&#39;quote&#39;"} {
		if !strings.Contains(got, want) {
			t.Fatalf("missing escape %q in %q", want, got)
		}
	}
	if strings.ContainsAny(got, "\r\n") {
		t.Fatal("newlines can inject blocks")
	}
}
