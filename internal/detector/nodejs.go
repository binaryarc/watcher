package detector

import (
	"fmt"
	"os/exec"
	"strings"
)

type NodeDetector struct{}

func (d *NodeDetector) Name() string {
	return "node"
}

func (d *NodeDetector) Detect() (*Runtime, error) {
	runtime := &Runtime{
		Name:  "node",
		Found: false,
	}

	nodePath, err := exec.LookPath("node")
	if err != nil {
		return runtime, nil
	}

	runtime.Path = nodePath
	runtime.Found = true

	cmd := exec.Command("node", "--version")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return runtime, fmt.Errorf("failed to execute node --version: %w", err)
	}

	runtime.Version = parseNodeVersion(string(output))
	return runtime, nil
}

func parseNodeVersion(output string) string {
	// v18.16.0
	version := strings.TrimSpace(output)
	version = strings.TrimPrefix(version, "v")
	return version
}
