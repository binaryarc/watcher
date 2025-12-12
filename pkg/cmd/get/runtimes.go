package get

import (
	"context"
	"fmt"

	"github.com/binaryarc/watcher/internal/detector"
	"github.com/binaryarc/watcher/internal/grpcclient"
	"github.com/binaryarc/watcher/internal/output"
	"github.com/spf13/cobra"
)

var runtimesCmd = &cobra.Command{
	Use:   "runtimes",
	Short: "Get all detected runtimes",
	Long:  `Scan and display all detected runtime versions on the system`,
	Run:   runGetRuntimes,
}

func init() {
	GetCmd.AddCommand(runtimesCmd)
	runtimesCmd.Flags().String("host", "", "Remote server address (e.g., server:9090)")
}

func runGetRuntimes(cmd *cobra.Command, args []string) {
	outputFormat, _ := cmd.Flags().GetString("output")
	host, _ := cmd.Flags().GetString("host")

	var runtimes []*detector.Runtime
	var err error

	// Remote vs Local detection
	if host != "" {
		// Remote observation via gRPC
		runtimes, err = observeRemoteRuntimes(host, outputFormat)
		if err != nil {
			fmt.Printf("❌ Failed to observe remote server: %v\n", err)
			return
		}
	} else {
		// Local detection
		runtimes = observeLocalRuntimes(outputFormat)
	}

	if len(runtimes) == 0 {
		if outputFormat == "table" {
			fmt.Println("❌ No runtimes detected.")
		}
		return
	}

	// Output formatting
	switch outputFormat {
	case "json":
		if err := output.PrintRuntimesJSON(runtimes); err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	case "yaml":
		if err := output.PrintRuntimesYAML(runtimes); err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	case "table":
		output.PrintRuntimesTable(runtimes)
		fmt.Printf("\n📊 Total: %d runtime(s) detected\n", len(runtimes))
	default:
		fmt.Printf("Unknown output format: %s\n", outputFormat)
		fmt.Println("Supported formats: table, json, yaml")
	}
}

func observeLocalRuntimes(outputFormat string) []*detector.Runtime {
	if outputFormat == "table" {
		fmt.Println("👁️  Observing local runtimes...\n")
	}

	detectors := detector.GetAllDetectors()
	var runtimes []*detector.Runtime

	for _, det := range detectors {
		runtime, err := det.Detect()
		if err != nil {
			if outputFormat == "table" {
				fmt.Printf("⚠️  Error detecting %s: %v\n", det.Name(), err)
			}
			continue
		}

		if runtime.Found {
			runtimes = append(runtimes, runtime)
		}
	}

	return runtimes
}

// observeRemoteRuntimes fetches runtime info from remote server via gRPC
func observeRemoteRuntimes(host string, outputFormat string) ([]*detector.Runtime, error) {
	if outputFormat == "table" {
		fmt.Printf("🌐 Connecting to remote server: %s...\n\n", host)
	}

	// Create gRPC client
	client, err := grpcclient.NewClient(host)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	// Fetch runtimes
	ctx := context.Background()
	runtimes, err := client.ObserveRuntimes(ctx)
	if err != nil {
		return nil, err
	}

	return runtimes, nil
}
