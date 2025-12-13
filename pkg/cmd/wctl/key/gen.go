package key

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/binaryarc/watcher/internal/keymanager"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "key",
	Short: "Manage API key",
	Long:  `Generate and view API key for authentication with watcher servers`,
}

var genCmd = &cobra.Command{
	Use:   "gen",
	Short: "Generate a new API key",
	Long:  `Generate a new API key (replaces existing key)`,
	RunE:  runGenerate,
}

func init() {
	Cmd.AddCommand(genCmd)
}

func getKeyManager() (*keymanager.Manager, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	keysDir := filepath.Join(homeDir, ".watcher", "keys")
	return keymanager.NewManager(keysDir)
}

func runGenerate(cmd *cobra.Command, args []string) error {
	manager, err := getKeyManager()
	if err != nil {
		return err
	}

	apiKey, err := manager.Generate()
	if err != nil {
		return fmt.Errorf("failed to generate key: %w", err)
	}

	if err := manager.Save(keymanager.DefaultKeyName, apiKey); err != nil {
		return fmt.Errorf("failed to save key: %w", err)
	}

	fmt.Println("Generated API key:")
	fmt.Println()
	fmt.Println(apiKey)
	fmt.Println()
	fmt.Println("Register this key on the server:")
	fmt.Printf("wsctl add key %s \"<description>\"\n", apiKey)

	return nil
}
