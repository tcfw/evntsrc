package cmd

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tcfw/evntsrc/internal/ttlworker"
)

//NewStartCmd provides a version command
func NewStartCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "Start the worker",
		Run: func(cmd *cobra.Command, args []string) {
			port, _ := cmd.Flags().GetInt("port")
			worker, err := ttlworker.NewWorker(port)
			if err != nil {
				log.Fatalf("Failed to consturct worker: %s", err.Error())
			}
			worker.StartAndWait()
		},
	}

	cmd.Flags().String("tracer", "jaeger-agent:5775", "endpoint of the jaeger-agent. Set to 'false' to disable tracing")
	cmd.Flags().IntP("port", "p", 443, "listening port for GRPC")
	cmd.Flags().BoolP("verbose", "v", false, "Display status every 30 seconds")
	viper.BindPFlag("tracer", cmd.Flags().Lookup("tracer"))

	return cmd
}
