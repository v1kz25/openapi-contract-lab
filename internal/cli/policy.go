package cli

import (
	"fmt"

	"github.com/v1kz25/openapi-contract-lab/internal/oasdiff"
)

// Policy only decides the gate; it never reclassifies or filters findings.
type severityPolicy int

func parsePolicy(value string) (severityPolicy, error) {
	switch value {
	case "none":
		return 0, nil
	case "info":
		return 1, nil
	case "warning":
		return 2, nil
	case "error":
		return 3, nil
	default:
		return 0, fmt.Errorf("--fail-on must be error, warning, info or none")
	}
}

func (p severityPolicy) blocks(findings []oasdiff.Finding) bool {
	if p == 0 {
		return false
	}
	for _, finding := range findings {
		var level severityPolicy
		switch finding.Severity {
		case "info":
			level = 1
		case "warning":
			level = 2
		case "error":
			level = 3
		default:
			// An unrecognized severity is not evidence of compatibility.
			return true
		}
		if level >= p {
			return true
		}
	}
	return false
}
