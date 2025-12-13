package get

import "github.com/spf13/cobra"

var Cmd = &cobra.Command{
	Use:   "get",
	Short: "Get resources",
	Long:  `Get runtime information or API key`,
}

func init() {
	Cmd.AddCommand(runtimesCmd)
	Cmd.AddCommand(runtimeCmd)
	Cmd.AddCommand(keyCmd)
}
