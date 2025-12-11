package get

import (
	"fmt"

	"github.com/binaryarc/watcher/internal/detector"
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
}

func runGetRuntimes(cmd *cobra.Command, args []string) {
	fmt.Println("👁️  Observing all runtimes...\n")

	detectors := detector.GetAllDetectors()
	var runtimes []*detector.Runtime

	for _, det := range detectors {
		runtime, err := det.Detect()
		if err != nil {
			fmt.Printf("⚠️  Error detecting %s: %v\n", det.Name(), err)
			continue
		}

		if runtime.Found {
			runtimes = append(runtimes, runtime)
		}
	}

	if len(runtimes) == 0 {
		fmt.Println("❌ No runtimes detected on this system.")
		return
	}

	output.PrintRuntimesTable(runtimes)
	fmt.Printf("\n📊 Total: %d runtime(s) detected\n", len(runtimes))
}
