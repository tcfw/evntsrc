package cmd

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tcfw/evntsrc/internal/interconnect"
)

//NewConnectCmd provides a connect command
func NewConnectCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "connect",
		Short: "Connect to a region",
		Run: func(cmd *cobra.Command, args []string) {
			viper.BindPFlag("certificate", cmd.Flags().Lookup("certificate"))
			viper.BindPFlag("region", cmd.Flags().Lookup("region"))

			endpoint, _ := cmd.Flags().GetString("remote")
			nats, _ := cmd.Flags().GetString("nats")
			if err := interconnect.Connect(endpoint, nats); err != nil {
				log.Fatal(err)
			}
		},
	}

	cmd.Flags().StringP("remote", "r", "", "Endpoint for remote region")
	cmd.Flags().String("nats", "localhost:4222", "endpoint for NATS server")
	cmd.Flags().StringP("certificate", "c", "", "Certificate for endpoint")
	cmd.Flags().String("region", "", "Region/Zone identifier")

	cmd.MarkFlagRequired("certificate")
	cmd.MarkFlagRequired("remote")
	cmd.MarkFlagRequired("region")

	return cmd
}
