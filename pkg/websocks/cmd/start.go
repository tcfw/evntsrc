package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tcfw/evntsrc/pkg/websocks"
)

//NewServeCmd provides a version command
func NewServeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Start the service",
		Run: func(cmd *cobra.Command, args []string) {
			port, _ := cmd.Flags().GetInt("port")
			websocks.Run(port)
		},
	}

	cmd.Flags().String("tracer", "jaeger-agent:5775", "endpoint of the jaeger-agent. Set to 'false' to disable tracing")
	cmd.Flags().Int("port", 12345, "listening port for GRPC")
	viper.BindPFlag("tracer", cmd.Flags().Lookup("tracer"))

	return cmd
}
