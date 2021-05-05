package service

import (
	"customerservice/internal/biz"
	"customerservice/internal/utl/log"
	"encoding/json"
	"fmt"
	"net/http"
)

type SendParam struct {
	AppId    string      `json:"app_id" validate:"required"`
	ClientId string      `json:"client_id" validate:"required"`
	Data     interface{} `json:"data" validate:"required"`
	UserId   int         `json:"user_id"`
}

// 通过HTTP接口发送消息
func Send(writer http.ResponseWriter, request *http.Request) {
	defer request.Body.Close()
	var in SendParam
	_ = json.NewDecoder(request.Body).Decode(&in)

	// 验证参数
	if err := Validate(in); err != nil {
		fmt.Println("参数错误：", err)
		Respond(writer, API_FAIL, err.Error(), []string{})
		return
	}

	log.Info("发送一条http消息：", in.Data)

	// 发送消息
	md := biz.NewMessage(in.AppId, in.ClientId, biz.Common, in.Data)
	md.Send2Client()

	Respond(writer, API_SUCCESS, "success", []string{})
}
