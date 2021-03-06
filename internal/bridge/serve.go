package bridge

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	pb "github.com/tcfw/evntsrc/internal/bridge/protos"
	"github.com/tcfw/evntsrc/internal/tracing"
	"google.golang.org/grpc"
)

//RunGRPC starts the GRPC server
func RunGRPC(port int, natsEndpoint string) {

	connectNats(natsEndpoint)
	defer natsConn.Close()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		panic(err)
	}

	tracing.InitGlobalTracer("Bridge")

	grpcServer := grpc.NewServer(tracing.GRPCServerOptions()...)
	pb.RegisterBridgeServiceServer(grpcServer, newServer())

	log.Println("Starting gRPC server...")
	go func() {
		grpcServer.Serve(lis)
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan

	log.Println("Shutdown signal received, exiting...")
	grpcServer.GracefulStop()
}
