package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	pb "github.com/tcfw/evntsrc/internal/bridge/protos"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func main() {
	conn, err := grpc.Dial("localhost:1235", tracing.GRPCClientOptions()...)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	client := pb.NewBridgeServiceClient(conn)

	ctx := context.Background()
	md := metadata.New(map[string]string{"qid": "test"})

	stream, err := client.RelayEvents(metadata.NewOutgoingContext(ctx, md))
	if err != nil {
		panic(err)
	}

	go func() {
		fmt.Println("Opened stream")
		for {
			event, err := stream.Recv()
			if err != nil {
				return
			}
			fmt.Printf("%v\n", event)
			fmt.Printf("Latency: %s\n", time.Since(*event.GetTime()))
		}
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan
}
