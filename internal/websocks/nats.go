package websocks

import (
	"os"

	nats "github.com/nats-io/nats.go"
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
