package bridge

import (
	"fmt"
	"log"
	"net"

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

	grpcServer := grpc.NewServer() //tracing.GRPCServerOptions()...)
	pb.RegisterBridgeServiceServer(grpcServer, newServer())

	log.Println("Starting gRPC server...")
	grpcServer.Serve(lis)
}
