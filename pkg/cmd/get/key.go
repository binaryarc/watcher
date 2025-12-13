package get

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/binaryarc/watcher/internal/keymanager"
	"github.com/spf13/cobra"
)

var keyCmd = &cobra.Command{
	Use:   "key",
	Short: "Get API key",
	Long:  `Display the current API key`,
	RunE:  runGetKey,
}

func runGetKey(cmd *cobra.Command, args []string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	keysDir := filepath.Join(homeDir, ".watcher", "keys")
	manager, err := keymanager.NewManager(keysDir)
	if err != nil {
		return err
	}

	apiKey, err := manager.Load(keymanager.DefaultKeyName)
	if err != nil {
		return fmt.Errorf("no API key found\n\nGenerate a new key with:\n   wctl key generate")
	}

	fmt.Println(apiKey)

	return nil
}
