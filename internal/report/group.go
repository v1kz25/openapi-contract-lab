package report

import (
	"sort"
	"strings"
)

type severityGroup struct {
	Severity string
	Findings []finding
}

type operationGroup struct {
	Number     int
	Title      string
	Counts     counts
	Severities []severityGroup
}

// Exact path/method pairs avoid conflating distinct upstream values. Incomplete
// pairs share a final group, but each finding still shows its original endpoint.
func groupFindings(items []finding) []operationGroup {
	type endpoint struct{ path, operation string }
	buckets := make(map[endpoint][]finding)
	var keys []endpoint
	for _, item := range items {
		key := endpoint{item.Path, item.Operation}
		if strings.TrimSpace(key.path) == "" || strings.TrimSpace(key.operation) == "" {
			key = endpoint{}
		}
		if _, exists := buckets[key]; !exists {
			keys = append(keys, key)
		}
		buckets[key] = append(buckets[key], item)
	}
	sort.Slice(keys, func(i, j int) bool {
		if keys[i] == (endpoint{}) {
			return false
		}
		if keys[j] == (endpoint{}) {
			return true
		}
		if keys[i].path != keys[j].path {
			return keys[i].path < keys[j].path
		}
		return keys[i].operation < keys[j].operation
	})
	groups := make([]operationGroup, 0, len(keys))
	for _, key := range keys {
		group := operationGroup{Number: len(groups) + 1, Title: key.operation + " " + key.path}
		if key == (endpoint{}) {
			group.Title = "Sin ruta u operación"
		}
		bySeverity := make(map[string][]finding)
		var severities []string
		for _, item := range buckets[key] {
			if _, exists := bySeverity[item.Severity]; !exists {
				severities = append(severities, item.Severity)
			}
			bySeverity[item.Severity] = append(bySeverity[item.Severity], item)
			group.Counts.Total++
			switch item.Tone {
			case "error":
				group.Counts.Errors++
			case "warning":
				group.Counts.Warnings++
			case "info":
				group.Counts.Info++
			default:
				group.Counts.Unknown++
			}
		}
		sort.Slice(severities, func(i, j int) bool {
			a, b := severityRank(severities[i]), severityRank(severities[j])
			if a != b {
				return a < b
			}
			return severities[i] < severities[j]
		})
		for _, severity := range severities {
			group.Severities = append(group.Severities, severityGroup{severity, bySeverity[severity]})
		}
		groups = append(groups, group)
	}
	return groups
}

func severityRank(severity string) int {
	switch severity {
	case "error":
		return 0
	case "warning":
		return 1
	case "info":
		return 2
	default:
		return 3
	}
}
