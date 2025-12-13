package serve

import (
	"fmt"

	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add resources",
	Long:  `Add API keys`,
}

var addKeyCmd = &cobra.Command{
	Use:   "key <api-key> [description]",
	Short: "Add a new API key",
	Long:  `Add a new API key to allow clients to authenticate`,
	Args:  cobra.RangeArgs(1, 2),
	RunE:  runAddKey,
}

func init() {
	addCmd.AddCommand(addKeyCmd)
}

func runAddKey(cmd *cobra.Command, args []string) error {
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

	fmt.Println("API key added successfully")
	if description != "" {
		fmt.Printf("Description: %s\n", description)
	}
	fmt.Printf("Key: %s\n", maskKey(apiKey))

	return nil
}
