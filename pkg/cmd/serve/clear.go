package serve

import (
	"fmt"

	"github.com/spf13/cobra"
)

var clearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Clear resources",
	Long:  `Clear all API keys`,
}

var clearKeysCmd = &cobra.Command{
	Use:   "keys",
	Short: "Clear all API keys",
	Long:  `Remove all registered API keys`,
	RunE:  runClearKeys,
}

func init() {
	clearCmd.AddCommand(clearKeysCmd)
}

func runClearKeys(cmd *cobra.Command, args []string) error {
	store, err := getKeyStore()
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
