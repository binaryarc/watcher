package compare

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/binaryarc/watcher/internal/detector"
	"github.com/binaryarc/watcher/internal/grpcclient"
	"github.com/binaryarc/watcher/internal/output"
	"github.com/spf13/cobra"
)

var runtimesCmd = &cobra.Command{
	Use:   "runtimes",
	Short: "Compare runtime versions across multiple servers",
	Long: `Compare runtime versions across multiple servers to identify version inconsistencies.

This command queries multiple servers in parallel and displays a comparison table
showing which runtimes have different versions across your infrastructure.`,
	Run: runCompareRuntimes,
}

func init() {
	CompareCmd.AddCommand(runtimesCmd)
	runtimesCmd.Flags().StringSlice("hosts", []string{}, "Comma-separated list of server addresses (required)")
	runtimesCmd.MarkFlagRequired("hosts")
}

// ServerRuntimes holds runtime information for a single server
type ServerRuntimes struct {
	Host     string
	Runtimes map[string]*detector.Runtime
	Error    error
}

func runCompareRuntimes(cmd *cobra.Command, args []string) {
	hosts, _ := cmd.Flags().GetStringSlice("hosts")
	outputFmt, _ := cmd.Flags().GetString("output")

	if len(hosts) == 0 {
		fmt.Println("❌ Error: --hosts flag is required")
		fmt.Println("Example: wctl compare runtimes --hosts server1:9090,server2:9090")
		return
	}

	if outputFmt == "table" {
		fmt.Printf("🌐 Comparing runtimes across %d server(s)...\n\n", len(hosts))
	}

	// API 키 가져오기
	apiKey, _ := cmd.Root().PersistentFlags().GetString("api-key")

	serverResults := fetchAllServers(hosts, apiKey, outputFmt)

	var successfulServers []ServerRuntimes
	for _, result := range serverResults {
		if result.Error != nil {
			if outputFmt == "table" {
				fmt.Printf("⚠️  Failed to connect to %s: %v\n", result.Host, result.Error)
			}
		} else {
			successfulServers = append(successfulServers, result)
		}
	}

	if len(successfulServers) == 0 {
		fmt.Println("\n❌ Failed to connect to all servers")
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

			ctx := context.Background()
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
				versions[i] = "-"
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

	nonEmptyVersions := make(map[string]int)
	emptyCount := 0
	errorCount := 0

	for _, v := range versions {
		if v == "-" {
			emptyCount++
		} else if v == "ERROR" {
			errorCount++
		} else {
			nonEmptyVersions[v]++
		}
	}

	if errorCount > 0 {
		return "ERROR"
	}

	if emptyCount == len(versions) {
		return "MISSING"
	}

	if emptyCount > 0 {
		return "PARTIAL"
	}

	if len(nonEmptyVersions) == 1 {
		return "SAME"
	}

	return "DIFF"
}

func printSummary(comparison *output.ComparisonData) {
	sameCount := 0
	diffCount := 0
	partialCount := 0

	for _, rt := range comparison.Runtimes {
		switch rt.Status {
		case "SAME":
			sameCount++
		case "DIFF":
			diffCount++
		case "PARTIAL":
			partialCount++
		}
	}

	fmt.Println("\n📊 Summary:")
	fmt.Printf("   • %d server(s) compared\n", len(comparison.Hosts))
	fmt.Printf("   • %d runtime(s) with differences\n", diffCount)
	if partialCount > 0 {
		fmt.Printf("   • %d runtime(s) partially installed\n", partialCount)
	}
	fmt.Printf("   • %d runtime(s) consistent\n", sameCount)
}
