package biz

import (
	"context"
	"customerservice/api/pb"
	"customerservice/internal/utl"
	"customerservice/internal/utl/log"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type MsgData struct {
	AppId     string      `json:"app_id"`
	ClientId  string      `json:"client_id"`
	Type      MessageType `json:"type"`
	Data      interface{} `json:"data"`
	MessageId string      `json:"message_id"`
}

// 消息通道
var MessageChan = make(chan *MsgData, 10000)

func NewMessage(appId, clientId string, t MessageType, data interface{}) *MsgData {
	return &MsgData{
		AppId:    appId,
		ClientId: clientId,
		Type:     t,
		Data:     data,
	}
}

//单发消息
func (md *MsgData) Send2Client() error {
	var (
		addr    string
		isLocal bool
		err     error
	)
	// 判断接收方是否在本机
	addr, isLocal, err = utl.GetHostByClientId(md.ClientId)
	if err != nil {
		log.Error("clientId is invalid: ", err)
		return err
	}
	if isLocal {
		md.SendLocal()
	} else {
		//通过grpc发送到对应机器
		log.Info("sending a gRPC message:", md.Data, addr)
		err = md.SendBygRPC(addr)
		if err != nil {
			log.Error("failed to send a gRPC message：", err)
			return err
		}
		log.Info("send a gRPC message success.")
	}
	return nil
}

// 发送到本机通道中
func (md *MsgData) SendLocal() {
	if md.Type != Heartbeat{
		log.WithFields(logrus.Fields{
			"app_id":    md.AppId,
			"client_id": md.ClientId,
		}).Debug("sending a local message:", md.Data)
	}


	//生成一个message id
	mid := utl.GenUUID()
	md.MessageId = mid

	//持久化
	//....

	//塞入通道，goroutine异步处理推送消息
	MessageChan <- md
}

// 通过gRPC发送消息
func (md *MsgData) SendBygRPC(addr string) error {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return err
	}

	cli := pb.NewCustomerClient(conn)
	in := pb.Message{
		AppId:    md.AppId,
		ClientId: md.ClientId,
		Type:     int32(md.Type),
		Data:     md.Data.(string),
	}

	_, err = cli.SendMessage(context.TODO(), &in)
	if err != nil {
		return err
	}
	return nil
}
