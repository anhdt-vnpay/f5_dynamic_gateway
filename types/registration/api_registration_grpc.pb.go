// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.6.1
// source: registration/api_registration.proto

package registration

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// ApiRegistrationServiceClient is the client API for ApiRegistrationService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ApiRegistrationServiceClient interface {
	// Register
	Register(ctx context.Context, in *ApiRegistrationRequest, opts ...grpc.CallOption) (*ApiRegistrationResponse, error)
}

type apiRegistrationServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewApiRegistrationServiceClient(cc grpc.ClientConnInterface) ApiRegistrationServiceClient {
	return &apiRegistrationServiceClient{cc}
}

func (c *apiRegistrationServiceClient) Register(ctx context.Context, in *ApiRegistrationRequest, opts ...grpc.CallOption) (*ApiRegistrationResponse, error) {
	out := new(ApiRegistrationResponse)
	err := c.cc.Invoke(ctx, "/registration.ApiRegistrationService/Register", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ApiRegistrationServiceServer is the server API for ApiRegistrationService service.
// All implementations must embed UnimplementedApiRegistrationServiceServer
// for forward compatibility
type ApiRegistrationServiceServer interface {
	// Register
	Register(context.Context, *ApiRegistrationRequest) (*ApiRegistrationResponse, error)
	mustEmbedUnimplementedApiRegistrationServiceServer()
}

// UnimplementedApiRegistrationServiceServer must be embedded to have forward compatible implementations.
type UnimplementedApiRegistrationServiceServer struct {
}

func (UnimplementedApiRegistrationServiceServer) Register(context.Context, *ApiRegistrationRequest) (*ApiRegistrationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Register not implemented")
}
func (UnimplementedApiRegistrationServiceServer) mustEmbedUnimplementedApiRegistrationServiceServer() {
}

// UnsafeApiRegistrationServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ApiRegistrationServiceServer will
// result in compilation errors.
type UnsafeApiRegistrationServiceServer interface {
	mustEmbedUnimplementedApiRegistrationServiceServer()
}

func RegisterApiRegistrationServiceServer(s grpc.ServiceRegistrar, srv ApiRegistrationServiceServer) {
	s.RegisterService(&ApiRegistrationService_ServiceDesc, srv)
}

func _ApiRegistrationService_Register_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ApiRegistrationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ApiRegistrationServiceServer).Register(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/registration.ApiRegistrationService/Register",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ApiRegistrationServiceServer).Register(ctx, req.(*ApiRegistrationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ApiRegistrationService_ServiceDesc is the grpc.ServiceDesc for ApiRegistrationService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ApiRegistrationService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "registration.ApiRegistrationService",
	HandlerType: (*ApiRegistrationServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Register",
			Handler:    _ApiRegistrationService_Register_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "registration/api_registration.proto",
}
