package etcd

import (
	"context"
	"customerservice/internal/pkg/setting"
	"customerservice/internal/utl/log"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"net"
	"strconv"
	"time"
)

type ServiceData struct {
	Client             *clientv3.Client
	LeaseGrantResp     *clientv3.LeaseGrantResponse
	LeaseKeepAliveChan <-chan *clientv3.LeaseKeepAliveResponse
	CancelFunc         context.CancelFunc
	Lease              clientv3.Lease
}

var Services = ServiceData{}

// 注册服务到etcd
// 创建租约，监测服务存活情况
func NewServiceReg(addr []string, keepAliveTime int64) error {
	// 建立连接
	cfg := clientv3.Config{
		Endpoints:   addr,
		DialTimeout: time.Second * 2,
	}
	client, err := clientv3.New(cfg)
	if err != nil {
		return err
	}
	// 创建租约
	Services.Client = client
	err = Services.createLease(keepAliveTime)

	go Services.listenLeaseKeepAlive()

	// 获取本机IP + 端口号
	host := net.JoinHostPort(setting.Server.LocalHost, strconv.Itoa(setting.Config.Common.RPCPort))

	// 注册服务
	err = Services.PutServices(setting.Config.Common.ServerPrefix+host, host)
	if err != nil {
		log.Panic("注册服务失败：", err)
	}

	// 获取所有服务节点地址
	err = Services.GetServices(setting.Config.Common.ServerPrefix)
	if err != nil {
		log.Panic("初始化服务失败：", err)
	}
	return nil
}

//建立租约
func (s *ServiceData) createLease(keepAliveTime int64) error {
	lease := clientv3.NewLease(s.Client)

	ctx := context.TODO()
	leaseGrantResp, err := lease.Grant(ctx, keepAliveTime)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())
	leaseKeepAliveChan, err := lease.KeepAlive(ctx, leaseGrantResp.ID)
	if err != nil {
		return err
	}

	// 保存租约信息
	s.LeaseGrantResp = leaseGrantResp
	s.Lease = lease
	s.LeaseKeepAliveChan = leaseKeepAliveChan
	s.CancelFunc = cancel
	return nil
}

// 监听续租情况
func (s *ServiceData) listenLeaseKeepAlive() {
	for {
		select {
		case leaseResp := <-s.LeaseKeepAliveChan:
			if leaseResp == nil {
				log.Info("续租失败")
				return
			}
		}
	}
}

// 保存服务地址到etcd中
func (s *ServiceData) PutServices(k, v string) error {
	//这个key绑定了租约ID，租约过期后也会触发删除这个key
	_, err := s.Client.Put(context.TODO(), k, v, clientv3.WithLease(s.LeaseGrantResp.ID))
	if err != nil {
		return err
	}
	return nil
}

// 获取etcd中所有服务器host
func (s *ServiceData) GetServices(prefix string) error {
	resp, err := s.Client.Get(context.TODO(), prefix, clientv3.WithPrefix())
	if err != nil {
		return err
	}

	for _, v := range resp.Kvs {
		// 保存所有服务器地址
		SetServerListMap(string(v.Key), string(v.Value))
	}

	//监听服务地址变化
	go s.watcher()

	return nil
}

// 更新服务地址列表
func SetServerListMap(k, v string) {
	setting.Server.ServerListLock.Lock()
	defer setting.Server.ServerListLock.Unlock()
	setting.Server.ServerList[k] = v
	log.Info("发现服务地址: ", v)
}

// 删除服务地址
func DelServerListMap(k string) {
	setting.Server.ServerListLock.Lock()
	defer setting.Server.ServerListLock.Unlock()
	delete(setting.Server.ServerList, k)
	log.Info("删除服务地址: ", k)
}

func (s *ServiceData) watcher() {
	watchChan := s.Client.Watch(context.TODO(), setting.Config.Common.ServerPrefix, clientv3.WithPrefix())
	for {
		select {
		case resp := <-watchChan:
			for _, v := range resp.Events {
				switch v.Type {
				case mvccpb.DELETE:
					DelServerListMap(string(v.Kv.Key))
				case mvccpb.PUT:
					SetServerListMap(string(v.Kv.Key), string(v.Kv.Value))
				}
			}
		}
	}
}
