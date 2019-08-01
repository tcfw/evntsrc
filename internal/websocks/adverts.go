package websocks

import (
	"fmt"
	"os"

	"github.com/tcfw/evntsrc/internal/event"
)

func (c *Client) publishBroadcast(eventType string, data []byte) {
	rEvent := event.NewEvent()
	rEvent.Stream = c.authKey.Stream
	rEvent.Source = "ws"
	rEvent.Subject = "advertisement"
	rEvent.Type = eventType
	rEvent.Data = data
	rEvent.Metadata = map[string]string{
		"source_ip":     c.conn.RemoteAddr().String(),
		"connection_id": c.connectionID,
	}
	hostname, err := os.Hostname()
	if err == nil {
		rEvent.Metadata["ev_host"] = hostname
	}

	channel := fmt.Sprintf("_STREAM.%d.%s", c.authKey.Stream, rEvent.Subject)

	c.publisher.Publish(channel, rEvent)
}

func (c *Client) broadcastConnect() {
	c.publishBroadcast("connect", nil)
}

func (c *Client) broadcastDisconnect() {
	c.publishBroadcast("disconnect", nil)
}

func (c *Client) broadcastSub(subject string) {
	c.publishBroadcast("subscribe", []byte(subject))
}

func (c *Client) broadcastUnsub(subject string) {
	c.publishBroadcast("unsubscribe", []byte(subject))
}
