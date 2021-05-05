package biz

import (
	"customerservice/internal/pkg/etcd"
	"encoding/json"
	"time"
)

const (
	KEYAPPID = "appid"
)

type AppInfo struct {
	AppID   string `json:"app_id"`
	RegTime int64  `json:"reg_time"`
}

//保存APPID到ETCD
func Save(appID string) error {
	t := time.Now().Unix()
	appInfo := AppInfo{
		AppID:   appID,
		RegTime: t,
	}
	appInfoJson, _ := json.Marshal(appInfo)
	err := etcd.Services.Put(KEYAPPID+appID, string(appInfoJson))
	if err != nil{
		return err
	}
	return nil
}
