package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tcfw/evntsrc/pkg/stsmetrics"
)

//NewJobCmd triggers of metrics requests
func NewJobCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "job",
		Short: "Trigger metrics jobs for streams",
		Run: func(cmd *cobra.Command, args []string) {
			stream, _ := cmd.Flags().GetInt32("stream")
			nats, _ := cmd.Flags().GetString("nats")
			stsmetrics.JobRequest(nats, stream)
		},
	}

	cmd.Flags().String("tracer", "jaeger-agent:5775", "endpoint of the jaeger-agent. Set to 'false' to disable tracing")
	cmd.Flags().Int32("stream", 0, "stream ID to trigger job for (default all streams)")
	cmd.Flags().String("nats", "localhost:4222", "endpoint for NATS server")
	viper.BindPFlag("tracer", cmd.Flags().Lookup("tracer"))

	return cmd
}
