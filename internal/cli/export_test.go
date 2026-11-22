package cli

import (
	"bytes"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
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

type failingReportFile struct {
	file               *os.File
	writeErr, closeErr error
	short, closed      bool
}

func (f *failingReportFile) Write(data []byte) (int, error) {
	if f.writeErr != nil || f.short {
		n, err := f.file.Write(data[:min(7, len(data))])
		if err != nil {
			return n, err
		}
		return n, f.writeErr
	}
	return f.file.Write(data)
}

func (f *failingReportFile) Close() error {
	f.closed = true
	return errors.Join(f.file.Close(), f.closeErr)
}

func TestNewReportFileFailureRecovery(t *testing.T) {
	writeErr, closeErr := errors.New("synthetic disk write failure"), errors.New("synthetic close failure")
	data := []byte(`{"findings":[]}` + "\n")
	for _, test := range []struct {
		name, stage        string
		writeErr, closeErr error
		short              bool
		want               error
	}{
		{"partial-write", "write", writeErr, nil, false, writeErr},
		{"short-write", "write", nil, nil, true, io.ErrShortWrite},
		{"close", "close", nil, closeErr, false, closeErr},
		{"write-and-close", "write", writeErr, closeErr, false, writeErr},
	} {
		t.Run(test.name, func(t *testing.T) {
			dir := t.TempDir()
			path := filepath.Join(dir, "new.json")
			file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
			if err != nil {
				t.Fatal(err)
			}
			defer file.Close()
			failing := &failingReportFile{file: file, writeErr: test.writeErr, closeErr: test.closeErr, short: test.short}
			err = writeAndClose(path, failing, data)
			if !errors.Is(err, test.want) || !strings.HasPrefix(err.Error(), test.stage+" ") || !failing.closed {
				t.Fatalf("error=%v, closed=%v", err, failing.closed)
			}
			want := data
			if test.stage == "write" {
				want = data[:7]
			}
			actual, err := os.ReadFile(path)
			if err != nil || !bytes.Equal(actual, want) {
				t.Fatalf("failed new file contents=%q err=%v", actual, err)
			}
			// Retrying the same name must not overwrite even a partial report.
			if err := writeNewFile(path, data); err == nil {
				t.Fatal("partial/existing file overwritten on retry")
			}
			after, err := os.ReadFile(path)
			if err != nil || !bytes.Equal(after, actual) {
				t.Fatal("failed file changed during retry")
			}
			// Recovery uses a fresh destination; the failed file remains untouched.
			if err := writeNewFile(filepath.Join(dir, "recovered.json"), data); err != nil {
				t.Fatal(err)
			}
		})
	}
}
