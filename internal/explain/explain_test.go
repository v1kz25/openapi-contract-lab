package explain

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/v1kz25/openapi-contract-lab/internal/oasdiff"
)

type messageSnapshot struct {
	RuleID      string      `json:"ruleId"`
	Explanation Explanation `json:"explanation"`
}

func TestVerifiedMessages(t *testing.T) {
	var snapshots []messageSnapshot
	if len(ruleCases()) != 10 || len(templates) != 10 {
		t.Fatal("expected ten documented rules and cases")
	}
	for _, test := range ruleCases() {
		t.Run(test.id, func(t *testing.T) {
			if !json.Valid([]byte(test.base)) || !json.Valid([]byte(test.head)) {
				t.Fatal("invalid synthetic JSON fixture")
			}
			report := oasdiff.Report{Engine: "oasdiff", EngineVersion: oasdiff.Version, Findings: []oasdiff.Finding{{RuleID: test.id, UpstreamLevel: test.level, Message: "original upstream message", Upstream: json.RawMessage(`{"untouched":true}`)}}}
			before, _ := json.Marshal(report)
			enriched := Enrich(report)
			after, _ := json.Marshal(report)
			if !bytes.Equal(before, after) || !reflect.DeepEqual(enriched.Findings[0].Finding, report.Findings[0]) {
				t.Fatal("upstream finding was changed")
			}
			detail := enriched.Findings[0].Explanation
			if detail.Kind != "template" {
				t.Fatalf("no verified explanation for %s", test.id)
			}
			snapshots = append(snapshots, messageSnapshot{test.id, detail})
		})
	}
	assertSnapshot(t, "messages", snapshots)
}

func TestNeutralFallback(t *testing.T) {
	for _, test := range []struct {
		id      string
		level   int
		version string
	}{
		{"future-rule", 3, oasdiff.Version},
		{"request-parameter-became-required", 99, oasdiff.Version},
		{"request-parameter-became-required", 1, oasdiff.Version},
		{"request-parameter-became-required", 2, oasdiff.Version},
		{"request-parameter-became-required", 3, "future-version"},
	} {
		report := oasdiff.Report{Engine: "oasdiff", EngineVersion: test.version, Findings: []oasdiff.Finding{{RuleID: test.id, UpstreamLevel: test.level}}}
		enriched := Enrich(report)
		if !reflect.DeepEqual(enriched.Findings[0].Finding, report.Findings[0]) {
			t.Fatal("fallback changed finding")
		}
		assertSnapshot(t, "fallback", enriched.Findings[0].Explanation)
	}
	empty := Enrich(oasdiff.Report{Engine: "oasdiff", EngineVersion: oasdiff.Version})
	if empty.Findings == nil || len(empty.Findings) != 0 {
		t.Fatal("empty changelog must remain an empty array")
	}
}

func assertSnapshot(t *testing.T, name string, value any) {
	t.Helper()
	actual, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	want, err := os.ReadFile(filepath.Join("testdata", name+".golden.json"))
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(bytes.TrimSpace(actual), bytes.TrimSpace(bytes.ReplaceAll(want, []byte("\r\n"), []byte("\n")))) {
		t.Fatalf("%s snapshot mismatch:\n%s", name, actual)
	}
}
