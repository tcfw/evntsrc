package ingress

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	streamauth "github.com/tcfw/evntsrc/pkg/streamauth/protos"
	"google.golang.org/grpc"
)

//ValidateAuth verifies the sent HTTP realm auth
func ValidateAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		stream, err := getStream(r)
		if err != nil {
			http.Error(w, "Unable to determine stream", http.StatusBadRequest)
			return
		}

		key, secret, ok := r.BasicAuth()
		if !ok || key == "" || secret == "" {
			log.Printf("Auth failed: %v, %v, %v", key, secret, ok)
			w.Header().Set("WWW-Authenticate", `Basic realm="API Login Required"`)
			http.Error(w, "Unauthorised", 401)
			return
		}

		conn, err := grpc.Dial("streamauth:443", grpc.WithInsecure())
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Failed to verify auth", http.StatusInternalServerError)
			return
		}
		cli := streamauth.NewStreamAuthServiceClient(conn)

		_, err = cli.ValidateKeySecret(context.Background(), &streamauth.KSRequest{Stream: *stream, Key: key, Secret: secret})
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		conn.Close()

		next.ServeHTTP(w, r)
	})
}

func getStream(r *http.Request) (*int32, error) {
	parts := mux.Vars(r)

	expectedStream := parts["stream"]
	testedStream, err := strconv.Atoi(expectedStream)
	if err != nil {
		return nil, err
	}
	stream := int32(testedStream)

	return &stream, nil
}
