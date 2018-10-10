package ingress

import "net/http"

//HandlePub publishes events to NATS
func HandlePub(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK!"))
}
