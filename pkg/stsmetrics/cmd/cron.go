package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tcfw/evntsrc/pkg/stsmetrics"
)

//NewCronCmd triggers of metrics requests
func NewCronCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cron",
		Short: "Trigger metrics jobs on a cron like basis",
		Run: func(cmd *cobra.Command, args []string) {
			nats, _ := cmd.Flags().GetString("nats")
			cronExpr, _ := cmd.Flags().GetString("cron")
			stsmetrics.Cron(nats, cronExpr)
		},
	}

	cmd.Flags().String("tracer", "jaeger-agent:5775", "endpoint of the jaeger-agent. Set to 'false' to disable tracing")
	cmd.Flags().String("nats", "localhost:4222", "endpoint for NATS server")
	cmd.Flags().String("cron", "0 */5 * * * *", "cron frequency")
	viper.BindPFlag("tracer", cmd.Flags().Lookup("tracer"))

	return cmd
}
