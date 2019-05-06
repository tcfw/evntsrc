package emails

import (
	"fmt"
	"log"
	"net"

	pb "github.com/tcfw/evntsrc/internal/emails/protos"
	"github.com/tcfw/evntsrc/internal/tracing"
	"google.golang.org/grpc"
)

//RunGRPC starts main server & work queue
func RunGRPC(port int) {
	go startWorker()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	tracing.InitGlobalTracer("Emails")

	grpcServer := grpc.NewServer(tracing.GRPCServerOptions()...)
	pb.RegisterEmailServiceServer(grpcServer, NewServer())

	log.Printf("Starting gRPC server (port %d)\n", port)
	grpcServer.Serve(lis)
}
