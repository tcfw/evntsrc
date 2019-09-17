package interconnect

import (
	"io"
	"log"

	nats "github.com/nats-io/nats.go"
	event "github.com/tcfw/evntsrc/internal/event/protos"
	pb "github.com/tcfw/evntsrc/internal/interconnect/protos"
)

type srv struct {
	natsConn *nats.Conn
}

const (
	mdForwarded     string = "forwarded"
	mdForwardedFrom string = "forwarded_from"
)

//Replay opens bi-direction NAT streaming
func (s *srv) Relay(stream pb.InterconnectService_RelayServer) error {
	log.Println("Opening relay [REMOTE]...")

	toRelay := make(chan *event.Event, 1024)
	closed := false

	relay := &relay{natsConn: s.natsConn}

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
			return nil
		}

		forwardedEventReq, err := stream.Recv()
		if err != nil {
			close(writeClose)
			log.Println("Closing relay [REMOTE]...")

			if err == io.EOF {
				return nil
			}

			return err
		}

		if err := relay.publishedForwarded(forwardedEventReq); err != nil {
			close(writeClose)
			return err
		}
	}
}

//Init attempts to make a connection with the NATS server
func (s *srv) Init(natsEndpoint string) error {
	nc, err := nats.Connect(natsEndpoint)
	if err != nil {
		return err
	}

	s.natsConn = nc

	return nil
}
