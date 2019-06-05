package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

//NewSchedulerCmd provides a version command
func NewSchedulerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run",
		Short: "Start the scheduler",
		Run: func(cmd *cobra.Command, args []string) {
		},
	}

	cmd.Flags().Bool("once", false, "Only run one pass of the scheduler")

	cmd.Flags().String("tracer", "jaeger-agent:5775", "endpoint of the jaeger-agent. Set to 'false' to disable tracing")
	cmd.Flags().IntP("port", "p", 443, "listening port for GRPC")
	viper.BindPFlag("tracer", cmd.Flags().Lookup("tracer"))

	return cmd
}
