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
			nats, _ := cmd.Flags().GetString("nats")
			websocks.Run(port, nats)
		},
	}

	cmd.Flags().String("tracer", "jaeger-agent:5775", "endpoint of the jaeger-agent. Set to 'false' to disable tracing")
	cmd.Flags().Int("port", 12345, "listening port for GRPC")
	cmd.Flags().String("nats", "localhost:4222", "endpoint for NATS server")
	viper.BindPFlag("tracer", cmd.Flags().Lookup("tracer"))

	return cmd
}