package biz

import (
	"sync"
)

type ClientManager struct {
	ClientMap     map[string]*Client
	ClientMapLock sync.RWMutex

	// 把每个客户端信息跟应用id关联起来
	AppClientMap     map[string][]string
	AppClientMapLock sync.RWMutex
}

var CliManager = NewClientManager() // 管理者

func NewClientManager() *ClientManager {
	return &ClientManager{
		ClientMap:    make(map[string]*Client, 100),
		AppClientMap: make(map[string][]string, 100),
	}
}

//保存客户端信息
func (m *ClientManager) SaveClient(cli *Client) {
	m.ClientMapLock.Lock()
	defer m.ClientMapLock.Unlock()
	m.ClientMap[cli.ClientID] = cli
}

//删除客户端
func (m *ClientManager) DeleteClient(cli *Client) {
	m.ClientMapLock.Lock()
	defer m.ClientMapLock.Unlock()
	delete(m.ClientMap, cli.ClientID)
}

//获取客户端信息
func (m *ClientManager) GetClient(clientID string) (*Client, bool) {
	m.ClientMapLock.RLock()
	defer m.ClientMapLock.RUnlock()

	res, ok := m.ClientMap[clientID]
	if !ok {
		return nil, ok
	} else {
		return res, ok
	}

}

// 获取本机所有client信息
func (m *ClientManager) GetAllClient() map[string]*Client {
	m.ClientMapLock.RLock()
	defer m.ClientMapLock.RUnlock()

	return m.ClientMap
}

// 添加到指定APP_ID下的客户端列表
// app_id : []client_id
func (m *ClientManager) AddClient2AppClient(appId string, client *Client) {
	m.AppClientMapLock.Lock()
	defer m.AppClientMapLock.Unlock()
	m.AppClientMap[appId] = append(m.AppClientMap[appId], client.ClientID)
}

// 删除指定APP_ID下的指定客户端
// app_id : []client_id
func (m *ClientManager) DelAppClient(client *Client) {
	m.AppClientMapLock.Lock()
	defer m.AppClientMapLock.Unlock()

	for index, clientId := range m.AppClientMap[client.AppID] {
		if clientId == client.ClientID {
			m.AppClientMap[client.AppID] = append(m.AppClientMap[client.AppID][:index], m.AppClientMap[client.AppID][index+1:]...)
		}
	}
}

// 获取指定app_id下的所有客户端
func (m *ClientManager) GetClientIdsByAppId(appId string) ([]string, bool) {
	m.AppClientMapLock.RLock()
	defer m.AppClientMapLock.RUnlock()

	r, ok := m.AppClientMap[appId]
	return r, ok
}
