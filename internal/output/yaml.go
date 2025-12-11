package output

import (
	"fmt"

	"github.com/binaryarc/watcher/internal/detector"
	"gopkg.in/yaml.v3"
)

// PrintRuntimesYAML prints runtimes in YAML format
func PrintRuntimesYAML(runtimes []*detector.Runtime) error {
	data, err := yaml.Marshal(runtimes)
	if err != nil {
		return fmt.Errorf("failed to marshal YAML: %w", err)
	}
	fmt.Println(string(data))
	return nil
}

// PrintRuntimeYAML prints a single runtime in YAML format
func PrintRuntimeYAML(runtime *detector.Runtime) error {
	data, err := yaml.Marshal(runtime)
	if err != nil {
		return fmt.Errorf("failed to marshal YAML: %w", err)
	}
	fmt.Println(string(data))
	return nil
}
