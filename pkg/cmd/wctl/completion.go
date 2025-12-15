package wctl

import (
	"fmt"

	"github.com/spf13/cobra"
)

var completionCmd = &cobra.Command{
	Use:   "completion [bash|zsh|fish|powershell]",
	Short: "Generate shell completion scripts",
	Long: `Emit completion scripts for bash, zsh, fish, and PowerShell.

Examples:
  wctl completion bash > /etc/bash_completion.d/wctl
  wctl completion zsh >> ~/.zshrc`,
	Args:      cobra.ExactValidArgs(1),
	ValidArgs: []string{"bash", "zsh", "fish", "powershell"},
	RunE: func(cmd *cobra.Command, args []string) error {
		root := cmd.Root()
		if root == nil {
			return fmt.Errorf("root command not found")
		}

		writer := cmd.OutOrStdout()
		switch args[0] {
		case "bash":
			return root.GenBashCompletion(writer)
		case "zsh":
			return root.GenZshCompletion(writer)
		case "fish":
			return root.GenFishCompletion(writer, true)
		case "powershell":
			return root.GenPowerShellCompletionWithDesc(writer)
		default:
			return fmt.Errorf("unsupported shell: %s", args[0])
		}
	},
}

func init() {
	rootCmd.AddCommand(completionCmd)
}
