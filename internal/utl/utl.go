package utl

import (
	"customerservice/internal/pkg/setting"
	"customerservice/internal/utl/crypto"
	uuid "github.com/satori/go.uuid"
	"net"
	"strconv"
	"strings"
)

//生成clientID，用户标识客户端连接了哪台机器
func CreateClientID() (string, error) {
	//对称加密IP和端口，当做clientId
	raw := []byte(setting.Server.LocalHost + ":" + strconv.Itoa(setting.Config.Common.RPCPort))
	str, err := crypto.Encrypt(raw, []byte(setting.Config.Common.CryptoKey))
	if err != nil {
		return "", err
	}
	return str, nil
}

func GetHostByClientId(client_id string) (string, bool, error) {
	addr, err := crypto.Decrypt(client_id, []byte(setting.Config.Common.CryptoKey))
	if err != nil {
		return "", false, err
	}

	localhost := GetIntranetIp()
	socket := net.JoinHostPort(localhost, strconv.Itoa(setting.Config.Common.RPCPort))
	islocal := false
	if addr == socket {
		islocal = true
	}

	return addr, islocal, nil
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

//GenUUID 生成uuid
func GenUUID() string {
	uuidFunc := uuid.NewV4()
	uuidStr := uuidFunc.String()
	uuidStr = strings.Replace(uuidStr, "-", "", -1)
	uuidByt := []rune(uuidStr)
	return string(uuidByt[8:24])
}

