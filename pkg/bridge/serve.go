package bridge

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	pb "github.com/tcfw/evntsrc/pkg/bridge/protos"
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

	grpcServer := grpc.NewServer()
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
