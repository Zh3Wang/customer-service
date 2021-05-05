package service

import (
	"customerservice/internal/biz"
	"customerservice/internal/utl/log"
	"encoding/json"
	"net/http"
)

type broadcastParam struct {
	AppID        string `json:"app_id" validate:"required"`
	Data         string `json:"data" validate:"required"`
	FromClientId string `json:"from_client_id" validate:"required"`
}

// 发送广播接口
func BroadCast(writer http.ResponseWriter, request *http.Request) {
	defer request.Body.Close()
	var (
		err error
		in  broadcastParam
	)
	_ = json.NewDecoder(request.Body).Decode(&in)

	//验证参数
	if err = Validate(in); err != nil {
		log.Error("参数错误: ", err)
		Respond(writer, API_FAIL, err.Error(), []string{})
		return
	}

	log.Info("Receive a broadcast: ", in.Data)

	bc := biz.NewBroadCast(in.AppID, biz.Broadcast, in.Data, in.FromClientId)
	if err = bc.BroadCast(); err != nil {
		log.Error("broadcast failed :", err)
		Respond(writer, API_FAIL, "failed", err)
	} else {
		Respond(writer, API_SUCCESS, "success", []string{})
	}

}
