package cli

import (
	"bytes"
	"encoding/json"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/v1kz25/openapi-contract-lab/internal/explain"
	"github.com/v1kz25/openapi-contract-lab/internal/oasdiff"
)

// Read the wire representation as a consumer, independently of production structs.
func contractObject(t *testing.T, data json.RawMessage) map[string]json.RawMessage {
	t.Helper()
	var object map[string]json.RawMessage
	if err := json.Unmarshal(data, &object); err != nil || object == nil {
		t.Fatalf("expected object: %s (%v)", data, err)
	}
	return object
}

func contractValue[T any](t *testing.T, object map[string]json.RawMessage, key string) T {
	t.Helper()
	data, ok := object[key]
	if !ok || bytes.Equal(data, []byte("null")) {
		t.Fatalf("missing or null required field %s", key)
	}
	var value T
	if err := json.Unmarshal(data, &value); err != nil {
		t.Fatalf("invalid %s: %s (%v)", key, data, err)
	}
	return value
}

func contractAbsent(t *testing.T, object map[string]json.RawMessage, keys ...string) {
	t.Helper()
	for _, key := range keys {
		if value, ok := object[key]; ok {
			t.Errorf("optional empty field %s must be omitted, got %s", key, value)
		}
	}
}

func TestJSONContractConsumer(t *testing.T) {
	// Repeated events and nested arrays must survive without sorting/deduplication.
	full := `{"id":"request-parameter-became-required","level":3,"text":"acción","operation":"GET","operationId":"list","path":"/items","section":"paths","fingerprint":"opaque","baseSource":{"file":"base.yaml","line":12,"column":11},"revisionSource":{"file":"head.yaml","line":13,"column":10},"future":{"values":[2,1,2],"precise":9007199254740993,"note":null}}`
	unknown := `{"id":"future-rule","level":99,"baseSource":{"file":"unknown.yaml"}}`
	for _, test := range []struct{ name, events string }{
		{"empty", `[]`},
		{"mixed", `[` + full + `,` + unknown + `,` + full + `]`},
	} {
		t.Run(test.name, func(t *testing.T) {
			report, err := oasdiff.Decode([]byte(test.events))
			if err != nil {
				t.Fatal(err)
			}
			report.Base, err = filepath.Abs("base.yaml")
			if err != nil {
				t.Fatal(err)
			}
			report.Head, err = filepath.Abs("head.yaml")
			if err != nil {
				t.Fatal(err)
			}
			var output bytes.Buffer
			if err := writeJSON(&output, explain.Enrich(report)); err != nil {
				t.Fatal(err)
			}
			root := contractObject(t, output.Bytes())
			if contractValue[string](t, root, "engine") != "oasdiff" || contractValue[string](t, root, "engineVersion") != "1.39.0" || contractValue[string](t, root, "base") != report.Base || contractValue[string](t, root, "head") != report.Head {
				t.Fatal("report identity or input paths changed")
			}
			findings := contractValue[[]json.RawMessage](t, root, "findings")
			var events []json.RawMessage
			if err := json.Unmarshal([]byte(test.events), &events); err != nil {
				t.Fatal(err)
			}
			if len(findings) != len(events) {
				t.Fatal("events lost or deduplicated")
			}
			for i, raw := range findings {
				finding := contractObject(t, raw)
				original := contractObject(t, events[i])
				if contractValue[string](t, finding, "ruleId") != contractValue[string](t, original, "id") {
					t.Fatal("upstream order or rule ID changed")
				}
				// UseNumber prevents a consumer assertion from hiding precision loss.
				decode := func(data []byte) any {
					decoder := json.NewDecoder(bytes.NewReader(data))
					decoder.UseNumber()
					var value any
					if err := decoder.Decode(&value); err != nil {
						t.Fatal(err)
					}
					return value
				}
				if !reflect.DeepEqual(decode(finding["upstream"]), decode(events[i])) {
					t.Fatal("upstream values, precision or array order changed")
				}
				explanation := contractObject(t, finding["explanation"])
				for _, field := range []string{"summary", "impact", "action"} {
					if contractValue[string](t, explanation, field) == "" {
						t.Fatalf("empty explanation %s", field)
					}
				}
				locations := contractValue[[]json.RawMessage](t, explanation, "locations")
				if i == 1 {
					if contractValue[int](t, finding, "upstreamLevel") != 99 || contractValue[string](t, finding, "severity") != "unknown" || contractValue[string](t, finding, "compatibility") != "unknown" || contractValue[string](t, explanation, "kind") != "fallback" {
						t.Fatal("unknown level/fallback contract changed")
					}
					contractAbsent(t, finding, "message", "operation", "operationId", "path", "section", "headSource", "fingerprint")
					if len(locations) != 1 {
						t.Fatal("invented or lost source")
					}
					source := contractObject(t, finding["baseSource"])
					location := contractObject(t, locations[0])
					if contractValue[string](t, source, "file") != "unknown.yaml" || contractValue[string](t, location, "file") != "unknown.yaml" || contractValue[string](t, location, "side") != "base" {
						t.Fatal("source-only location changed")
					}
					contractAbsent(t, source, "line", "column")
					contractAbsent(t, location, "line", "column", "url")
					continue
				}
				if contractValue[int](t, finding, "upstreamLevel") != 3 || contractValue[string](t, finding, "severity") != "error" || contractValue[string](t, finding, "compatibility") != "breaking" || contractValue[string](t, explanation, "kind") != "template" {
					t.Fatal("template contract changed")
				}
				for field, want := range map[string]string{"message": "acción", "operation": "GET", "operationId": "list", "path": "/items", "section": "paths", "fingerprint": "opaque"} {
					if contractValue[string](t, finding, field) != want {
						t.Fatalf("field %s lost", field)
					}
				}
				if len(locations) != 2 {
					t.Fatal("expected base and head locations")
				}
				for j, field := range []string{"baseSource", "headSource"} {
					source := contractObject(t, finding[field])
					location := contractObject(t, locations[j])
					if contractValue[string](t, location, "side") != []string{"base", "head"}[j] || contractValue[string](t, location, "file") != contractValue[string](t, source, "file") || contractValue[int](t, source, "line") != 12+j || contractValue[int](t, source, "column") != 11-j || contractValue[int](t, location, "line") != 12+j || contractValue[int](t, location, "column") != 11-j || contractValue[string](t, location, "url") == "" {
						t.Fatal("source fields/order/link changed")
					}
				}
			}
		})
	}
}

