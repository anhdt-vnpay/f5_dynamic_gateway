package runtime

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	pb "github.com/anhdt-vnpay/f5_dynamic_gateway/types/registration"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc/metadata"
)

const (
	default_gateway_name = "gateway0"
	gateway_key          = "x-gateway-name"
)

type gatewayService struct {
	name       string
	apiService pb.ApiRegistrationServiceServer
	gwRouter   *runtime.ServeMux
}

func NewGatewayService(connService ConnectionService, mux *runtime.ServeMux, logger ILogger) GatewayService {
	return NewGatewayServiceWithName(default_gateway_name, connService, mux, logger)
}

func NewGatewayServiceWithName(gatewayName string, connService ConnectionService, mux *runtime.ServeMux, logger ILogger) GatewayService {
	apiService := newApiRegistrationServiceServer(connService, logger)
	return &gatewayService{
		name:       gatewayName,
		apiService: apiService,
		gwRouter:   mux,
	}
}

func (gateway *gatewayService) register(router *runtime.ServeMux) error {
	ctx := context.Background()
	return pb.RegisterApiRegistrationServiceHandlerServer(ctx, router, gateway.apiService)
}

func (gateway *gatewayService) setupGatewayHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set headers
		key := fmt.Sprintf("%s%s", runtime.MetadataHeaderPrefix, gateway_key)
		r.Header.Set(key, gateway.name)
		// Next
		next.ServeHTTP(w, r)
	})
}

func (gateway *gatewayService) Handle(ctx context.Context, pattern string) *HandlerResult {
	router := http.NewServeMux()
	gwRouter := gateway.gwRouter
	err := gateway.register(gwRouter)
	router.Handle(pattern, gateway.setupGatewayHeader(gwRouter))
	return &HandlerResult{
		MuxServer: router,
		Err:       err,
	}
}

func GetGatewayName(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}
	names := md.Get(gateway_key)
	if len(names) > 0 {
		gatewayName := strings.Join(names, ",")
		return gatewayName
	}
	return ""
}
