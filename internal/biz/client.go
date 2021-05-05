package biz

import (
	"github.com/gorilla/websocket"
	"sync"
	"time"
)


type Client struct {
	ClientID string
	AppID    string
	Conn     *websocket.Conn
	ConnLock sync.RWMutex
	Time     int64
	Data     interface{}
}

func NewClient(clientID, AppID string, conn *websocket.Conn, data interface{}) *Client {
	cli := &Client{
		ClientID: clientID,
		AppID:    AppID,
		Conn:     conn,
		Time:     time.Now().Unix(),
		Data:     data,
	}
	return cli
}