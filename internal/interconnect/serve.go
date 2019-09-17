package interconnect

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/spf13/viper"
	pb "github.com/tcfw/evntsrc/internal/interconnect/protos"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	//Compressor encoding between client/server
	_ "google.golang.org/grpc/encoding/gzip"
)

//RunGRPC starts GRPC endpoint
func RunGRPC(port int, nats string) error {
	cert := viper.GetString("certificate")
	key := viper.GetString("key")
	if _, err := os.Stat(cert); os.IsNotExist(err) {
		return fmt.Errorf("Certificate not found: %s", err.Error())
	}
	if _, err := os.Stat(key); os.IsNotExist(err) {
		return fmt.Errorf("Key not found: %s", err.Error())
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	creds, err := credentials.NewServerTLSFromFile(cert, key)
	if err != nil {
		return fmt.Errorf("Failed to load TLS keys: %s", err)
	}

	server := grpc.NewServer(grpc.Creds(creds))

	srv := &srv{}

	if err := srv.Init(nats); err != nil {
		return fmt.Errorf("Failed to init: %s", err.Error())
	}

	pb.RegisterInterconnectServiceServer(server, srv)

	log.Printf("Starting interconnect server (%d): region: %s", port, viper.GetString("region"))

	if err := server.Serve(lis); err != nil {
		return fmt.Errorf("grpc serve error: %s", err)
	}

	return nil
}
