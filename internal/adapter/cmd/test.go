package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
	adapter "github.com/tcfw/evntsrc/internal/adapter"
	adapterPb "github.com/tcfw/evntsrc/internal/adapter/protos"
	"github.com/tcfw/evntsrc/internal/event/protos"
)

//NewTestCmd executes a test V8 adapter
func NewTestCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "test",
		Short: "Start a test execute",
		Run: func(cmd *cobra.Command, args []string) {
			adapterInst := &adapterPb.Adapter{
				Engine: adapterPb.Adapter_V8,
				Code: []byte(`
				log("hi");
				`),
			}

			event := &evntsrc_event.Event{}

			server := adapter.NewServer()
			server.StartPools()

			start := time.Now()
			resp, err := server.Execute(context.Background(), &adapterPb.ExecuteRequest{Adapter: adapterInst, Event: event})
			if err == nil {
				fmt.Printf("Log: %v\n", resp.Log)
				fmt.Printf("Took: %s\n", time.Since(start).String())
			} else {
				fmt.Printf("FATAL: %v\n", err.Error())
			}

			start = time.Now()
			resp, err = server.Execute(context.Background(), &adapterPb.ExecuteRequest{Adapter: adapterInst, Event: event})
			if err == nil {
				fmt.Printf("Log: %v\n", resp.Log)
				fmt.Printf("Took: %s\n", time.Since(start).String())
			} else {
				fmt.Printf("FATAL: %v\n", err.Error())
			}

			server.StopPools()

		},
	}

	return cmd
}
