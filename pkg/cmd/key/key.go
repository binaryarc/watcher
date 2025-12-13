package key

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/binaryarc/watcher/internal/keymanager"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var (
	copyFlag bool
)

// KeyCmd is the root command for key management
var KeyCmd = &cobra.Command{
	Use:   "key",
	Short: "Manage API keys",
	Long:  `Generate, list, and manage API keys for authentication with watcher servers`,
}

var generateCmd = &cobra.Command{
	Use:   "generate [name]",
	Short: "Generate a new API key",
	Long:  `Generate a new API key and save it with the given name (default: "default")`,
	Args:  cobra.MaximumNArgs(1),
	RunE:  runGenerate,
}

var showCmd = &cobra.Command{
	Use:   "show [name]",
	Short: "Show an API key",
	Long:  `Display a saved API key (default: "default")`,
	Args:  cobra.MaximumNArgs(1),
	RunE:  runShow,
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all API keys",
	Long:  `List all saved API keys`,
	RunE:  runList,
}

var deleteCmd = &cobra.Command{
	Use:   "delete <name>",
	Short: "Delete an API key",
	Long:  `Delete a saved API key by name`,
	Args:  cobra.ExactArgs(1),
	RunE:  runDelete,
}

func init() {
	showCmd.Flags().BoolVar(&copyFlag, "copy", false, "Copy the key to clipboard (requires xclip or pbcopy)")

	KeyCmd.AddCommand(generateCmd)
	KeyCmd.AddCommand(showCmd)
	KeyCmd.AddCommand(listCmd)
	KeyCmd.AddCommand(deleteCmd)
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
	name := keymanager.DefaultKeyName
	if len(args) > 0 {
		name = args[0]
	}

	manager, err := getKeyManager()
	if err != nil {
		return err
	}

	apiKey, err := manager.Generate()
	if err != nil {
		return fmt.Errorf("failed to generate key: %w", err)
	}

	if err := manager.Save(name, apiKey); err != nil {
		return fmt.Errorf("failed to save key: %w", err)
	}

	fmt.Printf("✅ Generated API key: %s\n", name)
	fmt.Printf("\n%s\n\n", apiKey)
	fmt.Println("📝 Register this key on the server with:")
	fmt.Printf("   watcher-server key add %s \"<description>\"\n", apiKey)

	return nil
}

func runShow(cmd *cobra.Command, args []string) error {
	name := keymanager.DefaultKeyName
	if len(args) > 0 {
		name = args[0]
	}

	manager, err := getKeyManager()
	if err != nil {
		return err
	}

	apiKey, err := manager.Load(name)
	if err != nil {
		return err
	}

	fmt.Printf("API Key (%s):\n%s\n", name, apiKey)

	if copyFlag {
		fmt.Println("\n⚠️  Clipboard copy not implemented yet")
	}

	return nil
}

func runList(cmd *cobra.Command, args []string) error {
	manager, err := getKeyManager()
	if err != nil {
		return err
	}

	keys, err := manager.List()
	if err != nil {
		return err
	}

	if len(keys) == 0 {
		fmt.Println("No API keys found")
		fmt.Println("\nGenerate a new key with:")
		fmt.Println("   wctl key generate")
		return nil
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.Header([]string{"Name", "Preview"})

	for _, keyName := range keys {
		apiKey, err := manager.Load(keyName)
		if err != nil {
			continue
		}

		preview := apiKey
		if len(preview) > 40 {
			preview = preview[:40] + "..."
		}

		table.Append([]string{keyName, preview})
	}

	table.Render()
	return nil
}

func runDelete(cmd *cobra.Command, args []string) error {
	name := args[0]

	manager, err := getKeyManager()
	if err != nil {
		return err
	}

	if err := manager.Delete(name); err != nil {
		return err
	}

	fmt.Printf("✅ Deleted API key: %s\n", name)
	return nil
}
