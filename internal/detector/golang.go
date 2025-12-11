package detector

import (
	"fmt"
	"os/exec"
	"strings"
)

type GoDetector struct{}

func (d *GoDetector) Name() string {
	return "go"
}

func (d *GoDetector) Detect() (*Runtime, error) {
	runtime := &Runtime{
		Name:  "go",
		Found: false,
	}

	goPath, err := exec.LookPath("go")
	if err != nil {
		return runtime, nil
	}

	runtime.Path = goPath
	runtime.Found = true

	cmd := exec.Command("go", "version")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return runtime, fmt.Errorf("failed to execute go version: %w", err)
	}

	runtime.Version = parseGoVersion(string(output))
	return runtime, nil
}

func parseGoVersion(output string) string {
	// go version go1.21.5 linux/amd64
	parts := strings.Fields(output)
	if len(parts) >= 3 {
		version := strings.TrimPrefix(parts[2], "go")
		return version
	}
	return "unknown"
}
