package wsctl

import (
	"os"

	"github.com/binaryarc/watcher/pkg/cmd/wsctl/add"
	"github.com/binaryarc/watcher/pkg/cmd/wsctl/clear"
	"github.com/binaryarc/watcher/pkg/cmd/wsctl/delete"
	"github.com/binaryarc/watcher/pkg/cmd/wsctl/get"
	"github.com/binaryarc/watcher/pkg/cmd/wsctl/run"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "wsctl",
	Short: "Watcher Server Control",
	Long:  `Control tool for Watcher gRPC server`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(run.Cmd)
	rootCmd.AddCommand(get.Cmd)
	rootCmd.AddCommand(add.Cmd)
	rootCmd.AddCommand(delete.Cmd)
	rootCmd.AddCommand(clear.Cmd)
}
