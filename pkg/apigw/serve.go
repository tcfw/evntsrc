package apigw

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/go-http-utils/logger"
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

	// defaultMarshaler := runtime.WithMarshalerOption(runtime.MIMEWildcard, &protoutil.JSONPb{OrigName: true})
	protobufMarshaler := runtime.WithMarshalerOption("application/protobuf", &runtime.ProtoMarshaller{})

	mux := runtime.NewServeMux(protobufMarshaler)
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
	}

	registerPassport(ctx, mux, opts)
	registerUsers(ctx, mux, opts)
	registerStreams(ctx, mux, opts)
	registerStreamAuth(ctx, mux, opts)

	handler := tracingWrapper(mux)
	handler = authGuard(handler)
	handler = logger.Handler(handler, os.Stdout, logger.CommonLoggerType)
	handler = cors.AllowAll().Handler(handler)

	fmt.Printf("Starting API GW (port %d)\n", port)
	return http.ListenAndServe(fmt.Sprintf(":%d", port), handler)
}
