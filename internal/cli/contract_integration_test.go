package cli_test

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/v1kz25/openapi-contract-lab/internal/cli"
)

func TestCLIContractDefaultsIntegration(t *testing.T) {
	binary := os.Getenv("OASDIFF_BINARY")
	if binary == "" {
		t.Skip("set OASDIFF_BINARY to local oasdiff 1.39.0")
	}
	binary, err := filepath.Abs(binary)
	if err != nil {
		t.Fatal(err)
	}
	// Test PATH with the selected engine even if OASDIFF_BINARY was renamed.
	engineDir := t.TempDir()
	engineName := "oasdiff"
	if runtime.GOOS == "windows" {
		engineName += ".exe"
	}
	engineData, err := os.ReadFile(binary)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(engineDir, engineName), engineData, 0700); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", engineDir+string(os.PathListSeparator)+os.Getenv("PATH"))
	dir := t.TempDir()
	for _, name := range []string{"base", "required"} {
		data, err := os.ReadFile(filepath.Join("../../spikes/oasdiff/fixtures", name+".yaml"))
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(dir, "-"+name+".yaml"), data, 0600); err != nil {
			t.Fatal(err)
		}
	}
	t.Chdir(dir)
	for _, test := range []struct {
		name, output string
		args         []string
		code         int
	}{
		{"defaults", "", []string{"compare", "--", "-base.yaml", "-required.yaml"}, 3},
		{"equals", "", []string{"compare", "--oasdiff=" + binary, "--format=json", "--fail-on=none", "--", "-base.yaml", "-required.yaml"}, 0},
		{"equals-output", "report.json", []string{"compare", "--output=report.json", "--", "-base.yaml", "-required.yaml"}, 3},
	} {
		t.Run(test.name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			if code := cli.Run(test.args, &stdout, &stderr); code != test.code || (code == 3 && !strings.HasPrefix(stderr.String(), "policy:")) || (code == 0 && stderr.Len() != 0) {
				t.Fatalf("exit=%d stderr=%s", code, &stderr)
			}
			data := stdout.Bytes()
			if test.output != "" {
				if stdout.Len() != 0 {
					t.Fatal("--output polluted stdout")
				}
				data, err = os.ReadFile(test.output)
				if err != nil {
					t.Fatal(err)
				}
			}
			var report struct {
				Engine, EngineVersion, Base, Head string
				Findings                          []struct{ RuleID, Severity, Compatibility string }
			}
			if err := json.Unmarshal(data, &report); err != nil {
				t.Fatal(err)
			}
			if report.Engine != "oasdiff" || report.EngineVersion != "1.39.0" || report.Base != filepath.Join(dir, "-base.yaml") || report.Head != filepath.Join(dir, "-required.yaml") || len(report.Findings) != 1 || report.Findings[0].RuleID != "request-parameter-became-required" || report.Findings[0].Severity != "error" || report.Findings[0].Compatibility != "breaking" {
				t.Fatalf("default JSON/gate/PATH contract changed: %+v", report)
			}
		})
	}
}
