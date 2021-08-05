package services

import (
	"context"
	"crontab/global"
	"errors"
	"math/rand"
	"os/exec"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
)

//任务执行器
type Executor struct {
	KvOp     clientv3.KV
	LeaseOp  clientv3.Lease
	LockPath string
}

type CancelLock struct {
	CancelFunc context.CancelFunc
	LeaseId    clientv3.LeaseID
	JobId      string
	IsLock     bool
}

func (e *Executor) ExecuteJob(schedule *JobSchedule, jobId string, info *JobExecuteInfo) {
	go func() {
		result := &JobExecuteResult{}
		result.StartTime = time.Now()
		result.Output = make([]byte, 0)
		result.JobId = jobId
		result.ExecuteInfo = info
		//随机睡眠0-1秒
		time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
		lock, err := e.TryLock(jobId)
		defer e.Unlock(lock)
		if err != nil {
			result.Err = err
			result.EndTime = time.Now()
		} else {
			result.StartTime = time.Now()
			//执行shell命令
			cmd := exec.CommandContext(info.Ctx, "/bin/sh", "-c", info.Job.Command)
			//执行并捕获输出
			ouput, err := cmd.CombinedOutput()
			result.EndTime = time.Now()
			result.Output = ouput
			result.Err = err
		}
		schedule.PushJobResult(result)
	}()
}

//抢占分布式锁
func (e *Executor) TryLock(jobId string) (*CancelLock, error) {

	//创建租约5秒
	grantResp, err := e.LeaseOp.Grant(context.TODO(), 5)
	if err != nil {
		global.Lg.Info("create lease grant failed", zap.Error(err), zap.String("jobId", jobId))
		return nil, err
	}
	//自动续约
	//用户取消锁
	cancelCtx, cancelFunc := context.WithCancel(context.TODO())
	//租约id
	leaseId := grantResp.ID
	_, err = e.LeaseOp.KeepAlive(cancelCtx, leaseId)

	//处理续租成功的应答协程 (暂不处理)

	//失败时候需要调用的函数
	errFunc := func() {
		e.LeaseOp.Revoke(context.TODO(), leaseId) //释放租约
		cancelFunc()                              //取消自动续租
	}

	//续租是否成功
	if err != nil {
		global.Lg.Info("grant create keepalive failed", zap.Error(err), zap.String("jobId", jobId))
		errFunc()
		return nil, err
	}

	//创建事务txn
	txn := e.KvOp.Txn(context.TODO())
	lockName := e.LockPath + jobId
	txn.If(clientv3.Compare(clientv3.CreateRevision(lockName), "=", 0)).
		Then(clientv3.OpPut(lockName, "", clientv3.WithLease(leaseId))).
		Else(clientv3.OpGet(lockName))

	//提交事务
	txnResp, err := txn.Commit()
	if err != nil {
		global.Lg.Info("tnx commit failed", zap.Error(err), zap.String("jobId", jobId))
		errFunc()
		return nil, err
	}

	//成功返回 失败释放租约
	if !txnResp.Succeeded {
		errFunc()
		return nil, errors.New("lock job failed")
	}
	cancelLock := &CancelLock{
		CancelFunc: cancelFunc,
		LeaseId:    leaseId,
		JobId:      jobId,
		IsLock:     true,
	}
	//上锁成功
	return cancelLock, nil
}

//释放锁
func (e *Executor) Unlock(lock *CancelLock) {
	//上锁成功才会取消
	if lock.IsLock {
		lock.CancelFunc()
		e.LeaseOp.Revoke(context.TODO(), lock.LeaseId)
	}
}
