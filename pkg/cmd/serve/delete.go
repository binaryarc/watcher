package serve

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	keyName string
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete resources",
	Long:  `Delete API keys`,
}

var deleteKeyCmd = &cobra.Command{
	Use:   "key",
	Short: "Delete an API key",
	Long:  `Delete an existing API key by name`,
	RunE:  runDeleteKey,
}

func init() {
	deleteKeyCmd.Flags().StringVar(&keyName, "name", "", "API key to delete (required)")
	deleteKeyCmd.MarkFlagRequired("name")

	deleteCmd.AddCommand(deleteKeyCmd)
}

func runDeleteKey(cmd *cobra.Command, args []string) error {
	store, err := getKeyStore()
	if err != nil {
		return err
	}

	if err := store.Remove(keyName); err != nil {
		return fmt.Errorf("failed to remove key: %w", err)
	}

	fmt.Println("API key removed successfully")
	fmt.Printf("Key: %s\n", maskKey(keyName))

	return nil
}
