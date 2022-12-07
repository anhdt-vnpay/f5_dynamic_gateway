package runtime

import (
	"net/http"

	"google.golang.org/grpc"
)

type connectionInfo struct {
	serviceName string
	endpoints   []*EndpointInfo
}

type EndpointInfo struct {
	Id       string
	Endpoint string
	GrpcConn *grpc.ClientConn
}

type ConnectionConfig struct {
	MaxConn int
}

type HandlerResult struct {
	MuxServer *http.ServeMux
	Err       error
}
