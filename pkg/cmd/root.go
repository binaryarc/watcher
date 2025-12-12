package cmd

import (
	"os"

	"github.com/binaryarc/watcher/pkg/cmd/compare"
	"github.com/binaryarc/watcher/pkg/cmd/get"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "wctl",
	Short: "👁️  Watcher - Observe your infrastructure",
	Long:  `Watcher is a kubectl-style CLI tool for observing runtime versions and services across your infrastructure.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {

	rootCmd.PersistentFlags().StringP("output", "o", "table", "Output format (table|json|yaml)")
	rootCmd.AddCommand(get.GetCmd)
	rootCmd.AddCommand(compare.CompareCmd)
}
