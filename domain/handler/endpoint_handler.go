package handler

import "google.golang.org/grpc"

type IEndpointHandler interface {
	RegisterEndpoint(conn *grpc.ClientConn, endpoint string, serviceName string) error
	UnRegisterEndpoint(conn *grpc.ClientConn, endpoint string, serviceName string) error
}
