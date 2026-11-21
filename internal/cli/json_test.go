package cli

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"reflect"
	"testing"
)

func TestJSONPreservesUnicodeAcrossLegacyPipelines(t *testing.T) {
	value := map[string]any{
		"summary":  "El parámetro pasa de opcional a obligatorio.",
		"path":     `C:\contrato español #1.yaml`,
		"upstream": json.RawMessage(`{"text":"acción 🧪 漢字"}`),
		"literal":  `\u00e1`,
		"controls": "quote\" backslash\\ newline\n tab\t nul\x00",
	}
	var output bytes.Buffer
	if err := writeJSON(&output, value); err != nil {
		t.Fatal(err)
	}
	for _, b := range output.Bytes() {
		if b > 0x7f {
			t.Fatalf("non-ASCII output byte: %x", b)
		}
	}
	var actual, want any
	if err := json.Unmarshal(output.Bytes(), &actual); err != nil {
		t.Fatal(err)
	}
	original, err := json.Marshal(value)
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(original, &want); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(actual, want) {
		t.Fatalf("JSON values changed: got=%#v want=%#v", actual, want)
	}
}

type jsonFailingWriter struct{ short bool }

func (w jsonFailingWriter) Write([]byte) (int, error) {
	if w.short {
		return 0, nil
	}
	return 0, io.ErrClosedPipe
}

func TestJSONOutputErrors(t *testing.T) {
	if err := writeJSON(jsonFailingWriter{}, "parámetro"); !errors.Is(err, io.ErrClosedPipe) {
		t.Fatalf("error=%v", err)
	}
	if err := writeJSON(jsonFailingWriter{short: true}, "parámetro"); !errors.Is(err, io.ErrShortWrite) {
		t.Fatalf("error=%v", err)
	}
	var output bytes.Buffer
	if err := writeJSON(&output, make(chan int)); err == nil || output.Len() != 0 {
		t.Fatal("unsupported value should fail without output")
	}
}
