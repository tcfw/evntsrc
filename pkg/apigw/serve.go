package apigw

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/go-http-utils/logger"
	protoutil "github.com/gogo/gateway"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/rs/cors"
	"google.golang.org/grpc"
)

//Run starts the JSON gw
func Run(port int) error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	runtime.HTTPError = CustomHTTPError

	defaultMarshaler := runtime.WithMarshalerOption(runtime.MIMEWildcard, &protoutil.JSONPb{OrigName: true})
	protobufMarshaler := runtime.WithMarshalerOption("application/protobuf", &runtime.ProtoMarshaller{})

	mux := runtime.NewServeMux(defaultMarshaler, protobufMarshaler)
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
	}

	registerStreams(ctx, mux, opts)

	handler := logger.Handler(mux, os.Stdout, logger.DevLoggerType)
	handler = tracingWrapper(handler)
	handler = cors.Default().Handler(handler)

	fmt.Printf("Starting API GW (port %d)\n", port)
	return http.ListenAndServe(fmt.Sprintf(":%d", port), handler)
}
