package detector

import (
	"fmt"
	"os/exec"
	"regexp"
)

type RedisDetector struct{}

func (d *RedisDetector) Name() string {
	return "redis"
}

func (d *RedisDetector) Detect() (*Runtime, error) {
	runtime := &Runtime{
		Name:  "redis",
		Found: false,
	}

	// Try redis-server first, then redis-cli
	redisPath, err := exec.LookPath("redis-server")
	if err != nil {
		redisPath, err = exec.LookPath("redis-cli")
		if err != nil {
			return runtime, nil
		}
	}

	runtime.Path = redisPath
	runtime.Found = true

	cmd := exec.Command(redisPath, "--version")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return runtime, fmt.Errorf("failed to execute redis --version: %w", err)
	}

	runtime.Version = parseRedisVersion(string(output))
	return runtime, nil
}

func parseRedisVersion(output string) string {
	// Redis server v=7.0.12 sha=00000000:0 malloc=jemalloc-5.2.1 bits=64 build=7206160ac829380f
	re := regexp.MustCompile(`v=(\d+\.\d+\.\d+)`)
	matches := re.FindStringSubmatch(output)

	if len(matches) > 1 {
		return matches[1]
	}

	return "unknown"
}
