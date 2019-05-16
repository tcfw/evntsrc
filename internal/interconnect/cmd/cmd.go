package cmd

import (
	"github.com/spf13/cobra"
)

// NewDefaultCommand creates the `wui` command and its nested children.
func NewDefaultCommand() *cobra.Command {
	// Parent command to which all subcommands are added.
	cmds := &cobra.Command{
		Use:   "interconnect",
		Short: "interconnect",
		Long:  `interconnect provide a GRPC bi-directional stream between federated zones`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	cmds.AddCommand(NewServeCmd())
	cmds.AddCommand(NewConnectCmd())

	return cmds
}
