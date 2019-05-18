package interconnect

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"google.golang.org/grpc/credentials"

	nats "github.com/nats-io/go-nats"
	"github.com/spf13/viper"
	event "github.com/tcfw/evntsrc/internal/event/protos"
	pb "github.com/tcfw/evntsrc/internal/interconnect/protos"
	"google.golang.org/grpc"
)

//Connect connects to a remote endpoint
func Connect(endpoint string, nats string) error {
	cert := viper.GetString("certificate")
	if _, err := os.Stat(cert); os.IsNotExist(err) {
		return fmt.Errorf("Certificate not found: %s", err.Error())
	}

	creds, err := credentials.NewClientTLSFromFile(cert, "")
	if err != nil {
		return err
	}

	log.Printf("Starting interconnect client: region: %s: remote: %s", viper.GetString("region"), endpoint)

	conn, err := grpc.Dial(endpoint, grpc.WithTransportCredentials(creds))
	if err != nil {
		return fmt.Errorf("Failed to connect to remote: %s", err.Error())
	}

	client := pb.NewInterconnectServiceClient(conn)
	stream, err := client.Relay(context.Background())
	if err != nil {
		return fmt.Errorf("Failed to start relay: %s", err.Error())
	}

	return startRelay(stream, nats)
}

func startRelay(stream pb.InterconnectService_RelayClient, natsEndpoint string) error {
	log.Println("Opening relay [LOCAL]...")
	toRelay := make(chan *event.Event, 1000)
	closed := false

	nc, err := nats.Connect(natsEndpoint)
	if err != nil {
		return fmt.Errorf("Failed to connect to NATS: %s", err.Error())
	}

	relay := &relay{natsConn: nc}

	writeClose, err := relay.WritePipe(toRelay)
	if err != nil {
		return err
	}

	go func() {
		for {
			ev, ok := <-toRelay
			if !ok {
				closed = true
				break
			}
			stream.Send(&pb.ForwardingRequest{Event: ev})
		}
	}()

	for {
		if closed {
			stream.CloseSend()
			return nil
		}

		forwardedEventReq, err := stream.Recv()
		if err == io.EOF {
			log.Println("Closing relay [LOCAL]...")
			stream.CloseSend()
			return nil
		}
		if err != nil {
			close(writeClose)
			return err
		}

		if err := relay.publishedForwarded(forwardedEventReq); err != nil {
			close(writeClose)
			return err
		}
	}

}
