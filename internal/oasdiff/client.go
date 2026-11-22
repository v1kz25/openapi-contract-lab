package oasdiff

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// Error identifies the stage that failed and preserves its underlying error.
type Error struct {
	Stage string
	Err   error
}

func (e *Error) Error() string { return "oasdiff " + e.Stage + ": " + e.Err.Error() }
func (e *Error) Unwrap() error { return e.Err }

type Client struct {
	binary string
	run    func(context.Context, string, ...string) ([]byte, error)
}

func New(binary string) *Client { return &Client{binary: binary, run: runCommand} }

// Compare accepts readable local regular files only. It enforces a 30 second
// deadline, the engine version, isolated configuration and disabled external refs.
func (c *Client) Compare(ctx context.Context, base, head string) (Report, error) {
	basePath, err := localFile(base)
	if err != nil {
		return Report{}, &Error{"base input", err}
	}
	headPath, err := localFile(head)
	if err != nil {
		return Report{}, &Error{"head input", err}
	}
	binary, err := exec.LookPath(c.binary)
	if err != nil {
		return Report{}, &Error{"executable", fmt.Errorf("locate %q (install oasdiff %s or use --oasdiff): %w", c.binary, Version, err)}
	}
	binary, err = filepath.Abs(binary)
	if err != nil {
		return Report{}, &Error{"executable", err}
	}
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	config, err := os.CreateTemp("", "openapi-contract-lab-oasdiff-*.yaml")
	if err != nil {
		return Report{}, &Error{"configuration", err}
	}
	defer os.Remove(config.Name())
	_, writeErr := config.WriteString("{}\n")
	closeErr := config.Close()
	if writeErr != nil {
		return Report{}, &Error{"configuration", writeErr}
	}
	if closeErr != nil {
		return Report{}, &Error{"configuration", closeErr}
	}
	version, err := c.run(ctx, binary, "--config", config.Name(), "--version")
	if err != nil {
		return Report{}, &Error{"version", err}
	}
	if strings.TrimSpace(string(version)) != "oasdiff version "+Version {
		return Report{}, &Error{"version", fmt.Errorf("expected oasdiff %s, received %q", Version, strings.TrimSpace(string(version)))}
	}
	output, err := c.run(ctx, binary, "changelog", "--config", config.Name(),
		"--format", "json", "--lang", "en", "--level", "INFO",
		"--allow-external-refs=false", "--open=false", "--fetch=false", "--", basePath, headPath)
	if err != nil {
		return Report{}, &Error{"comparison", err}
	}
	report, err := Decode(output)
	if err != nil {
		return Report{}, &Error{"decoding", err}
	}
	report.Base, report.Head = basePath, headPath
	return report, nil
}

func localFile(path string) (string, error) {
	if path == "" || path == "-" || strings.Contains(path, "://") || strings.HasPrefix(path, `\\`) || strings.HasPrefix(path, "//") {
		return "", fmt.Errorf("expected a local file path, received %q", path)
	}
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	f, err := os.Open(abs)
	if err != nil {
		return "", fmt.Errorf("open %q: %w", path, err)
	}
	defer f.Close()
	info, err := f.Stat()
	if err != nil {
		return "", err
	}
	if !info.Mode().IsRegular() {
		return "", fmt.Errorf("%q is not a regular file", path)
	}
	return abs, nil
}

type limitedBuffer struct {
	buffer bytes.Buffer
	limit  int
	err    error
}

func (b *limitedBuffer) Write(p []byte) (int, error) {
	if len(p) > b.limit-b.buffer.Len() {
		b.err = fmt.Errorf("engine output exceeds %d bytes", b.limit)
		return 0, b.err
	}
	return b.buffer.Write(p)
}

func runCommand(ctx context.Context, binary string, args ...string) ([]byte, error) {
	cmd := exec.CommandContext(ctx, binary, args...)
	cmd.WaitDelay = 2 * time.Second
	// Do not inherit upstream configuration overrides from the caller's shell.
	for _, variable := range os.Environ() {
		if !strings.HasPrefix(strings.ToUpper(variable), "OASDIFF_") {
			cmd.Env = append(cmd.Env, variable)
		}
	}
	stdout := &limitedBuffer{limit: 16 << 20}
	stderr := &limitedBuffer{limit: 64 << 10}
	cmd.Stdout, cmd.Stderr = stdout, stderr
	if err := cmd.Run(); err != nil {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		// Cmd.Wait may prefer a nonzero exit status over pipe-copy errors.
		// Retain the limit diagnosis as well as the original process failure.
		if stdout.err != nil {
			return nil, errors.Join(fmt.Errorf("stdout: %w", stdout.err), err)
		}
		if stderr.err != nil {
			return nil, errors.Join(fmt.Errorf("stderr: %w", stderr.err), err)
		}
		if message := strings.TrimSpace(stderr.buffer.String()); message != "" {
			return nil, fmt.Errorf("%w: %s", err, message)
		}
		return nil, err
	}
	return stdout.buffer.Bytes(), nil
}
