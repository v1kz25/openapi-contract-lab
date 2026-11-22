package oasdiff

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func buildTestEngine(t *testing.T) string {
	t.Helper()
	binary := filepath.Join(t.TempDir(), "synthetic-engine.exe")
	cmd := exec.Command("go", "build", "-o", binary, "./testdata/engine/main.go")
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("build synthetic engine: %v\n%s", err, output)
	}
	return binary
}

func TestProcessOutputFailures(t *testing.T) {
	binary := buildTestEngine(t)
	base, head := localInputs(t)
	for _, test := range []struct {
		mode string
		size int
	}{
		{"stdout-at-limit", 16 << 20},
		{"stderr-at-limit", 2},
	} {
		t.Run(test.mode, func(t *testing.T) {
			t.Setenv("LAB_ENGINE_MODE", test.mode)
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			output, err := runCommand(ctx, binary, "changelog")
			if err != nil || len(output) != test.size {
				t.Fatalf("exact capture boundary rejected: size=%d err=%v", len(output), err)
			}
		})
	}
	for _, test := range []struct{ mode, stage, diagnostic string }{
		{"stdout-limit", "comparison", "16777216 bytes"},
		{"stdout-limit-exit", "comparison", "16777216 bytes"},
		{"stderr-limit", "comparison", "65536 bytes"},
		{"stderr-limit-exit", "comparison", "65536 bytes"},
		{"bad-json", "decoding", "invalid oasdiff changelog JSON"},
		{"exit", "comparison", "synthetic engine failed"},
	} {
		t.Run(test.mode, func(t *testing.T) {
			t.Setenv("LAB_ENGINE_MODE", test.mode)
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			report, err := New(binary).Compare(ctx, base, head)
			var adapterErr *Error
			if !errors.As(err, &adapterErr) || adapterErr.Stage != test.stage || !strings.Contains(err.Error(), test.diagnostic) {
				// Avoid dumping the synthetic 64 KiB stderr into CI logs.
				t.Fatalf("mode %s: missing stage %s or diagnostic %q", test.mode, test.stage, test.diagnostic)
			}
			if report.Engine != "" || report.Findings != nil {
				t.Fatal("failed engine returned a partial report")
			}
			if strings.HasSuffix(test.mode, "-exit") || test.mode == "exit" {
				var processErr *exec.ExitError
				if !errors.As(err, &processErr) || processErr.ExitCode() != 7 {
					t.Fatal("original process failure lost")
				}
			}
			for _, path := range []string{base, head} {
				data, err := os.ReadFile(path)
				if err != nil || string(data) != "synthetic input" {
					t.Fatal("input changed during failed comparison")
				}
			}
		})
	}
}
