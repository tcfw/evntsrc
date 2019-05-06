package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tcfw/evntsrc/pkg/storer"
)

//NewStartCmd provides a version command
func NewStartCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "Start monitoring",
		Run: func(cmd *cobra.Command, args []string) {
			nats, _ := cmd.Flags().GetString("nats")
			port, _ := cmd.Flags().GetInt("port")
			storer.Start(nats, port)
		},
	}

	cmd.Flags().String("tracer", "jaeger-agent:5775", "endpoint of the jaeger-agent. Set to 'false' to disable tracing")
	cmd.Flags().String("nats", "localhost:4222", "endpoint for NATS server")
	cmd.Flags().IntP("port", "p", 443, "listening port for GRPC")
	cmd.Flags().Bool("migrate", false, "If migrations should run")
	viper.BindPFlag("tracer", cmd.Flags().Lookup("tracer"))
	viper.BindPFlag("migrate", cmd.Flags().Lookup("migrate"))

	return cmd
}
