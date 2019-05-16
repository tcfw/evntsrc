package cmd

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tcfw/evntsrc/internal/interconnect"
)

//NewServeCmd provides a serve command
func NewServeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Start the service",
		Run: func(cmd *cobra.Command, args []string) {
			viper.BindPFlag("certificate", cmd.Flags().Lookup("certificate"))
			viper.BindPFlag("key", cmd.Flags().Lookup("key"))
			viper.BindPFlag("region", cmd.Flags().Lookup("region"))

			port, _ := cmd.Flags().GetInt("port")
			nats, _ := cmd.Flags().GetString("nats")
			if err := interconnect.RunGRPC(port, nats); err != nil {
				log.Fatal(err)
			}
		},
	}

	cmd.Flags().IntP("port", "p", 443, "listening port for GRPC")
	cmd.Flags().String("nats", "localhost:4222", "endpoint for NATS server")
	cmd.Flags().StringP("certificate", "c", "", "Certificate for endpoint")
	cmd.Flags().StringP("key", "k", "", "TLS key for endpoint")
	cmd.Flags().String("region", "", "Region/Zone identifier")

	cmd.MarkFlagRequired("certificate")
	cmd.MarkFlagRequired("key")
	cmd.MarkFlagRequired("region")

	return cmd
}
