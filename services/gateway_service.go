package services

import (
	"context"
	"net/http"

	"github.com/anhdt-vnpay/f5_dynamic_gateway/domain/logger"
	"github.com/anhdt-vnpay/f5_dynamic_gateway/domain/services"
	pb "github.com/anhdt-vnpay/f5_dynamic_gateway/types/registration"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

type gatewayService struct {
	apiService pb.ApiRegistrationServiceServer
	gwRouter   *runtime.ServeMux
}

func NewGatewayService(connService services.ConnectionService, mux *runtime.ServeMux, logger logger.ILogger) services.GatewayService {
	apiService := NewApiRegistrationServiceServer(connService, logger)
	return &gatewayService{
		apiService: apiService,
		gwRouter:   mux,
	}
}

func (gateway *gatewayService) register(router *runtime.ServeMux) error {
	ctx := context.Background()
	return pb.RegisterApiRegistrationServiceHandlerServer(ctx, router, gateway.apiService)
}

func (gateway *gatewayService) Handle(ctx context.Context, pattern string) (*http.ServeMux, error) {
	router := http.NewServeMux()
	gwRouter := gateway.gwRouter
	err := gateway.register(gwRouter)
	router.Handle(pattern, gwRouter)
	return router, err
}
