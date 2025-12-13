package get

import (
	"fmt"
	"os"

	"github.com/binaryarc/watcher/pkg/cmd/wsctl/common"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "get",
	Short: "Get resources",
	Long:  `Get registered API keys`,
}

var keysCmd = &cobra.Command{
	Use:   "keys",
	Short: "Get all registered API keys",
	Long:  `List all registered API keys`,
	RunE:  runGetKeys,
}

func init() {
	Cmd.AddCommand(keysCmd)
}

func runGetKeys(cmd *cobra.Command, args []string) error {
	store, err := common.KeyStore()
	if err != nil {
		return err
	}

	keys := store.List()

	if len(keys) == 0 {
		fmt.Println("No API keys registered")
		fmt.Println()
		fmt.Println("Add a key with:")
		fmt.Println("   wsctl add key <api-key> \"<description>\"")
		return nil
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.Header([]string{"Key (masked)", "Description", "Created At"})

	for _, keyInfo := range keys {
		table.Append([]string{
			common.MaskKey(keyInfo.Key),
			keyInfo.Description,
			keyInfo.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	fmt.Printf("Registered API Keys (%d):\n", len(keys))
	table.Render()

	return nil
}
