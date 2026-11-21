package cli

import (
	"testing"

	"github.com/v1kz25/openapi-contract-lab/internal/oasdiff"
)

func TestSeverityPolicy(t *testing.T) {
	cases := []struct {
		severity string
		blocked  [4]bool // error, warning, info, none
	}{
		{"error", [4]bool{true, true, true, false}},
		{"warning", [4]bool{false, true, true, false}},
		{"info", [4]bool{false, false, true, false}},
		{"unknown", [4]bool{true, true, true, false}},
		{"future", [4]bool{true, true, true, false}},
		{"", [4]bool{true, true, true, false}},
	}
	for i, name := range []string{"error", "warning", "info", "none"} {
		t.Run(name, func(t *testing.T) {
			policy, err := parsePolicy(name)
			if err != nil {
				t.Fatal(err)
			}
			if policy.blocks(nil) {
				t.Fatal("empty report blocked")
			}
			for _, c := range cases {
				findings := []oasdiff.Finding{{Severity: "info"}, {Severity: c.severity}}
				if got := policy.blocks(findings); got != c.blocked[i] {
					t.Errorf("severity %q: blocked=%v, want %v", c.severity, got, c.blocked[i])
				}
			}
		})
	}
	for _, value := range []string{"", "ERROR", "warn", " error", "3"} {
		if _, err := parsePolicy(value); err == nil {
			t.Errorf("accepted %q", value)
		}
	}
}
