package cmd

import (
	"github.com/spf13/cobra"
)

// NewDefaultCommand creates the `wui` command and its nested children.
func NewDefaultCommand() *cobra.Command {
	// Parent command to which all subcommands are added.
	cmds := &cobra.Command{
		Use:   "metrics",
		Short: "metrics",
		Long:  `metrics gateways stream timeseries metrics`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	cmds.AddCommand(NewServeCmd())

	return cmds
}
