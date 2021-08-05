package services

import (
	"context"
	"crontab/global"
	"crontab/models"
	"strings"

	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type JobKill struct {
	Watcher    clientv3.Watcher
	BuildEvent *JobBuildEvent
	Schedule   *JobSchedule
}

func (k *JobKill) watchKiller() {
	go func() {
		watchChan := k.Watcher.Watch(context.TODO(), global.CRON_KILL_PATH, clientv3.WithPrefix())
		for watch := range watchChan {
			for _, event := range watch.Events {
				key := strings.TrimPrefix(string(event.Kv.Key), global.CRON_KILL_PATH)
				//获取到投递的强杀任务
				if event.Type == mvccpb.PUT {
					job := &models.Job{}
					jobEvent := k.BuildEvent.buildEventJob(global.JOB_EVENT_TYPE_KILL, job, key)
					k.Schedule.PushJob(jobEvent)
				}
			}
		}
	}()
}
