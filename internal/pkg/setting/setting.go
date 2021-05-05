package setting

import (
	"flag"
	"github.com/spf13/viper"
	"log"
	"net"
	"sync"
)

// 全局配置信息
type GlobalConfig struct {
	Common  common
	Etcd    etcd
	Mysql   mysql
	LogConf logConf
	AppMode string
}

// 通用配置
type common struct {
	HttpPort     int
	RPCPort      int
	CryptoKey    string
	LogPath      string
	LogSave      int
	ServerPrefix string
}

// etcd地址配置
type etcd struct {
	Addr []string
}

// mysql配置
type mysql struct {
	Ip       string
	Port     int
	Root     string
	Password string
}

type logConf struct {
	LogPath  string
	LogSave  int
	HideKeys bool
	Level    string
}

type ServerConfig struct {
	LocalHost      string
	ServerList     map[string]string
	ServerListLock sync.RWMutex
}

var (
	Config = &GlobalConfig{}
	Server = &ServerConfig{}
)

func Init() {
	viper.SetConfigName("conf")
	viper.SetConfigType("toml")
	viper.AddConfigPath("../configs")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("read config error: ", err)
	}

	// 解析到结构体
	_ = viper.Unmarshal(Config)

	//获取本机IP
	Server.LocalHost = GetIntranetIp()
	Server.ServerList = make(map[string]string)

	HttpPort := flag.Int("http", Config.Common.HttpPort, "http port")
	RPCPort := flag.Int("rpc", Config.Common.RPCPort, "rpc port")
	flag.Parse()
	Config.Common.HttpPort = *HttpPort
	Config.Common.RPCPort = *RPCPort
}

//获取本机内网IP
func GetIntranetIp() string {
	addrs, _ := net.InterfaceAddrs()

	for _, addr := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}

	return ""
}
