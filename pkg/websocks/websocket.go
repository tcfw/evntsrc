package websocks

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
	metrics "github.com/rcrowley/go-metrics"
	"github.com/tcfw/evntsrc/pkg/tracing"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true }, //TODO verify origins?
	// EnableCompression: true,
}

func serveWs(w http.ResponseWriter, r *http.Request) {
	span := tracing.GlobalTracer().StartSpan("serveWs")

	var apiKey string
	var apiSec string
	useAuthHeader := false

	if authHeader := r.Header.Get("Authorization"); authHeader != "" {
		fmt.Println("Attempting basic auth login")
		u, p, ok := r.BasicAuth()
		if !ok {
			fmt.Println("Failed to obtain basic auth")
			http.Error(w, "Invalid auth", http.StatusForbidden)
			return
		}
		apiKey = u
		apiSec = p
		useAuthHeader = true
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := NewClient(conn)

	if useAuthHeader {
		streamint, _ := strconv.ParseInt(r.URL.Path[len("/v1/"):], 10, 64)
		fmt.Printf("Attempting to use auth (stream: %v)\n", streamint)
		err := client.authFromHeader(apiKey, apiSec, int32(streamint))
		if err != nil {
			fmt.Println("Attempting to use auth: failed. Closing connection")
			conn.WriteControl(websocket.CloseMessage, []byte("Auth Failed"), time.Now().Add(5*time.Second))
			conn.Close()
			return
		}
		client.sendStruct(&ConnectionInfo{
			Ref:          "conn",
			ConnectionID: client.connectionID,
		})
		go client.broadcastConnect()
	}

	m := metrics.GetOrRegisterCounter("wsConnections", nil)
	m.Inc(1)

	span.Finish()

	go client.writePump()
	go client.readPump()
}
