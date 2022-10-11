package handler

import (
	"context"

	"github.com/anhdt-vnpay/f5_dynamic_gateway/v1/domain/handler"
	pb "github.com/anhdt-vnpay/f5_dynamic_gateway/examples/types/ping"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type exampleEndpointHandler struct {
	ctx context.Context
	mux *runtime.ServeMux
}

func NewExampleEndpointHandler(ctx context.Context, mux *runtime.ServeMux) handler.IEndpointHandler {
	return &exampleEndpointHandler{
		ctx: ctx,
		mux: mux,
	}
}

func (ex *exampleEndpointHandler) RegisterEndpoint(conn *grpc.ClientConn, endpoint string, serviceName string) error {
	pb.RegisterPingServiceHandler(ex.ctx, ex.mux, conn)
	return nil
}

func (*exampleEndpointHandler) UnRegisterEndpoint(conn *grpc.ClientConn, endpoint string, serviceName string) error {
	return nil
}
