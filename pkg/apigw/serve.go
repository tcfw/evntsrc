package apigw

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/go-http-utils/logger"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/rs/cors"
	"github.com/tcfw/evntsrc/pkg/tracing"
)

//Run starts the JSON gw
func Run(port int) error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	tracing.InitGlobalTracer("APIGW")

	go startMetrics()

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

	handler := tracingWrapper(mux)
	handler = metricsMiddleware(handler)
	handler = authGuard(handler)
	handler = logger.Handler(handler, os.Stdout, logger.CommonLoggerType)
	handler = cors.AllowAll().Handler(handler)

	fmt.Printf("Starting API GW (port %d)\n", port)
	return http.ListenAndServe(fmt.Sprintf(":%d", port), handler)
}
