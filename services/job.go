package services

import (
	"context"
	"crontab/dao"
	"crontab/global"
	"crontab/models"
	"crontab/utils"
	"encoding/json"
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
)

type JobServiceInterface interface {
	SaveJob(ctx *gin.Context, job *models.Job) (string, error)
	DestoryJob(ctx *gin.Context, jobId string) error
	JobList(ctx *gin.Context, page, pageSize int) ([]models.JobModel, int64)
	KillJob(ctx *gin.Context, jobId string) bool
	UpdateJob(ctx *gin.Context, jobId string, job *models.Job) error
	Show(ctx *gin.Context, jobId string) (*models.JobModel, error)
	WorkNode(ctx *gin.Context) ([]string, error)
}

type jobService struct {
	KvOp       clientv3.KV
	LeaseOp    clientv3.Lease
	SavePrefix string
	KillPrefix string
}

func NewJobServiceHandler() JobServiceInterface {
	return &jobService{
		KvOp:       clientv3.NewKV(global.EtcdClient),
		LeaseOp:    clientv3.NewLease(global.EtcdClient),
		SavePrefix: global.CRON_SAVE_PATH,
		KillPrefix: global.CRON_KILL_PATH,
	}
}

//保存数据
func (j *jobService) SaveJob(ctx *gin.Context, job *models.Job) (string, error) {
	jobDao := dao.NewJobHandler()
	db := global.DB.Begin()

	jobId, err := jobDao.SaveJob(db, job)
	if err != nil {
		db.Rollback()
		return "", err
	}
	jobStr, err := json.Marshal(job)
	if err != nil {
		db.Rollback()
		return "", err
	}
	//还需保存到etcd中
	if _, err := j.KvOp.Put(ctx, j.SavePrefix+jobId, string(jobStr)); err != nil {
		db.Rollback()
		return "", err
	}
	db.Commit()
	return jobId, nil
}

//删除数据
func (j *jobService) DestoryJob(ctx *gin.Context, jobId string) error {
	jobDao := dao.NewJobHandler()
	db := global.DB.Begin()
	err := jobDao.DestoryJobByJobId(db, jobId)
	if err != nil {
		db.Rollback()
		return err
	}
	if resp, err := j.KvOp.Delete(ctx, j.SavePrefix+jobId); err != nil || resp.Deleted < 1 {
		db.Rollback()
		if err != nil {
			return err
		}
		return errors.New("job deleted failed in etcd")
	}
	db.Commit()
	return nil
}

//获取job列表
func (j *jobService) JobList(ctx *gin.Context, page, pageSize int) ([]models.JobModel, int64) {
	jobDao := dao.NewJobHandler()
	total, ret := jobDao.GetJobList(page, pageSize)
	return ret, total
}

//强杀job
func (j *jobService) KillJob(ctx *gin.Context, jobId string) bool {
	leaseResp, err := j.LeaseOp.Grant(context.TODO(), 1)
	if err != nil {
		utils.LoggerFor(ctx).Info("create etcd lease failed", zap.String("error: ", err.Error()))
		return false
	}
	leaseId := leaseResp.ID
	if _, err := j.KvOp.Put(ctx, j.KillPrefix+jobId, "", clientv3.WithLease(leaseId)); err != nil {
		utils.LoggerFor(ctx).Info("put kill key failed", zap.String("error: ", err.Error()))
		return false
	}
	return true
}

func (j *jobService) UpdateJob(ctx *gin.Context, jobId string, job *models.Job) error {
	jobDao := dao.NewJobHandler()
	db := global.DB.Begin()
	jobData := map[string]interface{}{
		"name":        job.Name,
		"command":     job.Command,
		"cron_expire": job.CronExpire,
	}
	if err := jobDao.UpdateJobByJobId(db, jobId, jobData); err != nil {
		db.Rollback()
		return err
	}
	jobStr, err := json.Marshal(job)
	if err != nil {
		db.Rollback()
		return err
	}
	if _, err := j.KvOp.Put(ctx, j.SavePrefix+jobId, string(jobStr)); err != nil {
		db.Rollback()
		return err
	}
	db.Commit()
	return nil
}

func (j *jobService) Show(ctx *gin.Context, jobId string) (*models.JobModel, error) {
	jobDao := dao.NewJobHandler()
	return jobDao.GetJobByJobId(jobId)
}

func (j *jobService) WorkNode(ctx *gin.Context) ([]string, error) {
	resp, err := j.KvOp.Get(ctx, global.WORKER_SAVE_PATH, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}
	result := make([]string, 0)
	for _, work := range resp.Kvs {
		key := strings.TrimPrefix(string(work.Key), global.WORKER_SAVE_PATH)
		result = append(result, key)
	}
	return result, nil
}
