package reportpb

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	ReportService_ServiceName                   = "report.ReportService"
	ReportService_Ping_FullMethodName          = "/report.ReportService/Ping"
	ReportService_CreateReport_FullMethodName  = "/report.ReportService/CreateReport"
	ReportService_GetReport_FullMethodName     = "/report.ReportService/GetReport"
	ReportService_ListReports_FullMethodName   = "/report.ReportService/ListReports"
	ReportService_UpdateReport_FullMethodName  = "/report.ReportService/UpdateReport"
	ReportService_DeleteReport_FullMethodName  = "/report.ReportService/DeleteReport"
	ReportService_CreateEntry_FullMethodName   = "/report.ReportService/CreateEntry"
	ReportService_GetEntry_FullMethodName      = "/report.ReportService/GetEntry"
	ReportService_ListEntries_FullMethodName   = "/report.ReportService/ListEntries"
	ReportService_UpdateEntry_FullMethodName   = "/report.ReportService/UpdateEntry"
	ReportService_DeleteEntry_FullMethodName   = "/report.ReportService/DeleteEntry"
	ReportService_CreateComment_FullMethodName = "/report.ReportService/CreateComment"
	ReportService_GetComment_FullMethodName    = "/report.ReportService/GetComment"
	ReportService_ListComments_FullMethodName  = "/report.ReportService/ListComments"
	ReportService_UpdateComment_FullMethodName = "/report.ReportService/UpdateComment"
	ReportService_DeleteComment_FullMethodName = "/report.ReportService/DeleteComment"
)

type ReportServiceClient interface {
	Ping(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error)
	CreateReport(ctx context.Context, in *CreateDailyReportRequest, opts ...grpc.CallOption) (*DailyReport, error)
	GetReport(ctx context.Context, in *GetDailyReportRequest, opts ...grpc.CallOption) (*DailyReport, error)
	ListReports(ctx context.Context, in *ListDailyReportsRequest, opts ...grpc.CallOption) (*ListDailyReportsResponse, error)
	UpdateReport(ctx context.Context, in *UpdateDailyReportRequest, opts ...grpc.CallOption) (*DailyReport, error)
	DeleteReport(ctx context.Context, in *DeleteDailyReportRequest, opts ...grpc.CallOption) (*Empty, error)
	CreateEntry(ctx context.Context, in *CreateDailyReportEntryRequest, opts ...grpc.CallOption) (*DailyReportEntry, error)
	GetEntry(ctx context.Context, in *GetDailyReportEntryRequest, opts ...grpc.CallOption) (*DailyReportEntry, error)
	ListEntries(ctx context.Context, in *ListDailyReportEntriesRequest, opts ...grpc.CallOption) (*ListDailyReportEntriesResponse, error)
	UpdateEntry(ctx context.Context, in *UpdateDailyReportEntryRequest, opts ...grpc.CallOption) (*DailyReportEntry, error)
	DeleteEntry(ctx context.Context, in *DeleteDailyReportEntryRequest, opts ...grpc.CallOption) (*Empty, error)
	CreateComment(ctx context.Context, in *CreateDailyReportCommentRequest, opts ...grpc.CallOption) (*DailyReportComment, error)
	GetComment(ctx context.Context, in *GetDailyReportCommentRequest, opts ...grpc.CallOption) (*DailyReportComment, error)
	ListComments(ctx context.Context, in *ListDailyReportCommentsRequest, opts ...grpc.CallOption) (*ListDailyReportCommentsResponse, error)
	UpdateComment(ctx context.Context, in *UpdateDailyReportCommentRequest, opts ...grpc.CallOption) (*DailyReportComment, error)
	DeleteComment(ctx context.Context, in *DeleteDailyReportCommentRequest, opts ...grpc.CallOption) (*Empty, error)
}

type reportServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewReportServiceClient(cc grpc.ClientConnInterface) ReportServiceClient {
	return &reportServiceClient{cc: cc}
}

func invoke[Resp any](ctx context.Context, cc grpc.ClientConnInterface, method string, in any, out Resp, opts ...grpc.CallOption) (Resp, error) {
	if err := cc.Invoke(ctx, method, in, out, opts...); err != nil {
		var zero Resp
		return zero, err
	}
	return out, nil
}

