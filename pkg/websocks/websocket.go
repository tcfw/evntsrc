package websocks

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	metrics "github.com/rcrowley/go-metrics"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:    1024,
	WriteBufferSize:   1024,
	CheckOrigin:       func(r *http.Request) bool { return true }, //TODO verify origins?
	EnableCompression: true,
}

func serveWs(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := NewClient(conn)

	m := metrics.GetOrRegisterCounter("wsConnections", nil)
	m.Inc(1)

	go client.writePump()
	go client.readPump()
}
