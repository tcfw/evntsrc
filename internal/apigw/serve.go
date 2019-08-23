package apigw

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/rs/cors"
	"github.com/tcfw/evntsrc/internal/tracing"
	"go.uber.org/zap"
)

//Run starts the JSON gw
func Run(port int) error {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}
	defer logger.Sync()
	zap.ReplaceGlobals(logger)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tracing.InitGlobalTracer("APIGW")

	runtime.HTTPError = CustomHTTPError

	// defaultMarshaler := runtime.WithMarshalerOption(runtime.MIMEWildcard, &protoutil.JSONPb{OrigName: true})
	protobufMarshaler := runtime.WithMarshalerOption("application/protobuf", &runtime.ProtoMarshaller{})

	mux := runtime.NewServeMux(protobufMarshaler)
	opts := tracing.GRPCClientOptions()

	registerPassport(ctx, mux, opts)
	registerUsers(ctx, mux, opts)
	registerStreams(ctx, mux, opts)
	registerStreamAuth(ctx, mux, opts)
	registerBilling(ctx, mux, opts)
	registerMetrics(ctx, mux, opts)

	handler := authGuard(mux)
	handler = tracingWrapper(handler)
	handler = loggingWrapper(handler)
	handler = cors.AllowAll().Handler(handler)

	log.Printf("Starting API GW (port %d)\n", port)
	return http.ListenAndServe(fmt.Sprintf(":%d", port), handler)
}