func (c *reportServiceClient) Ping(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error) {
	return invoke(ctx, c.cc, ReportService_Ping_FullMethodName, in, new(Empty), opts...)
}
func (c *reportServiceClient) CreateReport(ctx context.Context, in *CreateDailyReportRequest, opts ...grpc.CallOption) (*DailyReport, error) {
	return invoke(ctx, c.cc, ReportService_CreateReport_FullMethodName, in, new(DailyReport), opts...)
}
func (c *reportServiceClient) GetReport(ctx context.Context, in *GetDailyReportRequest, opts ...grpc.CallOption) (*DailyReport, error) {
	return invoke(ctx, c.cc, ReportService_GetReport_FullMethodName, in, new(DailyReport), opts...)
}
func (c *reportServiceClient) ListReports(ctx context.Context, in *ListDailyReportsRequest, opts ...grpc.CallOption) (*ListDailyReportsResponse, error) {
	return invoke(ctx, c.cc, ReportService_ListReports_FullMethodName, in, new(ListDailyReportsResponse), opts...)
}
func (c *reportServiceClient) UpdateReport(ctx context.Context, in *UpdateDailyReportRequest, opts ...grpc.CallOption) (*DailyReport, error) {
	return invoke(ctx, c.cc, ReportService_UpdateReport_FullMethodName, in, new(DailyReport), opts...)
}
func (c *reportServiceClient) DeleteReport(ctx context.Context, in *DeleteDailyReportRequest, opts ...grpc.CallOption) (*Empty, error) {
	return invoke(ctx, c.cc, ReportService_DeleteReport_FullMethodName, in, new(Empty), opts...)
}
func (c *reportServiceClient) CreateEntry(ctx context.Context, in *CreateDailyReportEntryRequest, opts ...grpc.CallOption) (*DailyReportEntry, error) {
	return invoke(ctx, c.cc, ReportService_CreateEntry_FullMethodName, in, new(DailyReportEntry), opts...)
}
func (c *reportServiceClient) GetEntry(ctx context.Context, in *GetDailyReportEntryRequest, opts ...grpc.CallOption) (*DailyReportEntry, error) {
	return invoke(ctx, c.cc, ReportService_GetEntry_FullMethodName, in, new(DailyReportEntry), opts...)
}
func (c *reportServiceClient) ListEntries(ctx context.Context, in *ListDailyReportEntriesRequest, opts ...grpc.CallOption) (*ListDailyReportEntriesResponse, error) {
	return invoke(ctx, c.cc, ReportService_ListEntries_FullMethodName, in, new(ListDailyReportEntriesResponse), opts...)
}
func (c *reportServiceClient) UpdateEntry(ctx context.Context, in *UpdateDailyReportEntryRequest, opts ...grpc.CallOption) (*DailyReportEntry, error) {
	return invoke(ctx, c.cc, ReportService_UpdateEntry_FullMethodName, in, new(DailyReportEntry), opts...)
}
func (c *reportServiceClient) DeleteEntry(ctx context.Context, in *DeleteDailyReportEntryRequest, opts ...grpc.CallOption) (*Empty, error) {
	return invoke(ctx, c.cc, ReportService_DeleteEntry_FullMethodName, in, new(Empty), opts...)
}
func (c *reportServiceClient) CreateComment(ctx context.Context, in *CreateDailyReportCommentRequest, opts ...grpc.CallOption) (*DailyReportComment, error) {
	return invoke(ctx, c.cc, ReportService_CreateComment_FullMethodName, in, new(DailyReportComment), opts...)
}
func (c *reportServiceClient) GetComment(ctx context.Context, in *GetDailyReportCommentRequest, opts ...grpc.CallOption) (*DailyReportComment, error) {
	return invoke(ctx, c.cc, ReportService_GetComment_FullMethodName, in, new(DailyReportComment), opts...)
}
func (c *reportServiceClient) ListComments(ctx context.Context, in *ListDailyReportCommentsRequest, opts ...grpc.CallOption) (*ListDailyReportCommentsResponse, error) {
	return invoke(ctx, c.cc, ReportService_ListComments_FullMethodName, in, new(ListDailyReportCommentsResponse), opts...)
}
func (c *reportServiceClient) UpdateComment(ctx context.Context, in *UpdateDailyReportCommentRequest, opts ...grpc.CallOption) (*DailyReportComment, error) {
	return invoke(ctx, c.cc, ReportService_UpdateComment_FullMethodName, in, new(DailyReportComment), opts...)
}
func (c *reportServiceClient) DeleteComment(ctx context.Context, in *DeleteDailyReportCommentRequest, opts ...grpc.CallOption) (*Empty, error) {
	return invoke(ctx, c.cc, ReportService_DeleteComment_FullMethodName, in, new(Empty), opts...)
}

