package initialize

import (
	"crontab/global"
	"time"

	"github.com/fatih/color"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func InitEtcd() {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   global.Settings.EtcdInfo.Hosts,
		DialTimeout: time.Duration(global.Settings.EtcdInfo.DialTimeout) * time.Microsecond,
	})
	if err != nil {
		color.Red("[initEtcd failed] 初始化etcd连接失败")
		color.Yellow(err.Error())
		panic(err)
	}
	global.EtcdClient = client
}
