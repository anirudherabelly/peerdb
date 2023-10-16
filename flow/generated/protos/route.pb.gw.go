// Code generated by protoc-gen-grpc-gateway. DO NOT EDIT.
// source: route.proto

/*
Package protos is a reverse proxy.

It translates gRPC into RESTful JSON APIs.
*/
package protos

import (
	"context"
	"io"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/grpc-ecosystem/grpc-gateway/v2/utilities"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

// Suppress "imported and not used" errors
var _ codes.Code
var _ io.Reader
var _ status.Status
var _ = runtime.String
var _ = utilities.NewDoubleArray
var _ = metadata.Join

func request_FlowService_ValidatePeer_0(ctx context.Context, marshaler runtime.Marshaler, client FlowServiceClient, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var protoReq ValidatePeerRequest
	var metadata runtime.ServerMetadata

	newReader, berr := utilities.IOReaderFactory(req.Body)
	if berr != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "%v", berr)
	}
	if err := marshaler.NewDecoder(newReader()).Decode(&protoReq); err != nil && err != io.EOF {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "%v", err)
	}

	msg, err := client.ValidatePeer(ctx, &protoReq, grpc.Header(&metadata.HeaderMD), grpc.Trailer(&metadata.TrailerMD))
	return msg, metadata, err

}

func local_request_FlowService_ValidatePeer_0(ctx context.Context, marshaler runtime.Marshaler, server FlowServiceServer, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var protoReq ValidatePeerRequest
	var metadata runtime.ServerMetadata

	newReader, berr := utilities.IOReaderFactory(req.Body)
	if berr != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "%v", berr)
	}
	if err := marshaler.NewDecoder(newReader()).Decode(&protoReq); err != nil && err != io.EOF {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "%v", err)
	}

	msg, err := server.ValidatePeer(ctx, &protoReq)
	return msg, metadata, err

}

func request_FlowService_CreatePeer_0(ctx context.Context, marshaler runtime.Marshaler, client FlowServiceClient, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var protoReq CreatePeerRequest
	var metadata runtime.ServerMetadata

	newReader, berr := utilities.IOReaderFactory(req.Body)
	if berr != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "%v", berr)
	}
	if err := marshaler.NewDecoder(newReader()).Decode(&protoReq); err != nil && err != io.EOF {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "%v", err)
	}

	msg, err := client.CreatePeer(ctx, &protoReq, grpc.Header(&metadata.HeaderMD), grpc.Trailer(&metadata.TrailerMD))
	return msg, metadata, err

}

func local_request_FlowService_CreatePeer_0(ctx context.Context, marshaler runtime.Marshaler, server FlowServiceServer, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var protoReq CreatePeerRequest
	var metadata runtime.ServerMetadata

	newReader, berr := utilities.IOReaderFactory(req.Body)
	if berr != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "%v", berr)
	}
	if err := marshaler.NewDecoder(newReader()).Decode(&protoReq); err != nil && err != io.EOF {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "%v", err)
	}

	msg, err := server.CreatePeer(ctx, &protoReq)
	return msg, metadata, err

}

func request_FlowService_CreateCDCFlow_0(ctx context.Context, marshaler runtime.Marshaler, client FlowServiceClient, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var protoReq CreateCDCFlowRequest
	var metadata runtime.ServerMetadata

	newReader, berr := utilities.IOReaderFactory(req.Body)
	if berr != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "%v", berr)
	}
	if err := marshaler.NewDecoder(newReader()).Decode(&protoReq); err != nil && err != io.EOF {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "%v", err)
	}

	msg, err := client.CreateCDCFlow(ctx, &protoReq, grpc.Header(&metadata.HeaderMD), grpc.Trailer(&metadata.TrailerMD))
	return msg, metadata, err

}

func local_request_FlowService_CreateCDCFlow_0(ctx context.Context, marshaler runtime.Marshaler, server FlowServiceServer, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var protoReq CreateCDCFlowRequest
	var metadata runtime.ServerMetadata

	newReader, berr := utilities.IOReaderFactory(req.Body)
	if berr != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "%v", berr)
	}
	if err := marshaler.NewDecoder(newReader()).Decode(&protoReq); err != nil && err != io.EOF {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "%v", err)
	}

	msg, err := server.CreateCDCFlow(ctx, &protoReq)
	return msg, metadata, err

}

