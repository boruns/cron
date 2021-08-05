package services

import (
	"crontab/global"
	"crontab/models"
	"crontab/utils"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

type JobLogServiceInterface interface {
	GetLog(ctx *gin.Context, jobId string, page, pageSize int) (int64, []*models.JobLog, error)
}

type jobLog struct {
	Collection *mongo.Collection
}

type JobLogFilter struct {
	JobId string `bson:"job_id"`
}

type JobLogSort struct {
	SortOrder int `bson:"start_time"`
}

func NewJobLogService() JobLogServiceInterface {
	return &jobLog{
		Collection: utils.GetMongoDbDatabase().Collection("jobs"),
	}
}

//获取任务日志
func (l *jobLog) GetLog(ctx *gin.Context, jobId string, page, pageSize int) (int64, []*models.JobLog, error) {
	logFilter := &JobLogFilter{
		JobId: jobId,
	}
	// //按照任务开始时间倒叙
	logSort := &JobLogSort{
		SortOrder: -1,
	}
	skip := int64(page - 1)
	limit := int64(pageSize)
	total, err := l.Collection.CountDocuments(ctx, logFilter)
	if err != nil {
		return 0, nil, err
	}
	cursor, err := l.Collection.Find(ctx, logFilter, &options.FindOptions{
		Skip:  &skip,
		Limit: &limit,
		Sort:  logSort,
	})
	if err != nil {
		utils.LoggerFor(ctx).Info("get log failed", zap.String("jobId", jobId), zap.Error(err))
		return 0, nil, err
	}
	data := make([]*models.JobLog, 0)
	//延迟释放游标
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		jobLog := &models.JobLog{}
		if err := cursor.Decode(jobLog); err != nil {
			global.Lg.Info("bson to josn failed", zap.String("jobId", jobId), zap.Error(err))
			continue
		}
		data = append(data, jobLog)
	}
	return total, data, nil
}
