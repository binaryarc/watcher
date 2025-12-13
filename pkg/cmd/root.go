package cmd

import (
	"os"
	"path/filepath"

	"github.com/binaryarc/watcher/internal/keymanager"
	"github.com/binaryarc/watcher/pkg/cmd/compare"
	"github.com/binaryarc/watcher/pkg/cmd/get"
	"github.com/binaryarc/watcher/pkg/cmd/key"
	"github.com/spf13/cobra"
)

var (
	apiKey string
)

var rootCmd = &cobra.Command{
	Use:               "wctl",
	Short:             "👁️  Watcher - Observe your infrastructure",
	Long:              `Watcher is a kubectl-style CLI tool for observing runtime versions and services across your infrastructure.`,
	PersistentPreRunE: loadAPIKey,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringP("output", "o", "table", "Output format (table|json|yaml)")
	rootCmd.PersistentFlags().StringVar(&apiKey, "api-key", "", "API key for authentication")

	rootCmd.AddCommand(get.GetCmd)
	rootCmd.AddCommand(compare.CompareCmd)
	rootCmd.AddCommand(key.KeyCmd)
}

// loadAPIKey loads API key with priority: flag > env > file
func loadAPIKey(cmd *cobra.Command, args []string) error {
	// Skip for key commands
	if cmd.Parent() != nil && cmd.Parent().Use == "key" {
		return nil
	}

	// 1. Check if flag is set
	if apiKey != "" {
		return nil
	}

	// 2. Check environment variable
	if envKey := os.Getenv("WATCHER_API_KEY"); envKey != "" {
		apiKey = envKey
		return nil
	}

	// 3. Try to load from default file
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil
	}

	keysDir := filepath.Join(homeDir, ".watcher", "keys")
	manager, err := keymanager.NewManager(keysDir)
	if err != nil {
		return nil
	}

	key, err := manager.Load(keymanager.DefaultKeyName)
	if err != nil {
		return nil
	}

	apiKey = key
	return nil
}

// GetAPIKey returns the loaded API key
func GetAPIKey() string {
	return apiKey
}
