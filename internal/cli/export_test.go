package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestWriteNewFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "report.html")
	want := []byte("<meta charset=utf-8>acción\n")
	if err := writeNewFile(path, want); err != nil {
		t.Fatal(err)
	}
	if err := writeNewFile(path, []byte("replacement")); err == nil {
		t.Fatal("existing file overwritten")
	}
	actual, err := os.ReadFile(path)
	if err != nil || !bytes.Equal(actual, want) {
		t.Fatalf("file changed: %q (%v)", actual, err)
	}
	if err := writeNewFile(filepath.Join(dir, "missing-parent", "report.html"), want); err == nil {
		t.Fatal("missing directory should fail")
	}
}
