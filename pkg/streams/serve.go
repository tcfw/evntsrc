package streams

import (
	"fmt"
	"log"
	"net"

	pb "github.com/tcfw/evntsrc/pkg/streams/protos"
	"google.golang.org/grpc"
)

//RunGRPC starts the GRPC server
func RunGRPC(port int) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterStreamsServiceServer(grpcServer, newServer())

	log.Println("Starting gRPC server")
	grpcServer.Serve(lis)
}
