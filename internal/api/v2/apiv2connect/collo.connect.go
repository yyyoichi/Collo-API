// Code generated by protoc-gen-connect-go. DO NOT EDIT.
//
// Source: api/v2/collo.proto

package apiv2connect

import (
	connect "connectrpc.com/connect"
	context "context"
	errors "errors"
	http "net/http"
	strings "strings"
	v2 "yyyoichi/Collo-API/internal/api/v2"
)

// This is a compile-time assertion to ensure that this generated file and the connect package are
// compatible. If you get a compiler error that this constant is not defined, this code was
// generated with a version of connect newer than the one compiled into your binary. You can fix the
// problem by either regenerating this code with an older version of connect or updating the connect
// version compiled into your binary.
const _ = connect.IsAtLeastVersion0_1_0

const (
	// ColloNetworkServiceName is the fully-qualified name of the ColloNetworkService service.
	ColloNetworkServiceName = "api.v2.ColloNetworkService"
	// ColloWebServiceName is the fully-qualified name of the ColloWebService service.
	ColloWebServiceName = "api.v2.ColloWebService"
)

// These constants are the fully-qualified names of the RPCs defined in this package. They're
// exposed at runtime as Spec.Procedure and as the final two segments of the HTTP route.
//
// Note that these are different from the fully-qualified method names used by
// google.golang.org/protobuf/reflect/protoreflect. To convert from these constants to
// reflection-formatted method names, remove the leading slash and convert the remaining slash to a
// period.
const (
	// ColloNetworkServiceColloNetworkStreamProcedure is the fully-qualified name of the
	// ColloNetworkService's ColloNetworkStream RPC.
	ColloNetworkServiceColloNetworkStreamProcedure = "/api.v2.ColloNetworkService/ColloNetworkStream"
	// ColloWebServiceColloWebInitStreamProcedure is the fully-qualified name of the ColloWebService's
	// ColloWebInitStream RPC.
	ColloWebServiceColloWebInitStreamProcedure = "/api.v2.ColloWebService/ColloWebInitStream"
	// ColloWebServiceColloWebStreamProcedure is the fully-qualified name of the ColloWebService's
	// ColloWebStream RPC.
	ColloWebServiceColloWebStreamProcedure = "/api.v2.ColloWebService/ColloWebStream"
)

// ColloNetworkServiceClient is a client for the api.v2.ColloNetworkService service.
type ColloNetworkServiceClient interface {
	ColloNetworkStream(context.Context) *connect.BidiStreamForClient[v2.ColloNetworkStreamRequest, v2.ColloNetworkStreamResponse]
}

// NewColloNetworkServiceClient constructs a client for the api.v2.ColloNetworkService service. By
// default, it uses the Connect protocol with the binary Protobuf Codec, asks for gzipped responses,
// and sends uncompressed requests. To use the gRPC or gRPC-Web protocols, supply the
// connect.WithGRPC() or connect.WithGRPCWeb() options.
//
// The URL supplied here should be the base URL for the Connect or gRPC server (for example,
// http://api.acme.com or https://acme.com/grpc).
func NewColloNetworkServiceClient(httpClient connect.HTTPClient, baseURL string, opts ...connect.ClientOption) ColloNetworkServiceClient {
	baseURL = strings.TrimRight(baseURL, "/")
	return &colloNetworkServiceClient{
		colloNetworkStream: connect.NewClient[v2.ColloNetworkStreamRequest, v2.ColloNetworkStreamResponse](
			httpClient,
			baseURL+ColloNetworkServiceColloNetworkStreamProcedure,
			opts...,
		),
	}
}

// colloNetworkServiceClient implements ColloNetworkServiceClient.
type colloNetworkServiceClient struct {
	colloNetworkStream *connect.Client[v2.ColloNetworkStreamRequest, v2.ColloNetworkStreamResponse]
}

// ColloNetworkStream calls api.v2.ColloNetworkService.ColloNetworkStream.
func (c *colloNetworkServiceClient) ColloNetworkStream(ctx context.Context) *connect.BidiStreamForClient[v2.ColloNetworkStreamRequest, v2.ColloNetworkStreamResponse] {
	return c.colloNetworkStream.CallBidiStream(ctx)
}

// ColloNetworkServiceHandler is an implementation of the api.v2.ColloNetworkService service.
type ColloNetworkServiceHandler interface {
	ColloNetworkStream(context.Context, *connect.BidiStream[v2.ColloNetworkStreamRequest, v2.ColloNetworkStreamResponse]) error
}

// NewColloNetworkServiceHandler builds an HTTP handler from the service implementation. It returns
// the path on which to mount the handler and the handler itself.
//
// By default, handlers support the Connect, gRPC, and gRPC-Web protocols with the binary Protobuf
// and JSON codecs. They also support gzip compression.
func NewColloNetworkServiceHandler(svc ColloNetworkServiceHandler, opts ...connect.HandlerOption) (string, http.Handler) {
	colloNetworkServiceColloNetworkStreamHandler := connect.NewBidiStreamHandler(
		ColloNetworkServiceColloNetworkStreamProcedure,
		svc.ColloNetworkStream,
		opts...,
	)
	return "/api.v2.ColloNetworkService/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case ColloNetworkServiceColloNetworkStreamProcedure:
			colloNetworkServiceColloNetworkStreamHandler.ServeHTTP(w, r)
		default:
			http.NotFound(w, r)
		}
	})
}

// UnimplementedColloNetworkServiceHandler returns CodeUnimplemented from all methods.
type UnimplementedColloNetworkServiceHandler struct{}

