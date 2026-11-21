package explain

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/v1kz25/openapi-contract-lab/internal/oasdiff"
)

// Each documented rule is checked against the actual pinned engine, offline.
func TestIntegration(t *testing.T) {
	binary := os.Getenv("OASDIFF_BINARY")
	if binary == "" {
		t.Skip("set OASDIFF_BINARY to local oasdiff 1.39.0")
	}
	client := oasdiff.New(binary)
	for _, test := range ruleCases() {
		t.Run(test.id, func(t *testing.T) {
			dir := t.TempDir()
			base, head := filepath.Join(dir, "base.json"), filepath.Join(dir, "head.json")
			if err := os.WriteFile(base, []byte(test.base), 0600); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(head, []byte(test.head), 0600); err != nil {
				t.Fatal(err)
			}
			report, err := client.Compare(context.Background(), base, head)
			if err != nil {
				t.Fatal(err)
			}
			enriched := Enrich(report)
			matched := 0
			for _, finding := range enriched.Findings {
				if finding.RuleID != test.id {
					continue
				}
				matched++
				if finding.UpstreamLevel != test.level || finding.Explanation.Kind != "template" {
					t.Fatalf("unverified classification: %+v", finding)
				}
				if len(finding.Explanation.Locations) == 0 {
					t.Fatal("no reported locations")
				}
				for _, location := range finding.Explanation.Locations {
					if location.URL == "" || location.Line < 1 {
						t.Fatalf("unresolved reported source: %+v", location)
					}
				}
			}
			if matched != 1 {
				t.Fatalf("expected one %s event; got %d in %+v", test.id, matched, report.Findings)
			}
		})
	}
}
