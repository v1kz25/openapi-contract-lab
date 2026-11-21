package oasdiff

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestGoldenReports(t *testing.T) {
	for _, test := range []struct{ name, input string }{
		{"required", "../../spikes/oasdiff/evidence/required.changelog.json"},
		{"description", "../../spikes/oasdiff/evidence/description.changelog.json"},
		{"unknown", "testdata/unknown.json"},
	} {
		t.Run(test.name, func(t *testing.T) {
			input, err := os.ReadFile(test.input)
			if err != nil {
				t.Fatal(err)
			}
			report, err := Decode(input)
			if err != nil {
				t.Fatal(err)
			}
			assertGolden(t, test.name, report)
		})
	}
}

func assertGolden(t *testing.T, name string, report Report) {
	t.Helper()
	actual, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	want, err := os.ReadFile(filepath.Join("testdata", name+".golden.json"))
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(bytes.TrimSpace(actual), bytes.TrimSpace(bytes.ReplaceAll(want, []byte("\r\n"), []byte("\n")))) {
		t.Fatalf("normalized report differs from %s golden:\n%s", name, actual)
	}
}

func TestMalformedEvents(t *testing.T) {
	for _, input := range []string{
		``, `null`, `{}`, `{"changes":[]}`, `[] []`, `[`, `[null]`, `[7]`,
		`[{"level":3}]`, `[{"id":" " ,"level":3}]`, `[{"id":"x"}]`,
		`[{"id":"x","level":null}]`, `[{"id":"x","level":"ERR"}]`,
		`[{"id":"x","level":1.5}]`, `[{"id":"x","level":3,"path":7}]`,
		`[{"id":"x","level":3,"baseSource":{"file":"base","line":-1}}]`,
		`[{"id":"x","level":3,"revisionSource":{"file":""}}]`,
	} {
		t.Run(input, func(t *testing.T) {
			if _, err := Decode([]byte(input)); err == nil {
				t.Fatalf("accepted malformed input: %s", input)
			}
		})
	}
}

func TestUnknownIDUsesOnlyUpstreamLevel(t *testing.T) {
	for _, test := range []struct {
		level                   int
		severity, compatibility string
	}{
		{3, "error", "breaking"}, {2, "warning", "potentially-breaking"},
		{1, "info", "non-breaking"}, {0, "unknown", "unknown"}, {99, "unknown", "unknown"},
	} {
		input, _ := json.Marshal([]map[string]any{{"id": "future-rule", "level": test.level, "extra": map[string]any{"newField": true}}})
		report, err := Decode(input)
		if err != nil {
			t.Fatal(err)
		}
		finding := report.Findings[0]
		if finding.Severity != test.severity || finding.Compatibility != test.compatibility || finding.UpstreamLevel != test.level || finding.RuleID != "future-rule" {
			t.Fatalf("unexpected mapping: %+v", finding)
		}
		var preserved map[string]any
		if err := json.Unmarshal(finding.Upstream, &preserved); err != nil {
			t.Fatal(err)
		}
		if preserved["extra"].(map[string]any)["newField"] != true {
			t.Fatal("unknown field lost")
		}
	}
}
