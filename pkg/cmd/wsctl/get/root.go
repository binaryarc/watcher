package get

import "github.com/spf13/cobra"

var Cmd = &cobra.Command{
	Use:   "get",
	Short: "Get resources",
	Long:  `Get registered API keys`,
}

func init() {
	Cmd.AddCommand(keysCmd)
}
