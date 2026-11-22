package oasdiff

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// Hold a Windows process handle while the helper is alive; Wait then verifies
// termination without depending on PID reuse or shell-specific process listings.
func TestActiveProcessCancellation(t *testing.T) {
	binary := buildTestEngine(t)
	base, head := localInputs(t)
	for _, test := range []struct {
		name, mode, stage string
		deadline          bool
	}{
		{"cancel-comparison", "wait", "comparison", false},
		{"deadline-version", "version-wait", "version", true},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Setenv("LAB_ENGINE_MODE", test.mode)
			ready := filepath.Join(t.TempDir(), "ready.json")
			t.Setenv("LAB_ENGINE_READY", ready)
			duration := 10 * time.Second
			if test.deadline {
				duration = 3 * time.Second
			}
			ctx, cancel := context.WithTimeout(context.Background(), duration)
			done := make(chan error, 1)
			go func() {
				report, err := New(binary).Compare(ctx, base, head)
				if report.Findings != nil {
					err = errors.New("canceled comparison returned findings")
				}
				done <- err
			}()
			consumed := false
			defer func() {
				cancel()
				if !consumed {
					select {
					case <-done:
					case <-time.After(8 * time.Second):
						t.Error("process did not stop during cleanup")
					}
				}
			}()
			var state struct {
				PID    int
				Config string
			}
			for state.PID == 0 {
				data, _ := os.ReadFile(ready)
				if json.Unmarshal(data, &state) == nil && state.PID > 0 {
					break
				}
				select {
				case <-ctx.Done():
					t.Fatal("helper did not start before deadline")
				case <-time.After(10 * time.Millisecond):
				}
			}
			process, err := os.FindProcess(state.PID)
			if err != nil || ctx.Err() != nil {
				t.Fatalf("capture live helper handle: %v", err)
			}
			defer process.Release()
			stopped := make(chan error, 1)
			go func() {
				_, err := process.Wait()
				stopped <- err
			}()
			started := time.Now()
			want := error(context.DeadlineExceeded)
			if !test.deadline {
				want = context.Canceled
				cancel() // Only after the helper published readiness and a PID.
			}
			select {
			case err := <-done:
				consumed = true
				var adapterErr *Error
				if !errors.Is(err, want) || !errors.As(err, &adapterErr) || adapterErr.Stage != test.stage {
					t.Fatalf("unexpected cancellation error: %v", err)
				}
			case <-time.After(8 * time.Second):
				t.Fatal("cancellation exceeded 8 second tolerance")
			}
			select {
			case err := <-stopped:
				if err != nil {
					t.Fatalf("wait on helper handle: %v", err)
				}
			case <-time.After(time.Second):
				t.Fatal("helper remained alive after runCommand returned")
			}
			if time.Since(started) > 8*time.Second {
				t.Fatal("process shutdown exceeded tolerance")
			}
			if _, err := os.Stat(state.Config); !errors.Is(err, os.ErrNotExist) {
				t.Fatalf("temporary engine config retained: %v", err)
			}
			for _, path := range []string{base, head} {
				data, err := os.ReadFile(path)
				if err != nil || string(data) != "synthetic input" {
					t.Fatal("cancellation changed an input")
				}
			}
		})
	}
}
