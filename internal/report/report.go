// Package report renders deterministic local reports without external assets.
package report

import (
	"bytes"
	_ "embed"
	"fmt"
	"html"
	htmltemplate "html/template"
	"net/url"
	"strings"
	texttemplate "text/template"

	"github.com/v1kz25/openapi-contract-lab/internal/explain"
)

//go:embed report.html.tmpl
var htmlSource string

//go:embed report.md.tmpl
var markdownSource string

var htmlReport = htmltemplate.Must(htmltemplate.New("report").Parse(htmlSource))
var markdownReport = texttemplate.Must(texttemplate.New("report").Funcs(texttemplate.FuncMap{"md": markdownText}).Parse(markdownSource))

type counts struct{ Total, Errors, Warnings, Info, Unknown int }

type location struct {
	explain.Location
	Link htmltemplate.URL
}

type finding struct {
	explain.Finding
	Number    int
	Tone      string
	Endpoint  string
	Locations []location
}

type view struct {
	Engine, EngineVersion, Base, Head string
	Counts                            counts
	Groups                            []operationGroup
}

// Render groups only the human view; the input report and its order are untouched.
func Render(format string, report explain.Report) ([]byte, error) {
	v := view{Engine: report.Engine, EngineVersion: report.EngineVersion, Base: report.Base, Head: report.Head}
	var items []finding
	v.Counts.Total = len(report.Findings)
	for i, entry := range report.Findings {
		tone := "unknown"
		switch entry.Severity {
		case "error":
			v.Counts.Errors++
			tone = "error"
		case "warning":
			v.Counts.Warnings++
			tone = "warning"
		case "info":
			v.Counts.Info++
			tone = "info"
		default:
			v.Counts.Unknown++
		}
		item := finding{Finding: entry, Number: i + 1, Tone: tone, Endpoint: strings.TrimSpace(entry.Operation + " " + entry.Path)}
		for _, source := range entry.Explanation.Locations {
			item.Locations = append(item.Locations, location{Location: source, Link: localLink(source.URL)})
		}
		items = append(items, item)
	}
	v.Groups = groupFindings(items)
	var output bytes.Buffer
	var err error
	switch format {
	case "html":
		err = htmlReport.Execute(&output, v)
	case "markdown":
		err = markdownReport.Execute(&output, v)
	default:
		return nil, fmt.Errorf("unsupported report format %q", format)
	}
	if err != nil {
		return nil, fmt.Errorf("render %s report: %w", format, err)
	}
	data := bytes.ReplaceAll(output.Bytes(), []byte("\r\n"), []byte("\n"))
	return append(bytes.TrimRight(data, "\n"), '\n'), nil
}

// A template.URL cast is allowed only after checking the complete destination.
// External URLs, executable schemes and UNC hosts never become active links.
func localLink(raw string) htmltemplate.URL {
	u, err := url.Parse(raw)
	if err != nil || u.Scheme != "file" || u.Host != "" || u.User != nil || u.Opaque != "" || u.RawQuery != "" || u.ForceQuery || !strings.HasPrefix(u.Path, "/") || strings.HasPrefix(u.Path, "//") || strings.ContainsAny(u.Path, "\\\x00\r\n") {
		return ""
	}
	if u.Fragment != "" {
		if len(u.Fragment) < 2 || u.Fragment[0] != 'L' {
			return ""
		}
		for _, r := range u.Fragment[1:] {
			if r < '0' || r > '9' {
				return ""
			}
		}
	}
	return htmltemplate.URL(u.String())
}

func markdownText(value string) string {
	var result strings.Builder
	for _, r := range value {
		if r == '\n' || r == '\r' {
			result.WriteByte(' ')
			continue
		}
		if strings.ContainsRune("\\`*_{}[]()!|~", r) {
			result.WriteByte('\\')
		}
		result.WriteRune(r)
	}
	return html.EscapeString(result.String())
}
