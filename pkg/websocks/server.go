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
func Run(webPort int, natsEndpoint string) {
	connectNats(natsEndpoint)
	defer natsConn.Close()

	startHTTPServer(webPort)
}
