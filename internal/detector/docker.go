package detector

import (
	"fmt"
	"os/exec"
	"regexp"
)

type DockerDetector struct{}

func (d *DockerDetector) Name() string {
	return "docker"
}

func (d *DockerDetector) Detect() (*Runtime, error) {
	runtime := &Runtime{
		Name:  "docker",
		Found: false,
	}

	dockerPath, err := exec.LookPath("docker")
	if err != nil {
		return runtime, nil
	}

	runtime.Path = dockerPath
	runtime.Found = true

	cmd := exec.Command("docker", "--version")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return runtime, fmt.Errorf("failed to execute docker --version: %w", err)
	}

	runtime.Version = parseDockerVersion(string(output))
	return runtime, nil
}

func parseDockerVersion(output string) string {
	// Docker version 24.0.5, build ced0996
	re := regexp.MustCompile(`version (\d+\.\d+\.\d+)`)
	matches := re.FindStringSubmatch(output)

	if len(matches) > 1 {
		return matches[1]
	}

	return "unknown"
}
