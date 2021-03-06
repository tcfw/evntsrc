// Code generated by protoc-gen-grpc-gateway. DO NOT EDIT.
// source: metrics.proto

/*
Package evntsrc_stsmetrics is a reverse proxy.

It translates gRPC into RESTful JSON APIs.
*/
package evntsrc_stsmetrics

import (
	"context"
	"io"
	"net/http"

	"github.com/golang/protobuf/proto"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/grpc-ecosystem/grpc-gateway/utilities"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/status"
)

var _ codes.Code
var _ io.Reader
var _ status.Status
var _ = runtime.String
var _ = utilities.NewDoubleArray

func request_MetricsService_EventsCount_0(ctx context.Context, marshaler runtime.Marshaler, client MetricsServiceClient, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var protoReq MetricsRequest
	var metadata runtime.ServerMetadata

	var (
		val string
		e   int32
		ok  bool
		err error
		_   = err
	)

	val, ok = pathParams["stream"]
	if !ok {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "missing parameter %s", "stream")
	}

	protoReq.Stream, err = runtime.Int32(val)

	if err != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "type mismatch, parameter: %s, error: %v", "stream", err)
	}

	val, ok = pathParams["interval"]
	if !ok {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "missing parameter %s", "interval")
	}

	e, err = runtime.Enum(val, Interval_value)

	if err != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "type mismatch, parameter: %s, error: %v", "interval", err)
	}

	protoReq.Interval = Interval(e)

	msg, err := client.EventsCount(ctx, &protoReq, grpc.Header(&metadata.HeaderMD), grpc.Trailer(&metadata.TrailerMD))
	return msg, metadata, err

}

func request_MetricsService_EventsSize_0(ctx context.Context, marshaler runtime.Marshaler, client MetricsServiceClient, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var protoReq MetricsRequest
	var metadata runtime.ServerMetadata

	var (
		val string
		e   int32
		ok  bool
		err error
		_   = err
	)

	val, ok = pathParams["stream"]
	if !ok {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "missing parameter %s", "stream")
	}

	protoReq.Stream, err = runtime.Int32(val)

	if err != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "type mismatch, parameter: %s, error: %v", "stream", err)
	}

	val, ok = pathParams["interval"]
	if !ok {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "missing parameter %s", "interval")
	}

	e, err = runtime.Enum(val, Interval_value)

	if err != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "type mismatch, parameter: %s, error: %v", "interval", err)
	}

	protoReq.Interval = Interval(e)

	msg, err := client.EventsSize(ctx, &protoReq, grpc.Header(&metadata.HeaderMD), grpc.Trailer(&metadata.TrailerMD))
	return msg, metadata, err

}

func request_MetricsService_Connections_0(ctx context.Context, marshaler runtime.Marshaler, client MetricsServiceClient, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var protoReq MetricsRequest
	var metadata runtime.ServerMetadata

	var (
		val string
		e   int32
		ok  bool
		err error
		_   = err
	)

	val, ok = pathParams["stream"]
	if !ok {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "missing parameter %s", "stream")
	}

	protoReq.Stream, err = runtime.Int32(val)

	if err != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "type mismatch, parameter: %s, error: %v", "stream", err)
	}

	val, ok = pathParams["interval"]
	if !ok {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "missing parameter %s", "interval")
	}

	e, err = runtime.Enum(val, Interval_value)

	if err != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "type mismatch, parameter: %s, error: %v", "interval", err)
	}

	protoReq.Interval = Interval(e)

	msg, err := client.Connections(ctx, &protoReq, grpc.Header(&metadata.HeaderMD), grpc.Trailer(&metadata.TrailerMD))
	return msg, metadata, err

}

// RegisterMetricsServiceHandlerFromEndpoint is same as RegisterMetricsServiceHandler but
// automatically dials to "endpoint" and closes the connection when "ctx" gets done.
func RegisterMetricsServiceHandlerFromEndpoint(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error) {
	conn, err := grpc.Dial(endpoint, opts...)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			if cerr := conn.Close(); cerr != nil {
				grpclog.Infof("Failed to close conn to %s: %v", endpoint, cerr)
			}
			return
		}
		go func() {
			<-ctx.Done()
			if cerr := conn.Close(); cerr != nil {
				grpclog.Infof("Failed to close conn to %s: %v", endpoint, cerr)
			}
		}()
	}()

	return RegisterMetricsServiceHandler(ctx, mux, conn)
}

