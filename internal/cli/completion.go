package cli

import (
	"github.com/spf13/cobra"
)

func NewCompletionCommand(app *App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "completion [bash|zsh|fish|powershell]",
		Short: "Generate shell completion scripts",
		ValidArgs: []string{"bash", "zsh", "fish", "powershell"},
		Args:      cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			root := cmd.Root()
			switch args[0] {
			case "bash":
				return root.GenBashCompletion(app.Out)
			case "zsh":
				return root.GenZshCompletion(app.Out)
			case "fish":
				return root.GenFishCompletion(app.Out, true)
			case "powershell":
				return root.GenPowerShellCompletionWithDesc(app.Out)
			default:
				return cmd.Usage()
			}
		},
	}
	return cmd
}
