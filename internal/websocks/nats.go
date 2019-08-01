package websocks

import (
	"fmt"
	"os"

	"github.com/gogo/protobuf/proto"
	nats "github.com/nats-io/nats.go"
	"github.com/tcfw/evntsrc/internal/event"
)

var natsConn *nats.Conn

func connectNats(addr string) {
	envHost, exists := os.LookupEnv("NATS_HOST")
	if exists {
		addr = envHost
	}
	nc, err := nats.Connect(addr)
	if err != nil {
		panic(err)
	}

	natsConn = nc
}

//NatsPublisher sends an event into a nats cluster
type NatsPublisher struct{}

//Publish takes in a channel and an event, converts it into protobuf and sends it into a nats cluster
func (n *NatsPublisher) Publish(channel string, event *event.Event) error {
	eventJSONBytes, _ := proto.Marshal(event.ToProtobuf())

	labels := []string{
		fmt.Sprintf("%d", event.Stream),
	}

	if event.Subject == "advertisement" {
		labels = append(labels, "control")
	}

	bytePublishCounter.WithLabelValues(labels...).Add(float64(len(eventJSONBytes)))

	return natsConn.Publish(channel, eventJSONBytes)
}
