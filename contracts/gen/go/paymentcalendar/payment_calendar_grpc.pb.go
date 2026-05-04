package paymentcalendarpb

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	PaymentCalendarService_ServiceName                         = "paymentcalendar.PaymentCalendarService"
	PaymentCalendarService_Ping_FullMethodName                = "/paymentcalendar.PaymentCalendarService/Ping"
	PaymentCalendarService_CreatePayment_FullMethodName       = "/paymentcalendar.PaymentCalendarService/CreatePayment"
	PaymentCalendarService_GetPayment_FullMethodName          = "/paymentcalendar.PaymentCalendarService/GetPayment"
	PaymentCalendarService_ListPayments_FullMethodName        = "/paymentcalendar.PaymentCalendarService/ListPayments"
	PaymentCalendarService_UpdatePayment_FullMethodName       = "/paymentcalendar.PaymentCalendarService/UpdatePayment"
	PaymentCalendarService_DeletePayment_FullMethodName       = "/paymentcalendar.PaymentCalendarService/DeletePayment"
	PaymentCalendarService_MarkPaymentPaid_FullMethodName     = "/paymentcalendar.PaymentCalendarService/MarkPaymentPaid"
	PaymentCalendarService_GetProjectSummary_FullMethodName   = "/paymentcalendar.PaymentCalendarService/GetProjectSummary"
)

type PaymentCalendarServiceClient interface {
	Ping(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error)
	CreatePayment(ctx context.Context, in *CreatePaymentRequest, opts ...grpc.CallOption) (*Payment, error)
	GetPayment(ctx context.Context, in *GetPaymentRequest, opts ...grpc.CallOption) (*Payment, error)
	ListPayments(ctx context.Context, in *ListPaymentsRequest, opts ...grpc.CallOption) (*ListPaymentsResponse, error)
	UpdatePayment(ctx context.Context, in *UpdatePaymentRequest, opts ...grpc.CallOption) (*Payment, error)
	DeletePayment(ctx context.Context, in *DeletePaymentRequest, opts ...grpc.CallOption) (*Empty, error)
	MarkPaymentPaid(ctx context.Context, in *MarkPaymentPaidRequest, opts ...grpc.CallOption) (*Payment, error)
	GetProjectSummary(ctx context.Context, in *GetProjectSummaryRequest, opts ...grpc.CallOption) (*ProjectFinancialSummary, error)
}

type paymentCalendarServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewPaymentCalendarServiceClient(cc grpc.ClientConnInterface) PaymentCalendarServiceClient {
	return &paymentCalendarServiceClient{cc: cc}
}

func invoke[Resp any](ctx context.Context, cc grpc.ClientConnInterface, method string, in any, out Resp, opts ...grpc.CallOption) (Resp, error) {
	if err := cc.Invoke(ctx, method, in, out, opts...); err != nil {
		var zero Resp
		return zero, err
	}
	return out, nil
}

func (c *paymentCalendarServiceClient) Ping(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error) {
	return invoke(ctx, c.cc, PaymentCalendarService_Ping_FullMethodName, in, new(Empty), opts...)
}
func (c *paymentCalendarServiceClient) CreatePayment(ctx context.Context, in *CreatePaymentRequest, opts ...grpc.CallOption) (*Payment, error) {
	return invoke(ctx, c.cc, PaymentCalendarService_CreatePayment_FullMethodName, in, new(Payment), opts...)
}
func (c *paymentCalendarServiceClient) GetPayment(ctx context.Context, in *GetPaymentRequest, opts ...grpc.CallOption) (*Payment, error) {
	return invoke(ctx, c.cc, PaymentCalendarService_GetPayment_FullMethodName, in, new(Payment), opts...)
}
func (c *paymentCalendarServiceClient) ListPayments(ctx context.Context, in *ListPaymentsRequest, opts ...grpc.CallOption) (*ListPaymentsResponse, error) {
	return invoke(ctx, c.cc, PaymentCalendarService_ListPayments_FullMethodName, in, new(ListPaymentsResponse), opts...)
}
func (c *paymentCalendarServiceClient) UpdatePayment(ctx context.Context, in *UpdatePaymentRequest, opts ...grpc.CallOption) (*Payment, error) {
	return invoke(ctx, c.cc, PaymentCalendarService_UpdatePayment_FullMethodName, in, new(Payment), opts...)
}
func (c *paymentCalendarServiceClient) DeletePayment(ctx context.Context, in *DeletePaymentRequest, opts ...grpc.CallOption) (*Empty, error) {
	return invoke(ctx, c.cc, PaymentCalendarService_DeletePayment_FullMethodName, in, new(Empty), opts...)
}
func (c *paymentCalendarServiceClient) MarkPaymentPaid(ctx context.Context, in *MarkPaymentPaidRequest, opts ...grpc.CallOption) (*Payment, error) {
	return invoke(ctx, c.cc, PaymentCalendarService_MarkPaymentPaid_FullMethodName, in, new(Payment), opts...)
}
func (c *paymentCalendarServiceClient) GetProjectSummary(ctx context.Context, in *GetProjectSummaryRequest, opts ...grpc.CallOption) (*ProjectFinancialSummary, error) {
	return invoke(ctx, c.cc, PaymentCalendarService_GetProjectSummary_FullMethodName, in, new(ProjectFinancialSummary), opts...)
}

