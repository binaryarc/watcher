package compare

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/binaryarc/watcher/internal/detector"
	"github.com/binaryarc/watcher/internal/grpcclient"
	"github.com/binaryarc/watcher/internal/output"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "compare",
	Short: "Compare runtime versions across multiple servers",
	Long: `Compare runtime versions across multiple servers to identify inconsistencies.

Examples:
  # Compare all runtimes across three servers
  wctl compare runtimes --hosts server1:9090,server2:9090,server3:9090

  # Compare with JSON output
  wctl compare runtimes --hosts server1:9090,server2:9090 -o json`,
}

var runtimesCmd = &cobra.Command{
	Use:   "runtimes",
	Short: "Compare runtime versions across multiple servers",
	Long: `Compare runtime versions across multiple servers to identify version inconsistencies.

This command queries multiple servers in parallel and displays a comparison table
showing which runtimes have different versions across your infrastructure.`,
	Run: runCompareRuntimes,
}

func init() {
	Cmd.AddCommand(runtimesCmd)
	runtimesCmd.Flags().StringSlice("hosts", []string{}, "Comma-separated list of server addresses (required)")
	runtimesCmd.MarkFlagRequired("hosts")
}

type ServerRuntimes struct {
	Host     string
	Runtimes map[string]*detector.Runtime
	Error    error
}

const hostRequestTimeout = 10 * time.Second

func runCompareRuntimes(cmd *cobra.Command, args []string) {
	hosts, _ := cmd.Flags().GetStringSlice("hosts")
	for i, host := range hosts {
		hosts[i] = strings.TrimSpace(host)
	}
	outputFmt, _ := cmd.Flags().GetString("output")

	if len(hosts) == 0 {
		fmt.Println("Error: --hosts flag is required")
		fmt.Println("Example: wctl compare runtimes --hosts server1:9090,server2:9090")
		return
	}

	if outputFmt == "table" {
		fmt.Printf("Comparing runtimes across %d server(s)...\n\n", len(hosts))
	}

	apiKey, _ := cmd.Root().PersistentFlags().GetString("api-key")

	serverResults := fetchAllServers(hosts, apiKey, outputFmt)

	var successfulServers []ServerRuntimes
	for _, result := range serverResults {
		if result.Error != nil {
			if outputFmt == "table" {
				fmt.Printf("Warning: Failed to connect to %s: %v\n", result.Host, result.Error)
			}
		} else {
			successfulServers = append(successfulServers, result)
		}
	}

	if len(successfulServers) == 0 {
		fmt.Println("\nFailed to connect to all servers")
		return
	}

	if outputFmt == "table" && len(successfulServers) < len(hosts) {
		fmt.Println()
	}

	comparison := buildComparison(successfulServers)

	switch outputFmt {
	case "json":
		if err := output.PrintComparisonJSON(comparison); err != nil {
			fmt.Printf("Error printing JSON: %v\n", err)
		}
	case "yaml":
		if err := output.PrintComparisonYAML(comparison); err != nil {
			fmt.Printf("Error printing YAML: %v\n", err)
		}
	case "table":
		output.PrintComparisonTable(comparison)
		printSummary(comparison)
	default:
		fmt.Printf("Unknown output format: %s\n", outputFmt)
	}
}

func fetchAllServers(hosts []string, apiKey string, outputFmt string) []ServerRuntimes {
	var wg sync.WaitGroup
	results := make([]ServerRuntimes, len(hosts))

	for i, host := range hosts {
		wg.Add(1)
		go func(index int, hostAddr string) {
			defer wg.Done()

			client, err := grpcclient.NewClient(hostAddr, apiKey)
			if err != nil {
				results[index] = ServerRuntimes{
					Host:  hostAddr,
					Error: err,
				}
				return
			}
			defer client.Close()

			ctx, cancel := context.WithTimeout(context.Background(), hostRequestTimeout)
			defer cancel()

			runtimes, err := client.ObserveRuntimes(ctx)
			if err != nil {
				results[index] = ServerRuntimes{
					Host:  hostAddr,
					Error: err,
				}
				return
			}

			runtimeMap := make(map[string]*detector.Runtime)
			for _, rt := range runtimes {
				runtimeMap[rt.Name] = rt
			}

			results[index] = ServerRuntimes{
				Host:     hostAddr,
				Runtimes: runtimeMap,
				Error:    nil,
			}
		}(i, host)
	}

	wg.Wait()
	return results
}

func buildComparison(serverResults []ServerRuntimes) *output.ComparisonData {
	runtimeNames := make(map[string]bool)
	for _, server := range serverResults {
		if server.Error == nil {
			for name := range server.Runtimes {
				runtimeNames[name] = true
			}
		}
	}

	var runtimeComparisons []output.RuntimeComparison
	for name := range runtimeNames {
		versions := make([]string, len(serverResults))
		for i, server := range serverResults {
			if server.Error != nil {
				versions[i] = "ERROR"
			} else if rt, found := server.Runtimes[name]; found {
				versions[i] = rt.Version
			} else {
				versions[i] = "x"
			}
		}

		status := determineStatus(versions)

		runtimeComparisons = append(runtimeComparisons, output.RuntimeComparison{
			Name:     name,
			Versions: versions,
			Status:   status,
		})
	}

	hosts := make([]string, len(serverResults))
	for i, server := range serverResults {
		hostParts := strings.Split(server.Host, ":")
		if server.Error != nil {
			hosts[i] = hostParts[0] + " (ERR)"
		} else {
			hosts[i] = hostParts[0]
		}
	}

	return &output.ComparisonData{
		Hosts:    hosts,
		Runtimes: runtimeComparisons,
	}
}

func determineStatus(versions []string) string {
	if len(versions) == 0 {
		return "UNKNOWN"
	}

	baseVersion := ""
	for _, version := range versions {
		if version == "ERROR" {
			return "UNKNOWN"
		}
		if version == "x" {
			continue
		}
		if baseVersion == "" {
			baseVersion = version
			continue
		}
		if version != baseVersion {
			return "MISMATCH"
		}
	}

	if baseVersion == "" {
		return "UNKNOWN"
	}

	return "MATCH"
}

func printSummary(comparison *output.ComparisonData) {
	fmt.Println()
	fmt.Println("Summary:")
	totalRuntimes := len(comparison.Runtimes)
	mismatched := 0
	unknown := 0

	for _, runtime := range comparison.Runtimes {
		switch runtime.Status {
		case "MISMATCH":
			mismatched++
		case "UNKNOWN":
			unknown++
		}
	}

	fmt.Printf("- Total runtimes compared: %d\n", totalRuntimes)
	fmt.Printf("- Matching runtimes: %d\n", totalRuntimes-mismatched-unknown)
	fmt.Printf("- Mismatched runtimes: %d\n", mismatched)
	fmt.Printf("- Unknown status: %d\n", unknown)

	if mismatched > 0 {
		fmt.Println("\nRecommendation:")
		fmt.Println("- Review mismatched runtimes and plan upgrades/downgrades to standardize versions.")
	}
}
