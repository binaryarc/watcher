package get

import (
	"github.com/spf13/cobra"
)

var GetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get resources",
	Long:  `Get runtime information or API key`,
}

func init() {
	GetCmd.AddCommand(runtimesCmd)
	GetCmd.AddCommand(runtimeCmd)
	GetCmd.AddCommand(keyCmd)
}