type PaymentCalendarServiceServer interface {
	Ping(context.Context, *Empty) (*Empty, error)
	CreatePayment(context.Context, *CreatePaymentRequest) (*Payment, error)
	GetPayment(context.Context, *GetPaymentRequest) (*Payment, error)
	ListPayments(context.Context, *ListPaymentsRequest) (*ListPaymentsResponse, error)
	UpdatePayment(context.Context, *UpdatePaymentRequest) (*Payment, error)
	DeletePayment(context.Context, *DeletePaymentRequest) (*Empty, error)
	MarkPaymentPaid(context.Context, *MarkPaymentPaidRequest) (*Payment, error)
	GetProjectSummary(context.Context, *GetProjectSummaryRequest) (*ProjectFinancialSummary, error)
}

type UnimplementedPaymentCalendarServiceServer struct{}

func (UnimplementedPaymentCalendarServiceServer) Ping(context.Context, *Empty) (*Empty, error) {
	return nil, status.Error(codes.Unimplemented, "method Ping not implemented")
}
func (UnimplementedPaymentCalendarServiceServer) CreatePayment(context.Context, *CreatePaymentRequest) (*Payment, error) {
	return nil, status.Error(codes.Unimplemented, "method CreatePayment not implemented")
}
func (UnimplementedPaymentCalendarServiceServer) GetPayment(context.Context, *GetPaymentRequest) (*Payment, error) {
	return nil, status.Error(codes.Unimplemented, "method GetPayment not implemented")
}
func (UnimplementedPaymentCalendarServiceServer) ListPayments(context.Context, *ListPaymentsRequest) (*ListPaymentsResponse, error) {
	return nil, status.Error(codes.Unimplemented, "method ListPayments not implemented")
}
func (UnimplementedPaymentCalendarServiceServer) UpdatePayment(context.Context, *UpdatePaymentRequest) (*Payment, error) {
	return nil, status.Error(codes.Unimplemented, "method UpdatePayment not implemented")
}
func (UnimplementedPaymentCalendarServiceServer) DeletePayment(context.Context, *DeletePaymentRequest) (*Empty, error) {
	return nil, status.Error(codes.Unimplemented, "method DeletePayment not implemented")
}
func (UnimplementedPaymentCalendarServiceServer) MarkPaymentPaid(context.Context, *MarkPaymentPaidRequest) (*Payment, error) {
	return nil, status.Error(codes.Unimplemented, "method MarkPaymentPaid not implemented")
}
func (UnimplementedPaymentCalendarServiceServer) GetProjectSummary(context.Context, *GetProjectSummaryRequest) (*ProjectFinancialSummary, error) {
	return nil, status.Error(codes.Unimplemented, "method GetProjectSummary not implemented")
}

