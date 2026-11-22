package oasdiff

import (
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"
)

func localInputs(t *testing.T) (string, string) {
	t.Helper()
	dir := filepath.Join(t.TempDir(), "contracts with spaces")
	if err := os.Mkdir(dir, 0700); err != nil {
		t.Fatal(err)
	}
	base, head := filepath.Join(dir, "base.yaml"), filepath.Join(dir, "head.yaml")
	for _, name := range []string{base, head} {
		if err := os.WriteFile(name, []byte("synthetic input"), 0600); err != nil {
			t.Fatal(err)
		}
	}
	return base, head
}

func TestClientInvocation(t *testing.T) {
	base, head := localInputs(t)
	binary, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	client := New(binary)
	var config string
	var deadline time.Time
	calls := 0
	client.run = func(ctx context.Context, exe string, args ...string) ([]byte, error) {
		calls++
		if exe != binary {
			t.Fatalf("unexpected executable: %q", exe)
		}
		currentDeadline, ok := ctx.Deadline()
		if !ok || time.Until(currentDeadline) <= 0 || time.Until(currentDeadline) > 30*time.Second {
			t.Fatal("missing or incorrect 30 second deadline")
		}
		if calls == 1 {
			deadline = currentDeadline
		} else if !currentDeadline.Equal(deadline) {
			t.Fatal("version and comparison must share one deadline")
		}
		index := slices.Index(args, "--config")
		if index < 0 || index+1 >= len(args) {
			t.Fatal("missing explicit configuration")
		}
		config = args[index+1]
		contents, err := os.ReadFile(config)
		if err != nil || string(contents) != "{}\n" {
			t.Fatalf("config: %q, %v", contents, err)
		}
		if slices.Contains(args, "--version") {
			return []byte("oasdiff version 1.39.0\n"), nil
		}
		want := []string{"changelog", "--config", config, "--format", "json", "--lang", "en", "--level", "INFO", "--allow-external-refs=false", "--open=false", "--fetch=false", "--", base, head}
		if !slices.Equal(args, want) {
			t.Fatalf("args: %q", args)
		}
		return []byte("[]"), nil
	}
	report, err := client.Compare(context.Background(), base, head)
	if err != nil || report.Findings == nil || len(report.Findings) != 0 || calls != 2 {
		t.Fatalf("report=%+v calls=%d err=%v", report, calls, err)
	}
	if _, err := os.Stat(config); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("temporary configuration not removed: %v", err)
	}
}

func TestInvalidLocalInputs(t *testing.T) {
	base, head := localInputs(t)
	for _, path := range []string{"", "-", "https://example.invalid/spec.yaml", `\\server\share\spec.yaml`, "//server/share/spec.yaml", filepath.Join(t.TempDir(), "missing"), t.TempDir()} {
		client := New("unused")
		client.run = func(context.Context, string, ...string) ([]byte, error) {
			t.Fatal("engine invoked for invalid input")
			return nil, nil
		}
		_, err := client.Compare(context.Background(), path, head)
		var adapterErr *Error
		if !errors.As(err, &adapterErr) || adapterErr.Stage != "base input" {
			t.Fatalf("path=%q err=%v", path, err)
		}
	}
	_, err := New("unused").Compare(context.Background(), base, "missing-head.yaml")
	var adapterErr *Error
	if !errors.As(err, &adapterErr) || adapterErr.Stage != "head input" {
		t.Fatal(err)
	}
	_, err = New(filepath.Join(t.TempDir(), "missing-oasdiff")).Compare(context.Background(), base, head)
	if !errors.As(err, &adapterErr) || adapterErr.Stage != "executable" {
		t.Fatal(err)
	}
}

func TestClientFailures(t *testing.T) {
	base, head := localInputs(t)
	binary, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	for _, test := range []struct {
		name, version, output, stage string
		failure                      error
	}{
		{"wrong version", "oasdiff version 1.38.0", "[]", "version", nil},
		{"process failure", "oasdiff version 1.39.0", "", "comparison", errors.New("synthetic engine failure")},
		{"timeout", "oasdiff version 1.39.0", "", "comparison", context.DeadlineExceeded},
		{"bad json", "oasdiff version 1.39.0", "not JSON", "decoding", nil},
	} {
		t.Run(test.name, func(t *testing.T) {
			client := New(binary)
			client.run = func(_ context.Context, _ string, args ...string) ([]byte, error) {
				if slices.Contains(args, "--version") {
					return []byte(test.version), nil
				}
				return []byte(test.output), test.failure
			}
			_, err := client.Compare(context.Background(), base, head)
			var adapterErr *Error
			if !errors.As(err, &adapterErr) || adapterErr.Stage != test.stage {
				t.Fatalf("error=%v", err)
			}
			if test.failure != nil && !errors.Is(err, test.failure) {
				t.Fatal("underlying error lost")
			}
		})
	}
}

func TestCommandCancellationAndOutputLimit(t *testing.T) {
	binary, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := runCommand(ctx, binary, "--version"); !errors.Is(err, context.Canceled) {
		t.Fatalf("error=%v", err)
	}
	buffer := &limitedBuffer{limit: 3}
	if _, err := io.Copy(buffer, strings.NewReader("four")); err == nil || buffer.buffer.Len() > 3 {
		t.Fatal("output limit bypassed")
	}
}
