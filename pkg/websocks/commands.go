package websocks

import "github.com/simplereach/timeutils"

const (
	commandSubscribe   = "sub"
	commandUnsubscribe = "unsub"
	commandPublish     = "pub"
	commandAuth        = "auth"
	commandReplay      = "replay"
)

//InboundCommand is the basic struct for all commands coming from browser
type InboundCommand struct {
	Command string `json:"cmd"`
}

//AuthCommand provides streaming information and verification
type AuthCommand struct {
	*InboundCommand
	Stream int32  `json:"stream"`
	Key    string `json:"key"`
	Secret string `json:"secret"`
}

//SubscribeCommand starts a subscription
type SubscribeCommand struct {
	*InboundCommand
	Subject string `json:"subject"`
}

//PublishCommand sends data through to NATS
type PublishCommand struct {
	*SubscribeCommand
	Data        string `json:"data"`
	Source      string `json:"source"`
	Type        string `json:"type"`
	TypeVersion string `json:"typeVersion"`
	ContentType string `json:"contentType"`
}

//UnsubscribeCommand starts a subscription
type UnsubscribeCommand struct {
	*InboundCommand
	Subject string `json:"subject"`
}

//AckSubUnSucCommand provides error responses to WS clients
type AckSubUnSucCommand struct {
	Acktype string `json:"acktype"`
	Channel string `json:"string"`
	Error   string `json:"error,omitempty"`
}

//ReplayCommand instructs events to rebroadcast all events stored since time
type ReplayCommand struct {
	*SubscribeCommand
	Time   timeutils.Time `json:"startTime"`
	Stream int32          `json:"stream"`
	JustMe bool           `json:"justme"`
}
