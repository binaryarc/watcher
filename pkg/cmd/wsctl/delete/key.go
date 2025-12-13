package delete

import (
	"fmt"

	"github.com/binaryarc/watcher/pkg/cmd/wsctl/common"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete resources",
	Long:  `Delete API keys`,
}

var (
	keyName string
)

var keyCmd = &cobra.Command{
	Use:   "key",
	Short: "Delete an API key",
	Long:  `Delete an existing API key by name`,
	RunE:  runDeleteKey,
}

func init() {
	keyCmd.Flags().StringVar(&keyName, "name", "", "API key to delete (required)")
	keyCmd.MarkFlagRequired("name")
	Cmd.AddCommand(keyCmd)
}

func runDeleteKey(cmd *cobra.Command, args []string) error {
	store, err := common.KeyStore()
	if err != nil {
		return err
	}

	if err := store.Remove(keyName); err != nil {
		return fmt.Errorf("failed to remove key: %w", err)
	}

	fmt.Println("API key removed successfully")
	fmt.Printf("Key: %s\n", common.MaskKey(keyName))

	return nil
}
