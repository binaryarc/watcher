package serve

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/binaryarc/watcher/internal/keystore"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var keyCmd = &cobra.Command{
	Use:   "key",
	Short: "Manage API keys",
	Long:  `Add, remove, and list API keys for server authentication`,
}

var keyAddCmd = &cobra.Command{
	Use:   "add <api-key> [description]",
	Short: "Add a new API key",
	Long:  `Add a new API key to allow clients to authenticate`,
	Args:  cobra.RangeArgs(1, 2),
	RunE:  runKeyAdd,
}

var keyRemoveCmd = &cobra.Command{
	Use:   "remove <api-key>",
	Short: "Remove an API key",
	Long:  `Remove an existing API key`,
	Args:  cobra.ExactArgs(1),
	RunE:  runKeyRemove,
}

var keyListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all API keys",
	Long:  `List all registered API keys`,
	RunE:  runKeyList,
}

var keyClearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Clear all API keys",
	Long:  `Remove all registered API keys`,
	RunE:  runKeyClear,
}

func init() {
	keyCmd.AddCommand(keyAddCmd)
	keyCmd.AddCommand(keyRemoveCmd)
	keyCmd.AddCommand(keyListCmd)
	keyCmd.AddCommand(keyClearCmd)
}

func getKeyStore() (*keystore.Store, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	keysDir := filepath.Join(homeDir, ".watcher", "server")
	if err := os.MkdirAll(keysDir, 0700); err != nil {
		return nil, fmt.Errorf("failed to create keys directory: %w", err)
	}

	keystorePath := filepath.Join(keysDir, "keys.json")
	return keystore.NewStore(keystorePath)
}

func runKeyAdd(cmd *cobra.Command, args []string) error {
	apiKey := args[0]
	description := ""
	if len(args) > 1 {
		description = args[1]
	}

	store, err := getKeyStore()
	if err != nil {
		return err
	}

	if err := store.Add(apiKey, description); err != nil {
		return fmt.Errorf("failed to add key: %w", err)
	}

	fmt.Println("✅ API key added successfully")
	if description != "" {
		fmt.Printf("   Description: %s\n", description)
	}
	fmt.Printf("   Key: %s\n", maskKey(apiKey))

	return nil
}

func runKeyRemove(cmd *cobra.Command, args []string) error {
	apiKey := args[0]

	store, err := getKeyStore()
	if err != nil {
		return err
	}

	if err := store.Remove(apiKey); err != nil {
		return fmt.Errorf("failed to remove key: %w", err)
	}

	fmt.Println("✅ API key removed successfully")
	fmt.Printf("   Key: %s\n", maskKey(apiKey))

	return nil
}

func runKeyList(cmd *cobra.Command, args []string) error {
	store, err := getKeyStore()
	if err != nil {
		return err
	}

	keys := store.List()

	if len(keys) == 0 {
		fmt.Println("No API keys registered")
		fmt.Println("\nAdd a key with:")
		fmt.Println("   watcher-server key add <api-key> \"<description>\"")
		return nil
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.Header([]string{"Key (masked)", "Description", "Created At"})

	for _, keyInfo := range keys {
		table.Append([]string{
			maskKey(keyInfo.Key),
			keyInfo.Description,
			keyInfo.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	fmt.Printf("Registered API Keys (%d):\n", len(keys))
	table.Render()

	return nil
}

func runKeyClear(cmd *cobra.Command, args []string) error {
	store, err := getKeyStore()
	if err != nil {
		return err
	}

	// Confirmation
	fmt.Print("⚠️  Are you sure you want to remove ALL API keys? (yes/no): ")
	var response string
	fmt.Scanln(&response)

	if response != "yes" {
		fmt.Println("Cancelled")
		return nil
	}

	if err := store.Clear(); err != nil {
		return fmt.Errorf("failed to clear keys: %w", err)
	}

	fmt.Println("✅ All API keys removed")

	return nil
}

// maskKey masks the API key for display (shows first 10 and last 4 chars)
func maskKey(key string) string {
	if len(key) <= 14 {
		return key
	}
	return key[:10] + "..." + key[len(key)-4:]
}
