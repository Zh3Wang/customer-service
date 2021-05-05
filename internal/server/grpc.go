package server

import (
	"customerservice/api/pb"
	"customerservice/internal/service"
	"customerservice/internal/utl/log"
	"google.golang.org/grpc"
	"net"
)

type GRPCServerInfo struct {
	Port string
}

func NewGRPCServer(port string) *GRPCServerInfo {
	return &GRPCServerInfo{
		Port: port,
	}
}

//监听rpc服务
func (g *GRPCServerInfo) ListenRPCServer() {
	lis, err := net.Listen("tcp", ":"+g.Port)
	if err != nil {
		log.Error("启动gRPC失败：", err, g.Port)
		return
	}

	s := grpc.NewServer()
	pb.RegisterCustomerServer(s, &service.CustomerService{})

	log.Info("启动gRP服务，监听端口: " + g.Port)
	if err := s.Serve(lis); err != nil {
		log.Error("注册gRPC失败：", err)
		return
	}
}
