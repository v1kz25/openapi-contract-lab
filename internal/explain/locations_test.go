package explain

import (
	"net/url"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/v1kz25/openapi-contract-lab/internal/oasdiff"
)

func TestSourceLinks(t *testing.T) {
	input := filepath.Join(t.TempDir(), "contrato español #1.yaml")
	source := &oasdiff.Source{File: input, Line: 12, Column: 11}
	link := sourceLocation("head", source, input)
	parsed, err := url.Parse(link.URL)
	if err != nil || parsed.Scheme != "file" || parsed.Host != "" || parsed.Fragment != "L12" || !strings.Contains(link.URL, "%20") || !strings.Contains(link.URL, "%23") {
		t.Fatalf("invalid escaped local link: %+v (%v)", link, err)
	}
	if link.File != source.File || link.Line != 12 || link.Column != 11 || link.Side != "head" {
		t.Fatalf("location changed: %+v", link)
	}
	escaped := (&url.URL{Path: filepath.ToSlash(input)}).EscapedPath()
	encodedSource := &oasdiff.Source{File: escaped, Line: 12, Column: 11}
	if got := sourceLocation("head", encodedSource, input); got.URL != link.URL || got.File != escaped {
		t.Fatalf("escaped engine source lost: %+v", got)
	}
	for _, test := range []struct {
		source *oasdiff.Source
		input  string
	}{
		{&oasdiff.Source{File: filepath.Join(t.TempDir(), filepath.Base(input)), Line: 12}, input},
		{&oasdiff.Source{File: "https://example.invalid/spec.yaml", Line: 12}, input},
		{source, "relative.yaml"}, {source, ""}, {&oasdiff.Source{File: input, Line: -1}, input},
	} {
		if got := sourceLocation("base", test.source, test.input); got.URL != "" {
			t.Fatalf("invented link: %+v", got)
		}
	}
	fileOnly := sourceLocation("base", &oasdiff.Source{File: input}, input)
	parsed, err = url.Parse(fileOnly.URL)
	if err != nil || fileOnly.URL == "" || parsed.Fragment != "" {
		t.Fatal("invented line number")
	}
}

func TestWindowsSourceWithoutDrive(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("Windows path semantics")
	}
	for _, drive := range []string{"C:", "D:"} {
		input := drive + `\synthetic\contracts\base.yaml`
		source := &oasdiff.Source{File: "/synthetic/contracts/base.yaml", Line: 12, Column: 11}
		got := sourceLocation("base", source, input)
		want := "file:///" + drive + "/synthetic/contracts/base.yaml#L12"
		if got.URL != want || got.File != source.File {
			t.Fatalf("got=%+v want=%s", got, want)
		}
	}
}

func TestNoInventedLocations(t *testing.T) {
	report := oasdiff.Report{Engine: "oasdiff", EngineVersion: oasdiff.Version, Base: filepath.Join(t.TempDir(), "base.yaml"), Findings: []oasdiff.Finding{{RuleID: "future-rule", UpstreamLevel: 99}}}
	got := Enrich(report)
	if got.Findings[0].Explanation.Locations == nil || len(got.Findings[0].Explanation.Locations) != 0 {
		t.Fatal("invented source location")
	}
}
