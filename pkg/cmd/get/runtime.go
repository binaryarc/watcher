package get

import (
	"fmt"

	"github.com/binaryarc/watcher/internal/detector"
	"github.com/binaryarc/watcher/internal/output"
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
	outputFormat, _ := cmd.Flags().GetString("output")

	if outputFormat == "table" {
		fmt.Printf("👁️  Observing %s runtime...\n\n", runtimeName)
	}

	var det detector.Detector

	switch runtimeName {
	case "java":
		det = &detector.JavaDetector{}
	case "python":
		det = &detector.PythonDetector{}
	case "node", "nodejs":
		det = &detector.NodeDetector{}
	case "go", "golang":
		det = &detector.GoDetector{}
	case "docker":
		det = &detector.DockerDetector{}
	case "mysql", "mariadb":
		det = &detector.MySQLDetector{}
	case "redis":
		det = &detector.RedisDetector{}
	case "nginx":
		det = &detector.NginxDetector{}
	default:
		fmt.Printf("❌ Runtime '%s' is not supported yet.\n", runtimeName)
		fmt.Println("\nSupported runtimes:")
		fmt.Println("  - java")
		fmt.Println("  - python")
		fmt.Println("  - node")
		fmt.Println("  - go")
		fmt.Println("  - docker")
		fmt.Println("  - mysql")
		fmt.Println("  - redis")
		fmt.Println("  - nginx")
		return
	}

	runtime, err := det.Detect()
	if err != nil {
		fmt.Printf("❌ Error detecting %s: %v\n", runtimeName, err)
		return
	}

	if !runtime.Found {
		if outputFormat == "table" {
			fmt.Printf("❌ %s is not installed on this system.\n", runtime.Name)
		}
		return
	}

	// 출력 형식에 따라 분기
	switch outputFormat {
	case "json":
		if err := output.PrintRuntimeJSON(runtime); err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	case "yaml":
		if err := output.PrintRuntimeYAML(runtime); err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	case "table":
		fmt.Printf("✅ %s detected!\n\n", runtime.Name)
		output.PrintRuntimeTable(runtime)
	default:
		fmt.Printf("Unknown output format: %s\n", outputFormat)
		fmt.Println("Supported formats: table, json, yaml")
	}
}
