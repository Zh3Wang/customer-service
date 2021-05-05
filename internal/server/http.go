package server

import (
	"customerservice/internal/middleware"
	"customerservice/internal/pkg/setting"
	"customerservice/internal/service"
	"customerservice/internal/utl/log"
	"net/http"
	"strconv"
)

func ListenHttpServer() {

	http.HandleFunc("/register", middleware.AccessMiddleWare(service.Register))
	http.HandleFunc("/send", middleware.AccessMiddleWare(service.Send))
	http.HandleFunc("/broadcast", middleware.AccessMiddleWare(service.BroadCast))
	http.HandleFunc("/ws", Connect)

	// 客户端连接管理
	go Connection.Start()
	//监听消息通道
	go WriteMessage()
	go WatchConnectionStatus()

	//心跳
	go Heartbeat()

	//启动HTTP
	log.Info("启动HTTP服务：", setting.Config.Common.HttpPort)

	if err := http.ListenAndServe(":"+strconv.Itoa(setting.Config.Common.HttpPort), nil); err != nil {
		log.Panic("HTTP启动失败, ", err)
	}
}
