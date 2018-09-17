package cmd

import (
	"github.com/spf13/cobra"
)

// NewDefaultCommand creates the `wui` command and its nested children.
func NewDefaultCommand() *cobra.Command {
	// Parent command to which all subcommands are added.
	cmds := &cobra.Command{
		Use:   "stsmetrics",
		Short: "stsmetrics",
		Long:  `stsmetrics calculates stream time series metrics`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	cmds.AddCommand(NewWatchCmd())
	cmds.AddCommand(NewJobCmd())

	return cmds
}