func (UnimplementedColloNetworkServiceHandler) ColloNetworkStream(context.Context, *connect.BidiStream[v2.ColloNetworkStreamRequest, v2.ColloNetworkStreamResponse]) error {
	return connect.NewError(connect.CodeUnimplemented, errors.New("api.v2.ColloNetworkService.ColloNetworkStream is not implemented"))
}

// ColloWebServiceClient is a client for the api.v2.ColloWebService service.
type ColloWebServiceClient interface {
	ColloWebInitStream(context.Context, *connect.Request[v2.ColloWebInitStreamRequest]) (*connect.ServerStreamForClient[v2.ColloWebInitStreamResponse], error)
	ColloWebStream(context.Context, *connect.Request[v2.ColloWebStreamRequest]) (*connect.ServerStreamForClient[v2.ColloWebStreamResponse], error)
}

// NewColloWebServiceClient constructs a client for the api.v2.ColloWebService service. By default,
// it uses the Connect protocol with the binary Protobuf Codec, asks for gzipped responses, and
// sends uncompressed requests. To use the gRPC or gRPC-Web protocols, supply the connect.WithGRPC()
// or connect.WithGRPCWeb() options.
//
// The URL supplied here should be the base URL for the Connect or gRPC server (for example,
// http://api.acme.com or https://acme.com/grpc).
func NewColloWebServiceClient(httpClient connect.HTTPClient, baseURL string, opts ...connect.ClientOption) ColloWebServiceClient {
	baseURL = strings.TrimRight(baseURL, "/")
	return &colloWebServiceClient{
		colloWebInitStream: connect.NewClient[v2.ColloWebInitStreamRequest, v2.ColloWebInitStreamResponse](
			httpClient,
			baseURL+ColloWebServiceColloWebInitStreamProcedure,
			opts...,
		),
		colloWebStream: connect.NewClient[v2.ColloWebStreamRequest, v2.ColloWebStreamResponse](
			httpClient,
			baseURL+ColloWebServiceColloWebStreamProcedure,
			opts...,
		),
	}
}

// colloWebServiceClient implements ColloWebServiceClient.
type colloWebServiceClient struct {
	colloWebInitStream *connect.Client[v2.ColloWebInitStreamRequest, v2.ColloWebInitStreamResponse]
	colloWebStream     *connect.Client[v2.ColloWebStreamRequest, v2.ColloWebStreamResponse]
}

// ColloWebInitStream calls api.v2.ColloWebService.ColloWebInitStream.
func (c *colloWebServiceClient) ColloWebInitStream(ctx context.Context, req *connect.Request[v2.ColloWebInitStreamRequest]) (*connect.ServerStreamForClient[v2.ColloWebInitStreamResponse], error) {
	return c.colloWebInitStream.CallServerStream(ctx, req)
}

// ColloWebStream calls api.v2.ColloWebService.ColloWebStream.
func (c *colloWebServiceClient) ColloWebStream(ctx context.Context, req *connect.Request[v2.ColloWebStreamRequest]) (*connect.ServerStreamForClient[v2.ColloWebStreamResponse], error) {
	return c.colloWebStream.CallServerStream(ctx, req)
}

// ColloWebServiceHandler is an implementation of the api.v2.ColloWebService service.
type ColloWebServiceHandler interface {
	ColloWebInitStream(context.Context, *connect.Request[v2.ColloWebInitStreamRequest], *connect.ServerStream[v2.ColloWebInitStreamResponse]) error
	ColloWebStream(context.Context, *connect.Request[v2.ColloWebStreamRequest], *connect.ServerStream[v2.ColloWebStreamResponse]) error
}

// NewColloWebServiceHandler builds an HTTP handler from the service implementation. It returns the
// path on which to mount the handler and the handler itself.
//
// By default, handlers support the Connect, gRPC, and gRPC-Web protocols with the binary Protobuf
// and JSON codecs. They also support gzip compression.
func NewColloWebServiceHandler(svc ColloWebServiceHandler, opts ...connect.HandlerOption) (string, http.Handler) {
	colloWebServiceColloWebInitStreamHandler := connect.NewServerStreamHandler(
		ColloWebServiceColloWebInitStreamProcedure,
		svc.ColloWebInitStream,
		opts...,
	)
	colloWebServiceColloWebStreamHandler := connect.NewServerStreamHandler(
		ColloWebServiceColloWebStreamProcedure,
		svc.ColloWebStream,
		opts...,
	)
	return "/api.v2.ColloWebService/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case ColloWebServiceColloWebInitStreamProcedure:
			colloWebServiceColloWebInitStreamHandler.ServeHTTP(w, r)
		case ColloWebServiceColloWebStreamProcedure:
			colloWebServiceColloWebStreamHandler.ServeHTTP(w, r)
		default:
			http.NotFound(w, r)
		}
	})
}

// UnimplementedColloWebServiceHandler returns CodeUnimplemented from all methods.
type UnimplementedColloWebServiceHandler struct{}

func (UnimplementedColloWebServiceHandler) ColloWebInitStream(context.Context, *connect.Request[v2.ColloWebInitStreamRequest], *connect.ServerStream[v2.ColloWebInitStreamResponse]) error {
	return connect.NewError(connect.CodeUnimplemented, errors.New("api.v2.ColloWebService.ColloWebInitStream is not implemented"))
}

func (UnimplementedColloWebServiceHandler) ColloWebStream(context.Context, *connect.Request[v2.ColloWebStreamRequest], *connect.ServerStream[v2.ColloWebStreamResponse]) error {
	return connect.NewError(connect.CodeUnimplemented, errors.New("api.v2.ColloWebService.ColloWebStream is not implemented"))
}
