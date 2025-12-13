package serve

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "wsctl",
	Short: "Watcher Server Control",
	Long:  `Control tool for Watcher gRPC server`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(RunCmd)
	rootCmd.AddCommand(getCmd)
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(deleteCmd)
	rootCmd.AddCommand(clearCmd)
}
