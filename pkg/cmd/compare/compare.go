package compare

import (
	"github.com/spf13/cobra"
)

// CompareCmd represents the compare command
var CompareCmd = &cobra.Command{
	Use:   "compare",
	Short: "Compare runtime versions across multiple servers",
	Long: `Compare runtime versions across multiple servers to identify inconsistencies.

Examples:
  # Compare all runtimes across three servers
  wctl compare runtimes --hosts server1:9090,server2:9090,server3:9090

  # Compare with JSON output
  wctl compare runtimes --hosts server1:9090,server2:9090 -o json`,
}

func init() {
	// Subcommands will be added here
}
