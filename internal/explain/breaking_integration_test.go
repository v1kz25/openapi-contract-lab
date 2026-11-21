package explain

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/v1kz25/openapi-contract-lab/internal/oasdiff"
)

func TestBreakingControlsIntegration(t *testing.T) {
	binary := os.Getenv("OASDIFF_BINARY")
	if binary == "" {
		t.Skip("set OASDIFF_BINARY to local oasdiff 1.39.0")
	}
	// These reverse changes were observed with the pinned engine, not inferred
	// from the forward rule's severity. None has a verified local template.
	reverseIDs := []string{
		"request-parameter-type-generalized",
		"request-body-became-optional",
		"response-required-property-added",
	}
	for i, test := range breakingRuleCases(t) {
		t.Run(test.id, func(t *testing.T) {
			dir := t.TempDir()
			base, head := filepath.Join(dir, "base.json"), filepath.Join(dir, "head.json")
			for path, content := range map[string]string{base: test.base, head: test.head} {
				if err := os.WriteFile(path, []byte(content), 0600); err != nil {
					t.Fatal(err)
				}
			}
			client := oasdiff.New(binary)
			identical, err := client.Compare(context.Background(), base, base)
			if err != nil {
				t.Fatal(err)
			}
			if len(Enrich(identical).Findings) != 0 {
				t.Fatal("identical input produced a finding")
			}
			reverse, err := client.Compare(context.Background(), head, base)
			if err != nil {
				t.Fatal(err)
			}
			findings := Enrich(reverse).Findings
			if len(findings) != 1 {
				t.Fatalf("expected one reverse event, got %+v", findings)
			}
			f := findings[0]
			if f.RuleID != reverseIDs[i] || f.UpstreamLevel != 1 || f.Severity != "info" || f.Compatibility != "non-breaking" || f.Explanation.Kind != "fallback" {
				t.Fatalf("unexpected reverse classification or template: %+v", f)
			}
		})
	}
}
