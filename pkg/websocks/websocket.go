package websocks

import (
	"log"
	"net/http"

	"github.com/globalsign/mgo/bson"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true }, //TODO verify origins?
}

func serveWs(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := &Client{
		conn:          conn,
		send:          make(chan []byte, 256),
		subscriptions: map[string]chan bool{},
		connectionID:  bson.NewObjectId().Hex(),
		seq:           map[string]int64{},
		closed:        false,
	}

	go client.writePump()
	go client.readPump()
}
