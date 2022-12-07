module github.com/anhdt-vnpay/f5_dynamic_gateway/examples

go 1.16

replace github.com/anhdt-vnpay/f5_dynamic_gateway => ../

require (
	github.com/anhdt-vnpay/f5_dynamic_gateway v0.0.0-00010101000000-000000000000
	github.com/golang/glog v1.0.0
	github.com/gorilla/handlers v1.5.1
	github.com/gorilla/mux v1.8.0
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.11.3
	github.com/labstack/echo/v4 v4.9.1
	github.com/spf13/cobra v1.5.0
	google.golang.org/grpc v1.50.0
	google.golang.org/protobuf v1.28.1
)
