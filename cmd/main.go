package main

import (
	"customerservice/internal/pkg/etcd"
	"customerservice/internal/pkg/setting"
	"customerservice/internal/server"
	"customerservice/internal/utl/log"
	"fmt"
	"strconv"
)

func main() {
	defer func() {
		if err := recover();err != nil{
			log.Error(err)
		}
	}()
	Init()
	select {}
}

func Init() {
	// 配置
	setting.Init()
	fmt.Println("配置启动成功")
	// 日志
	log.InitLog()
	fmt.Println("日志启动成功")
	// 注册服务
	regService()
	// 启动gRPC服务
	g := server.NewGRPCServer(strconv.Itoa(setting.Config.Common.RPCPort))
	go g.ListenRPCServer()
	// http服务
	server.ListenHttpServer()
}

func regService() {
	err := etcd.NewServiceReg(setting.Config.Etcd.Addr, 2)
	if err != nil {
		log.Panic("init etcd error: ", err)
	}
}
