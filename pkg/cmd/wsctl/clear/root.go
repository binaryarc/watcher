package clear

import "github.com/spf13/cobra"

var Cmd = &cobra.Command{
	Use:   "clear",
	Short: "Clear resources",
	Long:  `Clear all API keys`,
}

func init() {
	Cmd.AddCommand(keysCmd)
}
