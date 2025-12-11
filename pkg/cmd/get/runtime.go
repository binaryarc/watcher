package get

import (
	"fmt"

	"github.com/binaryarc/watcher/internal/detector"
	"github.com/spf13/cobra"
)

var runtimeCmd = &cobra.Command{
	Use:   "runtime [name]",
	Short: "Get specific runtime version",
	Long:  `Get version information for a specific runtime (java, python, node, go, etc.)`,
	Args:  cobra.ExactArgs(1),
	Run:   runGetRuntime,
}

func init() {
	GetCmd.AddCommand(runtimeCmd)
}

func runGetRuntime(cmd *cobra.Command, args []string) {
	runtimeName := args[0]

	fmt.Printf("👁️  Observing %s runtime...\n\n", runtimeName)

	if runtimeName != "java" {
		fmt.Printf("❌ Runtime '%s' is not supported yet.\n", runtimeName)
		fmt.Println("Supported runtimes: java")
		return
	}

	javaDetector := &detector.JavaDetector{}
	runtime, err := javaDetector.Detect()
	if err != nil {
		fmt.Printf("❌ Error detecting %s: %v\n", runtimeName, err)
		return
	}

	if !runtime.Found {
		fmt.Printf("❌ %s is not installed on this system.\n", runtime.Name)
		return
	}

	fmt.Printf("✅ %s detected!\n\n", runtime.Name)
	fmt.Printf("Name:    %s\n", runtime.Name)
	fmt.Printf("Version: %s\n", runtime.Version)
	fmt.Printf("Path:    %s\n", runtime.Path)
}
