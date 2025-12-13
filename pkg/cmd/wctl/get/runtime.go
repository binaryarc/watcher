package get

import (
	"context"
	"fmt"

	"github.com/binaryarc/watcher/internal/detector"
	"github.com/binaryarc/watcher/internal/grpcclient"
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
	Cmd.AddCommand(runtimeCmd)
	runtimeCmd.Flags().String("host", "", "Remote server address (e.g., server:9090)")
}

func runGetRuntime(cmd *cobra.Command, args []string) {
	runtimeName := args[0]
	outputFormat, _ := cmd.Flags().GetString("output")
	host, _ := cmd.Flags().GetString("host")

	var runtime *detector.Runtime
	var err error

	if host != "" {
		apiKey, _ := cmd.Root().PersistentFlags().GetString("api-key")
		runtime, err = observeRemoteRuntime(host, apiKey, runtimeName, outputFormat)
		if err != nil {
			fmt.Printf("Failed to observe remote server: %v\n", err)
			return
		}
	} else {
		runtime, err = observeLocalRuntime(runtimeName, outputFormat)
		if err != nil {
			return
		}
	}

	if runtime == nil || !runtime.Found {
		if outputFormat == "table" {
			fmt.Printf("%s is not installed.\n", runtimeName)
		}
		return
	}

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
		fmt.Printf("%s detected!\n\n", runtime.Name)
		output.PrintRuntimeTable(runtime)
	default:
		fmt.Printf("Unknown output format: %s\n", outputFormat)
		fmt.Println("Supported formats: table, json, yaml")
	}
}

func observeLocalRuntime(runtimeName string, outputFormat string) (*detector.Runtime, error) {
	if outputFormat == "table" {
		fmt.Printf("Observing %s runtime...\n\n", runtimeName)
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
		fmt.Printf("Runtime '%s' is not supported yet.\n", runtimeName)
		fmt.Println()
		fmt.Println("Supported runtimes:")
		fmt.Println("  - java")
		fmt.Println("  - python")
		fmt.Println("  - node")
		fmt.Println("  - go")
		fmt.Println("  - docker")
		fmt.Println("  - mysql")
		fmt.Println("  - redis")
		fmt.Println("  - nginx")
		return nil, fmt.Errorf("unsupported runtime: %s", runtimeName)
	}

	runtime, err := det.Detect()
	if err != nil {
		fmt.Printf("Error detecting %s: %v\n", runtimeName, err)
		return nil, err
	}

	return runtime, nil
}

func observeRemoteRuntime(host string, apiKey string, runtimeName string, outputFormat string) (*detector.Runtime, error) {
	if outputFormat == "table" {
		fmt.Printf("Connecting to remote server: %s...\n\n", host)
	}

	client, err := grpcclient.NewClient(host, apiKey)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	ctx := context.Background()
	runtime, err := client.ObserveRuntime(ctx, runtimeName)
	if err != nil {
		return nil, err
	}

	return runtime, nil
}
