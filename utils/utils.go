package utils

import (
	"context"
	"crontab/global"
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func GetNowFormatTodayTime() string {
	now := time.Now()
	dateStr := fmt.Sprintf("%02d-%02d-%02d", now.Year(), int(now.Month()),
		now.Day())
	return dateStr
}

func GetNow() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func GetUUID() (string, error) {
	uId, err := uuid.NewV4()
	if err != nil {
		return "", err
	}
	return uId.String(), nil
}

// 给指定的context添加字段（关键方法）
func LoggerNewContext(ctx *gin.Context, fields ...zapcore.Field) {
	ctx.Set(strconv.Itoa(global.LoggerKey), LoggerFor(ctx).With(fields...))
}

// 从指定的context返回一个zap实例（关键方法）
func LoggerFor(ctx *gin.Context) *zap.Logger {
	if ctx == nil {
		return global.Lg
	}
	l, _ := ctx.Get(strconv.Itoa(global.LoggerKey))
	ctxLogger, ok := l.(*zap.Logger)
	if ok {
		return ctxLogger
	}
	return global.Lg
}

func GetMongoDbDatabase() *mongo.Database {
	return global.MongoCli.Database(global.Settings.MongodbInfo.Database)

}

//分页
func Find(database *mongo.Database, collection string, filter interface{}, limit, index int64) (data []map[string]interface{}, err error) {
	ctx, cannel := context.WithTimeout(context.Background(), time.Minute)
	defer cannel()
	var findoptions *options.FindOptions
	if limit > 0 {
		findoptions.SetLimit(limit)
		findoptions.SetSkip(limit * index)
	}
	cur, err := database.Collection(collection).Find(ctx, filter, findoptions)
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.Background())
	err = cur.All(context.Background(), &data)
	return
}
