package key

import "github.com/spf13/cobra"

var Cmd = &cobra.Command{
	Use:   "key",
	Short: "Manage API key",
	Long:  `Generate and view API key for authentication with watcher servers`,
}

func init() {
	Cmd.AddCommand(genCmd)
}
