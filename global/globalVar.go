package global

import (
	"crontab/config"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-redis/redis"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type JOB_EVENT_TYPE int

const (
	JOB_EVENT_TYPE_PUT    JOB_EVENT_TYPE = 0
	JOB_EVENT_TYPE_DELETE JOB_EVENT_TYPE = 1
	JOB_EVENT_TYPE_KILL   JOB_EVENT_TYPE = 2
)

var (
	Settings   config.ServerConfig
	Lg         *zap.Logger
	Trans      ut.Translator
	DB         *gorm.DB
	Redis      *redis.Client
	EtcdClient *clientv3.Client
	TraceId    string
	MongoCli   *mongo.Client
)

const (
	CRON_SAVE_PATH   = "/cron/jobs/"
	CRON_KILL_PATH   = "/cron/kill/"
	TM_FMT_WITH_MS   = "2006-01-02 15:04:05.000"
	CRON_LOCK_PATH   = "/cron/lock/"
	WORKER_SAVE_PATH = "/cron/works/"
)

const LoggerKey = iota
