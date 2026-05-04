package exportpb

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	ExportService_ServiceName                  = "export.ExportService"
	ExportService_Ping_FullMethodName         = "/export.ExportService/Ping"
	ExportService_BuildExport_FullMethodName  = "/export.ExportService/BuildExport"
)

type ExportServiceClient interface {
	Ping(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error)
	BuildExport(ctx context.Context, in *BuildExportRequest, opts ...grpc.CallOption) (*BuildExportResponse, error)
}

type exportServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewExportServiceClient(cc grpc.ClientConnInterface) ExportServiceClient {
	return &exportServiceClient{cc: cc}
}

func invoke[Resp any](ctx context.Context, cc grpc.ClientConnInterface, method string, in any, out Resp, opts ...grpc.CallOption) (Resp, error) {
	if err := cc.Invoke(ctx, method, in, out, opts...); err != nil {
		var zero Resp
		return zero, err
	}
	return out, nil
}

func (c *exportServiceClient) Ping(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error) {
	return invoke(ctx, c.cc, ExportService_Ping_FullMethodName, in, new(Empty), opts...)
}

func (c *exportServiceClient) BuildExport(ctx context.Context, in *BuildExportRequest, opts ...grpc.CallOption) (*BuildExportResponse, error) {
	return invoke(ctx, c.cc, ExportService_BuildExport_FullMethodName, in, new(BuildExportResponse), opts...)
}

type ExportServiceServer interface {
	Ping(context.Context, *Empty) (*Empty, error)
	BuildExport(context.Context, *BuildExportRequest) (*BuildExportResponse, error)
}

type UnimplementedExportServiceServer struct{}

func (UnimplementedExportServiceServer) Ping(context.Context, *Empty) (*Empty, error) {
	return nil, status.Error(codes.Unimplemented, "method Ping not implemented")
}

func (UnimplementedExportServiceServer) BuildExport(context.Context, *BuildExportRequest) (*BuildExportResponse, error) {
	return nil, status.Error(codes.Unimplemented, "method BuildExport not implemented")
}

func RegisterExportServiceServer(s grpc.ServiceRegistrar, srv ExportServiceServer) {
	s.RegisterService(&ExportService_ServiceDesc, srv)
}

func handleUnary(ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor, fullMethod string, in any, srv interface{}, call func(context.Context, any) (any, error)) (any, error) {
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return call(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: fullMethod}
	handler := func(ctx context.Context, req any) (any, error) {
		return call(ctx, req)
	}
	return interceptor(ctx, in, info, handler)
}

func _ExportService_Ping_Handler(srv interface{}, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor) (any, error) {
	return handleUnary(ctx, dec, interceptor, ExportService_Ping_FullMethodName, new(Empty), srv, func(ctx context.Context, req any) (any, error) {
		return srv.(ExportServiceServer).Ping(ctx, req.(*Empty))
	})
}

func _ExportService_BuildExport_Handler(srv interface{}, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor) (any, error) {
	return handleUnary(ctx, dec, interceptor, ExportService_BuildExport_FullMethodName, new(BuildExportRequest), srv, func(ctx context.Context, req any) (any, error) {
		return srv.(ExportServiceServer).BuildExport(ctx, req.(*BuildExportRequest))
	})
}

var ExportService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: ExportService_ServiceName,
	HandlerType: (*ExportServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{MethodName: "Ping", Handler: _ExportService_Ping_Handler},
		{MethodName: "BuildExport", Handler: _ExportService_BuildExport_Handler},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "export/export.proto",
}
