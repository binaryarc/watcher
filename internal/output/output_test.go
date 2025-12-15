package output

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/binaryarc/watcher/internal/detector"
)

func captureOutput(t *testing.T, fn func()) string {
	t.Helper()

	var buf bytes.Buffer
	oldStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create pipe: %v", err)
	}
	os.Stdout = w

	fn()

	_ = w.Close()
	os.Stdout = oldStdout
	if _, err := io.Copy(&buf, r); err != nil {
		t.Fatalf("failed to read output: %v", err)
	}
	_ = r.Close()
	return buf.String()
}

func sampleRuntimes() []*detector.Runtime {
	return []*detector.Runtime{
		{Name: "java", Version: "17.0.8", Path: "/opt/java/bin/java", Found: true},
		{Name: "python", Version: "3.10.12", Path: "/usr/bin/python3", Found: true},
	}
}

func TestPrintRuntimesJSON(t *testing.T) {
	out := captureOutput(t, func() {
		if err := PrintRuntimesJSON(sampleRuntimes()); err != nil {
			t.Fatalf("PrintRuntimesJSON error: %v", err)
		}
	})

	if !strings.Contains(out, `"java"`) || !strings.Contains(out, `"17.0.8"`) {
		t.Fatalf("unexpected JSON output: %s", out)
	}
}

func TestPrintRuntimeJSON(t *testing.T) {
	out := captureOutput(t, func() {
		if err := PrintRuntimeJSON(sampleRuntimes()[0]); err != nil {
			t.Fatalf("PrintRuntimeJSON error: %v", err)
		}
	})

	if !strings.Contains(out, `"java"`) {
		t.Fatalf("unexpected JSON output: %s", out)
	}
}

func TestPrintRuntimesYAML(t *testing.T) {
	out := captureOutput(t, func() {
		if err := PrintRuntimesYAML(sampleRuntimes()); err != nil {
			t.Fatalf("PrintRuntimesYAML error: %v", err)
		}
	})

	if !strings.Contains(out, "- name: java") {
		t.Fatalf("unexpected YAML output: %s", out)
	}
}

func TestPrintRuntimeYAML(t *testing.T) {
	out := captureOutput(t, func() {
		if err := PrintRuntimeYAML(sampleRuntimes()[0]); err != nil {
			t.Fatalf("PrintRuntimeYAML error: %v", err)
		}
	})

	if !strings.Contains(out, "name: java") {
		t.Fatalf("unexpected YAML output: %s", out)
	}
}

func TestPrintRuntimesTable(t *testing.T) {
	out := captureOutput(t, func() {
		PrintRuntimesTable(sampleRuntimes())
	})

	if !strings.Contains(out, "RUNTIME") || !strings.Contains(out, "java") {
		t.Fatalf("unexpected table output: %s", out)
	}
}

func TestPrintRuntimeTable(t *testing.T) {
	out := captureOutput(t, func() {
		PrintRuntimeTable(sampleRuntimes()[0])
	})

	if !strings.Contains(out, "Name") || !strings.Contains(out, "java") {
		t.Fatalf("unexpected runtime table output: %s", out)
	}
}

func TestPrintComparisonFormats(t *testing.T) {
	comparison := &ComparisonData{
		Hosts: []string{"server1", "server2"},
		Runtimes: []RuntimeComparison{
			{
				Name:     "python",
				Versions: []string{"3.10.12", "3.11.0"},
				Status:   "DIFF",
			},
		},
	}

	tableOut := captureOutput(t, func() {
		PrintComparisonTable(comparison)
	})
	if !strings.Contains(tableOut, "server1") || !strings.Contains(tableOut, "python") {
		t.Fatalf("unexpected comparison table output: %s", tableOut)
	}

	jsonOut := captureOutput(t, func() {
		if err := PrintComparisonJSON(comparison); err != nil {
			t.Fatalf("PrintComparisonJSON error: %v", err)
		}
	})
	if !strings.Contains(jsonOut, `"python"`) || !strings.Contains(jsonOut, `"server1"`) {
		t.Fatalf("unexpected comparison JSON output: %s", jsonOut)
	}

	yamlOut := captureOutput(t, func() {
		if err := PrintComparisonYAML(comparison); err != nil {
			t.Fatalf("PrintComparisonYAML error: %v", err)
		}
	})
	if !strings.Contains(yamlOut, "hosts:") || !strings.Contains(yamlOut, "python") {
		t.Fatalf("unexpected comparison YAML output: %s", yamlOut)
	}
}
