package websocks

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/tcfw/evntsrc/pkg/event"
)

func (c *Client) publishAdvert(eventType string, data []byte) {
	rEvent := &event.Event{}
	rEvent.SetID()
	rEvent.Stream = c.auth.Stream
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

	channel := fmt.Sprintf("_CONN.%d.%s", c.auth.Stream, rEvent.Subject)

	eventJSONBytes, _ := json.Marshal(rEvent)

	natsConn.Publish(channel, eventJSONBytes)
}

func (c *Client) advertiseConnect() {
	c.publishAdvert("connect", nil)
}

func (c *Client) advertiseDisconnect() {
	c.publishAdvert("disconnect", nil)
}

func (c *Client) advertiseSub(subject string) {
	c.publishAdvert("subscribe", []byte(subject))
}

func (c *Client) advertiseUnsub(subject string) {
	c.publishAdvert("unsubscribe", []byte(subject))
}
