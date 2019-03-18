package websocks

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

//startHTTPServer provides a HTTP server for / and websockets
func startHTTPServer(port int) {
	mux := mux.NewRouter()

	mux.HandleFunc("/v1/", func(w http.ResponseWriter, r *http.Request) {
		serveWs(w, r)
	})

	mux.HandleFunc("/v1/{stream:[0-9]+}", func(w http.ResponseWriter, r *http.Request) {
		serveWs(w, r)
	})

	mux.Use(metricsMiddleware)

	addr := fmt.Sprintf(":%d", port)

	log.Println("Starting HTTP server...")
	err := http.ListenAndServe(addr, mux)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

//Run connects to NAT and start web server
func Run(webPort int, natsEndpoint string) {
	connectNats(natsEndpoint)
	defer natsConn.Close()

	startMetrics()

	startHTTPServer(webPort)
}
