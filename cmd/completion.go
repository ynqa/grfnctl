package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var completionCmd = &cobra.Command{
	Use:       "completion [shell]",
	Short:     "Generate shell completion script",
	Long:      "Generate an autocompletion script for supported shells.",
	Args:      cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	ValidArgs: []string{"bash", "zsh"},
	RunE: func(cmd *cobra.Command, args []string) error {
		root := cmd.Root()
		switch args[0] {
		case "bash":
			// Generate Bash completion script to STDOUT.
			return root.GenBashCompletion(os.Stdout)
		case "zsh":
			// Generate Zsh completion script to STDOUT.
			return root.GenZshCompletion(os.Stdout)
		default:
			return fmt.Errorf("unsupported shell: %s", args[0])
		}
	},
}
