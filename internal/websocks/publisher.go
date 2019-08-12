package websocks

import (
	"fmt"

	"github.com/tcfw/evntsrc/internal/event"
)

//NatsPublisher sends an event into a nats cluster
type NatsPublisher struct{}

//Publish takes in a channel and an event, converts it into protobuf and sends it into a nats cluster
func (n *NatsPublisher) Publish(channel string, event *event.Event) error {
	eventBytes, _ := event.ToProtobuf().Marshal()

	labels := []string{
		fmt.Sprintf("%d", event.Stream),
	}

	if event.Subject == "advertisement" {
		labels = append(labels, "control")
	}

	bytePublishCounter.WithLabelValues(labels...).Add(float64(len(eventBytes)))

	return natsConn.Publish(channel, eventBytes)
}
