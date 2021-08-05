package services

import (
	"context"
	"crontab/global"
	"errors"
	"net"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
)

type RegisterNode struct {
	KvOp    clientv3.KV
	LeaseOp clientv3.Lease
	IPV4    string
}

func InitWorkerRegiste() {
	r := &RegisterNode{
		KvOp:    clientv3.NewKV(global.EtcdClient),
		LeaseOp: clientv3.NewLease(global.EtcdClient),
	}
	r.GetLocalIp()
	r.KeepOnLine()
}

func (r *RegisterNode) KeepOnLine() {
	var (
		ctx           context.Context
		cancelFunc    context.CancelFunc
		keepAliveChan <-chan *clientv3.LeaseKeepAliveResponse
		leaseResp     *clientv3.LeaseGrantResponse
		err           error
	)
	for {
		cancelFunc = nil
		leaseResp, err = r.LeaseOp.Grant(context.TODO(), 10)
		if err != nil {
			global.Lg.Info("create grant failed", zap.Error(err))
			goto RETRY
		}
		ctx, cancelFunc = context.WithCancel(context.TODO())
		//自动续租
		keepAliveChan, err = r.LeaseOp.KeepAlive(ctx, leaseResp.ID)
		if err != nil {
			goto RETRY
		}

		if _, err = r.KvOp.Put(ctx, global.WORKER_SAVE_PATH+r.IPV4, "", clientv3.WithLease(leaseResp.ID)); err != nil {
			global.Lg.Info("put work ip failed", zap.Error(err))
			goto RETRY
		}

		//处理续租应答
		for {
			if ch := <-keepAliveChan; ch == nil {
				goto RETRY
			}
		}

	RETRY:
		time.Sleep(1 * time.Second)
		if cancelFunc != nil {
			cancelFunc()
		}
	}

}

func (r *RegisterNode) GetLocalIp() error {
	//获取所有网卡
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		global.Lg.Info("get local ip failed", zap.Error(err))
		return err
	}
	for _, addr := range addrs {
		//这个是ip地址 ipv4或者ipv6
		if ipNet, isIpNet := addr.(*net.IPNet); isIpNet && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				r.IPV4 = ipNet.IP.String()
				return nil
			}
		}
	}

	return errors.New("not found ip addr")
}
