package output

import (
	"encoding/json"
	"fmt"

	"github.com/binaryarc/watcher/internal/detector"
)

// PrintRuntimesJSON prints runtimes in JSON format
func PrintRuntimesJSON(runtimes []*detector.Runtime) error {
	data, err := json.MarshalIndent(runtimes, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}
	fmt.Println(string(data))
	return nil
}

// PrintRuntimeJSON prints a single runtime in JSON format
func PrintRuntimeJSON(runtime *detector.Runtime) error {
	data, err := json.MarshalIndent(runtime, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}
	fmt.Println(string(data))
	return nil
}