func TestJSONContractWithoutSources(t *testing.T) {
	report, err := oasdiff.Decode([]byte(`[{"id":"request-parameter-became-required","level":3},{"id":"future-rule","level":1}]`))
	if err != nil {
		t.Fatal(err)
	}
	var output bytes.Buffer
	if err := writeJSON(&output, explain.Enrich(report)); err != nil {
		t.Fatal(err)
	}
	root := contractObject(t, output.Bytes())
	// Internal normalized reports may lack context; CLI comparisons never do.
	contractAbsent(t, root, "base", "head")
	findings := contractValue[[]json.RawMessage](t, root, "findings")
	if len(findings) != 2 {
		t.Fatal("source-less findings lost")
	}
	for i, raw := range findings {
		finding := contractObject(t, raw)
		contractAbsent(t, finding, "baseSource", "headSource", "message")
		explanation := contractObject(t, finding["explanation"])
		if len(contractValue[[]json.RawMessage](t, explanation, "locations")) != 0 || contractValue[string](t, explanation, "kind") != []string{"template", "fallback"}[i] {
			t.Fatal("missing sources must yield [] for template and fallback")
		}
		if i == 1 && (contractValue[string](t, finding, "severity") != "info" || contractValue[string](t, finding, "compatibility") != "non-breaking") {
			t.Fatal("unknown ID must retain known-level classification")
		}
	}
}
