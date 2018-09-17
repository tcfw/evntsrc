package stsmetrics

import (
	"github.com/nats-io/go-nats"
)

var natsConn *nats.Conn

func connectNats(addr string) {
	nc, err := nats.Connect(addr)
	if err != nil {
		panic(err)
	}

	natsConn = nc
}
