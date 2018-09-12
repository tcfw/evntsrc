package websocks

import (
	"fmt"
	"log"
	"net/http"
)

//startHTTPServer provides a HTTP server for / and websockets
func startHTTPServer(port int) {
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(w, r)
	})

	addr := fmt.Sprintf(":%d", port)

	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

//Run connects to NAT and start web server
func Run(webPort int) {
	connectNats("127.0.0.1:30858")
	defer func() {
		natsConn.Close()
	}()

	startHTTPServer(webPort)
}
