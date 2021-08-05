package services

import (
	"context"
	"crontab/global"
	"crontab/models"
	"crontab/utils"
	"encoding/json"
	"strings"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
)

type JobWatchServiceInterface interface {
	Watch(ctx context.Context)
}

type jobWatchService struct {
	KvOp       clientv3.KV
	WatcherOp  clientv3.Watcher
	SavePrefix string
	KillPrefix string
	Schedule   *JobSchedule
	BuildEvent *JobBuildEvent
}

type JobBuildEvent struct {
}

func (b *JobBuildEvent) buildEventJob(event global.JOB_EVENT_TYPE, job *models.Job, jobId string) *JobEvent {
	return &JobEvent{
		EventType: event,
		Job:       job,
		JobId:     jobId,
	}
}

//统一调度入口
func InitWatchAndScheduleAndExecJob(ctx context.Context) {
	kvOp := clientv3.NewKV(global.EtcdClient)
	leaseOp := clientv3.NewLease(global.EtcdClient)
	watcherOp := clientv3.NewWatcher(global.EtcdClient)
	buildEvent := &JobBuildEvent{}
	logCollection := utils.GetMongoDbDatabase().Collection("jobs")
	logSink := &LogSink{
		LogCollection: logCollection,
		LogChan:       make(chan *models.JobLog, 10000),
		AutoCommit:    make(chan *LogBatch, 1000),
	}
	executor := &Executor{
		KvOp:     kvOp,
		LeaseOp:  leaseOp,
		LockPath: global.CRON_LOCK_PATH,
	}
	schedule := &JobSchedule{
		JobEventChan:    make(chan *JobEvent, 10000),
		JobPlanTable:    make(map[string]*JobSchedulePlan),
		JobExecuteTable: make(map[string]*JobExecuteInfo),
		JobExecutor:     executor,
		JobResultChan:   make(chan *JobExecuteResult, 10000),
		LogSink:         logSink,
	}
	go schedule.schedulerLoop()
	go logSink.WriteLoop()
	jobWatcher := &jobWatchService{
		KvOp:       kvOp,
		WatcherOp:  watcherOp,
		SavePrefix: global.CRON_SAVE_PATH,
		KillPrefix: global.CRON_KILL_PATH,
		Schedule:   schedule,
		BuildEvent: buildEvent,
	}
	jobKillWatcher := &JobKill{
		Watcher:    clientv3.NewWatcher(global.EtcdClient),
		Schedule:   schedule,
		BuildEvent: buildEvent,
	}
	//监测任务强杀
	jobKillWatcher.watchKiller()
	//监测任务投递
	jobWatcher.Watch(ctx)
	//监测工作节点注册
	InitWorkerRegiste()
}

func (j *jobWatchService) Watch(ctx context.Context) {
	resp, err := j.KvOp.Get(ctx, j.SavePrefix, clientv3.WithPrefix())
	if err != nil {
		global.Lg.Error("get jobList failed", zap.Error(err))
	}
	for _, kv := range resp.Kvs {
		jobMap := models.Job{}
		if err := json.Unmarshal(kv.Value, &jobMap); err != nil {
			global.Lg.Error("json unmarshal failed", zap.Error(err))
			continue
		}
		//同步job
		jobEvent := j.BuildEvent.buildEventJob(global.JOB_EVENT_TYPE_PUT, &jobMap, strings.TrimPrefix(string(kv.Key), global.CRON_SAVE_PATH))
		j.Schedule.PushJob(jobEvent)
	}
	go func() {
		//从get时刻的后续版本开始监听
		watchStartRevision := resp.Header.Revision + 1
		//监听后续的变化
		watchChan := j.WatcherOp.Watch(ctx, j.SavePrefix, clientv3.WithRev(watchStartRevision), clientv3.WithPrefix())
		for watchResp := range watchChan {
			for _, event := range watchResp.Events {
				eventKey := strings.TrimPrefix(string(event.Kv.Key), global.CRON_SAVE_PATH)
				opJob := &models.Job{}
				var eventType global.JOB_EVENT_TYPE
				switch event.Type {
				case clientv3.EventTypeDelete: //删除
					eventType = global.JOB_EVENT_TYPE_DELETE
				case clientv3.EventTypePut: //更新
					eventType = global.JOB_EVENT_TYPE_PUT
					if err := json.Unmarshal(event.Kv.Value, opJob); err != nil {
						continue
					}
				}
				//构建jobEvent
				jobEvent := j.BuildEvent.buildEventJob(eventType, opJob, eventKey)
				//推送给调度器
				j.Schedule.PushJob(jobEvent)
			}
		}
	}()
}
