package server

import (
	"customerservice/internal/biz"
	"customerservice/internal/utl"
	"customerservice/internal/utl/log"
	"github.com/gorilla/websocket"
	"net/http"
	"time"
)

const (
	ARGS_ERROR   = -1000
	APP_ID_ERROR = -1001
)

// 创建一个新的ws连接
func Connect(writer http.ResponseWriter, request *http.Request) {
	//创建ws实例
	var ws = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	conn, err := ws.Upgrade(writer, request, nil)
	if err != nil {
		log.Error("ws 连接失败：", err)
		http.NotFound(writer, request)
		return
	}

	appid := request.FormValue("app_id")
	if appid == "" {
		data := &PushData{
			Type: APP_ID_ERROR,
			Msg:  "appid参数错误",
			Data: []string{},
		}
		_ = data.Push(conn)
		_ = conn.Close()
		return
	}

	//生成一个clientID，用于标记服务器地址
	clientID, _ := utl.CreateClientID()

	//构造客户端信息
	data := request.FormValue("data") //附带的一些信息
	client := biz.NewClient(clientID, appid, conn, data)

	// 上线通知
	Connection.Connect <- client

	// 返回连接成功
	d := &PushData{
		Type:     biz.Online,
		Msg:      "connected",
		Data:     []string{},
		ClientId: clientID,
	}
	client.ConnLock.Lock()
	_ = d.Push(client.Conn)
	client.ConnLock.Unlock()

}

type PushData struct {
	MessageId string          `json:"messageId"`
	Type      biz.MessageType `json:"code"`
	Msg       string          `json:"msg"`
	Data      interface{}     `json:"data"`
	ClientId  string          `json:"client_id"`
}

// 推送消息至客户端
func (p *PushData) Push(conn *websocket.Conn) error {
	return conn.WriteJSON(p)
}

// 从消息通道中读取消息，并推送给客户端
func WriteMessage() {
	for {
		select {
		case d := <-biz.MessageChan:
			client, isExist := biz.CliManager.GetClient(d.ClientId)
			if isExist {
				data := &PushData{
					MessageId: d.MessageId,
					Type:      d.Type,
					Msg:       "success",
					Data:      d.Data,
					ClientId:  d.ClientId,
				}
				client.ConnLock.Lock()
				err := data.Push(client.Conn)
				client.ConnLock.Unlock()
				if err != nil {
					log.Error("Disconneted:", d.ClientId)
					Connection.Disconnet <- client
				}
			}
		}
	}
}

// 定时发送心跳
func Heartbeat() {
	for {
		clients := biz.CliManager.GetAllClient()
		for _, v := range clients {
			data := &PushData{
				Type: biz.Heartbeat,
				Msg:  "success",
				Data: "ping",
			}
			v.ConnLock.Lock()
			err := data.Push(v.Conn)
			v.ConnLock.Unlock()
			if err != nil {
				// ping不通的标记为下线
				log.Error("can't ping:", v.ClientID)
				Connection.Disconnet <- v
			}
		}
		time.Sleep(time.Second * 10)
	}
}

/**
	监测长连接状态，是否中断
 */
func WatchConnectionStatus(){
	for  {
		clients := biz.CliManager.GetAllClient()
		for _, v := range clients {
			_, _, err := v.Conn.ReadMessage()
			if err != nil {
				//出错，代表连接断开，客户端下线了
				Connection.Disconnet <- v
				log.Error("Client has disconnected: ", v.ClientID)
				break
			}
		}
		time.Sleep(time.Second*5)
	}
}
