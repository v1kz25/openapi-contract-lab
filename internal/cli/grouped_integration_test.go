package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/v1kz25/openapi-contract-lab/internal/explain"
	"github.com/v1kz25/openapi-contract-lab/internal/oasdiff"
)

func TestGroupedReportsIntegration(t *testing.T) {
	binary := os.Getenv("OASDIFF_BINARY")
	if binary == "" {
		t.Skip("set OASDIFF_BINARY to local oasdiff 1.39.0")
	}
	// Three events across two routes and three operations, both ERR and INFO.
	spec := func(revision bool) []byte {
		operation := func(required bool) any {
			return map[string]any{
				"parameters": []any{map[string]any{"name": "filter", "in": "query", "required": required, "schema": map[string]any{"type": "string"}}},
				"responses":  map[string]any{"200": map[string]any{"description": "OK"}},
			}
		}
		data, err := json.Marshal(map[string]any{"openapi": "3.0.3", "info": map[string]any{"title": "Synthetic grouped API", "version": "1"}, "paths": map[string]any{
			"/items":  map[string]any{"get": operation(revision), "post": operation(!revision)},
			"/orders": map[string]any{"get": operation(revision)},
		}})
		if err != nil {
			t.Fatal(err)
		}
		return data
	}
	dir := t.TempDir()
	base, head := filepath.Join(dir, "base.json"), filepath.Join(dir, "head.json")
	for path, data := range map[string][]byte{base: spec(false), head: spec(true)} {
		if err := os.WriteFile(path, data, 0600); err != nil {
			t.Fatal(err)
		}
	}
	raw, err := oasdiff.New(binary).Compare(context.Background(), base, head)
	if err != nil {
		t.Fatal(err)
	}
	want := explain.Enrich(raw)
	if len(want.Findings) != 3 {
		t.Fatalf("expected three engine events: %+v", want.Findings)
	}
	var expected bytes.Buffer
	if err := writeJSON(&expected, want); err != nil {
		t.Fatal(err)
	}
	for _, format := range []string{"json", "html", "markdown"} {
		var blocked, allowed, stderr bytes.Buffer
		args := []string{"compare", "--oasdiff", binary, "--format", format}
		if code := Run(append(append([]string{}, args...), base, head), &blocked, &stderr); code != 3 || !strings.Contains(stderr.String(), "policy:") {
			t.Fatalf("%s exit=%d %s", format, code, &stderr)
		}
		stderr.Reset()
		if code := Run(append(append([]string{}, args...), "--fail-on", "none", base, head), &allowed, &stderr); code != 0 || stderr.Len() != 0 {
			t.Fatalf("%s exit=%d %s", format, code, &stderr)
		}
		if !bytes.Equal(blocked.Bytes(), allowed.Bytes()) {
			t.Fatal("gate changed grouped report")
		}
		if format == "json" {
			if !bytes.Equal(expected.Bytes(), blocked.Bytes()) {
				t.Fatal("JSON fields or upstream order changed")
			}
			continue
		}
		for _, endpoint := range []string{"GET /items", "POST /items", "GET /orders"} {
			if !strings.Contains(blocked.String(), endpoint) {
				t.Fatalf("%s missing %s", format, endpoint)
			}
		}
		for id, n := range map[string]int{"request-parameter-became-required": 2, "request-parameter-became-optional": 1} {
			if strings.Count(blocked.String(), id) != n {
				t.Fatalf("%s lost or duplicated %s", format, id)
			}
		}
		for _, level := range []string{"Severidad: error", "Severidad: info"} {
			if !strings.Contains(blocked.String(), level) {
				t.Fatalf("missing %s", level)
			}
		}
	}
}
