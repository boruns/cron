package services

import (
	"context"
	"crontab/global"
	"crontab/models"
	"time"

	"github.com/gorhill/cronexpr"
	"go.uber.org/zap"
)

//任务事件
type JobEvent struct {
	EventType global.JOB_EVENT_TYPE
	Job       *models.Job
	JobId     string
}

//任务调度器
type JobSchedule struct {
	JobEventChan    chan *JobEvent              //etcd中的任务
	JobPlanTable    map[string]*JobSchedulePlan //任务调度计划表
	JobExecuteTable map[string]*JobExecuteInfo  //任务执行表
	JobExecutor     *Executor
	JobResultChan   chan *JobExecuteResult //任务执行结果
	LogSink         *LogSink
}

//任务调度计划
type JobSchedulePlan struct {
	Job      *models.Job
	Expr     *cronexpr.Expression //解析好的表达式
	NextTime time.Time            //执行时间
}

//任务执行计划
type JobExecuteInfo struct {
	Job       *models.Job
	PlanTime  time.Time //理论上任务的调度时间
	RealTime  time.Time //实际的调度时间
	Ctx       context.Context
	CancelFun context.CancelFunc
}

//生成正在执行的调度信息
func (s *JobSchedule) BuildJobExecuteInfo(jobPlan *JobSchedulePlan) *JobExecuteInfo {
	ctx, cancelFunc := context.WithCancel(context.TODO())
	return &JobExecuteInfo{
		Job:       jobPlan.Job,
		PlanTime:  jobPlan.NextTime,
		RealTime:  time.Now(),
		Ctx:       ctx,
		CancelFun: cancelFunc,
	}
}

//任务执行结果
type JobExecuteResult struct {
	ExecuteInfo *JobExecuteInfo //执行状态
	Output      []byte          //脚本输出
	Err         error           //脚本错误
	JobId       string          //脚本id
	StartTime   time.Time       //脚本真实的启动时间
	EndTime     time.Time       //脚本的结束时间
}

//调度任务到执行器
func (s *JobSchedule) TryStartJob(jobId string, jobPlan *JobSchedulePlan) {
	//调度和执行任务
	//判断是否正在执行
	if _, ok := s.JobExecuteTable[jobId]; ok {
		//如果正在执行  跳过
		// global.Lg.Info("任务尚未退出，跳过执行", zap.String("jobId", jobId))
		return
	}
	jobExecuteInfo := s.BuildJobExecuteInfo(jobPlan)
	s.JobExecuteTable[jobId] = jobExecuteInfo
	//执行任务
	// fmt.Printf("计划任务：%#v  计划开始时间: %#v  计划执行时间: %#v\n", jobExecuteInfo.Job, jobExecuteInfo.PlanTime, jobExecuteInfo.RealTime)
	s.JobExecutor.ExecuteJob(s, jobId, jobExecuteInfo)
}

//构建调度计划
func (s *JobSchedule) BuildJobSchedulePlan(job *models.Job) (*JobSchedulePlan, error) {
	expr, err := cronexpr.Parse(job.CronExpire)
	if err != nil {
		return nil, err
	}
	return &JobSchedulePlan{
		Job:      job,
		Expr:     expr,
		NextTime: expr.Next(time.Now()),
	}, nil
}

//开始调度计划
func (s *JobSchedule) TryScheduleJob() time.Duration {
	now := time.Now()
	var nearTime *time.Time
	if len(s.JobPlanTable) == 0 {
		return time.Second * 1
	}
	//遍历所有任务
	for jobId, jobPlan := range s.JobPlanTable {
		if jobPlan.NextTime.Before(now) || jobPlan.NextTime.Equal(now) {
			//TODO尝试执行任务
			s.TryStartJob(jobId, jobPlan)
			jobPlan.NextTime = jobPlan.Expr.Next(now) //更新下次执行时间
		}
		//统计最近要过期的任务时间
		if nearTime == nil || jobPlan.NextTime.Before(*nearTime) {
			nearTime = &jobPlan.NextTime
		}
	}

	//下次调度间隔
	return (*nearTime).Sub(now)
}

//处理任务事件
func (s *JobSchedule) HandlerJobEvent(jobEvent *JobEvent) {
	switch jobEvent.EventType {
	case global.JOB_EVENT_TYPE_PUT:
		jobPlan, err := s.BuildJobSchedulePlan(jobEvent.Job)
		if err != nil {
			global.Lg.Error("buildJobShedulePlan failed", zap.Error(err), zap.String("jobId", jobEvent.JobId))
			return
		}
		s.JobPlanTable[jobEvent.JobId] = jobPlan
	case global.JOB_EVENT_TYPE_DELETE:
		global.Lg.Info("delete job", zap.String("jobId:", jobEvent.JobId))
		delete(s.JobPlanTable, jobEvent.JobId)
	case global.JOB_EVENT_TYPE_KILL:
		//如果存在正在执行的任务
		if event, ok := s.JobExecuteTable[jobEvent.JobId]; ok {
			global.Lg.Info("kill job", zap.String("jobId", jobEvent.JobId))
			event.CancelFun()
		}
	}
}

//任务调度
func (s *JobSchedule) schedulerLoop() {
	var scheduleAfter time.Duration
	scheduleAfter = s.TryScheduleJob()
	scheduleTimer := time.NewTimer(scheduleAfter)

	for {
		select {
		case jobEvent := <-s.JobEventChan:
			s.HandlerJobEvent(jobEvent)
		case <-scheduleTimer.C: //最近的任务到期了
		case jobResult := <-s.JobResultChan: //监听任务执行结果
			s.HandleJobResult(jobResult)
		}
		//调度一次任务
		scheduleAfter = s.TryScheduleJob()
		//重置定时器
		scheduleTimer.Reset(scheduleAfter)
	}
}

//推送任务
func (s *JobSchedule) PushJob(jobEvent *JobEvent) {
	s.JobEventChan <- jobEvent
}

//回传任务执行结果
func (s *JobSchedule) PushJobResult(jobResult *JobExecuteResult) {
	s.JobResultChan <- jobResult
}

//处理任务结果
func (s *JobSchedule) HandleJobResult(jobResult *JobExecuteResult) {
	//删除执行状态
	delete(s.JobExecuteTable, jobResult.JobId)
	err := ""
	if jobResult.Err != nil {
		err = jobResult.Err.Error()
	}
	jobLog := &models.JobLog{
		JobId:        jobResult.JobId,
		JobName:      jobResult.ExecuteInfo.Job.Name,
		Command:      jobResult.ExecuteInfo.Job.Command,
		Err:          err,
		Output:       string(jobResult.Output),
		PlanTime:     jobResult.ExecuteInfo.PlanTime.UnixNano() / 1000 / 1000,
		ScheduleTime: jobResult.ExecuteInfo.RealTime.UnixNano() / 1000 / 1000,
		StartTime:    jobResult.StartTime.UnixNano() / 1000 / 1000,
		EndTime:      jobResult.EndTime.UnixNano() / 1000 / 1000,
	}
	s.LogSink.PushLog(jobLog)
}
