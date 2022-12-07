package handler

import (
	"context"

	pb "github.com/anhdt-vnpay/f5_dynamic_gateway/examples/types/ping"
	dRuntime "github.com/anhdt-vnpay/f5_dynamic_gateway/runtime"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

type exampleEndpointHandler struct {
	ctx context.Context
	mux *runtime.ServeMux
}

func NewExampleEndpointHandler(ctx context.Context, mux *runtime.ServeMux) dRuntime.IEndpointHandler {
	return &exampleEndpointHandler{
		ctx: ctx,
		mux: mux,
	}
}

func (ex *exampleEndpointHandler) RegisterEndpoint(serviceName string, endpointInfo *dRuntime.EndpointInfo) error {
	pb.RegisterPingServiceHandler(ex.ctx, ex.mux, endpointInfo.GrpcConn)
	return nil
}

func (*exampleEndpointHandler) UnRegisterEndpoint(serviceName string, endpointInfo *dRuntime.EndpointInfo) error {
	return nil
}
