package clear

import (
	"fmt"

	"github.com/binaryarc/watcher/pkg/cmd/wsctl/common"
	"github.com/spf13/cobra"
)

var keysCmd = &cobra.Command{
	Use:   "keys",
	Short: "Clear all API keys",
	Long:  `Remove all registered API keys`,
	RunE:  runClearKeys,
}

func runClearKeys(cmd *cobra.Command, args []string) error {
	store, err := common.KeyStore()
	if err != nil {
		return err
	}

	fmt.Print("Are you sure you want to remove ALL API keys? (yes/no): ")
	var response string
	fmt.Scanln(&response)

	if response != "yes" {
		fmt.Println("Cancelled")
		return nil
	}

	if err := store.Clear(); err != nil {
		return fmt.Errorf("failed to clear keys: %w", err)
	}

	fmt.Println("All API keys removed")

	return nil
}
