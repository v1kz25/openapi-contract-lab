package explain

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

type ruleCase struct {
	id         string
	level      int
	base, head string
}

func parameterSpec(required bool, schema string) string {
	return fmt.Sprintf(`{"openapi":"3.0.3","info":{"title":"Synthetic API","version":"1.0.0"},"paths":{"/items":{"get":{"parameters":[{"name":"filter","in":"query","required":%t,"schema":%s}],"responses":{"200":{"description":"Synthetic response"}}}}}}`, required, schema)
}

func responseSpec(enum string) string {
	color := map[string]any{"type": "string", "enum": json.RawMessage(enum)}
	schema := map[string]any{"type": "object", "properties": map[string]any{"color": color}}
	response := map[string]any{"description": "Synthetic response", "content": map[string]any{"application/json": map[string]any{"schema": schema}}}
	operation := map[string]any{"responses": map[string]any{"200": response}}
	spec := map[string]any{"openapi": "3.0.3", "info": map[string]any{"title": "Synthetic API", "version": "1.0.0"}, "paths": map[string]any{"/items": map[string]any{"get": operation}}}
	data, err := json.Marshal(spec)
	if err != nil {
		panic(err)
	}
	return string(data)
}

func ruleCases(t *testing.T) []ruleCase {
	t.Helper()
	plain := parameterSpec(false, `{"type":"string"}`)
	required := parameterSpec(true, `{"type":"string"}`)
	enumOne := parameterSpec(false, `{"type":"string","enum":["red"]}`)
	enumTwo := parameterSpec(false, `{"type":"string","enum":["red","blue"]}`)
	nullable := parameterSpec(false, `{"type":"string","nullable":true}`)
	return append([]ruleCase{
		{"request-parameter-became-required", 3, plain, required},
		{"request-parameter-became-optional", 1, required, plain},
		{"request-parameter-became-enum", 3, plain, enumOne},
		{"request-parameter-became-not-enum", 1, enumOne, plain},
		{"request-parameter-enum-value-added", 1, enumOne, enumTwo},
		{"request-parameter-enum-value-removed", 3, enumTwo, enumOne},
		{"request-parameter-became-nullable", 1, plain, nullable},
		{"request-parameter-became-not-nullable", 3, nullable, plain},
		{"response-property-enum-value-added", 3, responseSpec(`["red"]`), responseSpec(`["red","blue"]`)},
		{"response-property-enum-value-removed", 1, responseSpec(`["red","blue"]`), responseSpec(`["red"]`)},
	}, breakingRuleCases(t)...)
}

func breakingRuleCases(t *testing.T) []ruleCase {
	t.Helper()
	var cases []ruleCase
	for _, test := range []struct{ name, id string }{
		{"parameter", "request-parameter-type-changed"},
		{"body", "request-body-became-required"},
		{"response", "response-required-property-removed"},
	} {
		read := func(side string) string {
			data, err := os.ReadFile(filepath.Join("testdata", "breaking", test.name+"-"+side+".json"))
			if err != nil {
				t.Fatal(err)
			}
			return string(data)
		}
		cases = append(cases, ruleCase{test.id, 3, read("base"), read("head")})
	}
	return cases
}