func RegisterPaymentCalendarServiceServer(s grpc.ServiceRegistrar, srv PaymentCalendarServiceServer) {
	s.RegisterService(&PaymentCalendarService_ServiceDesc, srv)
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

func _PaymentCalendarService_Ping_Handler(srv interface{}, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor) (any, error) {
	return handleUnary(ctx, dec, interceptor, PaymentCalendarService_Ping_FullMethodName, new(Empty), srv, func(ctx context.Context, req any) (any, error) {
		return srv.(PaymentCalendarServiceServer).Ping(ctx, req.(*Empty))
	})
}
func _PaymentCalendarService_CreatePayment_Handler(srv interface{}, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor) (any, error) {
	return handleUnary(ctx, dec, interceptor, PaymentCalendarService_CreatePayment_FullMethodName, new(CreatePaymentRequest), srv, func(ctx context.Context, req any) (any, error) {
		return srv.(PaymentCalendarServiceServer).CreatePayment(ctx, req.(*CreatePaymentRequest))
	})
}
func _PaymentCalendarService_GetPayment_Handler(srv interface{}, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor) (any, error) {
	return handleUnary(ctx, dec, interceptor, PaymentCalendarService_GetPayment_FullMethodName, new(GetPaymentRequest), srv, func(ctx context.Context, req any) (any, error) {
		return srv.(PaymentCalendarServiceServer).GetPayment(ctx, req.(*GetPaymentRequest))
	})
}
func _PaymentCalendarService_ListPayments_Handler(srv interface{}, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor) (any, error) {
	return handleUnary(ctx, dec, interceptor, PaymentCalendarService_ListPayments_FullMethodName, new(ListPaymentsRequest), srv, func(ctx context.Context, req any) (any, error) {
		return srv.(PaymentCalendarServiceServer).ListPayments(ctx, req.(*ListPaymentsRequest))
	})
}
func _PaymentCalendarService_UpdatePayment_Handler(srv interface{}, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor) (any, error) {
	return handleUnary(ctx, dec, interceptor, PaymentCalendarService_UpdatePayment_FullMethodName, new(UpdatePaymentRequest), srv, func(ctx context.Context, req any) (any, error) {
		return srv.(PaymentCalendarServiceServer).UpdatePayment(ctx, req.(*UpdatePaymentRequest))
	})
}
func _PaymentCalendarService_DeletePayment_Handler(srv interface{}, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor) (any, error) {
	return handleUnary(ctx, dec, interceptor, PaymentCalendarService_DeletePayment_FullMethodName, new(DeletePaymentRequest), srv, func(ctx context.Context, req any) (any, error) {
		return srv.(PaymentCalendarServiceServer).DeletePayment(ctx, req.(*DeletePaymentRequest))
	})
}
func _PaymentCalendarService_MarkPaymentPaid_Handler(srv interface{}, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor) (any, error) {
	return handleUnary(ctx, dec, interceptor, PaymentCalendarService_MarkPaymentPaid_FullMethodName, new(MarkPaymentPaidRequest), srv, func(ctx context.Context, req any) (any, error) {
		return srv.(PaymentCalendarServiceServer).MarkPaymentPaid(ctx, req.(*MarkPaymentPaidRequest))
	})
}
func _PaymentCalendarService_GetProjectSummary_Handler(srv interface{}, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor) (any, error) {
	return handleUnary(ctx, dec, interceptor, PaymentCalendarService_GetProjectSummary_FullMethodName, new(GetProjectSummaryRequest), srv, func(ctx context.Context, req any) (any, error) {
		return srv.(PaymentCalendarServiceServer).GetProjectSummary(ctx, req.(*GetProjectSummaryRequest))
	})
}

var PaymentCalendarService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: PaymentCalendarService_ServiceName,
	HandlerType: (*PaymentCalendarServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{MethodName: "Ping", Handler: _PaymentCalendarService_Ping_Handler},
		{MethodName: "CreatePayment", Handler: _PaymentCalendarService_CreatePayment_Handler},
		{MethodName: "GetPayment", Handler: _PaymentCalendarService_GetPayment_Handler},
		{MethodName: "ListPayments", Handler: _PaymentCalendarService_ListPayments_Handler},
		{MethodName: "UpdatePayment", Handler: _PaymentCalendarService_UpdatePayment_Handler},
		{MethodName: "DeletePayment", Handler: _PaymentCalendarService_DeletePayment_Handler},
		{MethodName: "MarkPaymentPaid", Handler: _PaymentCalendarService_MarkPaymentPaid_Handler},
		{MethodName: "GetProjectSummary", Handler: _PaymentCalendarService_GetProjectSummary_Handler},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "paymentcalendar/payment_calendar.proto",
}
