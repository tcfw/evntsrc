package websocks

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/websocket"
	nats "github.com/nats-io/go-nats"
)

func BenchmarkConnSub(b *testing.B) {
	b.StopTimer()
	s := httptest.NewServer(http.HandlerFunc(echo))
	defer s.Close()
	// Convert http://127.0.0.1 to ws://127.0.0.
	u := "ws" + strings.TrimPrefix(s.URL, "http")

	natsConn, _ = nats.Connect(nats.DefaultURL)

	b.StartTimer()
	for n := 0; n < b.N; n++ {
		// Connect to the server
		ws, _, err := websocket.DefaultDialer.Dial(u, nil)
		if err != nil {
			return
		}
		c := NewClient(ws)
		c.conn.Close()
		c.closeConnSub <- true
		if err := c.ConnSub(); err != nil {
			return
		}
		ws.Close()
	}
}

func echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer c.Close()
	for {
		_, _, err := c.ReadMessage()
		if err != nil {
			break
		}

	}
}
