package interconnect

import (
	"fmt"

	"github.com/gogo/protobuf/proto"

	nats "github.com/nats-io/go-nats"
	"github.com/spf13/viper"
	event "github.com/tcfw/evntsrc/internal/event/protos"
	pb "github.com/tcfw/evntsrc/internal/interconnect/protos"
)

type relay struct {
	natsConn *nats.Conn
}

func (s *relay) WritePipe(relayOut chan *event.Event) (chan struct{}, error) {
	writeClose := make(chan struct{})
	natsIn := make(chan *nats.Msg)
	regionName := viper.GetString("region")

	sub, err := s.natsConn.ChanQueueSubscribe("_USER.>", "interconnect", natsIn)
	if err != nil {
		return nil, err
	}

	go func() {
		for {
			select {
			case <-writeClose:
				sub.Unsubscribe()
				close(natsIn)
				return
			case msg := <-natsIn:
				ev := &event.Event{}
				err := proto.Unmarshal(msg.Data, ev)
				if err != nil {
					fmt.Printf("Failed to relay event: %s\n", err.Error())
				}
				//Ignore replays
				if _, ok := ev.Metadata["replay"]; ok {
					continue
				}
				//Ignore alreadyed forwarded
				if _, ok := ev.Metadata[mdForwarded]; ok {
					continue
				}

				ev.Metadata[mdForwarded] = "true"
				ev.Metadata[mdForwardedFrom] = regionName

				relayOut <- ev
			}
		}
	}()

	return writeClose, nil
}

func (s *relay) publishedForwarded(forwardedEventReq *pb.ForwardingRequest) error {
	if _, ok := forwardedEventReq.Event.Metadata[mdForwarded]; !ok {
		return fmt.Errorf("Missing forwarded metadata [%s]", mdForwarded)
	}
	if _, ok := forwardedEventReq.Event.Metadata[mdForwardedFrom]; !ok {
		return fmt.Errorf("Missing forwarded metadata [%s]", mdForwardedFrom)
	}

	channel := fmt.Sprintf("_USER.%d.%s", forwardedEventReq.Event.Stream, forwardedEventReq.Event.Subject)
	eventBytes, err := proto.Marshal(forwardedEventReq.Event)
	if err != nil {
		fmt.Printf("failed to forward event: %s\n", err.Error())
	}

	return s.natsConn.Publish(channel, eventBytes)
}
