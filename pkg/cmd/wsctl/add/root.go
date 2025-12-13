package add

import "github.com/spf13/cobra"

var Cmd = &cobra.Command{
	Use:   "add",
	Short: "Add resources",
	Long:  `Add API keys`,
}

func init() {
	Cmd.AddCommand(keyCmd)
}