func request_FlowService_MirrorStatus_0(ctx context.Context, marshaler runtime.Marshaler, client FlowServiceClient, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var protoReq MirrorStatusRequest
	var metadata runtime.ServerMetadata

	var (
		val string
		ok  bool
		err error
		_   = err
	)

	val, ok = pathParams["flow_job_name"]
	if !ok {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "missing parameter %s", "flow_job_name")
	}

	protoReq.FlowJobName, err = runtime.String(val)
	if err != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "type mismatch, parameter: %s, error: %v", "flow_job_name", err)
	}

	msg, err := client.MirrorStatus(ctx, &protoReq, grpc.Header(&metadata.HeaderMD), grpc.Trailer(&metadata.TrailerMD))
	return msg, metadata, err

}

func local_request_FlowService_MirrorStatus_0(ctx context.Context, marshaler runtime.Marshaler, server FlowServiceServer, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var protoReq MirrorStatusRequest
	var metadata runtime.ServerMetadata

	var (
		val string
		ok  bool
		err error
		_   = err
	)

	val, ok = pathParams["flow_job_name"]
	if !ok {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "missing parameter %s", "flow_job_name")
	}

	protoReq.FlowJobName, err = runtime.String(val)
	if err != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "type mismatch, parameter: %s, error: %v", "flow_job_name", err)
	}

	msg, err := server.MirrorStatus(ctx, &protoReq)
	return msg, metadata, err

}

// RegisterFlowServiceHandlerServer registers the http handlers for service FlowService to "mux".
// UnaryRPC     :call FlowServiceServer directly.
// StreamingRPC :currently unsupported pending https://github.com/grpc/grpc-go/issues/906.
// Note that using this registration option will cause many gRPC library features to stop working. Consider using RegisterFlowServiceHandlerFromEndpoint instead.
func RegisterFlowServiceHandlerServer(ctx context.Context, mux *runtime.ServeMux, server FlowServiceServer) error {

	mux.Handle("POST", pattern_FlowService_ValidatePeer_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		var stream runtime.ServerTransportStream
		ctx = grpc.NewContextWithServerTransportStream(ctx, &stream)
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		var err error
		var annotatedContext context.Context
		annotatedContext, err = runtime.AnnotateIncomingContext(ctx, mux, req, "/peerdb_route.FlowService/ValidatePeer", runtime.WithHTTPPathPattern("/v1/peers/validate"))
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := local_request_FlowService_ValidatePeer_0(annotatedContext, inboundMarshaler, server, req, pathParams)
		md.HeaderMD, md.TrailerMD = metadata.Join(md.HeaderMD, stream.Header()), metadata.Join(md.TrailerMD, stream.Trailer())
		annotatedContext = runtime.NewServerMetadataContext(annotatedContext, md)
		if err != nil {
			runtime.HTTPError(annotatedContext, mux, outboundMarshaler, w, req, err)
			return
		}

		forward_FlowService_ValidatePeer_0(annotatedContext, mux, outboundMarshaler, w, req, resp, mux.GetForwardResponseOptions()...)

	})

	mux.Handle("POST", pattern_FlowService_CreatePeer_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		var stream runtime.ServerTransportStream
		ctx = grpc.NewContextWithServerTransportStream(ctx, &stream)
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		var err error
		var annotatedContext context.Context
		annotatedContext, err = runtime.AnnotateIncomingContext(ctx, mux, req, "/peerdb_route.FlowService/CreatePeer", runtime.WithHTTPPathPattern("/v1/peers/create"))
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := local_request_FlowService_CreatePeer_0(annotatedContext, inboundMarshaler, server, req, pathParams)
		md.HeaderMD, md.TrailerMD = metadata.Join(md.HeaderMD, stream.Header()), metadata.Join(md.TrailerMD, stream.Trailer())
		annotatedContext = runtime.NewServerMetadataContext(annotatedContext, md)
		if err != nil {
			runtime.HTTPError(annotatedContext, mux, outboundMarshaler, w, req, err)
			return
		}

		forward_FlowService_CreatePeer_0(annotatedContext, mux, outboundMarshaler, w, req, resp, mux.GetForwardResponseOptions()...)

	})

	mux.Handle("POST", pattern_FlowService_CreateCDCFlow_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		var stream runtime.ServerTransportStream
		ctx = grpc.NewContextWithServerTransportStream(ctx, &stream)
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		var err error
		var annotatedContext context.Context
		annotatedContext, err = runtime.AnnotateIncomingContext(ctx, mux, req, "/peerdb_route.FlowService/CreateCDCFlow", runtime.WithHTTPPathPattern("/v1/flows/cdc/create"))
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := local_request_FlowService_CreateCDCFlow_0(annotatedContext, inboundMarshaler, server, req, pathParams)
		md.HeaderMD, md.TrailerMD = metadata.Join(md.HeaderMD, stream.Header()), metadata.Join(md.TrailerMD, stream.Trailer())
		annotatedContext = runtime.NewServerMetadataContext(annotatedContext, md)
		if err != nil {
			runtime.HTTPError(annotatedContext, mux, outboundMarshaler, w, req, err)
			return
		}

		forward_FlowService_CreateCDCFlow_0(annotatedContext, mux, outboundMarshaler, w, req, resp, mux.GetForwardResponseOptions()...)

	})

	mux.Handle("GET", pattern_FlowService_MirrorStatus_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		var stream runtime.ServerTransportStream
		ctx = grpc.NewContextWithServerTransportStream(ctx, &stream)
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		var err error
		var annotatedContext context.Context
		annotatedContext, err = runtime.AnnotateIncomingContext(ctx, mux, req, "/peerdb_route.FlowService/MirrorStatus", runtime.WithHTTPPathPattern("/v1/mirrors/{flow_job_name}"))
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := local_request_FlowService_MirrorStatus_0(annotatedContext, inboundMarshaler, server, req, pathParams)
		md.HeaderMD, md.TrailerMD = metadata.Join(md.HeaderMD, stream.Header()), metadata.Join(md.TrailerMD, stream.Trailer())
		annotatedContext = runtime.NewServerMetadataContext(annotatedContext, md)
		if err != nil {
			runtime.HTTPError(annotatedContext, mux, outboundMarshaler, w, req, err)
			return
		}

		forward_FlowService_MirrorStatus_0(annotatedContext, mux, outboundMarshaler, w, req, resp, mux.GetForwardResponseOptions()...)

	})

	return nil
}

