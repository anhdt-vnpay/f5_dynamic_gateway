package services

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	pb "github.com/anhdt-vnpay/f5_dynamic_gateway/examples/types/ping"
	"github.com/anhdt-vnpay/f5_dynamic_gateway/services"
)

type pingService struct {
	pb.UnimplementedPingServiceServer
}

func NewPingService() pb.PingServiceServer {
	return &pingService{}
}

func (*pingService) PingMe(ctx context.Context, in *pb.PingRequest) (*pb.PingReply, error) {
	gatewayName := services.GetGatewayName(ctx)
	fmt.Println(fmt.Sprintf("gateway name : %s", *gatewayName))
	return &pb.PingReply{Message: "PONG"}, nil
}

func (*pingService) SlowPing(ctx context.Context, in *pb.PingRequest) (*pb.PingReply, error) {
	duration, err := strconv.Atoi(in.Delay)
	if err != nil {
		log.Fatalln("Parameters not valid: ", err)
	}
	time.Sleep(time.Duration(duration) * time.Second)

	return &pb.PingReply{Message: "SLOW PONG " + in.Delay + " seconds"}, nil
}
