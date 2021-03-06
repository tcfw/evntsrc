package ingress

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	streamauth "github.com/tcfw/evntsrc/internal/streamauth/protos"
	"github.com/tcfw/evntsrc/internal/tracing"
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

		conn, err := grpc.Dial("streamauth:443", tracing.GRPCClientOptions()...)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Failed to verify auth", http.StatusInternalServerError)
			return
		}
		cli := streamauth.NewStreamAuthServiceClient(conn)

		if _, err = cli.ValidateKeySecret(context.Background(), &streamauth.KSRequest{Stream: *stream, Key: key, Secret: secret}); err != nil {
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

	expectedStream, ok := parts["stream"]
	if !ok {
		return nil, errors.New("Stream not found in request")
	}

	testedStream, err := strconv.Atoi(expectedStream)
	if err != nil {
		return nil, err
	}
	stream := int32(testedStream)

	return &stream, nil
}
