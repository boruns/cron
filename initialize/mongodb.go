package initialize

import (
	"context"
	"crontab/global"
	"time"

	"github.com/fatih/color"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func InitMongodb() {
	clientOption := options.Client().
		ApplyURI(global.Settings.MongodbInfo.Url).
		SetConnectTimeout(time.Duration(global.Settings.MongodbInfo.Timeout) * time.Millisecond)
	mongoCli, err := mongo.Connect(context.TODO(), clientOption)
	if err != nil {
		color.Red("client mongodb failed")
		color.Yellow("error: %s\n", err.Error())
		panic(err)
	}
	//检查连接
	if err := mongoCli.Ping(context.TODO(), nil); err != nil {
		panic(err)
	}
	global.MongoCli = mongoCli
}
