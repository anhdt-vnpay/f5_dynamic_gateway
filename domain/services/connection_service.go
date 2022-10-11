package services

type ConnectionConfig struct {
	MaxConn int
}

type ConnectionService interface {
	SetConfig(config ConnectionConfig, isForceApply bool) error
	Add(endpoint string, serviceName string) error
}
