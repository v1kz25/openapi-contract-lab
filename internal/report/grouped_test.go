package report

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/v1kz25/openapi-contract-lab/internal/explain"
	"github.com/v1kz25/openapi-contract-lab/internal/oasdiff"
)

func groupedReport() explain.Report {
	r := snapshotReport()
	r.Findings = append(r.Findings,
		explain.Enrich(oasdiff.Report{Findings: []oasdiff.Finding{
			{RuleID: "orders-info", Severity: "info", UpstreamLevel: 1, Operation: "GET", Path: "/orders", Message: "Synthetic orders info"},
			{RuleID: "post-error", Severity: "error", UpstreamLevel: 3, Operation: "POST", Path: "/items", Message: "Synthetic POST error"},
			{RuleID: "orders-error", Severity: "error", UpstreamLevel: 3, Operation: "GET", Path: "/orders", Message: "Synthetic orders error"},
			{RuleID: "path-only", Severity: "warning", UpstreamLevel: 2, Path: "/items", Message: "Synthetic path event"},
			{RuleID: "method-only", Severity: "future", UpstreamLevel: 42, Operation: "PATCH", Message: "Synthetic method event"},
		}}).Findings...)
	return r
}

func TestGroupingOrderAndConservation(t *testing.T) {
	r := groupedReport()
	// Include identical events, incomplete endpoints, an empty severity and
	// multiple unknown severities. They must remain separate occurrences.
	r.Findings = append(r.Findings, r.Findings[0], r.Findings[8], r.Findings[8])
	r.Findings[10].Severity = ""
	before, err := json.Marshal(r)
	if err != nil {
		t.Fatal(err)
	}
	var items []finding
	for i, f := range r.Findings {
		items = append(items, finding{Finding: f, Number: i + 1, Tone: map[string]string{"error": "error", "warning": "warning", "info": "info"}[f.Severity]})
	}
	groups := groupFindings(items)
	var titles []string
	var numbers []int
	seen := make(map[int]int)
	for _, g := range groups {
		titles = append(titles, g.Title)
		var total, errors, warnings, info, unknown int
		previous := -1
		for _, s := range g.Severities {
			if rank := severityRank(s.Severity); rank < previous {
				t.Fatal("severity order changed")
			} else {
				previous = rank
			}
			last := 0
			for _, f := range s.Findings {
				if f.Number <= last {
					t.Fatal("upstream order within severity changed")
				}
				last = f.Number
				if !reflect.DeepEqual(f.Finding, r.Findings[f.Number-1]) {
					t.Fatal("finding fields changed")
				}
				seen[f.Number]++
				numbers = append(numbers, f.Number)
				total++
				switch f.Severity {
				case "error":
					errors++
				case "warning":
					warnings++
				case "info":
					info++
				default:
					unknown++
				}
			}
		}
		if g.Counts != (counts{total, errors, warnings, info, unknown}) {
			t.Fatalf("wrong summary: %+v", g.Counts)
		}
	}
	if !reflect.DeepEqual(titles, []string{"GET /items", "POST /items", "GET /orders", "Sin ruta u operación"}) {
		t.Fatalf("groups: %v", titles)
	}
	if !reflect.DeepEqual(numbers, []int{1, 10, 2, 6, 3, 7, 5, 8, 11, 9, 12, 4}) {
		t.Fatalf("finding order: %v", numbers)
	}
	if len(seen) != len(r.Findings) {
		t.Fatal("findings lost")
	}
	for _, n := range seen {
		if n != 1 {
			t.Fatal("finding duplicated")
		}
	}
	for _, format := range []string{"html", "markdown"} {
		first, err := Render(format, r)
		if err != nil {
			t.Fatal(err)
		}
		again, err := Render(format, r)
		if err != nil || !bytes.Equal(first, again) {
			t.Fatal("nondeterministic rendering")
		}
		if !strings.Contains(string(first), "Sin ruta u operación") {
			t.Fatal("missing incomplete endpoint section")
		}
		for i := range r.Findings {
			prefix := fmt.Sprintf("##### %d. ", i+1)
			if format == "html" {
				prefix = fmt.Sprintf("<h5>%d. ", i+1)
			}
			if strings.Count(string(first), prefix) != 1 {
				t.Fatalf("%s lost or duplicated finding %d", format, i+1)
			}
		}
	}
	after, err := json.Marshal(r)
	if err != nil || !bytes.Equal(before, after) {
		t.Fatal("human rendering changed JSON fields or order")
	}
}