// RegisterFlowServiceHandlerFromEndpoint is same as RegisterFlowServiceHandler but
// automatically dials to "endpoint" and closes the connection when "ctx" gets done.
func RegisterFlowServiceHandlerFromEndpoint(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error) {
	conn, err := grpc.DialContext(ctx, endpoint, opts...)
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

	return RegisterFlowServiceHandler(ctx, mux, conn)
}

// RegisterFlowServiceHandler registers the http handlers for service FlowService to "mux".
// The handlers forward requests to the grpc endpoint over "conn".
func RegisterFlowServiceHandler(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error {
	return RegisterFlowServiceHandlerClient(ctx, mux, NewFlowServiceClient(conn))
}

// RegisterFlowServiceHandlerClient registers the http handlers for service FlowService
// to "mux". The handlers forward requests to the grpc endpoint over the given implementation of "FlowServiceClient".
// Note: the gRPC framework executes interceptors within the gRPC handler. If the passed in "FlowServiceClient"
// doesn't go through the normal gRPC flow (creating a gRPC client etc.) then it will be up to the passed in
// "FlowServiceClient" to call the correct interceptors.
func RegisterFlowServiceHandlerClient(ctx context.Context, mux *runtime.ServeMux, client FlowServiceClient) error {

	mux.Handle("POST", pattern_FlowService_ValidatePeer_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		var err error
		var annotatedContext context.Context
		annotatedContext, err = runtime.AnnotateContext(ctx, mux, req, "/peerdb_route.FlowService/ValidatePeer", runtime.WithHTTPPathPattern("/v1/peers/validate"))
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := request_FlowService_ValidatePeer_0(annotatedContext, inboundMarshaler, client, req, pathParams)
		annotatedContext = runtime.NewServerMetadataContext(annotatedContext, md)
		if err != nil {
			runtime.HTTPError(annotatedContext, mux, outboundMarshaler, w, req, err)
			return
		}

		forward_FlowService_ValidatePeer_0(annotatedContext, mux, outboundMarshaler, w, req, resp, mux.GetForwardResponseOptions()...)

	})

	mux.Handle("POST", pattern_FlowService_CreatePeer_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		var err error
		var annotatedContext context.Context
		annotatedContext, err = runtime.AnnotateContext(ctx, mux, req, "/peerdb_route.FlowService/CreatePeer", runtime.WithHTTPPathPattern("/v1/peers/create"))
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := request_FlowService_CreatePeer_0(annotatedContext, inboundMarshaler, client, req, pathParams)
		annotatedContext = runtime.NewServerMetadataContext(annotatedContext, md)
		if err != nil {
			runtime.HTTPError(annotatedContext, mux, outboundMarshaler, w, req, err)
			return
		}

		forward_FlowService_CreatePeer_0(annotatedContext, mux, outboundMarshaler, w, req, resp, mux.GetForwardResponseOptions()...)

	})

	mux.Handle("POST", pattern_FlowService_CreateCDCFlow_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		var err error
		var annotatedContext context.Context
		annotatedContext, err = runtime.AnnotateContext(ctx, mux, req, "/peerdb_route.FlowService/CreateCDCFlow", runtime.WithHTTPPathPattern("/v1/flows/cdc/create"))
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := request_FlowService_CreateCDCFlow_0(annotatedContext, inboundMarshaler, client, req, pathParams)
		annotatedContext = runtime.NewServerMetadataContext(annotatedContext, md)
		if err != nil {
			runtime.HTTPError(annotatedContext, mux, outboundMarshaler, w, req, err)
			return
		}

		forward_FlowService_CreateCDCFlow_0(annotatedContext, mux, outboundMarshaler, w, req, resp, mux.GetForwardResponseOptions()...)

	})

	mux.Handle("GET", pattern_FlowService_MirrorStatus_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		var err error
		var annotatedContext context.Context
		annotatedContext, err = runtime.AnnotateContext(ctx, mux, req, "/peerdb_route.FlowService/MirrorStatus", runtime.WithHTTPPathPattern("/v1/mirrors/{flow_job_name}"))
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := request_FlowService_MirrorStatus_0(annotatedContext, inboundMarshaler, client, req, pathParams)
		annotatedContext = runtime.NewServerMetadataContext(annotatedContext, md)
		if err != nil {
			runtime.HTTPError(annotatedContext, mux, outboundMarshaler, w, req, err)
			return
		}

		forward_FlowService_MirrorStatus_0(annotatedContext, mux, outboundMarshaler, w, req, resp, mux.GetForwardResponseOptions()...)

	})

	return nil
}

