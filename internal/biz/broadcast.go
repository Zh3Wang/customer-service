package biz

import (
	"context"
	"customerservice/api/pb"
	"customerservice/internal/pkg/setting"
	"customerservice/internal/utl/log"
	"encoding/json"
	"errors"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
	"strconv"
)

/**
广播消息
*/
type BroadCastData struct {
	Type         MessageType `json:"type"`
	AppId        string      `json:"app_id"`
	Data         interface{} `json:"data"`
	FromClientId string      `json:"from_client_id"`
}

var (
	ErrAppIdNonExistent = errors.New("app id is non-existent")
)

func NewBroadCast(appId string, t MessageType, data interface{}, fromClientId string) *BroadCastData {
	d, _ := json.Marshal(data)
	return &BroadCastData{
		Type:         t,
		AppId:        appId,
		Data:         string(d),
		FromClientId: fromClientId,
	}
}

/*
	发送广播消息
*/
func (b *BroadCastData) BroadCast() error {
	log.WithFields(logrus.Fields{
		"AppId": b.AppId,
		"Data":  b.Data,
		"Type":  b.Type,
	}).Info("Sending a broadcast: ", b)

	//先给本地的client发送消息
	if err := b.BroadCastLocal(); err != nil {
		return err
	}

	// 通知集群中其他机器
	// 获取集群机器地址
	localhost := net.JoinHostPort(setting.Server.LocalHost, strconv.Itoa(setting.Config.Common.RPCPort))
	for _, addr := range setting.Server.ServerList {
		// 跳过本机
		if localhost == addr {
			continue
		}
		if err := b.BroadCastBygRPC(addr); err != nil {
			log.Error("Broadcast failed : ", addr, err)
		}
	}

	return nil
}

func (b *BroadCastData) BroadCastLocal() error {
	// 找到本机下app_id对应的所有client_ids
	clientIds, ok := CliManager.GetClientIdsByAppId(b.AppId)
	if !ok {
		// 该APP_id不存在，返回错误
		return ErrAppIdNonExistent
	}

	// 给所有客户端发消息
	for _, clientId := range clientIds {
		if b.FromClientId == clientId {
			continue
		}
		client, ok := CliManager.GetClient(clientId)
		if !ok {
			continue
		}
		// 发送消息
		md := NewMessage(b.AppId, client.ClientID, b.Type, b.Data)
		md.SendLocal()
	}
	return nil
}

/*
	通过gRPC通知集群广播消息
*/
func (b *BroadCastData) BroadCastBygRPC(addr string) error {
	gConn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return err
	}

	cli := pb.NewCustomerClient(gConn)

	var in = &pb.BroadcastData{
		AppId: b.AppId,
		Type:  int32(b.Type),
		Data:  b.Data.(string),
	}
	_, err = cli.BroadCast(context.TODO(), in)
	if err != nil {
		return err
	}
	return nil
}
