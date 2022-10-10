package handler

type IDbHandler interface {
	LoadServiceName() ([]string, error)
	Load(serviceName string) ([]string, error)
	Store(endpoint string, serviceName string) error
	Delete(endpoint string, serviceName string) error
}
