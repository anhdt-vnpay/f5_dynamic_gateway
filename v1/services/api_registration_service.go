package services

import (
	"context"

	ilogger "github.com/anhdt-vnpay/f5_dynamic_gateway/v1/domain/logger"
	iservices "github.com/anhdt-vnpay/f5_dynamic_gateway/v1/domain/services"
	pb "github.com/anhdt-vnpay/f5_dynamic_gateway/v1/types/registration"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

type apiRegistrationService struct {
	pb.UnimplementedApiRegistrationServiceServer
	connService iservices.ConnectionService
	logger      ilogger.ILogger
}

func NewApiRegistrationServiceServer(connService iservices.ConnectionService, logger ilogger.ILogger) pb.ApiRegistrationServiceServer {
	return &apiRegistrationService{
		connService: connService,
		logger:      logger,
	}
}

func (api *apiRegistrationService) Register(ctx context.Context, req *pb.ApiRegistrationRequest) (*pb.ApiRegistrationResponse, error) {
	endpoint := req.Endpoint
	serviceName := req.ServiceName
	err := api.connService.Add(endpoint, serviceName)
	if err != nil {
		api.logger.Errorf("Add new endpoint failed with error %s", err.Error())
		return nil, &runtime.HTTPStatusError{
			HTTPStatus: 400,
			Err:        err,
		}
	}
	return &pb.ApiRegistrationResponse{
		Code:    200,
		Message: "Register Endpoint Successful",
	}, nil
}
