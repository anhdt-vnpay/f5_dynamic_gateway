package handler

import (
	"context"

	"github.com/anhdt-vnpay/f5_dynamic_gateway/examples/services"
	pb "github.com/anhdt-vnpay/f5_dynamic_gateway/examples/types/ping"
	dRuntime "github.com/anhdt-vnpay/f5_dynamic_gateway/runtime"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

type exampleEndpointHandler2 struct {
	ctx context.Context
	mux *runtime.ServeMux
}

func NewExampleEndpointHandler2(ctx context.Context, mux *runtime.ServeMux) dRuntime.IEndpointHandler {
	return &exampleEndpointHandler2{
		ctx: ctx,
		mux: mux,
	}
}

func (ex *exampleEndpointHandler2) RegisterEndpoint(serviceName string, endpointInfo *dRuntime.EndpointInfo) error {
	// pb.RegisterPingServiceHandler(ex.ctx, ex.mux, endpointInfo.GrpcConn)
	service := services.NewPingService()
	pb.RegisterPingServiceHandlerServer(ex.ctx, ex.mux, service)
	return nil
}

func (*exampleEndpointHandler2) UnRegisterEndpoint(serviceName string, endpointInfo *dRuntime.EndpointInfo) error {
	return nil
}
