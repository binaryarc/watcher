package get

import (
	"github.com/spf13/cobra"
)

var GetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get resources",
	Long:  `Get information about runtimes, services, and other resources`,
}
