package service

import (
	"customerservice/internal/biz"
	"customerservice/internal/utl/log"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
)

type registerParam struct {
	AppID string `json:"app_id" validate:"required"`
}

func Register(writer http.ResponseWriter, request *http.Request) {
	defer request.Body.Close()
	var in registerParam
	_ = json.NewDecoder(request.Body).Decode(&in)

	// 验证参数
	if err := Validate(in); err != nil {
		fmt.Println(err)
		Respond(writer, API_FAIL, err.Error(), []string{})
		return
	}

	// 保存APPID
	if err := biz.Save(in.AppID); err != nil {
		log.WithFields(logrus.Fields{
			"app_id": in.AppID,
		}).Error("app id 保存失败:", err)
		Respond(writer, API_FAIL, "success", err)
	}else{
		Respond(writer, API_SUCCESS, "success", "")
	}
}
