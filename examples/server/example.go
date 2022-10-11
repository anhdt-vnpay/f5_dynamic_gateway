package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	handler "github.com/anhdt-vnpay/f5_dynamic_gateway/examples/handler"
	exService "github.com/anhdt-vnpay/f5_dynamic_gateway/examples/services"
	pb "github.com/anhdt-vnpay/f5_dynamic_gateway/examples/types/ping"
	lHandler "github.com/anhdt-vnpay/f5_dynamic_gateway/handler"
	"github.com/anhdt-vnpay/f5_dynamic_gateway/services"
	"github.com/golang/glog"
	"github.com/gorilla/handlers"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
)

func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set headers
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "*")
		w.Header().Set("Access-Control-Expose-Headers", "Content-Range")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		// Next
		next.ServeHTTP(w, r)
	})
}

func GatewayServer(port int) {

	address := fmt.Sprintf(":%d", port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		msg := fmt.Sprintf("net.Listen: grpcServer error: %s", err)
		fmt.Println("grpcServer running ", msg)
		return
	}

	defer listener.Close()

	if err := httpHandlers(listener); err != nil {
		glog.Fatal(err)
	}
}

func httpHandlers(listener net.Listener) error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	mux := runtime.NewServeMux(
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
			MarshalOptions: protojson.MarshalOptions{
				EmitUnpopulated: true,
			},
		}))
	endpointHandler := handler.NewExampleEndpointHandler(ctx, mux)
	dbHandler, err := lHandler.NewDefaultDbHandler("./tmp/db")
	if err != nil {
		fmt.Println("dbhandler error " + err.Error())
	}
	connService := services.NewDefaultConnectionService(dbHandler, endpointHandler)
	gatewayService := services.NewGatewayService(connService, mux, nil)
	router, err := gatewayService.Handle(ctx, "/")
	if err != nil {
		fmt.Println("handler router error " + err.Error())
	}
	cors_handler := CORS(router)
	log_handler := handlers.LoggingHandler(os.Stdout, cors_handler)

	server := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 60 * time.Second,
		Handler:      log_handler}
	return server.Serve(listener)
}

func GatewayGrpcServer(port int) {
	// Register gRPC server endpoint
	// Note: Make sure the gRPC server is running properly and accessible
	pingService := exService.NewPingService()
	fmt.Println("grpcServer running on port ", port)
	grpclistener, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		msg := fmt.Sprintf("net.Listen: grpcServer error: %s", err)
		fmt.Println("grpcServer running ", msg)
		return
	}
	s := grpc.NewServer()
	pb.RegisterPingServiceServer(s, pingService)

	if err := s.Serve(grpclistener); err != nil {
		msg := fmt.Sprintf("net.Listen: grpcServer error: %s", err)
		fmt.Println("grpcServer running ", msg)
		return
	}
}
