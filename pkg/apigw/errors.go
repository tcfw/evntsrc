package apigw

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
)

type errorBody struct {
	Err string `json:"error,omitempty"`
}

//CustomHTTPError formats errors in a slightly nicer way
func CustomHTTPError(ctx context.Context, _ *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, _ *http.Request, err error) {
	log.Printf("Attempting to format error")
	const fallback = `{"error": "failed to marshal error message"}`

	w.Header().Set("Content-type", marshaler.ContentType())
	w.WriteHeader(runtime.HTTPStatusFromCode(grpc.Code(err)))
	jErr := json.NewEncoder(w).Encode(errorBody{
		Err: grpc.ErrorDesc(err),
	})

	if jErr != nil {
		w.Write([]byte(fallback))
	}
}
