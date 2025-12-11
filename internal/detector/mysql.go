package detector

import (
	"fmt"
	"os/exec"
	"regexp"
)

type MySQLDetector struct{}

func (d *MySQLDetector) Name() string {
	return "mysql"
}

func (d *MySQLDetector) Detect() (*Runtime, error) {
	runtime := &Runtime{
		Name:  "mysql",
		Found: false,
	}

	// Try mysql first, then mariadb
	mysqlPath, err := exec.LookPath("mysql")
	if err != nil {
		mysqlPath, err = exec.LookPath("mariadb")
		if err != nil {
			return runtime, nil
		}
		runtime.Name = "mariadb"
	}

	runtime.Path = mysqlPath
	runtime.Found = true

	cmd := exec.Command(mysqlPath, "--version")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return runtime, fmt.Errorf("failed to execute mysql --version: %w", err)
	}

	runtime.Version = parseMySQLVersion(string(output))
	return runtime, nil
}

func parseMySQLVersion(output string) string {
	// mysql  Ver 8.0.34 for Linux on x86_64
	// mysql  Ver 15.1 Distrib 10.11.4-MariaDB
	re := regexp.MustCompile(`Ver (\d+\.\d+\.\d+)`)
	matches := re.FindStringSubmatch(output)

	if len(matches) > 1 {
		return matches[1]
	}

	return "unknown"
}
