package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	handler "github.com/anhdt-vnpay/f5_dynamic_gateway/examples/handler"
	exService "github.com/anhdt-vnpay/f5_dynamic_gateway/examples/services"
	pb "github.com/anhdt-vnpay/f5_dynamic_gateway/examples/types/ping"
	lHandler "github.com/anhdt-vnpay/f5_dynamic_gateway/handler"
	dRuntime "github.com/anhdt-vnpay/f5_dynamic_gateway/runtime"
	"github.com/golang/glog"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
	endpointHandler := handler.NewExampleEndpointHandler2(ctx, mux)
	dbHandler, err := lHandler.NewDefaultDbHandler("./tmp/db")
	if err != nil {
		fmt.Println("dbhandler error " + err.Error())
	}
	connService := dRuntime.NewDefaultConnectionService(dbHandler, endpointHandler)
	gatewayService := dRuntime.NewGatewayServiceWithName("api1", connService, mux, nil)
	result := gatewayService.Handle(ctx, "/*")
	if result.Err != nil {
		fmt.Println("handler router error " + result.Err.Error())
	}
	router := result.MuxServer
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

	gServer := grpc.NewServer()
	pb.RegisterPingServiceServer(gServer, pingService)

	sigCh := make(chan os.Signal, 1)

	signal.Notify(sigCh, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	// signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		s := <-sigCh
		log.Printf("got signal %v, attempting graceful shutdown", s)
		// cancel()
		gServer.GracefulStop()
		// grpc.Stop() // leads to error while receiving stream response: rpc error: code = Unavailable desc = transport is closing
		wg.Done()
	}()

	if err := gServer.Serve(grpclistener); err != nil {
		msg := fmt.Sprintf("net.Listen: grpcServer error: %s", err)
		fmt.Println("grpcServer running ", msg)
		return
	}

	wg.Wait()
	log.Println("clean shutdown")
}

func GatewayEchoServer(port int) {
	e := echo.New()
	opts := []grpc.DialOption{
		grpc.WithAuthority("localhost:9000"),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	conn, _ := grpc.Dial("localhost:9000", opts...)
	defer conn.Close()
	// client := pb.NewPingServiceClient(conn)
	e.GET("/example2/ping", func(c echo.Context) error {
		res := &pb.PingReply{}
		err := conn.Invoke(c.Request().Context(), "/ping.PingService/PingMe", &pb.PingRequest{}, res)
		// res, err := client.PingMe(c.Request().Context(), &pb.PingRequest{})
		if err != nil {
			fmt.Println("error = " + err.Error())
			return err
		}
		// res := &pb.PingReply{
		// 	Message: "Pong",
		// }
		return c.JSON(http.StatusOK, res)
	})
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", port)))
}

func GatewayEchoServer2(port int) {

	address := fmt.Sprintf(":%d", port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		fmt.Printf("net.Listen: grpcServer error: %s \n", err)
		return
	}

	defer listener.Close()
	router := mux.NewRouter().StrictSlash(true)

	opts := []grpc.DialOption{
		grpc.WithAuthority("localhost:9000"),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	conn, _ := grpc.Dial("localhost:9000", opts...)
	defer conn.Close()

	router.HandleFunc("/example2/ping", func(w http.ResponseWriter, r *http.Request) {
		res := &pb.PingReply{}
		err := conn.Invoke(r.Context(), "/ping.PingService/PingMe", &pb.PingRequest{}, res)
		// res, err := client.PingMe(c.Request().Context(), &pb.PingRequest{})
		if err != nil {
			fmt.Println("error = " + err.Error())
			json.NewEncoder(w).Encode(err.Error())
		}
		json.NewEncoder(w).Encode(res)
	})
	server := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 60 * time.Second,
		Handler:      router}

	if err = server.Serve(listener); err != nil {
		glog.Fatal(err)
	}
}