// RegisterMetricsServiceHandler registers the http handlers for service MetricsService to "mux".
// The handlers forward requests to the grpc endpoint over "conn".
func RegisterMetricsServiceHandler(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error {
	return RegisterMetricsServiceHandlerClient(ctx, mux, NewMetricsServiceClient(conn))
}

// RegisterMetricsServiceHandlerClient registers the http handlers for service MetricsService
// to "mux". The handlers forward requests to the grpc endpoint over the given implementation of "MetricsServiceClient".
// Note: the gRPC framework executes interceptors within the gRPC handler. If the passed in "MetricsServiceClient"
// doesn't go through the normal gRPC flow (creating a gRPC client etc.) then it will be up to the passed in
// "MetricsServiceClient" to call the correct interceptors.
func RegisterMetricsServiceHandlerClient(ctx context.Context, mux *runtime.ServeMux, client MetricsServiceClient) error {

	mux.Handle("GET", pattern_MetricsService_EventsCount_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		rctx, err := runtime.AnnotateContext(ctx, mux, req)
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := request_MetricsService_EventsCount_0(rctx, inboundMarshaler, client, req, pathParams)
		ctx = runtime.NewServerMetadataContext(ctx, md)
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}

		forward_MetricsService_EventsCount_0(ctx, mux, outboundMarshaler, w, req, resp, mux.GetForwardResponseOptions()...)

	})

	mux.Handle("GET", pattern_MetricsService_EventsSize_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		rctx, err := runtime.AnnotateContext(ctx, mux, req)
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := request_MetricsService_EventsSize_0(rctx, inboundMarshaler, client, req, pathParams)
		ctx = runtime.NewServerMetadataContext(ctx, md)
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}

		forward_MetricsService_EventsSize_0(ctx, mux, outboundMarshaler, w, req, resp, mux.GetForwardResponseOptions()...)

	})

	mux.Handle("GET", pattern_MetricsService_Connections_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		rctx, err := runtime.AnnotateContext(ctx, mux, req)
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := request_MetricsService_Connections_0(rctx, inboundMarshaler, client, req, pathParams)
		ctx = runtime.NewServerMetadataContext(ctx, md)
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}

		forward_MetricsService_Connections_0(ctx, mux, outboundMarshaler, w, req, resp, mux.GetForwardResponseOptions()...)

	})

	return nil
}

var (
	pattern_MetricsService_EventsCount_0 = runtime.MustPattern(runtime.NewPattern(1, []int{2, 0, 2, 1, 1, 0, 4, 1, 5, 1, 2, 2, 2, 3, 1, 0, 4, 1, 5, 4}, []string{"v1", "stream", "metrics", "events", "interval"}, "", runtime.AssumeColonVerbOpt(true)))

	pattern_MetricsService_EventsSize_0 = runtime.MustPattern(runtime.NewPattern(1, []int{2, 0, 2, 1, 1, 0, 4, 1, 5, 1, 2, 2, 2, 3, 1, 0, 4, 1, 5, 4}, []string{"v1", "stream", "metrics", "storage", "interval"}, "", runtime.AssumeColonVerbOpt(true)))

	pattern_MetricsService_Connections_0 = runtime.MustPattern(runtime.NewPattern(1, []int{2, 0, 2, 1, 1, 0, 4, 1, 5, 1, 2, 2, 2, 3, 1, 0, 4, 1, 5, 4}, []string{"v1", "stream", "metrics", "conns", "interval"}, "", runtime.AssumeColonVerbOpt(true)))
)

var (
	forward_MetricsService_EventsCount_0 = runtime.ForwardResponseMessage

	forward_MetricsService_EventsSize_0 = runtime.ForwardResponseMessage

	forward_MetricsService_Connections_0 = runtime.ForwardResponseMessage
)
