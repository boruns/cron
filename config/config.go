package config

type ServerConfig struct {
	Name        string        `mapstructure:"name"`
	Port        int           `mapstructure:"port"`
	Locale      string        `mapstructure:"locale"`
	LogsAddress string        `mapstructure:"logsAddress"`
	MysqlInfo   MysqlConfig   `mapstructure:"mysql"`
	RedisInfo   RedisConfig   `mapstructure:"redis"`
	EtcdInfo    EtcdConfig    `mapstructure:"etcd"`
	JwtInfo     JwtConfig     `mapstructure:"jwt"`
	MinioInfo   MinioConfig   `mapstructure:"minio"`
	MongodbInfo MongodbConfig `mapstructure:"mongodb"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Db       int    `mapstructure:"db"`
	Password string `mapstructure:"password"`
}

type MysqlConfig struct {
	Host     string `mapstructure:"host"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Port     int    `mapstructure:"port"`
	DBName   string `mapstructure:"dbName"`
}

type EtcdConfig struct {
	Hosts       []string `mapstructure:"hosts"`
	DialTimeout int      `mapstructure:"dialTimeout"`
}

type JwtConfig struct {
	Key string `mapstructure:"key"`
}

type MinioConfig struct {
	AccessKey string `mapstructure:"access_key"`
	SecretKey string `mapstructure:"secret_key"`
}

type MongodbConfig struct {
	Url      string `mapstructure:"url"`
	Timeout  int    `mapstructure:"timeout"`
	Database string `mapstructure:"database"`
	Collect  string `mapstructure:"collect"`
}
