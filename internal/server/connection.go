package server

import (
	"customerservice/internal/biz"
	"customerservice/internal/utl/log"
)

type ClientConnect struct {
	Disconnet chan *biz.Client
	Connect   chan *biz.Client
}

var Connection = NewClientConnect() // 管理者

func NewClientConnect() *ClientConnect {
	return &ClientConnect{
		Disconnet: make(chan *biz.Client, 100),
		Connect:   make(chan *biz.Client, 100),
	}
}

func (m *ClientConnect) Start() {
	defer func() {
		if err := recover();err != nil{
			log.Error(err)
		}
	}()
	for {
		select {
		case client := <-m.Disconnet:
			//下线通知
			m.EventDisconnet(client)
		case client := <-m.Connect:
			//上线通知
			m.EventConnect(client)
		}
	}
}

func (m *ClientConnect) EventDisconnet(client *biz.Client) {
	log.Info("Client offline: ", client.ClientID)

	//通知整个应用，该客户端下线
	data := map[string]interface{}{
		"msg":      "disconnected",
		"from_cid": client.ClientID,
	}
	bc := biz.NewBroadCast(client.AppID, biz.Offline, data, client.ClientID)
	_ = bc.BroadCast()

	//关闭连接
	_ = client.Conn.Close()
	//删除本地相关客户端信息
	biz.CliManager.DeleteClient(client)
	biz.CliManager.DelAppClient(client)
}

func (m *ClientConnect) EventConnect(client *biz.Client) {
	log.Info("Client online：", client.ClientID)

	// 保存客户端信息
	biz.CliManager.SaveClient(client)

	// 关联appid和clientid
	biz.CliManager.AddClient2AppClient(client.AppID, client)

	//通知整个应用
	data := make(map[string]string)
	data["msg"] = "connected"
	data["from_cid"] = 	client.ClientID
	bc := biz.NewBroadCast(client.AppID, biz.Online, data, client.ClientID)
	_ = bc.BroadCast()
}
