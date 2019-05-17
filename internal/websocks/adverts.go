package websocks

import (
	"fmt"
	"os"
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/tcfw/evntsrc/internal/event"
)

func (c *Client) publishBroadcast(eventType string, data []byte) {
	rEvent := &event.Event{}
	rEvent.SetID()
	rEvent.Stream = c.authKey.Stream
	rEvent.Source = "ws"
	rEvent.Subject = "advertisement"
	rEvent.CEVersion = "0.1"
	rEvent.Type = eventType
	rEvent.TypeVersion = "0.1"
	rEvent.ContentType = ""
	rEvent.Time = event.ZeroableTime{Time: time.Now()}
	rEvent.Data = data
	rEvent.Metadata = map[string]string{}
	rEvent.Metadata["source_ip"] = c.conn.RemoteAddr().String()
	rEvent.Metadata["connection_id"] = c.connectionID
	hostname, err := os.Hostname()
	if err == nil {
		rEvent.Metadata["host"] = hostname
	}

	channel := fmt.Sprintf("_STREAM.%d.%s", c.authKey.Stream, rEvent.Subject)

	eventJSONBytes, _ := proto.Marshal(rEvent.ToProtobuf())
	natsConn.Publish(channel, eventJSONBytes)
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