// These reviewable synthetic examples are also the grouped snapshots. The before
// files were rendered with develop before changing the presentation, same data.
func TestGroupedExampleSnapshots(t *testing.T) {
	for _, format := range []string{"html", "markdown"} {
		actual, err := Render(format, groupedReport())
		if err != nil {
			t.Fatal(err)
		}
		extension := format
		if format == "markdown" {
			extension = "md"
		}
		want, err := os.ReadFile("../../docs/examples/issue-23-grouped." + extension)
		if err != nil {
			t.Fatal(err)
		}
		want = bytes.ReplaceAll(want, []byte("\r\n"), []byte("\n"))
		if !bytes.Equal(actual, want) {
			t.Fatalf("review grouped %s example snapshot", format)
		}
	}
}

func TestEndpointKeysStayDistinct(t *testing.T) {
	items := []finding{
		{Finding: explain.Finding{Finding: oasdiff.Finding{Path: "/z", Operation: "get"}}},
		{Finding: explain.Finding{Finding: oasdiff.Finding{Path: "/a", Operation: "CUSTOM"}}},
		{Finding: explain.Finding{Finding: oasdiff.Finding{Path: "/z", Operation: "GET"}}},
		{Finding: explain.Finding{Finding: oasdiff.Finding{Path: " ", Operation: "GET"}}},
	}
	var got []string
	for _, g := range groupFindings(items) {
		got = append(got, g.Title)
	}
	if !reflect.DeepEqual(got, []string{"CUSTOM /a", "GET /z", "get /z", "Sin ruta u operación"}) {
		t.Fatalf("keys conflated: %v", got)
	}
}

func TestGroupHeadersAndNavigationStayPassive(t *testing.T) {
	r := snapshotReport()
	r.Findings = r.Findings[:1]
	f := &r.Findings[0]
	f.Path = "/<img src=x onerror=alert(1)>\n# hostile"
	f.Operation = `GET"><script>alert(2)</script>`
	f.Severity = `future"><script>alert(3)</script>`
	for _, format := range []string{"html", "markdown"} {
		data, err := Render(format, r)
		if err != nil {
			t.Fatal(err)
		}
		text := string(data)
		if strings.Contains(text, "<script>") || strings.Contains(text, "<img ") || (format == "markdown" && strings.Contains(text, "\n# hostile")) {
			t.Fatal("active group heading")
		}
		if !strings.Contains(text, "&lt;script&gt;") || !strings.Contains(text, "&lt;img") {
			t.Fatal("group content lost")
		}
		if format == "html" {
			if !strings.Contains(text, `href="#group-1"`) || !strings.Contains(text, `id="group-1"`) {
				t.Fatal("unsafe or missing navigation ID")
			}
			if !strings.Contains(text, `content="default-src 'none'; style-src 'unsafe-inline'; base-uri 'none'; form-action 'none'"`) {
				t.Fatal("CSP changed")
			}
			if strings.Contains(text, `<script`) || strings.Contains(text, `<link`) || strings.Contains(text, `<iframe`) {
				t.Fatal("active resources")
			}
		}
	}
	// Empty severity stays visible, without affecting the original finding.
	f.Severity = ""
	for _, format := range []string{"html", "markdown"} {
		data, err := Render(format, r)
		if err != nil || !strings.Contains(string(data), "Severidad: No reportada") {
			t.Fatal("empty severity hidden")
		}
	}
}
