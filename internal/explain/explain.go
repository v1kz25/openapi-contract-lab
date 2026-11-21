// Package explain adds verified, deterministic explanations to oasdiff findings.
package explain

import "github.com/v1kz25/openapi-contract-lab/internal/oasdiff"

type Location struct {
	Side   string `json:"side"`
	File   string `json:"file"`
	Line   int    `json:"line,omitempty"`
	Column int    `json:"column,omitempty"`
	URL    string `json:"url,omitempty"`
}

type Explanation struct {
	Kind      string     `json:"kind"`
	Summary   string     `json:"summary"`
	Impact    string     `json:"impact"`
	Action    string     `json:"action"`
	Locations []Location `json:"locations"`
}

type Finding struct {
	oasdiff.Finding
	Explanation Explanation `json:"explanation"`
}

// Report retains the adapter's fields and adds an explanation to each finding.
type Report struct {
	Engine        string    `json:"engine"`
	EngineVersion string    `json:"engineVersion"`
	Base          string    `json:"base,omitempty"`
	Head          string    `json:"head,omitempty"`
	Findings      []Finding `json:"findings"`
}

// Enrich never changes upstream classification, messages or source locations.
// Templates apply only to the verified engine version and rule/level pair.
func Enrich(report oasdiff.Report) Report {
	result := Report{
		Engine: report.Engine, EngineVersion: report.EngineVersion, Base: report.Base, Head: report.Head,
		Findings: make([]Finding, 0, len(report.Findings)),
	}
	for _, finding := range report.Findings {
		detail := Explanation{
			Kind:      "fallback",
			Summary:   "No hay una explicación verificada para esta combinación de regla y nivel.",
			Impact:    "Consulte el mensaje y los datos originales de oasdiff; no se ha inferido una causa adicional.",
			Action:    "Revise manualmente el cambio y las ubicaciones disponibles.",
			Locations: make([]Location, 0, 2),
		}
		if template, ok := templates[finding.RuleID]; ok && report.Engine == "oasdiff" && report.EngineVersion == oasdiff.Version && finding.UpstreamLevel == template.level {
			detail.Kind, detail.Summary = "template", template.summary
			detail.Impact, detail.Action = template.impact, template.action
		}
		if finding.BaseSource != nil {
			detail.Locations = append(detail.Locations, sourceLocation("base", finding.BaseSource, report.Base))
		}
		if finding.HeadSource != nil {
			detail.Locations = append(detail.Locations, sourceLocation("head", finding.HeadSource, report.Head))
		}
		result.Findings = append(result.Findings, Finding{Finding: finding, Explanation: detail})
	}
	return result
}
