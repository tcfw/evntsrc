package cmd

import (
	"github.com/spf13/cobra"
)

// NewDefaultCommand creates the `wui` command and its nested children.
func NewDefaultCommand() *cobra.Command {
	// Parent command to which all subcommands are added.
	cmds := &cobra.Command{
		Use:   "ttlscheduler",
		Short: "ttlscheduler",
		Long:  `ttlscheduler schedules and manages ttl replay worker allocations`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	cmds.AddCommand(
		NewSchedulerCmd(),
	)

	return cmds
}
