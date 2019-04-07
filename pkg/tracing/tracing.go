package tracing

import (
	"context"
	"fmt"
	"log"
	"time"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	ot "github.com/opentracing/opentracing-go"
	viper "github.com/spf13/viper"
	jaeger "github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/zipkin"
	"github.com/uber/jaeger-lib/metrics/prometheus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const (
	prefixTracerState  = "x-b3-"
	zipkinTraceID      = prefixTracerState + "traceid"
	zipkinSpanID       = prefixTracerState + "spanid"
	zipkinParentSpanID = prefixTracerState + "parentspanid"
	zipkinSampled      = prefixTracerState + "sampled"
	zipkinFlags        = prefixTracerState + "flags"
)

var otHeaders = []string{
	zipkinTraceID,
	zipkinSpanID,
	zipkinParentSpanID,
	zipkinSampled,
	zipkinFlags,
}

var (
	defaultEndpoint = "jaeger-agent:5775"
)

//ExtractOTHeadersFromContext TODO
func ExtractOTHeadersFromContext(ctx context.Context) *metadata.MD {
	pairs := []string{}
	for _, h := range otHeaders {
		if v, ok := ctx.Value(h).(string); ok && len(v) > 0 {
			fmt.Printf("Found %v %v", h, v)
			pairs = append(pairs, h, v)
		}
	}

	md := metadata.Pairs(pairs...)

	return &md
}

//InitGlobalTracer applies a jaeger based global tracer
func InitGlobalTracer(name string) {

	tracingEndpoint := viper.GetString("tracer")
	if tracingEndpoint == "false" {
		log.Println("!! Tracing is disabled !!")
		return
	}

	zipkinPropagator := zipkin.NewZipkinB3HTTPHeaderPropagator()
	injector := jaeger.TracerOptions.Injector(ot.HTTPHeaders, zipkinPropagator)
	extractor := jaeger.TracerOptions.Extractor(ot.HTTPHeaders, zipkinPropagator)

	metricsFactory := prometheus.New()
	metricsTags := map[string]string{
		"service": name,
	}

	sender, _ := jaeger.NewUDPTransport(defaultEndpoint, 0)
	tracer, _ := jaeger.NewTracer(
		name,
		jaeger.NewConstSampler(true),
		jaeger.NewRemoteReporter(
			sender,
			jaeger.ReporterOptions.BufferFlushInterval(1*time.Second)),
		injector,
		extractor,
		jaeger.TracerOptions.ZipkinSharedRPCSpan(true),
		jaeger.TracerOptions.Metrics(jaeger.NewMetrics(metricsFactory, metricsTags)),
	)

	ot.SetGlobalTracer(tracer)
}

//GRPCClientOptions applies interceptors for tracing on GRPC clients
func GRPCClientOptions() []grpc.DialOption {
	var opts []grpc.DialOption

	opts = append(opts, grpc.WithInsecure())

	//Tracing
	opts = append(opts, grpc.WithStreamInterceptor(grpc_middleware.ChainStreamClient(
		grpc_prometheus.StreamClientInterceptor,
		grpc_opentracing.StreamClientInterceptor(
			grpc_opentracing.WithTracer(ot.GlobalTracer())))))

	opts = append(opts, grpc.WithUnaryInterceptor(grpc_middleware.ChainUnaryClient(
		grpc_prometheus.UnaryClientInterceptor,
		grpc_opentracing.UnaryClientInterceptor(
			grpc_opentracing.WithTracer(ot.GlobalTracer())))))

	return opts
}

//GRPCServerOptions applies interceptors for tracing on GRPC servers
func GRPCServerOptions() []grpc.ServerOption {
	var opts []grpc.ServerOption

	//Tracing
	opts = append(opts, grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
		grpc_prometheus.StreamServerInterceptor,
		grpc_opentracing.StreamServerInterceptor(
			grpc_opentracing.WithTracer(ot.GlobalTracer())))))
	opts = append(opts, grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
		grpc_prometheus.UnaryServerInterceptor,
		grpc_opentracing.UnaryServerInterceptor(
			grpc_opentracing.WithTracer(ot.GlobalTracer())))))

	return opts
}

//StartSpan starts a new span from a given context
func StartSpan(ctx context.Context, name string) ot.Span {
	span, _ := ot.StartSpanFromContext(ctx, name)
	return span
}

//StartChildSpan is a wrapper method to start a subspan
func StartChildSpan(childSpan ot.Span, name string) ot.Span {
	return ot.StartSpan(name, ot.ChildOf(childSpan.Context()))
}
