package get

import (
	"context"
	"fmt"

	"github.com/binaryarc/watcher/internal/detector"
	"github.com/binaryarc/watcher/internal/grpcclient"
	"github.com/binaryarc/watcher/internal/output"
	"github.com/spf13/cobra"
)

var (
	detectorsProvider = detector.GetAllDetectors
	grpcClientFactory = func(host, apiKey string) (grpcRuntimeClient, error) {
		return grpcclient.NewClient(host, apiKey)
	}
)

type grpcRuntimeClient interface {
	ObserveRuntimes(ctx context.Context) ([]*detector.Runtime, error)
	Close() error
}

var runtimesCmd = &cobra.Command{
	Use:   "runtimes",
	Short: "Get all detected runtimes",
	Long:  `Scan and display all detected runtime versions on the system`,
	Run:   runGetRuntimes,
}

func init() {
	Cmd.AddCommand(runtimesCmd)
	runtimesCmd.Flags().String("host", "", "Remote server address (e.g., server:9090)")
}

func runGetRuntimes(c *cobra.Command, args []string) {
	outputFormat, _ := c.Flags().GetString("output")
	host, _ := c.Flags().GetString("host")

	var runtimes []*detector.Runtime
	var err error

	if host != "" {
		apiKey, _ := c.Root().PersistentFlags().GetString("api-key")
		runtimes, err = observeRemoteRuntimes(host, apiKey, outputFormat)
		if err != nil {
			fmt.Printf("Failed to observe remote server: %v\n", err)
			return
		}
	} else {
		runtimes = observeLocalRuntimes(outputFormat)
	}

	if len(runtimes) == 0 {
		if outputFormat == "table" {
			fmt.Println("No runtimes detected.")
		}
		return
	}

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
		fmt.Printf("\nTotal: %d runtime(s) detected\n", len(runtimes))
	default:
		fmt.Printf("Unknown output format: %s\n", outputFormat)
		fmt.Println("Supported formats: table, json, yaml")
	}
}

func observeLocalRuntimes(outputFormat string) []*detector.Runtime {
	if outputFormat == "table" {
		fmt.Println("Observing local runtimes...")
		fmt.Println()
	}

	detectors := detectorsProvider()
	var runtimes []*detector.Runtime

	for _, det := range detectors {
		runtime, err := det.Detect()
		if err != nil {
			if outputFormat == "table" {
				fmt.Printf("Warning: Error detecting %s: %v\n", det.Name(), err)
			}
			continue
		}

		if runtime.Found {
			runtimes = append(runtimes, runtime)
		}
	}

	return runtimes
}

func observeRemoteRuntimes(host string, apiKey string, outputFormat string) ([]*detector.Runtime, error) {
	if outputFormat == "table" {
		fmt.Printf("Connecting to remote server: %s...\n\n", host)
	}

	client, err := grpcClientFactory(host, apiKey)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	ctx := context.Background()
	runtimes, err := client.ObserveRuntimes(ctx)
	if err != nil {
		return nil, err
	}

	return runtimes, nil
}
