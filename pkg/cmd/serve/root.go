package serve

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "watcher-server",
	Short: "Watcher Server",
	Long:  `Watcher gRPC server for remote runtime observation`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(ServeCmd)
	rootCmd.AddCommand(keyCmd)
}
