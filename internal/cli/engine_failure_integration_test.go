package cli

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestEngineFailuresPreserveFiles(t *testing.T) {
	binary := filepath.Join(t.TempDir(), "synthetic-engine.exe")
	cmd := exec.Command("go", "build", "-o", binary, "../oasdiff/testdata/engine/main.go")
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("build synthetic engine: %v\n%s", err, output)
	}
	dir := t.TempDir()
	base, head, existing := filepath.Join(dir, "base.yaml"), filepath.Join(dir, "head.yaml"), filepath.Join(dir, "existing.json")
	originals := map[string][32]byte{}
	for _, path := range []string{base, head, existing} {
		data := []byte("synthetic data to preserve\n")
		if err := os.WriteFile(path, data, 0600); err != nil {
			t.Fatal(err)
		}
		originals[path] = sha256.Sum256(data)
	}
	for _, mode := range []string{"stdout-limit", "stdout-limit-exit", "stderr-limit", "stderr-limit-exit", "bad-json", "exit"} {
		t.Run(mode, func(t *testing.T) {
			t.Setenv("LAB_ENGINE_MODE", mode)
			fresh := filepath.Join(t.TempDir(), "new-report")
			for _, test := range []struct{ policy, format, destination string }{
				{"error", "json", ""},
				{"none", "json", fresh},
				{"error", "html", existing},
				{"none", "markdown", base},
				{"error", "json", head},
			} {
				args := []string{"compare", "--oasdiff", binary, "--fail-on", test.policy, "--format", test.format}
				if test.destination != "" {
					args = append(args, "--output", test.destination)
				}
				args = append(args, base, head)
				var stdout, stderr bytes.Buffer
				if code := Run(args, &stdout, &stderr); code != 1 || stdout.Len() != 0 || !strings.HasPrefix(stderr.String(), "error: oasdiff ") || strings.Contains(stderr.String(), "policy:") {
					t.Fatalf("%s/%s: expected operational failure before output/gate", test.policy, test.format)
				}
				if _, err := os.Stat(fresh); !errors.Is(err, os.ErrNotExist) {
					t.Fatal("engine failure created a partial destination")
				}
				for path, hash := range originals {
					data, err := os.ReadFile(path)
					if err != nil || sha256.Sum256(data) != hash {
						t.Fatal("engine failure changed an original/existing file")
					}
				}
			}
		})
	}
}

type partialOutput struct{ bytes.Buffer }

func (w *partialOutput) Write(data []byte) (int, error) {
	n, _ := w.Buffer.Write(data[:min(7, len(data))])
	return n, errors.New("synthetic partial stdout failure")
}

func TestPartialStdoutFailurePrecedesGate(t *testing.T) {
	binary := os.Getenv("OASDIFF_BINARY")
	if binary == "" {
		t.Skip("set OASDIFF_BINARY to local oasdiff 1.39.0")
	}
	for _, policy := range []string{"error", "none"} {
		output := &partialOutput{}
		var stderr bytes.Buffer
		args := []string{"compare", "--oasdiff", binary, "--fail-on", policy, "../../spikes/oasdiff/fixtures/base.yaml", "../../spikes/oasdiff/fixtures/required.yaml"}
		if code := Run(args, output, &stderr); code != 1 || output.Len() != 7 || json.Valid(output.Bytes()) || !strings.Contains(stderr.String(), "could not write output") || strings.Contains(stderr.String(), "policy:") {
			t.Fatalf("partial stdout/gate regression, exit=%d", code)
		}
	}
}
