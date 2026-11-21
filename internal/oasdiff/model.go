// Package oasdiff runs the pinned local engine and normalizes its changelog.
package oasdiff

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

const Version = "1.39.0"

type Source struct {
	File   string `json:"file"`
	Line   int    `json:"line,omitempty"`
	Column int    `json:"column,omitempty"`
}

type Finding struct {
	RuleID        string  `json:"ruleId"`
	Severity      string  `json:"severity"`
	Compatibility string  `json:"compatibility"`
	UpstreamLevel int     `json:"upstreamLevel"`
	Message       string  `json:"message,omitempty"`
	Operation     string  `json:"operation,omitempty"`
	OperationID   string  `json:"operationId,omitempty"`
	Path          string  `json:"path,omitempty"`
	Section       string  `json:"section,omitempty"`
	BaseSource    *Source `json:"baseSource,omitempty"`
	HeadSource    *Source `json:"headSource,omitempty"`
	Fingerprint   string  `json:"fingerprint,omitempty"`
	// Upstream preserves the complete event, including fields we do not know yet.
	Upstream json.RawMessage `json:"upstream"`
}

type Report struct {
	Engine        string    `json:"engine"`
	EngineVersion string    `json:"engineVersion"`
	Base          string    `json:"base,omitempty"`
	Head          string    `json:"head,omitempty"`
	Findings      []Finding `json:"findings"`
}

// Decode accepts the bare changelog array emitted by the pinned engine.
// Unknown IDs are preserved; classification depends only on the upstream level.
func Decode(data []byte) (Report, error) {
	var events []json.RawMessage
	trimmed := bytes.TrimSpace(data)
	if len(trimmed) == 0 || trimmed[0] != '[' {
		return Report{}, fmt.Errorf("expected an oasdiff changelog JSON array")
	}
	if err := json.Unmarshal(trimmed, &events); err != nil {
		return Report{}, fmt.Errorf("invalid oasdiff changelog JSON: %w", err)
	}
	report := Report{Engine: "oasdiff", EngineVersion: Version, Findings: make([]Finding, 0, len(events))}
	for i, event := range events {
		var entry struct {
			ID             string  `json:"id"`
			Level          *int    `json:"level"`
			Text           string  `json:"text"`
			Operation      string  `json:"operation"`
			OperationID    string  `json:"operationId"`
			Path           string  `json:"path"`
			Section        string  `json:"section"`
			BaseSource     *Source `json:"baseSource"`
			RevisionSource *Source `json:"revisionSource"`
			Fingerprint    string  `json:"fingerprint"`
		}
		if err := json.Unmarshal(event, &entry); err != nil {
			return Report{}, fmt.Errorf("invalid oasdiff event %d: %w", i, err)
		}
		if strings.TrimSpace(entry.ID) == "" || entry.Level == nil {
			return Report{}, fmt.Errorf("oasdiff event %d requires id and integer level", i)
		}
		for _, source := range []*Source{entry.BaseSource, entry.RevisionSource} {
			if source != nil && (source.File == "" || source.Line < 0 || source.Column < 0) {
				return Report{}, fmt.Errorf("oasdiff event %d has an invalid source location", i)
			}
		}
		severity, compatibility := "unknown", "unknown"
		switch *entry.Level {
		case 3:
			severity, compatibility = "error", "breaking"
		case 2:
			severity, compatibility = "warning", "potentially-breaking"
		case 1:
			severity, compatibility = "info", "non-breaking"
		}
		report.Findings = append(report.Findings, Finding{
			RuleID: entry.ID, Severity: severity, Compatibility: compatibility,
			UpstreamLevel: *entry.Level, Message: entry.Text, Operation: entry.Operation,
			OperationID: entry.OperationID, Path: entry.Path, Section: entry.Section,
			BaseSource: entry.BaseSource, HeadSource: entry.RevisionSource,
			Fingerprint: entry.Fingerprint, Upstream: event,
		})
	}
	return report, nil
}
