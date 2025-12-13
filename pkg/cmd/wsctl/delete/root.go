package delete

import "github.com/spf13/cobra"

var Cmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete resources",
	Long:  `Delete API keys`,
}

func init() {
	Cmd.AddCommand(keyCmd)
}