var (
	pattern_FlowService_ValidatePeer_0 = runtime.MustPattern(runtime.NewPattern(1, []int{2, 0, 2, 1, 2, 2}, []string{"v1", "peers", "validate"}, ""))

	pattern_FlowService_CreatePeer_0 = runtime.MustPattern(runtime.NewPattern(1, []int{2, 0, 2, 1, 2, 2}, []string{"v1", "peers", "create"}, ""))

	pattern_FlowService_CreateCDCFlow_0 = runtime.MustPattern(runtime.NewPattern(1, []int{2, 0, 2, 1, 2, 2, 2, 3}, []string{"v1", "flows", "cdc", "create"}, ""))

	pattern_FlowService_MirrorStatus_0 = runtime.MustPattern(runtime.NewPattern(1, []int{2, 0, 2, 1, 1, 0, 4, 1, 5, 2}, []string{"v1", "mirrors", "flow_job_name"}, ""))
)

var (
	forward_FlowService_ValidatePeer_0 = runtime.ForwardResponseMessage

	forward_FlowService_CreatePeer_0 = runtime.ForwardResponseMessage

	forward_FlowService_CreateCDCFlow_0 = runtime.ForwardResponseMessage

	forward_FlowService_MirrorStatus_0 = runtime.ForwardResponseMessage
)
