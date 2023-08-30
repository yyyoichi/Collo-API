// Code generated by protoc-gen-connect-go. DO NOT EDIT.
//
// Source: proto/collo/v1/collo.proto

package collov1connect

import (
	connect "connectrpc.com/connect"
	context "context"
	errors "errors"
	http "net/http"
	strings "strings"
	v1 "yyyoichi/Collo-API/gen/proto/collo/v1"
)

// This is a compile-time assertion to ensure that this generated file and the connect package are
// compatible. If you get a compiler error that this constant is not defined, this code was
// generated with a version of connect newer than the one compiled into your binary. You can fix the
// problem by either regenerating this code with an older version of connect or updating the connect
// version compiled into your binary.
const _ = connect.IsAtLeastVersion0_1_0

const (
	// ColloServiceName is the fully-qualified name of the ColloService service.
	ColloServiceName = "collo.v1.ColloService"
)

// These constants are the fully-qualified names of the RPCs defined in this package. They're
// exposed at runtime as Spec.Procedure and as the final two segments of the HTTP route.
//
// Note that these are different from the fully-qualified method names used by
// google.golang.org/protobuf/reflect/protoreflect. To convert from these constants to
// reflection-formatted method names, remove the leading slash and convert the remaining slash to a
// period.
const (
	// ColloServiceColloStreamProcedure is the fully-qualified name of the ColloService's ColloStream
	// RPC.
	ColloServiceColloStreamProcedure = "/collo.v1.ColloService/ColloStream"
)

// ColloServiceClient is a client for the collo.v1.ColloService service.
type ColloServiceClient interface {
	ColloStream(context.Context, *connect.Request[v1.ColloRequest]) (*connect.ServerStreamForClient[v1.ColloStreamResponse], error)
}

// NewColloServiceClient constructs a client for the collo.v1.ColloService service. By default, it
// uses the Connect protocol with the binary Protobuf Codec, asks for gzipped responses, and sends
// uncompressed requests. To use the gRPC or gRPC-Web protocols, supply the connect.WithGRPC() or
// connect.WithGRPCWeb() options.
//
// The URL supplied here should be the base URL for the Connect or gRPC server (for example,
// http://api.acme.com or https://acme.com/grpc).
func NewColloServiceClient(httpClient connect.HTTPClient, baseURL string, opts ...connect.ClientOption) ColloServiceClient {
	baseURL = strings.TrimRight(baseURL, "/")
	return &colloServiceClient{
		colloStream: connect.NewClient[v1.ColloRequest, v1.ColloStreamResponse](
			httpClient,
			baseURL+ColloServiceColloStreamProcedure,
			opts...,
		),
	}
}

// colloServiceClient implements ColloServiceClient.
type colloServiceClient struct {
	colloStream *connect.Client[v1.ColloRequest, v1.ColloStreamResponse]
}

// ColloStream calls collo.v1.ColloService.ColloStream.
func (c *colloServiceClient) ColloStream(ctx context.Context, req *connect.Request[v1.ColloRequest]) (*connect.ServerStreamForClient[v1.ColloStreamResponse], error) {
	return c.colloStream.CallServerStream(ctx, req)
}

// ColloServiceHandler is an implementation of the collo.v1.ColloService service.
type ColloServiceHandler interface {
	ColloStream(context.Context, *connect.Request[v1.ColloRequest], *connect.ServerStream[v1.ColloStreamResponse]) error
}

// NewColloServiceHandler builds an HTTP handler from the service implementation. It returns the
// path on which to mount the handler and the handler itself.
//
// By default, handlers support the Connect, gRPC, and gRPC-Web protocols with the binary Protobuf
// and JSON codecs. They also support gzip compression.
func NewColloServiceHandler(svc ColloServiceHandler, opts ...connect.HandlerOption) (string, http.Handler) {
	colloServiceColloStreamHandler := connect.NewServerStreamHandler(
		ColloServiceColloStreamProcedure,
		svc.ColloStream,
		opts...,
	)
	return "/collo.v1.ColloService/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case ColloServiceColloStreamProcedure:
			colloServiceColloStreamHandler.ServeHTTP(w, r)
		default:
			http.NotFound(w, r)
		}
	})
}

// UnimplementedColloServiceHandler returns CodeUnimplemented from all methods.
type UnimplementedColloServiceHandler struct{}

func (UnimplementedColloServiceHandler) ColloStream(context.Context, *connect.Request[v1.ColloRequest], *connect.ServerStream[v1.ColloStreamResponse]) error {
	return connect.NewError(connect.CodeUnimplemented, errors.New("collo.v1.ColloService.ColloStream is not implemented"))
}
