package emails

import (
	"fmt"
	"log"
	"net"

	pb "github.com/tcfw/evntsrc/pkg/emails/protos"
	"google.golang.org/grpc"
)

func RunGRPC(port int) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterEmailServiceServer(grpcServer, NewServer())

	log.Printf("Starting gRPC server (port %d)\n", port)
	grpcServer.Serve(lis)
}