type ReportServiceServer interface {
	Ping(context.Context, *Empty) (*Empty, error)
	CreateReport(context.Context, *CreateDailyReportRequest) (*DailyReport, error)
	GetReport(context.Context, *GetDailyReportRequest) (*DailyReport, error)
	ListReports(context.Context, *ListDailyReportsRequest) (*ListDailyReportsResponse, error)
	UpdateReport(context.Context, *UpdateDailyReportRequest) (*DailyReport, error)
	DeleteReport(context.Context, *DeleteDailyReportRequest) (*Empty, error)
	CreateEntry(context.Context, *CreateDailyReportEntryRequest) (*DailyReportEntry, error)
	GetEntry(context.Context, *GetDailyReportEntryRequest) (*DailyReportEntry, error)
	ListEntries(context.Context, *ListDailyReportEntriesRequest) (*ListDailyReportEntriesResponse, error)
	UpdateEntry(context.Context, *UpdateDailyReportEntryRequest) (*DailyReportEntry, error)
	DeleteEntry(context.Context, *DeleteDailyReportEntryRequest) (*Empty, error)
	CreateComment(context.Context, *CreateDailyReportCommentRequest) (*DailyReportComment, error)
	GetComment(context.Context, *GetDailyReportCommentRequest) (*DailyReportComment, error)
	ListComments(context.Context, *ListDailyReportCommentsRequest) (*ListDailyReportCommentsResponse, error)
	UpdateComment(context.Context, *UpdateDailyReportCommentRequest) (*DailyReportComment, error)
	DeleteComment(context.Context, *DeleteDailyReportCommentRequest) (*Empty, error)
}

type UnimplementedReportServiceServer struct{}

func (UnimplementedReportServiceServer) Ping(context.Context, *Empty) (*Empty, error) { return nil, status.Error(codes.Unimplemented, "method Ping not implemented") }
func (UnimplementedReportServiceServer) CreateReport(context.Context, *CreateDailyReportRequest) (*DailyReport, error) {
	return nil, status.Error(codes.Unimplemented, "method CreateReport not implemented")
}
func (UnimplementedReportServiceServer) GetReport(context.Context, *GetDailyReportRequest) (*DailyReport, error) {
	return nil, status.Error(codes.Unimplemented, "method GetReport not implemented")
}
func (UnimplementedReportServiceServer) ListReports(context.Context, *ListDailyReportsRequest) (*ListDailyReportsResponse, error) {
	return nil, status.Error(codes.Unimplemented, "method ListReports not implemented")
}
func (UnimplementedReportServiceServer) UpdateReport(context.Context, *UpdateDailyReportRequest) (*DailyReport, error) {
	return nil, status.Error(codes.Unimplemented, "method UpdateReport not implemented")
}
func (UnimplementedReportServiceServer) DeleteReport(context.Context, *DeleteDailyReportRequest) (*Empty, error) {
	return nil, status.Error(codes.Unimplemented, "method DeleteReport not implemented")
}
func (UnimplementedReportServiceServer) CreateEntry(context.Context, *CreateDailyReportEntryRequest) (*DailyReportEntry, error) {
	return nil, status.Error(codes.Unimplemented, "method CreateEntry not implemented")
}
func (UnimplementedReportServiceServer) GetEntry(context.Context, *GetDailyReportEntryRequest) (*DailyReportEntry, error) {
	return nil, status.Error(codes.Unimplemented, "method GetEntry not implemented")
}
func (UnimplementedReportServiceServer) ListEntries(context.Context, *ListDailyReportEntriesRequest) (*ListDailyReportEntriesResponse, error) {
	return nil, status.Error(codes.Unimplemented, "method ListEntries not implemented")
}
func (UnimplementedReportServiceServer) UpdateEntry(context.Context, *UpdateDailyReportEntryRequest) (*DailyReportEntry, error) {
	return nil, status.Error(codes.Unimplemented, "method UpdateEntry not implemented")
}
func (UnimplementedReportServiceServer) DeleteEntry(context.Context, *DeleteDailyReportEntryRequest) (*Empty, error) {
	return nil, status.Error(codes.Unimplemented, "method DeleteEntry not implemented")
}
func (UnimplementedReportServiceServer) CreateComment(context.Context, *CreateDailyReportCommentRequest) (*DailyReportComment, error) {
	return nil, status.Error(codes.Unimplemented, "method CreateComment not implemented")
}
func (UnimplementedReportServiceServer) GetComment(context.Context, *GetDailyReportCommentRequest) (*DailyReportComment, error) {
	return nil, status.Error(codes.Unimplemented, "method GetComment not implemented")
}
func (UnimplementedReportServiceServer) ListComments(context.Context, *ListDailyReportCommentsRequest) (*ListDailyReportCommentsResponse, error) {
	return nil, status.Error(codes.Unimplemented, "method ListComments not implemented")
}
func (UnimplementedReportServiceServer) UpdateComment(context.Context, *UpdateDailyReportCommentRequest) (*DailyReportComment, error) {
	return nil, status.Error(codes.Unimplemented, "method UpdateComment not implemented")
}
func (UnimplementedReportServiceServer) DeleteComment(context.Context, *DeleteDailyReportCommentRequest) (*Empty, error) {
	return nil, status.Error(codes.Unimplemented, "method DeleteComment not implemented")
}

