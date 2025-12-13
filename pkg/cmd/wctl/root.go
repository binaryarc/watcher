package wctl

import (
	"os"
	"path/filepath"

	"github.com/binaryarc/watcher/internal/keymanager"
	"github.com/binaryarc/watcher/pkg/cmd/wctl/compare"
	"github.com/binaryarc/watcher/pkg/cmd/wctl/get"
	"github.com/binaryarc/watcher/pkg/cmd/wctl/key"
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
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringP("output", "o", "table", "Output format (table|json|yaml)")
	rootCmd.PersistentFlags().StringVar(&apiKey, "api-key", "", "API key for authentication")

	rootCmd.AddCommand(get.Cmd)
	rootCmd.AddCommand(compare.Cmd)
	rootCmd.AddCommand(key.Cmd)
}

func loadAPIKey(cmd *cobra.Command, args []string) error {
	if cmd.Parent() != nil && cmd.Parent().Use == "key" {
		return nil
	}

	if apiKey != "" {
		return nil
	}

	if envKey := os.Getenv("WATCHER_API_KEY"); envKey != "" {
		apiKey = envKey
		return nil
	}

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

func GetAPIKey() string {
	return apiKey
}
