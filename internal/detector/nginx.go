package detector

import (
	"fmt"
	"os/exec"
	"regexp"
)

type NginxDetector struct{}

func (d *NginxDetector) Name() string {
	return "nginx"
}

func (d *NginxDetector) Detect() (*Runtime, error) {
	runtime := &Runtime{
		Name:  "nginx",
		Found: false,
	}

	nginxPath, err := exec.LookPath("nginx")
	if err != nil {
		return runtime, nil
	}

	runtime.Path = nginxPath
	runtime.Found = true

	cmd := exec.Command("nginx", "-v")
	output, err := cmd.CombinedOutput()
	if err != nil {
		// nginx -v outputs to stderr even on success
		if len(output) == 0 {
			return runtime, fmt.Errorf("failed to execute nginx -v: %w", err)
		}
	}

	runtime.Version = parseNginxVersion(string(output))
	return runtime, nil
}

func parseNginxVersion(output string) string {
	// nginx version: nginx/1.24.0
	re := regexp.MustCompile(`nginx/(\d+\.\d+\.\d+)`)
	matches := re.FindStringSubmatch(output)

	if len(matches) > 1 {
		return matches[1]
	}

	return "unknown"
}