func RegisterReportServiceServer(s grpc.ServiceRegistrar, srv ReportServiceServer) {
	s.RegisterService(&ReportService_ServiceDesc, srv)
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

func _ReportService_Ping_Handler(srv interface{}, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor) (any, error) {
	return handleUnary(ctx, dec, interceptor, ReportService_Ping_FullMethodName, new(Empty), srv, func(ctx context.Context, req any) (any, error) {
		return srv.(ReportServiceServer).Ping(ctx, req.(*Empty))
	})
}
func _ReportService_CreateReport_Handler(srv interface{}, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor) (any, error) {
	return handleUnary(ctx, dec, interceptor, ReportService_CreateReport_FullMethodName, new(CreateDailyReportRequest), srv, func(ctx context.Context, req any) (any, error) {
		return srv.(ReportServiceServer).CreateReport(ctx, req.(*CreateDailyReportRequest))
	})
}
func _ReportService_GetReport_Handler(srv interface{}, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor) (any, error) {
	return handleUnary(ctx, dec, interceptor, ReportService_GetReport_FullMethodName, new(GetDailyReportRequest), srv, func(ctx context.Context, req any) (any, error) {
		return srv.(ReportServiceServer).GetReport(ctx, req.(*GetDailyReportRequest))
	})
}
func _ReportService_ListReports_Handler(srv interface{}, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor) (any, error) {
	return handleUnary(ctx, dec, interceptor, ReportService_ListReports_FullMethodName, new(ListDailyReportsRequest), srv, func(ctx context.Context, req any) (any, error) {
		return srv.(ReportServiceServer).ListReports(ctx, req.(*ListDailyReportsRequest))
	})
}
func _ReportService_UpdateReport_Handler(srv interface{}, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor) (any, error) {
	return handleUnary(ctx, dec, interceptor, ReportService_UpdateReport_FullMethodName, new(UpdateDailyReportRequest), srv, func(ctx context.Context, req any) (any, error) {
		return srv.(ReportServiceServer).UpdateReport(ctx, req.(*UpdateDailyReportRequest))
	})
}
func _ReportService_DeleteReport_Handler(srv interface{}, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor) (any, error) {
	return handleUnary(ctx, dec, interceptor, ReportService_DeleteReport_FullMethodName, new(DeleteDailyReportRequest), srv, func(ctx context.Context, req any) (any, error) {
		return srv.(ReportServiceServer).DeleteReport(ctx, req.(*DeleteDailyReportRequest))
	})
}
func _ReportService_CreateEntry_Handler(srv interface{}, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor) (any, error) {
	return handleUnary(ctx, dec, interceptor, ReportService_CreateEntry_FullMethodName, new(CreateDailyReportEntryRequest), srv, func(ctx context.Context, req any) (any, error) {
		return srv.(ReportServiceServer).CreateEntry(ctx, req.(*CreateDailyReportEntryRequest))
	})
}
func _ReportService_GetEntry_Handler(srv interface{}, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor) (any, error) {
	return handleUnary(ctx, dec, interceptor, ReportService_GetEntry_FullMethodName, new(GetDailyReportEntryRequest), srv, func(ctx context.Context, req any) (any, error) {
		return srv.(ReportServiceServer).GetEntry(ctx, req.(*GetDailyReportEntryRequest))
	})
}
func _ReportService_ListEntries_Handler(srv interface{}, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor) (any, error) {
	return handleUnary(ctx, dec, interceptor, ReportService_ListEntries_FullMethodName, new(ListDailyReportEntriesRequest), srv, func(ctx context.Context, req any) (any, error) {
		return srv.(ReportServiceServer).ListEntries(ctx, req.(*ListDailyReportEntriesRequest))
	})
}
func _ReportService_UpdateEntry_Handler(srv interface{}, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor) (any, error) {
	return handleUnary(ctx, dec, interceptor, ReportService_UpdateEntry_FullMethodName, new(UpdateDailyReportEntryRequest), srv, func(ctx context.Context, req any) (any, error) {
		return srv.(ReportServiceServer).UpdateEntry(ctx, req.(*UpdateDailyReportEntryRequest))
	})
}
func _ReportService_DeleteEntry_Handler(srv interface{}, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor) (any, error) {
	return handleUnary(ctx, dec, interceptor, ReportService_DeleteEntry_FullMethodName, new(DeleteDailyReportEntryRequest), srv, func(ctx context.Context, req any) (any, error) {
		return srv.(ReportServiceServer).DeleteEntry(ctx, req.(*DeleteDailyReportEntryRequest))
	})
}
func _ReportService_CreateComment_Handler(srv interface{}, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor) (any, error) {
	return handleUnary(ctx, dec, interceptor, ReportService_CreateComment_FullMethodName, new(CreateDailyReportCommentRequest), srv, func(ctx context.Context, req any) (any, error) {
		return srv.(ReportServiceServer).CreateComment(ctx, req.(*CreateDailyReportCommentRequest))
	})
}
func _ReportService_GetComment_Handler(srv interface{}, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor) (any, error) {
	return handleUnary(ctx, dec, interceptor, ReportService_GetComment_FullMethodName, new(GetDailyReportCommentRequest), srv, func(ctx context.Context, req any) (any, error) {
		return srv.(ReportServiceServer).GetComment(ctx, req.(*GetDailyReportCommentRequest))
	})
}
func _ReportService_ListComments_Handler(srv interface{}, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor) (any, error) {
	return handleUnary(ctx, dec, interceptor, ReportService_ListComments_FullMethodName, new(ListDailyReportCommentsRequest), srv, func(ctx context.Context, req any) (any, error) {
		return srv.(ReportServiceServer).ListComments(ctx, req.(*ListDailyReportCommentsRequest))
	})
}
func _ReportService_UpdateComment_Handler(srv interface{}, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor) (any, error) {
	return handleUnary(ctx, dec, interceptor, ReportService_UpdateComment_FullMethodName, new(UpdateDailyReportCommentRequest), srv, func(ctx context.Context, req any) (any, error) {
		return srv.(ReportServiceServer).UpdateComment(ctx, req.(*UpdateDailyReportCommentRequest))
	})
}
func _ReportService_DeleteComment_Handler(srv interface{}, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor) (any, error) {
	return handleUnary(ctx, dec, interceptor, ReportService_DeleteComment_FullMethodName, new(DeleteDailyReportCommentRequest), srv, func(ctx context.Context, req any) (any, error) {
		return srv.(ReportServiceServer).DeleteComment(ctx, req.(*DeleteDailyReportCommentRequest))
	})
}

var ReportService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: ReportService_ServiceName,
	HandlerType: (*ReportServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{MethodName: "Ping", Handler: _ReportService_Ping_Handler},
		{MethodName: "CreateReport", Handler: _ReportService_CreateReport_Handler},
		{MethodName: "GetReport", Handler: _ReportService_GetReport_Handler},
		{MethodName: "ListReports", Handler: _ReportService_ListReports_Handler},
		{MethodName: "UpdateReport", Handler: _ReportService_UpdateReport_Handler},
		{MethodName: "DeleteReport", Handler: _ReportService_DeleteReport_Handler},
		{MethodName: "CreateEntry", Handler: _ReportService_CreateEntry_Handler},
		{MethodName: "GetEntry", Handler: _ReportService_GetEntry_Handler},
		{MethodName: "ListEntries", Handler: _ReportService_ListEntries_Handler},
		{MethodName: "UpdateEntry", Handler: _ReportService_UpdateEntry_Handler},
		{MethodName: "DeleteEntry", Handler: _ReportService_DeleteEntry_Handler},
		{MethodName: "CreateComment", Handler: _ReportService_CreateComment_Handler},
		{MethodName: "GetComment", Handler: _ReportService_GetComment_Handler},
		{MethodName: "ListComments", Handler: _ReportService_ListComments_Handler},
		{MethodName: "UpdateComment", Handler: _ReportService_UpdateComment_Handler},
		{MethodName: "DeleteComment", Handler: _ReportService_DeleteComment_Handler},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "report/report.proto",
}
