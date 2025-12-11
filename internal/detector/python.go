package detector

import (
	"fmt"
	"os/exec"
	"regexp"
)

type PythonDetector struct{}

func (d *PythonDetector) Name() string {
	return "python"
}

func (d *PythonDetector) Detect() (*Runtime, error) {
	runtime := &Runtime{
		Name:  "python",
		Found: false,
	}

	// Try python3 first, then python
	pythonCmd := "python3"
	pythonPath, err := exec.LookPath(pythonCmd)
	if err != nil {
		// python3 없으면 python 시도
		pythonCmd = "python"
		pythonPath, err = exec.LookPath(pythonCmd)
		if err != nil {
			return runtime, nil // 둘 다 없음
		}
	}

	runtime.Path = pythonPath
	runtime.Found = true

	cmd := exec.Command(pythonCmd, "--version")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return runtime, fmt.Errorf("failed to execute python --version: %w", err)
	}

	runtime.Version = parsePythonVersion(string(output))
	return runtime, nil
}

func parsePythonVersion(output string) string {
	// Python 3.9.16
	re := regexp.MustCompile(`Python (\d+\.\d+\.\d+)`)
	matches := re.FindStringSubmatch(output)

	if len(matches) > 1 {
		return matches[1]
	}

	return "unknown"
}
