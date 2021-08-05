package services

import (
	"context"
	"crontab/global"
	"crontab/models"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type LogSink struct {
	LogCollection *mongo.Collection
	LogChan       chan *models.JobLog
	AutoCommit    chan *LogBatch
}

type LogBatch struct {
	Logs []interface{}
}

func (s *LogSink) WriteLoop() {
	var logBatch *LogBatch
	var timer *time.Timer
	for {
		select {
		case log := <-s.LogChan:
			//此处优化的重点是 防止每次日志都请求mongodb 浪费性能，转换为批次处理日志 减少io
			if logBatch == nil {
				logBatch = &LogBatch{}
			}
			//产生新的批次的时候，给一个1秒的超时时间，到时自动提交
			timer = time.AfterFunc(time.Second*2, func(logbatch *LogBatch) func() {
				return func() {
					s.AutoCommit <- logbatch
				}
			}(logBatch))
			//把新的日志添加到批次中去
			logBatch.Logs = append(logBatch.Logs, log)
			//如果批次满了 立即发送
			if len(logBatch.Logs) > 100 {
				s.SaveLogs(logBatch)
				//保存之后清空
				logBatch = nil
				//取消定时器
				timer.Stop()
			}
		case timeoutBatch := <-s.AutoCommit: //超时批次
			//判断过期批次是否依旧是当前批次 防止自动提交和超时提交共同提交了
			if timeoutBatch != logBatch { //如果不等于就说明已经提交过了
				continue //跳过提交的
			}
			s.SaveLogs(timeoutBatch)
			logBatch = nil
		}
	}
}

//批量写入日志
func (s *LogSink) SaveLogs(batch *LogBatch) {
	_, err := s.LogCollection.InsertMany(context.TODO(), batch.Logs)
	if err != nil {
		global.Lg.Info("insert db failed", zap.Any("data", batch), zap.Error(err))
	}
}

//往队列中推送日志  这地方需要考虑队列满了的情况
func (s *LogSink) PushLog(log *models.JobLog) {
	s.LogChan <- log
}
