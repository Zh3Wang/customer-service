package service

import (
	"context"
	"customerservice/api/pb"
	"customerservice/internal/biz"
	"customerservice/internal/utl/log"
)

type CustomerService struct {
	pb.UnimplementedCustomerServer
	md *biz.MsgData
	bc *biz.BroadCastData
}

// gRPC接口--发送消息
func (s *CustomerService) SendMessage(ctx context.Context, in *pb.Message) (*pb.Reply, error) {
	t := biz.MessageType(in.Type)
	s.md = biz.NewMessage(in.AppId, in.ClientId, t, in.Data)
	log.Info("Received a gRPC message：", in.Data)
	_ = s.md.Send2Client()
	rp := &pb.Reply{
		Code: 1,
		Msg:  "ok",
	}
	return rp, nil
}

// gRPC接口--广播消息
func (s *CustomerService) BroadCast(ctx context.Context, in *pb.BroadcastData) (*pb.Reply, error) {
	t := biz.MessageType(in.Type)
	s.bc = biz.NewBroadCast(in.AppId, t, in.Data, in.FromClientId)
	_ = s.bc.BroadCastLocal()
	rp := &pb.Reply{
		Code: 1,
		Msg:  "ok",
	}
	return rp, nil
}
