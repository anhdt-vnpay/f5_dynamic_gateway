package runtime

import "context"

type IDbHandler interface {
	LoadServiceName() ([]string, error)
	Load(serviceName string) ([]string, error)
	Store(endpoint string, serviceName string) error
	Delete(endpoint string, serviceName string) error
}

type IEndpointHandler interface {
	RegisterEndpoint(serviceName string, endpointInfo *EndpointInfo) error
	UnRegisterEndpoint(serviceName string, endpointInfo *EndpointInfo) error
}

type ILogger interface {
	Fatalf(message string, args ...interface{})
	Infof(message string, args ...interface{})
	Debugf(message string, args ...interface{})
	Warnf(message string, args ...interface{})
	Errorf(message string, args ...interface{})
}

type ConnectionService interface {
	SetConfig(config ConnectionConfig, isForceApply bool) error
	Add(serviceName string, endpoint string) error
}

type GatewayService interface {
	Handle(ctx context.Context, pattern string) *HandlerResult
}
